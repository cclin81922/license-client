package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os/exec"
	"strings"

	"github.com/cclin81922/license-client/license-client/samples/lib"
)

func main() {
	// case 1 :: to provide an explicitly delineated command and argument array
	dateCmd := exec.Command("date")
	dateOut, dateErr := dateCmd.Output()
	if dateErr != nil {
		panic(dateErr)
	}
	fmt.Println("> date")
	fmt.Println(string(dateOut))

	// case 2 :: to pipe data to the external process on its stdin and collect the results from its stdout
	// tag: async, pipe, stdin, stdout
	grepCmd := exec.Command("grep", "hello")
	grepIn, _ := grepCmd.StdinPipe()
	grepOut, _ := grepCmd.StdoutPipe()
	grepCmd.Start()
	grepIn.Write([]byte("hello grep\ngoodbye grep"))
	grepIn.Close()
	grepBytes, _ := ioutil.ReadAll(grepOut)
	grepCmd.Wait()
	fmt.Println("> grep hello")
	fmt.Println(string(grepBytes))

	// case 3 :: to run starts the specified command and waits for it to complete
	// tag: sync
	fmt.Println("> sleep 3")
	sleepCmd := exec.Command("sleep", "3")
	sleepErr := sleepCmd.Run()
	if sleepErr != nil {
		panic(sleepErr)
	}
	fmt.Println("")

	// case 4
	// tag: stdin, stdout
	/*
	   About wcCmd.Stdin = strings.NewReader("12345")

	   type of wcCmd.Stdin is interface io.Reader which has method Read(p []byte) (n int, err error)
	   function strings.NewReader returns pointer to struct strings.Reader (i.e. *Reader)
	   struct strings.Reader has method Read(b []byte) (n int, err error) so it is a kind of interface io.Reader

	   About wcCmd.Stdout = &wcOut

	   type of wcCmd.Stdout is interface io.Writer which has method Write(p []byte) (n int, err error)
	   struct bytes.Buffer has method Write(p []byte) (n int, err error) so it is a kind of interface io.Writer

	   About wcCmd.Stdout = wcOut vs. wcCmd.Stdout = &wcOut

	   see https://stackoverflow.com/questions/13511203/why-cant-i-assign-a-struct-to-an-interface
	*/
	wcCmd := exec.Command("wc", "-c")
	wcCmd.Stdin = strings.NewReader("12345")
	var wcOut bytes.Buffer
	wcCmd.Stdout = &wcOut
	wcErr := wcCmd.Run()
	if wcErr != nil {
		panic(wcErr)
	}
	fmt.Println("> wc -c")
	fmt.Println(wcOut.String())

	// case 5 :: to pass in one command-line string
	// tag: bash -c
	lsCmd := exec.Command("bash", "-c", "ls -a -l -h")
	lsOut, dateErr := lsCmd.Output()
	if dateErr != nil {
		panic(dateErr)
	}
	fmt.Println("> ls -a -l -h")
	fmt.Println(string(lsOut))

	// case 6 :: to pipe multiple exec.Command instances
	lib.Pipe()

	fmt.Println("End")

}
