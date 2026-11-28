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

	err := flagSet.Parse(args)
	if err != nil {
		return Args{}, err
	}

	if flagSet.NArg() != 1 {
		flagSet.Usage()
		return Args{}, &UnexpectedArgsError{LeftoverArgs: slices.Clone(flagSet.Args())}
	}

	host := flagSet.Arg(0)
	return Args{Host: host}, nil
}