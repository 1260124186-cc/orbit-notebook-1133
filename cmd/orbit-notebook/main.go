package main

import (
	"example.com/orbit-notebook/internal/command"
	"fmt"
	"os"
)

func main() {
	if err := command.Run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}
