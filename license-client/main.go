package main

import "os"
import "fmt"

func main() {
	fmt.Println("LICENSE_SERVER: ", os.Getenv("LICENSE_SERVER"))
}
