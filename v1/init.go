package httpx

import (
	"html"
	"net/http"
	"strings"
	"sync"

	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"golang.org/x/text/language"

	// locales
	"github.com/go-playground/locales/de"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/es"
	"github.com/go-playground/locales/fr"
	"github.com/go-playground/locales/it"
	"github.com/go-playground/locales/ja"
	"github.com/go-playground/locales/ko"
	"github.com/go-playground/locales/lv"
	"github.com/go-playground/locales/pt"
	"github.com/go-playground/locales/ru"
	"github.com/go-playground/locales/zh"

	// translations
	de_trans "github.com/go-playground/validator/v10/translations/de"
	en_trans "github.com/go-playground/validator/v10/translations/en"
	es_trans "github.com/go-playground/validator/v10/translations/es"
	fr_trans "github.com/go-playground/validator/v10/translations/fr"
	it_trans "github.com/go-playground/validator/v10/translations/it"
	ja_trans "github.com/go-playground/validator/v10/translations/ja"
	ko_trans "github.com/go-playground/validator/v10/translations/ko"
	pt_trans "github.com/go-playground/validator/v10/translations/pt"
	ru_trans "github.com/go-playground/validator/v10/translations/ru"
	zh_trans "github.com/go-playground/validator/v10/translations/zh"
)

var (
	initOnce    sync.Once
	V           *validator.Validate
	translators map[string]ut.Translator
)

func init() {
	initOnce.Do(func() {
		if V == nil {
			def := validator.New()
			_ = def.RegisterValidation("nohtml", func(fl validator.FieldLevel) bool {
				s := fl.Field().String()
				return html.EscapeString(s) == s
			})
			V = def
		}

		// Universal‑translator + регистрация языков
		uni := ut.New(
			en.New(),
			en.New(),
			ru.New(),
			de.New(),
			lv.New(),
			zh.New(),
			fr.New(),
			es.New(),
			it.New(),
			pt.New(),
			ja.New(),
			ko.New(),
		)

		translators = make(map[string]ut.Translator, 12)

		// helper
		add := func(code string, reg func(*validator.Validate, ut.Translator) error) {
			tr, _ := uni.GetTranslator(code)
			_ = reg(V, tr)
			translators[code] = tr
		}

		add("en", en_trans.RegisterDefaultTranslations)
		add("ru", ru_trans.RegisterDefaultTranslations)
		add("de", de_trans.RegisterDefaultTranslations)
		add("zh", zh_trans.RegisterDefaultTranslations)
		add("fr", fr_trans.RegisterDefaultTranslations)
		add("es", es_trans.RegisterDefaultTranslations)
		add("lv", en_trans.RegisterDefaultTranslations) // TODO: заменить, когда выйдет поддержка латышского языка
		add("it", it_trans.RegisterDefaultTranslations)
		add("pt", pt_trans.RegisterDefaultTranslations)
		add("ja", ja_trans.RegisterDefaultTranslations)
		add("ko", ko_trans.RegisterDefaultTranslations)
	})
}

func baseLocale(tag string) string { // helper
	if i := strings.IndexByte(tag, '-'); i > 0 {
		return tag[:i]
	}
	return tag
}

// TranslatorFor выбирает переводчик «на лету»
//
//  1. X-Request-Lang
//  2. Accept-Language
//  3. fallback -> "en"
func TranslatorFor(r *http.Request) ut.Translator {
	if lang := baseLocale(r.Header.Get("X-Request-Lang")); lang != "" {
		if tr, ok := translators[lang]; ok {
			return tr
		}
	}
	if al := r.Header.Get("Accept-Language"); al != "" {
		if tags, _, err := language.ParseAcceptLanguage(al); err == nil && len(tags) > 0 {
			if tr, ok := translators[baseLocale(tags[0].String())]; ok {
				return tr
			}
		}
	}
	return translators["en"]
}

// RegisterCustomValidator добавляет кастомное правило в валидатор + переводы.
//
// Пример:
//
//	httpx.RegisterCustomValidator("tz", validateTimeZone, map[string]string{
//	    "en": "Must be a valid IANA time zone",
//	    "ru": "Некорректная IANA таймзона (например, Europe/Moscow)",
//	})
func RegisterCustomValidator(tag string, fn validator.Func, messages map[string]string) error {
	if V == nil {
		return ErrValidatorUnset
	}

	if err := V.RegisterValidation(tag, fn); err != nil {
		return err
	}

	// Регистрируем переводы для всех подключённых Translator’ов
	for lang, msg := range messages {
		tr, ok := translators[lang]
		if !ok {
			continue // язык не инициализирован - пропускаем
		}

		_ = V.RegisterTranslation(tag, tr,
			func(ut ut.Translator) error {
				return ut.Add(tag, msg, true)
			},
			func(ut ut.Translator, fe validator.FieldError) string {
				t, _ := ut.T(tag)
				return t
			},
		)
	}

	return nil
}
