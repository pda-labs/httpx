package httpx

import (
	"net/http"
)

// redirect - базовая функция. Все остальные вызывают её.
//
// location - абсолютный (предпочтительно) или относительный URL.
// Если location == "", будет отправлен только статус без Location‑header.
func redirect(w http.ResponseWriter, r *http.Request, status int, location string) {
	if location != "" {
		w.Header().Set("Location", location)
	}

	// JSON‑обёртка полезна для SPA / fetch‑клиентов,
	// классические браузеры её игнорируют.
	JSON(w, r, status, map[string]string{
		"location": location,
	})
}

/* 300 */

// RedirectMultipleChoices - 300 MULTIPLE_CHOICES
//
// Сервер предлагает несколько вариантов ресурса (например, разные
// форматы). Передайте URL контента по умолчанию в location, чтобы
// клиент мог автоматически перейти.
//
// Status: 300 Multiple Choices
func RedirectMultipleChoices(w http.ResponseWriter, r *http.Request, location string) {
	redirect(w, r, http.StatusMultipleChoices, location)
}

/* 301 */

// RedirectMovedPermanently - 301 MOVED_PERMANENTLY
//
// Канонический URL изменился навсегда (SEO‑friendly).
// Браузеры кешируют; поисковики передают “link‑juice”.
//
// Status: 301 Moved Permanently
func RedirectMovedPermanently(w http.ResponseWriter, r *http.Request, location string) {
	redirect(w, r, http.StatusMovedPermanently, location)
}

/* 302 */

// RedirectFound - 302 FOUND
//
// Временный переезд (legacy). Современный эквивалент - 307.
// Используйте, если нужно сохранить метод GET / HEAD.
//
// Status: 302 Found
func RedirectFound(w http.ResponseWriter, r *http.Request, location string) {
	redirect(w, r, http.StatusFound, location)
}

/* 303 */

// RedirectSeeOther - 303 SEE_OTHER
//
// После успешного POST’а отправьте пользователя на GET‑страницу
// результата. Браузеры всегда выполнят последующий GET.
//
// Status: 303 See Other
func RedirectSeeOther(w http.ResponseWriter, r *http.Request, location string) {
	redirect(w, r, http.StatusSeeOther, location)
}

/* 304 */

// RedirectNotModified - 304 NOT_MODIFIED
//
// Особый случай: тело отсутствует; Location не нужен.
// Вызывайте, если ETag / If‑Modified‑Since совпали.
//
// Status: 304 Not Modified
func RedirectNotModified(w http.ResponseWriter, r *http.Request) {
	// 304 не должен иметь тела по RFC 7232.
	w.WriteHeader(http.StatusNotModified)
}

/* 307 */

// RedirectTemporary - 307 TEMPORARY_REDIRECT
//
// Временный переезд: **сохраняет** HTTP‑метод и тело запроса
// (в отличие от 302). Подходит для API, когда POST нужно
// временно проксировать на другую ноду.
//
// Status: 307 Temporary Redirect
func RedirectTemporary(w http.ResponseWriter, r *http.Request, location string) {
	redirect(w, r, http.StatusTemporaryRedirect, location)
}

/* 308 */

// RedirectPermanent - 308 PERMANENT_REDIRECT
//
// Постоянный переезд, но, как и 307, сохраняет метод + тело.
// Выбирайте для REST‑ресурсов, которые переместились навсегда.
//
// Status: 308 Permanent Redirect
func RedirectPermanent(w http.ResponseWriter, r *http.Request, location string) {
	redirect(w, r, http.StatusPermanentRedirect, location)
}
