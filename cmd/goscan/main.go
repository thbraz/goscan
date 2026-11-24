package main

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/thbraz/goscan/internal/cli"
	"github.com/thbraz/goscan/internal/scanner"
)

func main() {
	args, err := cli.ParseArgs(os.Args[1:], os.Stderr)
	if err != nil {
		os.Exit(exitCodeFor(err))
	}

	results := scanner.Scan(args)

	for _, v := range results {
		if v.Open {
			fmt.Printf("%+v \n", v)
		}
	}

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