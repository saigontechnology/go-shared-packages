package locale_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/datngo2sgtech/go-packages/locale"
)

func TestExtractLanguageCode(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		name           string
		locale         string
		expectedResult string
		expectedErr    error
	}{
		{
			name:           "Invalid locale",
			locale:         "ar-SA",
			expectedResult: "",
			expectedErr:    locale.ErrInvalidLocale,
		},
		{
			name:           "Valid locale",
			locale:         "ar_SA",
			expectedResult: "ar",
			expectedErr:    nil,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			actualResult, err := locale.ExtractLanguageCode(tc.locale)
			r := require.New(t)
			r.Equal(tc.expectedResult, actualResult)
			r.ErrorIs(err, tc.expectedErr)
		})
	}
}
