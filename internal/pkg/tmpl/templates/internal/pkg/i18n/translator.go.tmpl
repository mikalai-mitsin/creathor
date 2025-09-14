package i18n

import (
	"embed"

	"github.com/leonelquinteros/gotext"
)

//go:embed translations/**/*.po
var TranslationsFS embed.FS

type Translator struct {
	locales map[string]*gotext.Locale
}

func NewTranslator() (*Translator, error) {
	translator := &Translator{
		locales: make(map[string]*gotext.Locale),
	}
	langs, err := TranslationsFS.ReadDir("translations")
	if err != nil {
		return nil, err
	}
	for _, lang := range langs {
		if lang.IsDir() {
			l := gotext.NewLocaleFSWithPath(lang.Name(), TranslationsFS, "translations")
			l.AddDomain("default")
			translator.locales[lang.Name()] = l
		}
	}
	return translator, nil
}

func (t *Translator) Translate(lang string, key string) string {
	locale, ok := t.locales[lang]
	if !ok {
		return key
	}
	return locale.Get(key)
}
