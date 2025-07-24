package httpx

import "net/http"

// ErrorBadRequest - 400 BAD_REQUEST
// Некорректный синтаксис или формат входящего запроса.
// Типовые случаи:
//   - Клиент отправил повреждённый JSON / XML / form-body, который не парсится.
//   - Слэш-значения query-параметров не приводятся к нужному типу (например, ?age=abc).
//   - Отсутствует обязательный заголовок (Content-Type, Authorization и пр.).
//
// Если структура распарсилась, но значения невалидны, используйте ErrorValidation.
// Рекомендуется писать понятный `msg`, чтобы UI/CLI показал человеку причину.
//
// Status: 400 Bad Request
//
// Code:   BAD_REQUEST
func ErrorBadRequest(w http.ResponseWriter, r *http.Request, msg string) {
	Error(w, r, http.StatusBadRequest, "BAD_REQUEST", msg, nil)
}

// ErrorValidation - 400 VALIDATION
//
// Структура запроса корректна, но значения нарушают бизнес-правила:
//   - e-mail не соответствует RFC 5322; телефон вне E.164; age < 0 и т.п.
//   - Отсутствуют взаимозависимые поля (например, одновременно latitude + longitude).
//
// Передавайте `details` в формате `map[string]string`, где ключ - поле, значение - причина.
// Тогда фронтенд сможет подсветить ошибки конкретных полей.
//
// Status: 400 Bad Request
//
// Code:   VALIDATION
func ErrorValidation(w http.ResponseWriter, r *http.Request, details interface{}) {
	Error(w, r, http.StatusBadRequest, "VALIDATION", "Request failed validation", details)
}

/* 401 */

// ErrorUnauthorized - 401 UNAUTHORIZED
//
// К запросу не приложены корректные учётные данные:
//
//   - Отсутствует или просрочен JWT / OAuth-токен.
//   - Неверный Basic-auth логин/пароль.
//   - В cookies нет сессионного идентификатора.
//
// Не раскрывайте причину в деталях: злоумышленнику не нужно знать, просрочен ли токен или он просто неверен.
//
// Status: 401 Unauthorized
//
// Code:   UNAUTHORIZED
func ErrorUnauthorized(w http.ResponseWriter, r *http.Request, msg string) {
	Error(w, r, http.StatusUnauthorized, "UNAUTHORIZED", msg, nil)
}

/* 402 */

// ErrorPaymentRequired - 402 PAYMENT_REQUIRED
//
// Запрос отклонён из-за неоплаченного тарифа или превышения квоты.
//
// Примеры применений:
//   - SaaS-API блокирует вызовы, когда баланс 0 ₽.
//   - Исчерпан месячный лимит сообщений / хранилища.
//
// Status: 402 Payment Required
//
// Code:   PAYMENT_REQUIRED
func ErrorPaymentRequired(w http.ResponseWriter, r *http.Request, msg string) {
	Error(w, r, http.StatusPaymentRequired, "PAYMENT_REQUIRED", msg, nil)
}

/* 403 */

// ErrorForbidden - 403 FORBIDDEN
//
// Клиент аутентифицирован, но не имеет доступа к ресурсу/действию.
//
//   - Пользователь пытается изменить чужие данные.
//   - Токен без нужного scope/role.
//
// Отличается от 401 тем, что личность известна, но право запрещено.
//
// Status: 403 Forbidden
//
// Code:   FORBIDDEN
func ErrorForbidden(w http.ResponseWriter, r *http.Request, msg string) {
	Error(w, r, http.StatusForbidden, "FORBIDDEN", msg, nil)
}

/* 404 */

// ErrorNotFound - 404 NOT_FOUND
//
// Запрошенный объект или эндпоинт отсутствует.
//
// Передайте `res` («user», «file», «order») для генерации сообщения
// `user not found`; если пусто - вернётся универсальное «Resource not found».
//
// Status: 404 Not Found
//
// Code:   NOT_FOUND
func ErrorNotFound(w http.ResponseWriter, r *http.Request, res string) {
	txt := "Resource not found"
	if res != "" {
		txt = res + " not found"
	}
	Error(w, r, http.StatusNotFound, "NOT_FOUND", txt, nil)
}

/* 405 */

// ErrorMethodNotAllowed - 405 METHOD_NOT_ALLOWED
//
// Клиент использовал HTTP-метод, который не разрешён на данном URL.
// Например, попытка выполнить POST, когда разрешён только GET.
//
// Status: 405 Method Not Allowed
//
// Code:   METHOD_NOT_ALLOWED
func ErrorMethodNotAllowed(w http.ResponseWriter, r *http.Request) {
	Error(w, r, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "Method not allowed", nil)
}

/* 406 */

// ErrorNotAcceptable - 406 NOT_ACCEPTABLE
//
// Сервер не может предоставить ответ в формате, указанном в заголовке Accept.
// Пример: клиент хочет `text/csv`, а сервис поддерживает только `application/json`.
//
// Status: 406 Not Acceptable
//
// Code:   NOT_ACCEPTABLE
func ErrorNotAcceptable(w http.ResponseWriter, r *http.Request, msg string) {
	Error(w, r, http.StatusNotAcceptable, "NOT_ACCEPTABLE", msg, nil)
}

/* 407 */

// ErrorProxyAuthRequired - 407 PROXY_AUTH_REQUIRED
//
// Используется прокси-серверами; приложение-бэкенд почти не применяет.
// Отправляйте, если ваш сервис сам является промежуточным прокси
// и требует авторизацию у клиента.
//
// Status: 407 Proxy Authentication Required
//
// Code:   PROXY_AUTH_REQUIRED
func ErrorProxyAuthRequired(w http.ResponseWriter, r *http.Request, msg string) {
	Error(w, r, http.StatusProxyAuthRequired, "PROXY_AUTH_REQUIRED", msg, nil)
}

/* 408 */

// ErrorRequestTimeout - 408 REQUEST_TIMEOUT
//
// Сервер закрыл соединение, потому что клиент слишком долго не присылал тело
// (обычно > дескриптора ReadTimeout). Эффективно предотвращает удержание
// соединений «зависшими» клиентами.
//
// Status: 408 Request Timeout
//
// Code:   REQUEST_TIMEOUT
func ErrorRequestTimeout(w http.ResponseWriter, r *http.Request, msg string) {
	Error(w, r, http.StatusRequestTimeout, "REQUEST_TIMEOUT", msg, nil)
}

/* 409 */

// ErrorConflict - 409 CONFLICT
//
// Конфликт состояния/уникальности.
//
//   - Попытка создать ресурс с уже существующим уникальным полем (email).
//   - PATCH версии объекта, который успел изменить другой пользователь (ETag).
//
// Status: 409 Conflict
//
// Code:   CONFLICT
func ErrorConflict(w http.ResponseWriter, r *http.Request, msg string) {
	Error(w, r, http.StatusConflict, "CONFLICT", msg, nil)
}

/* 410 */

// ErrorGone - 410 GONE
//
// Ресурс существовал, но удалён необратимо и не будет доступен снова.
// Пример: запись архивирована и политика запрещает её восстановление.
//
// Status: 410 Gone
//
// Code:   GONE
func ErrorGone(w http.ResponseWriter, r *http.Request, res string) {
	txt := "Resource is gone"
	if res != "" {
		txt = res + " is gone"
	}
	Error(w, r, http.StatusGone, "GONE", txt, nil)
}

/* 411 */

// ErrorLengthRequired - 411 LENGTH_REQUIRED
//
// Отсутствует заголовок `Content-Length`, без которого сервер
// не желает принимать тело (например, для потоковой загрузки файлов).
//
// Status: 411 Length Required
//
// Code:   LENGTH_REQUIRED
func ErrorLengthRequired(w http.ResponseWriter, r *http.Request, msg string) {
	Error(w, r, http.StatusLengthRequired, "LENGTH_REQUIRED", msg, nil)
}

/* 412 */

// ErrorPreconditionFailed - 412 PRECONDITION_FAILED
//
// Предусловия (`If-Match`, `If-None-Match`, ETag) не выполнены.
// Часто используется для оптимистических блокировок.
//
// Status: 412 Precondition Failed
//
// Code:   PRECONDITION_FAILED
func ErrorPreconditionFailed(w http.ResponseWriter, r *http.Request, msg string) {
	Error(w, r, http.StatusPreconditionFailed, "PRECONDITION_FAILED", msg, nil)
}

/* 413 */

// ErrorPayloadTooLarge - 413 PAYLOAD_TOO_LARGE
//
// Тело запроса превышает установленный лимит (например, > 10 МБ).
// Укажите лимит в `msg`, чтобы пользователь мог исправиться.
//
// Status: 413 Payload Too Large
//
// Code:   PAYLOAD_TOO_LARGE
func ErrorPayloadTooLarge(w http.ResponseWriter, r *http.Request, msg string) {
	Error(w, r, http.StatusRequestEntityTooLarge, "PAYLOAD_TOO_LARGE", msg, nil)
}

/* 414 */

// ErrorURITooLong - 414 URI_TOO_LONG
//
// URL + query-string столь велики, что сервер их не обрабатывает.
// Обычно случается при GET со сверхбольшими параметрами.
//
// Status: 414 URI Too Long
//
// Code:   URI_TOO_LONG
func ErrorURITooLong(w http.ResponseWriter, r *http.Request, msg string) {
	Error(w, r, http.StatusRequestURITooLong, "URI_TOO_LONG", msg, nil)
}

/* 415 */

// ErrorUnsupportedMediaType - 415 UNSUPPORTED_MEDIA_TYPE
//
// Тип содержимого запроса (`Content-Type`) не поддерживается сервером.
// Типичный пример: загрузили `application/xml`, а API ждёт `application/json`.
//
// Status: 415 Unsupported Media Type
//
// Code:   UNSUPPORTED_MEDIA_TYPE
func ErrorUnsupportedMediaType(w http.ResponseWriter, r *http.Request, msg string) {
	Error(w, r, http.StatusUnsupportedMediaType, "UNSUPPORTED_MEDIA_TYPE", msg, nil)
}

/* 416 */

// ErrorRangeNotSatisfiable - 416 RANGE_NOT_SATISFIABLE
//
// Клиент запросил диапазон файла, выходящий за пределы размера ресурса.
//
// Status: 416 Range Not Satisfiable
//
// Code:   RANGE_NOT_SATISFIABLE
func ErrorRangeNotSatisfiable(w http.ResponseWriter, r *http.Request, msg string) {
	Error(w, r, http.StatusRequestedRangeNotSatisfiable, "RANGE_NOT_SATISFIABLE", msg, nil)
}

/* 417 */

// ErrorExpectationFailed - 417 EXPECTATION_FAILED
//
// Поле `Expect: 100-continue` или иное ожидание не было выполнено сервером.
//
// Status: 417 Expectation Failed
//
// Code:   EXPECTATION_FAILED
func ErrorExpectationFailed(w http.ResponseWriter, r *http.Request, msg string) {
	Error(w, r, http.StatusExpectationFailed, "EXPECTATION_FAILED", msg, nil)
}

/* 418 */

// ErrorTeapot - 418 TEAPOT
//
// Пасхалка из протокола «Hyper Text Coffee Pot Control».
// Используется исключительно в демо/health-check’ах.
//
// Status: 418 I'm a teapot
//
// Code:   TEAPOT
func ErrorTeapot(w http.ResponseWriter, r *http.Request) {
	Error(w, r, http.StatusTeapot, "TEAPOT", "I'm a teapot", nil)
}

/* 421 */

// ErrorMisdirectedRequest - 421 MISDIRECTED_REQUEST
//
// Запрос адресован к серверу, который не может с ним справиться
// (обычно при использовании HTTP/2 + SNI не на тот хост).
//
// Status: 421 Misdirected Request
//
// Code:   MISDIRECTED_REQUEST
func ErrorMisdirectedRequest(w http.ResponseWriter, r *http.Request, msg string) {
	Error(w, r, http.StatusMisdirectedRequest, "MISDIRECTED_REQUEST", msg, nil)
}

/* 422 */

// ErrorUnprocessableEntity - 422 UNPROCESSABLE
//
// Семантически неверный запрос: JSON корректен, но нарушает доменные правила.
// Укажите `det` (map / struct) с подробностями.
//
// Status: 422 Unprocessable Entity
//
// Code:   UNPROCESSABLE
func ErrorUnprocessableEntity(w http.ResponseWriter, r *http.Request, msg string, det interface{}) {
	Error(w, r, http.StatusUnprocessableEntity, "UNPROCESSABLE", msg, det)
}

/* 423 */

// ErrorLocked - 423 LOCKED
//
// Ресурс заблокирован другим пользователем/процессом; операция невозможна,
// пока лок не будет снят.
//
// Status: 423 Locked
//
// Code:   LOCKED
func ErrorLocked(w http.ResponseWriter, r *http.Request, msg string) {
	Error(w, r, http.StatusLocked, "LOCKED", msg, nil)
}

/* 424 */

// ErrorFailedDependency - 424 FAILED_DEPENDENCY
//
// Предыдущая операция в той же цепочке (например, batch) завершилась неудачей.
//
// Status: 424 Failed Dependency
//
// Code:   FAILED_DEPENDENCY
func ErrorFailedDependency(w http.ResponseWriter, r *http.Request, msg string) {
	Error(w, r, http.StatusFailedDependency, "FAILED_DEPENDENCY", msg, nil)
}

/* 425 */

// ErrorTooEarly - 425 TOO_EARLY
//
// Сервер отказывается обработать запрос, присланный слишком рано
// (0-RTT параметры небезопасны). Актуально только с TLS 1.3.
//
// Status: 425 Too Early
//
// Code:   TOO_EARLY
func ErrorTooEarly(w http.ResponseWriter, r *http.Request, msg string) {
	Error(w, r, http.StatusTooEarly, "TOO_EARLY", msg, nil)
}

/* 426 */

// ErrorUpgradeRequired - 426 UPGRADE_REQUIRED
//
// Клиент должен перейти на другой протокол (чаще всего WebSocket).
//
// Status: 426 Upgrade Required
//
// Code:   UPGRADE_REQUIRED
func ErrorUpgradeRequired(w http.ResponseWriter, r *http.Request, msg string) {
	Error(w, r, http.StatusUpgradeRequired, "UPGRADE_REQUIRED", msg, nil)
}

/* 428 */

// ErrorPreconditionRequired - 428 PRECONDITION_REQUIRED
//
// Сервер требует, чтобы запрос содержал `If-Match` / `If-Unmodified-Since`
// или подобные заголовки для безопасного изменения ресурса.
//
// Status: 428 Precondition Required
//
// Code:   PRECONDITION_REQUIRED
func ErrorPreconditionRequired(w http.ResponseWriter, r *http.Request, msg string) {
	Error(w, r, http.StatusPreconditionRequired, "PRECONDITION_REQUIRED", msg, nil)
}

/* 429 */

// ErrorTooManyRequests - 429 RATE_LIMIT
//
// Клиент превысил лимит запросов.
// Передавайте ретри-сообщение и (желательно) заголовок `Retry-After`.
//
// Status: 429 Too Many Requests
//
// Code:   RATE_LIMIT
func ErrorTooManyRequests(w http.ResponseWriter, r *http.Request, msg string) {
	Error(w, r, http.StatusTooManyRequests, "RATE_LIMIT", msg, nil)
}

/* 431 */

// ErrorHeaderFieldsTooLarge - 431 HEADER_FIELDS_TOO_LARGE
//
// Совокупный размер заголовков превышает лимит сервера (cookies, tokens).
//
// Status: 431 Request Header Fields Too Large
//
// Code:   HEADER_FIELDS_TOO_LARGE
func ErrorHeaderFieldsTooLarge(w http.ResponseWriter, r *http.Request, msg string) {
	Error(w, r, http.StatusRequestHeaderFieldsTooLarge, "HEADER_FIELDS_TOO_LARGE", msg, nil)
}

/* 451 */

// ErrorLegalReasons - 451 LEGAL_REASONS
//
// Контент недоступен по юридическим причинам (DMCA, GDPR и т.д.).
// В `msg` желательно добавить ссылку на нормативный акт или запрос.
//
// Status: 451 Unavailable For Legal Reasons
//
// Code:   LEGAL_REASONS
func ErrorLegalReasons(w http.ResponseWriter, r *http.Request, msg string) {
	Error(w, r, http.StatusUnavailableForLegalReasons, "LEGAL_REASONS", msg, nil)
}
