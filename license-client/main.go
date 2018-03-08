package main

import (
	"bytes"
	"log"
	"os"
	"os/exec"
)

func Activate() {
	log.Println("INFO | Activating")

	// TODO
	// mv /opt/src /path/to/target
	// []string{"mv", "/opt/src", "/path/to/target"}
}

func Destroy() {
	log.Println("INFO | Destroying")

	// TODO

	os.Exit(1)
}

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
	log.Println("INFO | Start")

	// Get license server location
	licenseServer := os.Getenv("LICENSE_SERVER")
	if len(licenseServer) == 0 {
		log.Println("ERROR | Missing LICENSE_SERVER")
		Destroy()
	}
	log.Println("DEBUG | licenseServer:", licenseServer)

	// Run licensing step 1
	//
	// shell: curl --cert ./data/pki/client.cert.pem --key ./data/pki/client.key.pem --cacert ./data/pki/ca.cert.pem ${LICENSE_SERVER}
	// golang: []string{"curl", "--cert", "./data/pki/client.cert.pem", "--key", "./data/pki/client.key.pem", "--cacert", "./data/pki/ca.cert.pem", os.Getenv("LICENSE_SERVER")}
	curlCmd := exec.Command("curl", "--cert", "./data/pki/client.cert.pem", "--key", "./data/pki/client.key.pem", "--cacert", "./data/pki/ca.cert.pem", licenseServer)

	curlOut, curlErr := curlCmd.Output()
	if curlErr != nil {
		log.Printf("ERROR | curlErr: %s\n", curlErr)
		Destroy()
	}
	dataSecret := string(curlOut)
	log.Println("DEBUG | dataSecret:", dataSecret)

	// Run licensing step 2
	//
	// shell: dd if=./data/src.des3 | openssl des3 -d -k "${PASS}" | tar zxf - -C /tmp
	// golang: []string{"dd", "if=./data/src.des3"}
	// golang: []string{"openssl", "des3", "-d", "-k", string(curlOut)}
	// golang: []string{"tar", "zxf", "-", "-C", "/tmp"}
	ddCmd := exec.Command("dd", "if=./data/src.des3")
	opensslCmd := exec.Command("openssl", "des3", "-d", "-k", dataSecret)
	tarCmd := exec.Command("tar", "zxf", "-", "-C", "/tmp")

	decryptOut, decryptStderr, decryptErr := Pipeline(ddCmd, opensslCmd, tarCmd)
	if decryptErr != nil {
		log.Printf("ERROR | decryptErr: %s\n", decryptErr)
		Destroy()
	}

	// Print the stdout, if any
	if len(decryptOut) > 0 {
		log.Printf("DEBUG | decryptOut: %s\n", decryptOut)
	}

	// Print the stderr, if any
	if len(decryptStderr) > 0 {
		log.Printf("DEBUG | decryptStderr: %s\n", decryptStderr)
	}

	Activate()

	log.Println("INFO | End")
}
