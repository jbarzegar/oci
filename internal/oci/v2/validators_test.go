package ociv2_test

import (
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

// Validate
// func TestValidateReference(t *testing.T) {
// 	valid := []string{
// 		"v1.0.0",
// 	}
// }
