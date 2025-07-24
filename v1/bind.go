package httpx

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
)

// errDecode  - повреждённый или непарсибельный JSON.
var errDecode = errors.New("httpx: cannot decode JSON body")

// errValidatorUnset - глобальный валидатор не сконфигурирован.
var errValidatorUnset = errors.New("httpx: validator is not set (call httpx.V = validator.New())")

const maxBodySize int64 = 1 << 23 // 8 MiB

// BindValidate читает JSON‑тело, валидирует dst и локализует ошибки.
//
// Возвращает:
//  1. details - map[field]translated msg; nil, если валидация прошла или ошибка другого типа.
//  2. err     - любая ошибка процесса (декодинг, отсутствие валидатора, validator.ValidationErrors).
//
// Использование:
//
//	var dto SignupDTO
//	if det, err := httpx.BindValidate(r, &dto); err != nil {
//	    if det != nil {
//	        httpx.ErrorValidation(w, r, det) // 400 + детали
//	    } else {
//	        httpx.ErrorBadRequest(w, r, err.Error())
//	    }
//	    return
//	}
func BindValidate[T any](r *http.Request, dst *T) (map[string]string, error) {
	// Читаем тело с учётом контекста + лимита
	if r.ContentLength != 0 {
		defer r.Body.Close()

		// MaxBytesReader: вернёт 4xx если тело превышает лимит.
		limited := http.MaxBytesReader(nil, r.Body, maxBodySize)

		decoder := json.NewDecoder(limited)
		decoder.DisallowUnknownFields()

		// Читаем первую структуру
		if err := decoder.Decode(dst); err != nil {
			select { // если ctx отменён - лучше вернуть context error
			case <-r.Context().Done():
				return nil, r.Context().Err()
			default:
			}
			return nil, fmt.Errorf("%w: %v", errDecode, err)
		}
		// Проверяем trailing garbage
		if decoder.More() {
			return nil, fmt.Errorf("%w: extra data after JSON object", errDecode)
		}
	}

	//  Валидатор
	if V == nil {
		return nil, errValidatorUnset
	}

	//  Валидация
	if err := V.Struct(dst); err != nil {
		var ve validator.ValidationErrors
		if !errors.As(err, &ve) {
			return nil, err
		}
		tr := TranslatorFor(r)
		details := make(map[string]string, len(ve))
		for _, fe := range ve {
			details[fe.Field()] = fe.Translate(tr)
		}
		return details, err
	}

	return nil, nil
}
