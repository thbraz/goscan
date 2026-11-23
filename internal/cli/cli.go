package cli

import (
	"flag"
	"fmt"
	"io"
	"slices"
	"strings"
)

const (
	binName = "goscan"
	defaultHost = "localhost"
)

type Args struct {
	Host string
}

type UnexpectedArgsError struct {
	LeftoverArgs []string
}

func (e *UnexpectedArgsError) Error() string {
    return fmt.Sprintf("unexpected arguments: %s", strings.Join(e.LeftoverArgs, ", "))
}

func ParseArgs(args []string, writer io.Writer) (Args, error) {
	flagSet := flag.NewFlagSet(binName, flag.ContinueOnError)
	flagSet.SetOutput(writer)
	
	host := flagSet.String("t", defaultHost, "Target specification")

	err := flagSet.Parse(args)
	if err != nil {
		return Args{}, err
	}

	if flagSet.NArg() > 0 {
		flagSet.Usage()
		return Args{}, &UnexpectedArgsError{LeftoverArgs: slices.Clone(flagSet.Args())}
	}

	return Args{Host: *host}, nil
}