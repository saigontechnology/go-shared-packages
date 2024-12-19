package test

import (
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/pelletier/go-toml/v2"
	"golang.org/x/text/language"

	"github.com/datngo2sgtech/go-packages/env"
	"github.com/datngo2sgtech/go-packages/localization"
	"github.com/datngo2sgtech/go-packages/must"
)

const (
	TypeHtOptionTranslationsDir = "HtOptionTranslationDir"
)

// HandlerTest is a utility to implement unit test for handlerFunc.
type HandlerTest interface {
	UT
	Recorder() *httptest.ResponseRecorder
	GinContext() *gin.Context
	GinEngine() *gin.Engine
	AssertRestResponse(statusCode int, responseBody string)
	RequireRestResponse(statusCode int, responseBody string)
}

type handlerTest struct {
	*ut
	cfg *htConfig
	rec *httptest.ResponseRecorder
	gc  *gin.Context
	ge  *gin.Engine
}

type HtOption interface {
	Type() string
	Value() string
}

type HtOptionTranslationsDir struct {
	TranslationDir string
}

func (o *HtOptionTranslationsDir) Type() string {
	return TypeHtOptionTranslationsDir
}

func (o *HtOptionTranslationsDir) Value() string {
	return o.TranslationDir
}

func NewHandlerTest(t *testing.T, options ...HtOption) HandlerTest {
	t.Helper()

	// nolint: tenv
	os.Setenv(env.EnvironmentVariable, "test")
	env.Load()

	cfg, err := newHtConfig()
	must.NotFail(err)
	parseHtOptions(cfg, options...)

	rec := httptest.NewRecorder()
	gc, ge := gin.CreateTestContext(rec)

	// Set up test localization
	bundle := i18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)
	bundle.MustLoadMessageFile(filepath.Join(cfg.TranslationsDir, "en-US.toml"))
	gc.Set(localization.ContextLocalizerKey, i18n.NewLocalizer(bundle, "en-US"))

	return &handlerTest{
		ut:  newUT(t),
		cfg: cfg,
		rec: rec,
		gc:  gc,
		ge:  ge,
	}
}

func parseHtOptions(cfg *htConfig, options ...HtOption) {
	for _, option := range options {
		if option.Type() == TypeHtOptionTranslationsDir {
			cfg.TranslationsDir = option.Value()
		}
	}
}

func (t *handlerTest) Recorder() *httptest.ResponseRecorder {
	return t.rec
}

func (t *handlerTest) GinContext() *gin.Context {
	return t.gc
}

func (t *handlerTest) GinEngine() *gin.Engine {
	return t.ge
}

func (t *handlerTest) AssertRestResponse(statusCode int, responseBody string) {
	t.assert.Equal(statusCode, t.rec.Code)
	t.assert.Equal(responseBody, t.rec.Body.String())
}

func (t *handlerTest) RequireRestResponse(statusCode int, responseBody string) {
	t.require.Equal(statusCode, t.rec.Code)
	t.require.Equal(responseBody, t.rec.Body.String())
}
