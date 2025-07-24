package httpx

import "net/http"

/* 500 */

// ErrorInternal - 500 INTERNAL
//
// «Мы что-то сломали». Универсальная внутренняя ошибка, когда
// причина скрыта от клиента, а разработчики увидят её в логах.
//
//   - Используйте при панике, ошибке базы, непредвиденном исключении.
//   - Не выдавайте детали (stacktrace) публично - только лаконичное msg.
//
// Status: 500 Internal Server Error
// Code:   INTERNAL
func ErrorInternal(w http.ResponseWriter, r *http.Request, msg string) {
	Error(w, r, http.StatusInternalServerError, "INTERNAL", msg, nil)
}

/* 501 */

// ErrorNotImplemented - 501 NOT_IMPLEMENTED
//
// Фича/эндпоинт задокументированы, но ещё не реализованы.
//
//   - Полезно оставлять заглушки: фронт или интегратор поймут, что
//     функционал появится позже, а не сломан.
//   - Передайте `feature`, чтобы указать, что именно недоступно.
//
// Status: 501 Not Implemented
// Code:   NOT_IMPLEMENTED
func ErrorNotImplemented(w http.ResponseWriter, r *http.Request, feature string) {
	txt := "Feature not implemented"
	if feature != "" {
		txt = feature + " not implemented"
	}
	Error(w, r, http.StatusNotImplemented, "NOT_IMPLEMENTED", txt, nil)
}

/* 502 */

// ErrorBadGateway - 502 BAD_GATEWAY
//
// Сервис-прокси получил ошибочный ответ **от вышестоящего бекенда**.
//
//   - API-шлюз не смог достучаться до микросервиса или получил от него 500.
//   - Ваш сервис сам является прокси к стороннему API.
//
// Status: 502 Bad Gateway
// Code:   BAD_GATEWAY
func ErrorBadGateway(w http.ResponseWriter, r *http.Request, msg string) {
	Error(w, r, http.StatusBadGateway, "BAD_GATEWAY", msg, nil)
}

/* 503 */

// ErrorServiceUnavailable - 503 SERVICE_UNAVAILABLE
//
// Сервис временно недоступен: ведутся работы, перегрузка, отключён
// по feature-флагу или находится в процессе деплоя.
//
//   - Верните заголовок `Retry-After`, если знаете время восстановления.
//   - Отличается от 500 тем, что это **ожидаемое** состояние сервиса.
//
// Status: 503 Service Unavailable
// Code:   SERVICE_UNAVAILABLE
func ErrorServiceUnavailable(w http.ResponseWriter, r *http.Request, msg string) {
	Error(w, r, http.StatusServiceUnavailable, "SERVICE_UNAVAILABLE", msg, nil)
}

/* 504 */

// ErrorTimeout - 504 TIMEOUT
//
// **Шлюз** (reverse-proxy, API Gateway) не дождался ответа
// от внутреннего сервиса или стороннего API.
//
//   - Укажите в `msg`, какой именно апстрим «завис», чтобы помочь Ops.
//   - Клиентам можно советовать повторить запрос позднее.
//
// Status: 504 Gateway Timeout
// Code:   TIMEOUT
func ErrorTimeout(w http.ResponseWriter, r *http.Request, msg string) {
	Error(w, r, http.StatusGatewayTimeout, "TIMEOUT", msg, nil)
}

/* 505 */

// ErrorHTTPVersionNotSupported - 505 VERSION_NOT_SUPPORTED
//
// Сервер не поддерживает версию протокола HTTP, указанную клиентом.
// В современных API встречается крайне редко.
//
// Status: 505 HTTP Version Not Supported
// Code:   VERSION_NOT_SUPPORTED
func ErrorHTTPVersionNotSupported(w http.ResponseWriter, r *http.Request, msg string) {
	Error(w, r, http.StatusHTTPVersionNotSupported, "VERSION_NOT_SUPPORTED", msg, nil)
}

/* 506 */

// ErrorVariantAlsoNegotiates - 506 VARIANT_NEGOTIATES
//
// Экзотика из RFC 2295 (Transparent Content Negotiation):
// круговая зависимость на этапе согласования варианта ресурса.
//
//   - Практически не используется, но функция есть для полноты.
//
// Status: 506 Variant Also Negotiates
// Code:   VARIANT_NEGOTIATES
func ErrorVariantAlsoNegotiates(w http.ResponseWriter, r *http.Request, msg string) {
	Error(w, r, http.StatusVariantAlsoNegotiates, "VARIANT_NEGOTIATES", msg, nil)
}

/* 507 */

// ErrorInsufficientStorage - 507 INSUFFICIENT_STORAGE
//
// Сервер не смог завершить операцию из-за нехватки места
// (например, заключительный /upload-chunk превысил квоту диска).
//
//   - Возвращайте в системах хранения, файловых сервисах, S3-совместимых API.
//
// Status: 507 Insufficient Storage
// Code:   INSUFFICIENT_STORAGE
func ErrorInsufficientStorage(w http.ResponseWriter, r *http.Request, msg string) {
	Error(w, r, http.StatusInsufficientStorage, "INSUFFICIENT_STORAGE", msg, nil)
}

/* 508 */

// ErrorLoopDetected - 508 LOOP_DETECTED
//
// При обработке WebDAV-запроса обнаружена бесконечная рекурсия (циклическая ссылка).
// Остаётся ради стандарта; обычным REST-API не нужно.
//
// Status: 508 Loop Detected
// Code:   LOOP_DETECTED
func ErrorLoopDetected(w http.ResponseWriter, r *http.Request, msg string) {
	Error(w, r, http.StatusLoopDetected, "LOOP_DETECTED", msg, nil)
}

/* 510 */

// ErrorNotExtended - 510 NOT_EXTENDED
//
// Расширение протокола, необходимое клиенту, не поддерживается сервером.
// Определён в RFC 2774. На практике почти не встречается.
//
// Status: 510 Not Extended
// Code:   NOT_EXTENDED
func ErrorNotExtended(w http.ResponseWriter, r *http.Request, msg string) {
	Error(w, r, http.StatusNotExtended, "NOT_EXTENDED", msg, nil)
}

/* 511 */

// ErrorNetworkAuthRequired - 511 NETWORK_AUTH_REQUIRED
//
// Клиент должен пройти сетевую аутентификацию (captive portal).
// Используется в публичных Wi-Fi до входа на портал.
//
// Status: 511 Network Authentication Required
// Code:   NETWORK_AUTH_REQUIRED
func ErrorNetworkAuthRequired(w http.ResponseWriter, r *http.Request, msg string) {
	Error(w, r, http.StatusNetworkAuthenticationRequired, "NETWORK_AUTH_REQUIRED", msg, nil)
}
