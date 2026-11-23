package cli

import (
	"bytes"
	"errors"
	"slices"
	"testing"
	"flag"
	 "io"
)

func TestParseArgsOk(t *testing.T) {
	args := []string{"-t", "127.0.0.1"}
	var buffer bytes.Buffer

	exp := Args{Host: "127.0.0.1"}
	res, err := ParseArgs(args, &buffer)
	if err != nil {
		t.Fatalf("ParseArgs(%q) unexpected return: %v", args, err)
	}

	if res != exp {
		t.Errorf("ParseArgs(%q): %+v; expected: %+v", args, res, exp)
	}

	if buffer.Len() > 0 {
		t.Errorf("ParseArgs(%q): wrote unexpected output: %q", args, buffer.String())
	}
}

func TestParseArgsOmittedFlag(t *testing.T) {
	args := []string{}
	var buffer bytes.Buffer

	exp := Args{Host: "localhost"}
	res, err := ParseArgs(args, &buffer)
	if err != nil {
		t.Fatalf("ParseArgs(%q) unexpected return: %v", args, err)
	}

	if res != exp {
		t.Errorf("ParseArgs(%q): %+v; expected: %+v", args, res, exp)
	}

	if buffer.Len() > 0 {
		t.Errorf("ParseArgs(%q): wrote unexpected output: %q", args, buffer.String())
	}
}

func TestParseArgsWithUnknownFlag(t *testing.T) {
	args := []string{"-x"}
	var buffer bytes.Buffer

	exp := Args{}
	res, err := ParseArgs(args, &buffer)

	if err == nil {
		t.Fatalf("ParseArgs(%q): expected error, got nil", args)
	}

	var unexpectedArgsErr *UnexpectedArgsError
	if errors.As(err, &unexpectedArgsErr) {
		t.Fatalf("ParseArgs(%q): expected parse error, got UnexpectedArgsError: %v", args, err)
	}

	if res != exp {
		t.Errorf("ParseArgs(%q): %+v; expected: %+v", args, res, exp)
	}

	if buffer.Len() == 0 {
		t.Errorf("ParseArgs(%q): expected usage output, got none", args)
	}
}

func TestParseArgsWithMissingFlagValue(t *testing.T) {
	args := []string{"-t"}
	var buffer bytes.Buffer

	exp := Args{}
	res, err := ParseArgs(args, &buffer)
	if err == nil {
		t.Fatalf("ParseArgs(%q): expected error, got nil", args)
	}

	var unexpectedArgsErr *UnexpectedArgsError
	if errors.As(err, &unexpectedArgsErr) {
		t.Fatalf("ParseArgs(%q): expected parse error, got UnexpectedArgsError: %v", args, err)
	}

	if res != exp {
		t.Errorf("ParseArgs(%q): %+v; expected: %+v", args, res, exp)
	}

	if buffer.Len() == 0 {
		t.Errorf("ParseArgs(%q): expected usage output, got none", args)
	}
}

func TestParseArgsWithOnlyPositionalArg(t *testing.T) {
	args := []string{"extra"}
	var buffer bytes.Buffer
	expLeftover := []string{"extra"}
	expArgs := Args{}

	res, err := ParseArgs(args, &buffer)
	if err == nil {
		t.Fatalf("ParseArgs(%q): expected error, got nil", args)
	}

	var unexpectedArgsErr *UnexpectedArgsError
	if !errors.As(err, &unexpectedArgsErr) {
		t.Fatalf("ParseArgs(%q): expected UnexpectedArgsError, got: %v", args, err)
	}

	if !slices.Equal(unexpectedArgsErr.LeftoverArgs, expLeftover) {
		t.Errorf("ParseArgs(%q): LeftoverArgs %q; expected: %q", args, unexpectedArgsErr.LeftoverArgs, expLeftover)
	}

	if res != expArgs {
		t.Errorf("ParseArgs(%q): %+v; expected: %+v", args, res, expArgs)
	}

	if buffer.Len() == 0 {
		t.Errorf("ParseArgs(%q): expected usage output, got none", args)
	}
}

func TestParseArgsWithExtraPositionalArgAfterFlag(t *testing.T) {
	args := []string{"-t", "127.0.0.1", "extra"}
	var buffer bytes.Buffer
	expLeftover := []string{"extra"}
	expArgs := Args{}

	res, err := ParseArgs(args, &buffer)
	if err == nil {
		t.Fatalf("ParseArgs(%q): expected error, got nil", args)
	}

	var unexpectedArgsErr *UnexpectedArgsError
	if !errors.As(err, &unexpectedArgsErr) {
		t.Fatalf("ParseArgs(%q): expected UnexpectedArgsError, got: %v", args, err)
	}

	if !slices.Equal(unexpectedArgsErr.LeftoverArgs, expLeftover) {
		t.Errorf("ParseArgs(%q): LeftoverArgs %q; expected: %q", args, unexpectedArgsErr.LeftoverArgs, expLeftover)
	}

	if res != expArgs {
		t.Errorf("ParseArgs(%q): %+v; expected: %+v", args, res, expArgs)
	}

	if buffer.Len() == 0 {
		t.Errorf("ParseArgs(%q): expected usage output, got none", args)
	}
}

func TestParseArgsWithExtraPositionalArgBeforeFlag(t *testing.T) {
	args := []string{"extra", "-t", "127.0.0.1"}
	var buffer bytes.Buffer
	expLeftover := []string{"extra", "-t", "127.0.0.1"}
	expArgs := Args{}

	res, err := ParseArgs(args, &buffer)
	if err == nil {
		t.Fatalf("ParseArgs(%q): expected error, got nil", args)
	}

	var unexpectedArgsErr *UnexpectedArgsError
	if !errors.As(err, &unexpectedArgsErr) {
		t.Fatalf("ParseArgs(%q): expected UnexpectedArgsError, got: %v", args, err)
	}

	if !slices.Equal(unexpectedArgsErr.LeftoverArgs, expLeftover) {
		t.Errorf("ParseArgs(%q): LeftoverArgs %q; expected: %q", args, unexpectedArgsErr.LeftoverArgs, expLeftover)
	}

	if res != expArgs {
		t.Errorf("ParseArgs(%q): %+v; expected: %+v", args, res, expArgs)
	}

	if buffer.Len() == 0 {
		t.Errorf("ParseArgs(%q): expected usage output, got none", args)
	}
}

func TestParseArgsWithTerminator(t *testing.T) {
	args := []string{"--", "-t", "127.0.0.1"}
	var buffer bytes.Buffer
	expLeftover := []string{"-t", "127.0.0.1"}
	expArgs := Args{}

	res, err := ParseArgs(args, &buffer)
	if err == nil {
		t.Fatalf("ParseArgs(%q): expected error, got nil", args)
	}

	var unexpectedArgsErr *UnexpectedArgsError
	if !errors.As(err, &unexpectedArgsErr) {
		t.Fatalf("ParseArgs(%q): expected UnexpectedArgsError, got: %v", args, err)
	}

	if !slices.Equal(unexpectedArgsErr.LeftoverArgs, expLeftover) {
		t.Errorf("ParseArgs(%q): LeftoverArgs %q; expected: %q", args, unexpectedArgsErr.LeftoverArgs, expLeftover)
	}

	if res != expArgs {
		t.Errorf("ParseArgs(%q): %+v; expected: %+v", args, res, expArgs)
	}

	if buffer.Len() == 0 {
		t.Errorf("ParseArgs(%q): expected usage output, got none", args)
	}
}

func TestParseArgsLeftoverArgsDoesNotAliasInput(t *testing.T) {
	args := []string{"-t", "127.0.0.1", "extra"}

	_, err := ParseArgs(args, io.Discard)

	var unexpectedArgsErr *UnexpectedArgsError
	if !errors.As(err, &unexpectedArgsErr) {
		t.Fatalf("ParseArgs(%q): expected UnexpectedArgsError, got: %v", args, err)
	}

	args[2] = "mutated"

	if got := unexpectedArgsErr.LeftoverArgs[0]; got != "extra" {
		t.Errorf("LeftoverArgs aliases input slice: got %q after mutating input", got)
	}
}

func TestParseArgsHelp(t *testing.T) {
	for _, args := range [][]string{{"-h"}, {"-help"}} {
		var buffer bytes.Buffer

		res, err := ParseArgs(args, &buffer)

		if !errors.Is(err, flag.ErrHelp) {
			t.Errorf("ParseArgs(%q): expected flag.ErrHelp, got: %v", args, err)
		}
		if res != (Args{}) {
			t.Errorf("ParseArgs(%q): %+v; expected zero Args", args, res)
		}
		if buffer.Len() == 0 {
			t.Errorf("ParseArgs(%q): expected usage output, got none", args)
		}
	}
}