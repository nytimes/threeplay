package threeplay

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCheckForAPIError(t *testing.T) {
	t.Parallel()
	var tests = []struct {
		name     string
		input    string
		expected error
	}{
		{
			"no error",
			"{}",
			nil,
		},
		{
			"auth error",
			`{"iserror": true, "errors": {"authentication": "unauthorized!"}}`,
			ErrUnauthorized,
		},
		{
			"not found",
			`{"iserror": true, "errors": {"not_found": "can't find what you asked for"}}`,
			ErrNotFound,
		},
		{
			"generic error",
			`{"iserror":true,"errors":{"internal_error":"something went wrong"}}`,
			errors.New(`api error: {"iserror":true,"errors":{"internal_error":"something went wrong"}}`),
		},
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			err := checkForAPIError([]byte(test.input))
			assert.Equal(t, err, test.expected)
		})
	}
}
