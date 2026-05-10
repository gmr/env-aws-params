package main

import (
	"fmt"
	"regexp"
	"sort"
	"strings"
)

var InvalidPattern = regexp.MustCompile(`[^a-zA-Z0-9_]`)

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
