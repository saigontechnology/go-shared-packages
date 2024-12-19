package locale

import (
	"github.com/gin-gonic/gin"
)

const (
	ContextLocaleKey = "x-locale"
	HeaderLocaleKey  = "x-locale"
)

// XLocaleExtractionMiddleware extract locale from header x-locale and set it to context
// Locale value can be got everytime, everywhere later if a function has a context.
func XLocaleExtractionMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		locale := DefaultLocale
		if xLocale := c.GetHeader(HeaderLocaleKey); xLocale != "" {
			locale = xLocale
		}

		// Locale can be got from a gin context only, usually in a handler
		c.Set(ContextLocaleKey, locale)

		c.Next()
	}
}
