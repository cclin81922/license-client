package main

import "syscall"
import "os"
import "os/exec"
import "fmt"

func main() {
    fmt.Println("LICENSE_SERVER: ", os.Getenv("LICENSE_SERVER"))

    binary, lookErr := exec.LookPath("curl")
    if lookErr != nil {
        panic(lookErr)
    }

    args := []string{"curl", "--cert", "/tmp/pki/client.cert.pem", "--key", "/tmp/pki/client.key.pem", "--cacert", "/tmp/pki/ca.cert.pem", os.Getenv("LICENSE_SERVER")}

    env := os.Environ()

    execErr := syscall.Exec(binary, args, env)
    if execErr != nil {
        panic(execErr)
    }
}
