package main

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/thbraz/goscan/internal/cli"
)

func main() {
	args, err := cli.ParseArgs(os.Args[1:], os.Stderr)
	if err != nil {
		os.Exit(exitCodeFor(err))
	}

	fmt.Printf("Scanning: %s\n", args.Host)
}

func exitCodeFor(err error) int {
	if errors.Is(err, flag.ErrHelp) {
		return 0
	}

	var unexpected *cli.UnexpectedArgsError
	if errors.As(err, &unexpected) {
		fmt.Fprintln(os.Stderr, err)
	}

	return 2
}