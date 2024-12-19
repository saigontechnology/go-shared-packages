package locale_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"

	"github.com/datngo2sgtech/go-packages/locale"
)

func TestLocaleMiddleware(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		name           string
		setupFunc      func() *http.Request
		expectedLocale string
	}{
		{
			name: "No x-locale in header",
			setupFunc: func() *http.Request {
				req, _ := http.NewRequestWithContext(
					context.Background(),
					http.MethodGet,
					"/v1/wishlists",
					nil,
				)
				return req
			},
			expectedLocale: locale.DefaultLocale,
		},
		{
			name: "No x-locale in header",
			setupFunc: func() *http.Request {
				req, _ := http.NewRequestWithContext(
					context.Background(),
					http.MethodGet,
					"/v1/wishlists",
					nil,
				)
				req.Header.Set(locale.HeaderLocaleKey, "en_US")
				return req
			},
			expectedLocale: "en_US",
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			rec := httptest.NewRecorder()
			gc, _ := gin.CreateTestContext(rec)
			gc.Request = tc.setupFunc()
			locale.XLocaleExtractionMiddleware()(gc)
			r := require.New(t)
			localeValue, exist := gc.Get(locale.ContextLocaleKey)
			r.True(exist)
			r.Equal(tc.expectedLocale, localeValue)
		})
	}
}
