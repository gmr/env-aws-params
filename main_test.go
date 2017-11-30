package main

import (
	"fmt"
	"sort"
	"testing"
)

func AssertEqual(t *testing.T, v1 []string, v2 []string) {
	if len(v1) != len(v2) {
		t.Error("Slices are not the same length", len(v1), len(v2))
	}
	v1 = sort.StringSlice(v1)
	v2 = sort.StringSlice(v2)
	for i, _ := range v1 {
		if v1[i] != v2[i] {
			t.Error(fmt.Sprintf("Values at offset %v do not match", i), v1[i], v2[i])
		}
	}
}

func TestBuildEnvVarsUpcaseFalse(t *testing.T) {
	var params map[string]string

	params = make(map[string]string)
	params["FOO"] = "bar"
	params["baz"] = "qux"

	expectation := []string{"FOO=bar", "baz=qux"}
	envvars := BuildEnvVars(params, false)
	AssertEqual(t, envvars, expectation)
}

func TestBuildEnvVarsUpcaseTrue(t *testing.T) {
	var params map[string]string

	params = make(map[string]string)
	params["FOO"] = "bar"
	params["baz"] = "qux"

	expectation := []string{"FOO=bar", "BAZ=qux"}
	envvars := BuildEnvVars(params, true)
	AssertEqual(t, envvars, expectation)
}
