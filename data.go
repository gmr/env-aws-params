package main

import (
	"fmt"
	"regexp"
	"sort"
	"strings"
)

var InvalidPattern = regexp.MustCompile(`[^a-zA-Z0-9_]`)

// MergeEnvVars combines the inherited environment with SSM-derived variables,
// letting Parameter Store win whenever a key is present in both.
//
// Precedence is expressed purely through ordering: os/exec de-duplicates
// cmd.Env before handing it to the child and keeps the last occurrence of each
// key, so ssmVars is appended second. Pass a nil environ to inherit nothing.
// The result is always non-nil, because a nil cmd.Env makes os/exec fall back
// to the parent's environment.
func MergeEnvVars(ssmVars []string, environ []string) []string {
	merged := make([]string, 0, len(environ)+len(ssmVars))
	merged = append(merged, environ...)
	merged = append(merged, ssmVars...)
	return merged
}

func BuildEnvVars(parameters map[string]string, sanitize bool, strip bool, upcase bool) []string {
	var vars []string

	for k, v := range parameters {
		if sanitize {
			k = InvalidPattern.ReplaceAllString(k, "_")
		}
		if strip {
			k = InvalidPattern.ReplaceAllString(k, "")
		}
		if upcase {
			k = strings.ToUpper(k)
		}
		vars = append(vars, fmt.Sprintf("%s=%s", k, v))
	}
	sort.Strings(vars)
	return vars
}
