package main

import "os"
import "os/exec"
import "fmt"

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
	// dd if=/opt/src.des3 | openssl des3 -d -k "\$PASS" | tar zxf - -C /opt

	// command 3
	// mv /opt/src /path/to/target
}
