package chronicle

import (
	"fmt"
	"testing"

	chronicle "github.com/form3tech-oss/terraform-provider-chronicle/client"
	"github.com/pkg/errors"
)

func TestIsChronicleAPIErrorWithCode(t *testing.T) {
	notFound := &chronicle.ChronicleAPIError{HTTPStatusCode: 404, Message: "not found"}

	cases := []struct {
		name string
		err  error
		code int
		want bool
	}{
		{"bare error", notFound, 404, true},
		{"pkg/errors wrapped", errors.Wrap(notFound, "failed reading feed"), 404, true},
		{"fmt %w wrapped", fmt.Errorf("outer: %w", notFound), 404, true},
		{"different code", notFound, 500, false},
		{"unrelated error", fmt.Errorf("boom"), 404, false},
		{"nil error", nil, 404, false},
	}

	for _, tc := range cases {
		if got := IsChronicleAPIErrorWithCode(tc.err, tc.code); got != tc.want {
			t.Errorf("%s: IsChronicleAPIErrorWithCode = %v, want %v", tc.name, got, tc.want)
		}
	}
}
