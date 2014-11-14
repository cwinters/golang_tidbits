package main

import "io/ioutil"
import "fmt"

func main() {
	// just try to read a non-existent file and see what error we get
	_, err := ioutil.ReadFile("lol.txt")
	fmt.Printf("Returned error: %s\n", err)
}
