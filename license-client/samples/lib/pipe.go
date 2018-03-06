package lib

import (
	"bytes"
	"log"
	"os"
	"os/exec"
)

// Pipeline strings together the given exec.Cmd commands in a similar fashion
// to the Unix pipeline.  Each command's standard output is connected to the
// standard input of the next command, and the output of the final command in
// the pipeline is returned, along with the collected standard error of all
// commands and the first error found (if any).
//
// To provide input to the pipeline, assign an io.Reader to the first's Stdin.
func Pipeline(cmds ...*exec.Cmd) (pipeLineOutput, collectedStandardError []byte, pipeLineError error) {
	// Require at least one command
	if len(cmds) < 1 {
		return nil, nil, nil
	}

	// Collect the output from the command(s)
	var output bytes.Buffer
	var stderr bytes.Buffer

	last := len(cmds) - 1
	for i, cmd := range cmds[:last] {
		var err error
		// Connect each command's stdin to the previous command's stdout
		if cmds[i+1].Stdin, err = cmd.StdoutPipe(); err != nil {
			return nil, nil, err
		}
		// Connect each command's stderr to a buffer
		cmd.Stderr = &stderr
	}

	// Connect the output and error for the last command
	cmds[last].Stdout, cmds[last].Stderr = &output, &stderr

	// Start each command
	for _, cmd := range cmds {
		if err := cmd.Start(); err != nil {
			return output.Bytes(), stderr.Bytes(), err
		}
	}

	// Wait for each command to complete
	for _, cmd := range cmds {
		if err := cmd.Wait(); err != nil {
			return output.Bytes(), stderr.Bytes(), err
		}
	}

	// Return the pipeline output and the collected standard error
	return output.Bytes(), stderr.Bytes(), nil
}

func Pipe() {
	// Collect directories from the command-line
	var dirs []string
	if len(os.Args) > 1 {
		dirs = os.Args[1:]
	} else {
		dirs = []string{"."}
	}

	// Run the command on each directory
	for _, dir := range dirs {
		// find $DIR -type f # Find all files
		ls := exec.Command("find", dir, "-type", "f")

		// | grep -v '/[._]' # Ignore hidden/temporary files
		visible := exec.Command("egrep", "-v", `/[._]`)

		// | sort -t. -k2 # Sort by file extension
		sort := exec.Command("sort", "-t.", "-k2")

		// Run the pipeline
		output, stderr, err := Pipeline(ls, visible, sort)
		if err != nil {
			log.Printf("dir %q: %s", dir, err)
		}

		// Print the stdout, if any
		if len(output) > 0 {
			log.Printf("%q:\n%s", dir, output)
		}

		// Print the stderr, if any
		if len(stderr) > 0 {
			log.Printf("%q: (stderr)\n%s", dir, stderr)
		}
	}
}
