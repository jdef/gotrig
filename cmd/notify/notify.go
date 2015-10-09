package main

import (
	"fmt"
	"os"

	"github.com/jdef/gotrig"
)

func usage() {
	fmt.Fprintln(os.Stderr, "notify fifodir message")
	os.Exit(1)
}

func main() {
	args := os.Args
	if len(args) < 3 {
		usage()
	}
	_, err := gotrig.Notify(args[1], args[2])
	if err != nil {
		fmt.Fprintf(os.Stderr, "notify: %v", err)
		os.Exit(1)
	}
}
