package main

import (
	"fmt"
	"os"
)

const (
	usageFmt = "%s proxy.pac url...\n"
)

const (
	_ = iota
	argumentErr
	runtimeErr
)

func main() {
	os.Exit(run())
}

func run() (exitCode int) {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, usageFmt, os.Args[0])
		exitCode = argumentErr
		return
	}

	proxyPac := os.Args[1]
	urls := os.Args[2:]

	err := process(proxyPac, urls)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		exitCode = runtimeErr
	}

	return
}
