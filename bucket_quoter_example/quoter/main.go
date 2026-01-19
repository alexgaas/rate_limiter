package main

import (
	"os"
	"quoter/cmd"
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
