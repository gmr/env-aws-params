package main

import (
	"errors"
	"flag"
	"io/ioutil"
	"testing"

	"github.com/urfave/cli"
)

func NewContext(t *testing.T, testArgs []string) *cli.Context {
	app := cli.NewApp()
	app.Writer = ioutil.Discard
	app.Flags = cliFlags()
	set := flag.NewFlagSet("test", 0)
	for _, f := range app.Flags {
		f.Apply(set)
	}
	set.Parse(testArgs)
	return cli.NewContext(app, set, nil)
}

func TestMissingPrefix(t *testing.T) {
	var testArgs []string

	testArgs = []string{"--upcase"}

	code, err := validateArgs(NewContext(t, testArgs))
	if code != 1 {
		t.Fatalf("expected code to be 1, got %v", code)
	}
	if err == nil {
		t.Fatalf("expected err to be set, got nil")
	}
}

func TestMissingCommand(t *testing.T) {
	var testArgs []string

	testArgs = []string{"--prefix", "/foo"}

	code, err := validateArgs(NewContext(t, testArgs))
	if code != 2 {
		t.Fatalf("expected code to be 2, got %v", code)
	}
	if err == nil {
		t.Fatalf("expected err to be set, got nil")
	}
}

func TestMissingStripAndSanitize(t *testing.T) {
	var testArgs []string

	testArgs = []string{"--prefix", "/foo", "--strip", "--sanitize", "/bin/bash"}

	code, err := validateArgs(NewContext(t, testArgs))
	if code != 3 {
		t.Fatalf("expected code to be 3, got %v", code)
	}
	if err == nil {
		t.Fatalf("expected err to be set, got nil")
	}
}

func TestValidCLIOptions(t *testing.T) {
	var testArgs []string

	testArgs = []string{"--prefix", "/foo", "--strip", "/bin/bash"}

	code, err := validateArgs(NewContext(t, testArgs))
	if code != 0 {
		t.Fatalf("expected code to be 0, got %v", code)
	}
	if err != nil {
		t.Fatalf("expected err to be nil, got %v", code)
	}
}

func TestErrorPrefix(t *testing.T) {
	testError := errors.New("foo bar")
	result := errorPrefix(testError)
	expectation := "ERROR: foo bar"
	if result != expectation {
		t.Fatalf("expected \"%v\", got \"%v\"", result, expectation)
	}
}
