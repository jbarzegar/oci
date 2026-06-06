package ociv2_test

import (
	"strings"
	"testing"

	ociv2 "github.com/jbarzegar/oci/internal/oci/v2"
	"github.com/stretchr/testify/assert"
)

// TestValidateManifestName tests that the "name" of a image is spec compliant
func TestValidateManifestNameValid(t *testing.T) {
	valid := []string{
		"abc",
		"a-b",
		"a--b",
		"a_b_c",
		"my.container",
		"my__container",
		"123image",
		"repo/image",
		"my/repo/sub/image",
		"a-b/c-d_e",
	}

	for _, v := range valid {
		assert.Equal(t, true, ociv2.ValidateManifestName(v), v)
	}
}

func TestValidateManifestNameInvalid(t *testing.T) {
	invalid := []string{
		"ABC",          // (contains uppercase letters)
		".abc",         // (starts with a separator)
		"abc-",         // (ends with a separator, requires alphanumeric after -)
		"image:latest", // (contains colon :)
		"my image",     // (contains space)
		"@symbol",      // (contains invalid special character)
		"",             // (empty string)
		"A.b",          // (starts with uppercase)
		"/image",       //(missing leading alphanumeric segment)
	}

	for _, v := range invalid {
		result := ociv2.ValidateManifestName(v)
		assert.Equal(t, false, result, v)
	}
}

func TestValidateReference(t *testing.T) {
	valid := []string{
		"v1.0.0",
		"0.0.0",
		"nginx",
		"123image",
		"_private",
		"my.image",
		"a-b_c.d",
		strings.Repeat("a", 128),
	}

	for _, v := range valid {
		result := ociv2.ValidateReference(v)
		assert.Equal(t, true, result, v)
	}
}

func TestInvalidReference(t *testing.T) {
	invalid := []string{
		"",                        // (Empty string - minimum length 1)
		"-start",                  // (Starts with hyphen)
		".start",                  // (Starts with dot)
		"my/image",                // (Slash / is not in allowed character set for this regex)
		"my image",                // (Space not allowed)
		"registry/repo/image:tag", // (Contains colon and slash)
		"my@name",                 // (Ampersand or special chars like @ not allowed)
		strings.Repeat("a", 129),  // (String exceeds 128 characters)
	}

	for _, v := range invalid {
		result := ociv2.ValidateReference(v)
		assert.Equal(t, false, result, v)
	}
}
