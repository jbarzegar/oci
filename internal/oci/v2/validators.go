package ociv2

import (
	"regexp"
	"strings"
)

func ValidateManifestName(name string) bool {
	// The regex pattern here differs slightly from the spec found at:
	// https://github.com/opencontainers/distribution-spec/blob/main/spec.md#pulling-manifests
	// The regex below differs slightly by the addition of ^ & $ on prefix and suffix respectively.
	// Regex provided in the spec produces a "Partial Match" see: regexr.com/8n9as
	// Meaning invalid strings would match as a false-positive
	// fixed regex to match full path: regexr.com/8n9av
	namePattern := regexp.MustCompile(
		`^[a-z0-9]+((\.|_|__|-+)[a-z0-9]+)*(\/[a-z0-9]+((\.|_|__|-+)[a-z0-9]+)*)*$`,
	)

	v := namePattern.MatchString(strings.TrimSpace(name))

	return v
}

func ValidateReference(ref string) bool {
	refPattern := regexp.MustCompile(
		`^[a-zA-Z0-9_][a-zA-Z0-9._-]{0,127}$`,
	)

	return refPattern.MatchString(ref)
}
