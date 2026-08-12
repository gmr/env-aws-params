package main

import (
	"maps"
	"regexp"
	"slices"
	"strings"

	log "github.com/sirupsen/logrus"
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
	// Transform in sorted parameter order so collisions resolve deterministically.
	values := make(map[string]string, len(parameters))
	sources := make(map[string]string, len(parameters))
	for _, name := range slices.Sorted(maps.Keys(parameters)) {
		key := name
		if sanitize {
			key = InvalidPattern.ReplaceAllString(key, "_")
		}
		if strip {
			key = InvalidPattern.ReplaceAllString(key, "")
		}
		if upcase {
			key = strings.ToUpper(key)
		}
		if prev, ok := sources[key]; ok {
			log.Warnf("Parameters %q and %q both map to %s; keeping the value of %q", prev, name, key, name)
		}
		sources[key] = name
		values[key] = parameters[name]
	}

	vars := make([]string, 0, len(values))
	for key, value := range values {
		vars = append(vars, key+"="+value)
	}
	slices.Sort(vars)
	return vars
}
