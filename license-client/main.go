package main

import (
	"bytes"
	"fmt"
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

func main() {
	fmt.Println("DEBUG | LICENSE_SERVER: ", os.Getenv("LICENSE_SERVER"))
	licenseServer := os.Getenv("LICENSE_SERVER")

	// command 1
	// curl --cert ./data/pki/client.cert.pem --key ./data/pki/client.key.pem --cacert ./data/pki/ca.cert.pem ${LICENSE_SERVER}
	// []string{"curl", "--cert", "./data/pki/client.cert.pem", "--key", "./data/pki/client.key.pem", "--cacert", "./data/pki/ca.cert.pem", os.Getenv("LICENSE_SERVER")}
	curlCmd := exec.Command("curl", "--cert", "./data/pki/client.cert.pem", "--key", "./data/pki/client.key.pem", "--cacert", "./data/pki/ca.cert.pem", licenseServer)

	curlOut, curlErr := curlCmd.Output()
	if curlErr != nil {
		panic(curlErr)
	}
	fmt.Println("DEBUG | curlOut: ", string(curlOut))

	// command 2
	// dd if=./data/src.des3 | openssl des3 -d -k "${PASS}" | tar zxf - -C /tmp
	// []string{"dd", "if=./data/src.des3"}
	// []string{"openssl", "des3", "-d", "-k", string(curlOut)}
	// []string{"tar", "zxf", "-", "-C", "/tmp"}
	ddCmd := exec.Command("dd", "if=./data/src.des3")
	opensslCmd := exec.Command("openssl", "des3", "-d", "-k", string(curlOut))
	tarCmd := exec.Command("tar", "zxf", "-", "-C", "/tmp")

	// Run the pipeline
	output, stderr, err := Pipeline(ddCmd, opensslCmd, tarCmd)
	if err != nil {
		log.Printf("%s", err)
	}

	// Print the stdout, if any
	if len(output) > 0 {
		log.Printf("%s", output)
	}

	// Print the stderr, if any
	if len(stderr) > 0 {
		log.Printf("%s", stderr)
	}

	// command 3
	// mv /opt/src /path/to/target
	// []string{"mv", "/opt/src", "/path/to/target"}
}
