package main

import (
	"fmt"
	"os/exec"
	"sort"
	"testing"
)

func AssertEqual(t *testing.T, value []string, expect []string) {
	if len(value) != len(expect) {
		t.Error("Slices are not the same length", len(value), len(expect))
	}
	sort.Strings(value)
	sort.Strings(expect)
	for i := range value {
		if value[i] != expect[i] {
			t.Error(fmt.Sprintf("Values at offset %v do not match", i), value[i], expect[i])
		}
	}
}

// AssertEqualOrdered is AssertEqual without the sort, for assertions where the
// ordering of the slice is itself the thing under test.
func AssertEqualOrdered(t *testing.T, value []string, expect []string) {
	t.Helper()
	if len(value) != len(expect) {
		t.Fatalf("expected %v entries, got %v: %v", len(expect), len(value), value)
	}
	for i := range expect {
		if value[i] != expect[i] {
			t.Errorf("Values at offset %v do not match: expected %q, got %q", i, expect[i], value[i])
		}
	}
}

func TestMergeEnvVarsAppendsSSMLast(t *testing.T) {
	ssmVars := []string{"ONLY_SSM=ssm", "SHARED=from-ssm"}
	environ := []string{"ONLY_ENV=env", "SHARED=from-environ"}

	expectation := []string{
		"ONLY_ENV=env", "SHARED=from-environ",
		"ONLY_SSM=ssm", "SHARED=from-ssm",
	}
	AssertEqualOrdered(t, MergeEnvVars(ssmVars, environ), expectation)
}

func TestMergeEnvVarsPristineInheritsNothing(t *testing.T) {
	ssmVars := []string{"ONLY_SSM=ssm"}
	AssertEqualOrdered(t, MergeEnvVars(ssmVars, nil), ssmVars)
}

func TestMergeEnvVarsIsNeverNil(t *testing.T) {
	// A nil result would make os/exec inherit the parent environment.
	if merged := MergeEnvVars(nil, nil); merged == nil {
		t.Fatal("expected a non-nil slice, got nil")
	}
}

// The precedence contract holds only because os/exec keeps the last occurrence
// of a duplicate key, so assert it against a real child process rather than
// against the ordering of the slice alone.
func TestMergeEnvVarsSSMWinsInChildProcess(t *testing.T) {
	cmd := exec.Command("/bin/sh", "-c", `printf '%s' "$SHARED"`)
	cmd.Env = MergeEnvVars([]string{"SHARED=from-ssm"}, []string{"SHARED=from-environ"})

	out, err := cmd.Output()
	if err != nil {
		t.Fatalf("running child: %v", err)
	}
	if got := string(out); got != "from-ssm" {
		t.Fatalf("expected child to see %q, got %q", "from-ssm", got)
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

func TestBuildEnvVarsCollisionKeepsLastSortedParameter(t *testing.T) {
	params := map[string]string{
		"db-host": "from-dash",
		"db_host": "from-underscore",
	}

	envVars := BuildEnvVars(params, true, false, true)
	AssertEqual(t, envVars, []string{"DB_HOST=from-underscore"})
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
