package main

import (
	"errors"
	"testing"
)

func TestValidateArgsNoPrefixIsValid(t *testing.T) {
	if err := validateArgs(1, false, false); err != nil {
		t.Fatalf("expected err to be nil, got %v", err)
	}
}

func TestValidateArgsMissingCommand(t *testing.T) {
	if err := validateArgs(0, false, false); err == nil {
		t.Fatalf("expected err to be set, got nil")
	}
}

func TestValidateArgsStripAndSanitize(t *testing.T) {
	if err := validateArgs(1, true, true); err == nil {
		t.Fatalf("expected err to be set, got nil")
	}
}

func TestValidateArgsValid(t *testing.T) {
	if err := validateArgs(1, false, true); err != nil {
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
