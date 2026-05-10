package main

import (
	"errors"
	"testing"
)

func TestValidateArgsNoPrefixIsValid(t *testing.T) {
	code, err := validateArgs(1, false, false)
	if code != 0 {
		t.Fatalf("expected code to be 0, got %v", code)
	}
	if err != nil {
		t.Fatalf("expected err to be nil, got %v", err)
	}
}

func TestValidateArgsMissingCommand(t *testing.T) {
	code, err := validateArgs(0, false, false)
	if code != 1 {
		t.Fatalf("expected code to be 1, got %v", code)
	}
	if err == nil {
		t.Fatalf("expected err to be set, got nil")
	}
}

func TestValidateArgsStripAndSanitize(t *testing.T) {
	code, err := validateArgs(1, true, true)
	if code != 2 {
		t.Fatalf("expected code to be 2, got %v", code)
	}
	if err == nil {
		t.Fatalf("expected err to be set, got nil")
	}
}

func TestValidateArgsValid(t *testing.T) {
	code, err := validateArgs(1, false, true)
	if code != 0 {
		t.Fatalf("expected code to be 0, got %v", code)
	}
	if err != nil {
		t.Fatalf("expected err to be nil, got %v", err)
	}
}

func TestErrorPrefix(t *testing.T) {
	testError := errors.New("foo bar")
	result := errorPrefix(testError)
	expectation := "ERROR: foo bar"
	if result != expectation {
		t.Fatalf("expected %q, got %q", expectation, result)
	}
}
