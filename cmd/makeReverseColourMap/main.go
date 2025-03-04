package main

import (
	"fmt"
	"os"
)

func main() {
	prog := NewProg()

	if len(os.Args) != 2 { //nolint:mnd
		fmt.Printf("Usage: %s filename", os.Args[0])
		os.Exit(1)
	}

	prog.filename = os.Args[1]

	prog.Run()
	os.Exit(prog.exitStatus)
}
