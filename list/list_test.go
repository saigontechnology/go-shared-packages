package list_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/datngo2sgtech/go-packages/list"
)

func TestIsSliceContainsPrefix(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		name     string
		strInput string
		arrInput []string
		result   bool
	}{
		{
			name:     "Normal Case",
			strInput: "/blah",
			arrInput: []string{"/blah?c=1", "/hi", "/yikes"},
			result:   true,
		},
		{
			name:     "No match case",
			strInput: "/blah",
			arrInput: []string{"/yikes", "/hi", "/no"},
			result:   false,
		},
		{
			name:     "Edge case",
			strInput: "",
			arrInput: []string{"/yikes", "/hi", "/no"},
			result:   false,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			actualResult := list.IsSliceContainsPrefix(tc.strInput, tc.arrInput)
			assert.Equal(t, tc.result, actualResult)
		})
	}
}
