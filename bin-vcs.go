package main

import (
	"os"
)

// The purpose of the main function here is take the command line input
// and dispatch it to the vcs-common resolver.
func main() {
	// Pull all the relevant arguments after the "./bin-vcs" and sent it
	// to the resolve function.
	Resolve(os.Args[1:])
}
