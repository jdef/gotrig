package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/jdef/gotrig"
)

func usage() {
	fmt.Fprintf(os.Stderr, "listen fifodir")
	os.Exit(1)
}

func main() {
	if len(os.Args) < 2 {
		usage()
	}

	x, err := gotrig.Listen(os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "listen: %v", err)
		os.Exit(1)
	}

	sigch := make(chan os.Signal)
	signal.Notify(sigch, syscall.SIGTERM, syscall.SIGINT, syscall.SIGHUP)

	rc := make(chan int, 1)

	go func() {
		exitCode := 0
		defer func() {
			rc <- exitCode
		}()
		for y := range x.B {
			fmt.Fprintf(os.Stderr, "received: %c\n", y)
		}
		if x.Err != nil {
			fmt.Fprintf(os.Stderr, "error: %v", x.Err)
			exitCode = 1
		}
	}()

	// wait for stop signal then terminate cleanly
	<-sigch
	x.Close()
	os.Exit(<-rc)
}
