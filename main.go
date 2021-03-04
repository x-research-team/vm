package main

import (
	"fmt"
	"os"

	"github.com/x-research-team/vm/internal/repl"
	"github.com/x-research-team/vm/vm"
)

func main() {
	args := os.Args[1:]
	//We must reset `os.Args`, or the `flag` module will not functioning correctly
	os.Args = os.Args[1:]
	if len(args) == 0 {
		fmt.Println("VM programming language REPL\n")
		repl.Start(os.Stdout, true)
	} else {
		if len(args) == 2 {
			if args[0] == "-d" || args[0] == "--debug" { // debug
				vm.RunProgram(true, args[1])
			} else {
				fmt.Println("Usage: vm -d file.vm")
				os.Exit(1)
			}
		} else {
			vm.RunProgram(false, args[0])
		}

	}
}
