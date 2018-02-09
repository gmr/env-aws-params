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
	params["baz"] = "qux"
	params["FOO"] = "bar"

	expectation := []string{"baz=qux", "FOO=bar"}
	envvars := BuildEnvVars(params, false, false, false)
	AssertEqual(t, envvars, expectation)
}

func TestBuildEnvVarsUpcaseTrue(t *testing.T) {
	var params map[string]string

	params = make(map[string]string)
	params["baz"] = "qux"
	params["FOO"] = "bar"

	expectation := []string{"BAZ=qux", "FOO=bar"}
	envVars := BuildEnvVars(params, false, false, true)
	AssertEqual(t, envVars, expectation)
}

func TestBuildEnvVarsUpperSanitize(t *testing.T) {
	var params map[string]string

	params = make(map[string]string)
	params["FOO"] = "bar"
	params["baz-corgie"] = "qux"
	params["wE_irD-kEY!"] = "zaphod"

	expectation := []string{"BAZ_CORGIE=qux", "FOO=bar", "WE_IRD_KEY_=zaphod"}
	envVars := BuildEnvVars(params, true, false, true)
	AssertEqual(t, envVars, expectation)
}

func TestBuildEnvVarsUpperStrip(t *testing.T) {
	var params map[string]string

	params = make(map[string]string)
	params["FOO"] = "bar"
	params["baz-corgie"] = "qux"
	params["wE_irD-kEY!"] = "zaphod"

	expectation := []string{"BAZCORGIE=qux", "FOO=bar", "WE_IRDKEY=zaphod"}
	envVars := BuildEnvVars(params, false, true, true)
	AssertEqual(t, envVars, expectation)
}
