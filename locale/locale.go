package locale

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

const (
	DefaultLocale = "ar_SA"
	ArabicLocale  = "ar_SA"
	EnglishLocale = "en_US"
)

var ErrInvalidLocale = errors.New("Invalid locale")

func GetLocaleFromGinContext(c *gin.Context) string {
	return c.GetString(ContextLocaleKey)
}

func ExtractLanguageCode(locale string) (string, error) {
	temp := strings.Split(locale, "_")
	if len(temp) != 2 {
		return "", errors.Wrap(ErrInvalidLocale, locale)
	}

	return temp[0], nil
}
