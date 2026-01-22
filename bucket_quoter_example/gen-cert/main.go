package main

import (
	"gen-cert/cmd"
	"os"
)

const (
	ExitCodeUnspecified = 1
)

func main() {
	err := cmd.Execute()
	if err != nil {
		os.Exit(ExitCodeUnspecified)
	}
}
