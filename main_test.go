package main

import (
	"fmt"
	"sort"
	"testing"
)

func AssertEqual(t *testing.T, value []string, expect []string) {
	if len(value) != len(expect) {
		t.Error("Slices are not the same length", len(value), len(expect))
	}
	sort.Strings(value)
	sort.Strings(expect)
	for i, _ := range value {
		if value[i] != expect[i] {
			t.Error(fmt.Sprintf("Values at offset %v do not match", i), value[i], expect[i])
		}
	}
}

func TestBuildEnvVarsUpcaseFalse(t *testing.T) {
	var params map[string]string

	params = make(map[string]string)
	params["FOO"] = "bar"
	params["baz"] = "qux"

	expectation := []string{"baz=qux", "FOO=bar"}
	envvars := BuildEnvVars(params, false)
	AssertEqual(t, envvars, expectation)
}

func TestBuildEnvVarsUpcaseTrue(t *testing.T) {
	var params map[string]string

	params = make(map[string]string)
	params["FOO"] = "bar"
	params["baz"] = "qux"

	expectation := []string{"BAZ=qux", "FOO=bar"}
	envvars := BuildEnvVars(params, true)
	AssertEqual(t, envvars, expectation)
}
