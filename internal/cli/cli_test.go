package cli

import (
	"bytes"
	"errors"
	"flag"
	"testing"
)

func TestParseArgsOk(t *testing.T) {
	tests := []struct {
		name string
		args []string
		exp  Args
	}{
		{
			name: "flag provided",
			args: []string{"127.0.0.1"},
			exp:  Args{Host: "127.0.0.1"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buffer bytes.Buffer

			res, err := ParseArgs(tt.args, &buffer)
			if err != nil {
				t.Fatalf("ParseArgs(%q) unexpected return: %v", tt.args, err)
			}

			if res != tt.exp {
				t.Errorf("ParseArgs(%q): %+v; expected: %+v", tt.args, res, tt.exp)
			}

			if buffer.Len() > 0 {
				t.Errorf("ParseArgs(%q): wrote unexpected output: %q", tt.args, buffer.String())
			}
		})
	}
}

func TestParseArgsParseError(t *testing.T) {
	tests := []struct {
		name   string
		args   []string
		expErr error
	}{
		{
			name: "unknown flag",
			args: []string{"-x"},
		},
		{
			name: "missing flag value",
			args: []string{"-t"},
		},
		{
			name:   "short help flag",
			args:   []string{"-h"},
			expErr: flag.ErrHelp,
		},
		{
			name:   "long help flag",
			args:   []string{"-help"},
			expErr: flag.ErrHelp,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buffer bytes.Buffer

			res, err := ParseArgs(tt.args, &buffer)
			if err == nil {
				t.Fatalf("ParseArgs(%q): expected error, got nil", tt.args)
			}

			var unexpectedArgsErr *UnexpectedArgsError
			if errors.As(err, &unexpectedArgsErr) {
				t.Fatalf("ParseArgs(%q): expected parse error, got UnexpectedArgsError: %v", tt.args, err)
			}

			if tt.expErr != nil && !errors.Is(err, tt.expErr) {
				t.Errorf("ParseArgs(%q): error %v; expected: %v", tt.args, err, tt.expErr)
			}

			if res != (Args{}) {
				t.Errorf("ParseArgs(%q): %+v; expected: %+v", tt.args, res, Args{})
			}

			if buffer.Len() == 0 {
				t.Errorf("ParseArgs(%q): expected usage output, got none", tt.args)
			}
		})
	}
}