package localization

import (
	"log"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/pelletier/go-toml/v2"
	"golang.org/x/text/language"
)

const (
	ContextLocalizerKey = "Health-localizer"
	DefaultLocale       = "ar_SA"
	HeaderLocaleKey     = "x-locale"
)

func InitBundle(translationsDir string) (*i18n.Bundle, error) {
	bundle := i18n.NewBundle(language.Arabic)
	bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)

	log.Printf("[LOCALIZER] Looking for translation files in \"%s\"\n", translationsDir)
	files, err := os.ReadDir(translationsDir)
	if err != nil {
		log.Printf("[LOCALIZER] Failed to read \"%s\" directory: %v\n", translationsDir, err)
		return nil, err
	}

	for _, f := range files {
		if !f.IsDir() && strings.HasSuffix(f.Name(), ".toml") {
			_, err := bundle.LoadMessageFile(translationsDir + "/" + f.Name())
			if err != nil {
				log.Printf("[LOCALIZER] Failed to load localization file: %v\n", err)

				return nil, err
			}

			log.Printf("[LOCALIZER] Loaded localization file: %s\n", f.Name())
		}
	}

	return bundle, nil
}

func InitLocalizerMiddleware(bundle *i18n.Bundle) gin.HandlerFunc {
	return func(c *gin.Context) {
		locale := DefaultLocale
		if xLocale := c.GetHeader(HeaderLocaleKey); xLocale != "" {
			locale = xLocale
		}

		lT := language.Make(locale)
		localizer := i18n.NewLocalizer(bundle, lT.String())
		c.Set(ContextLocalizerKey, localizer)

		c.Next()
	}
}

func Get(c *gin.Context) *i18n.Localizer {
	val, ok := c.MustGet(ContextLocalizerKey).(*i18n.Localizer)
	if !ok {
		return nil
	}

	return val
}
