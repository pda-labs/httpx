package httpx

import "net/http"

/* 200 */

// Ok - 200 Ok
//
// Универсальный шорткат для успешного ответа по умолчанию.
// Используйте, когда запрос выполнен успешно и сервер готов
// вернуть полезные данные (JSON) без каких‑либо побочных требований.
//
// Типовые сценарии:
//   - GET‑эндпойнт возвращает ресурс/список ресурсов;
//   - PUT/PATCH успешно обновил объект и возвращает его актуальное состояние;
//   - Health‑check, если хочется вернуть структурированное «status: ok».
//
// Rекомендации:
//   - В поле data передавайте DTO/структуру, пригодную для фронта;
//   - Если ответ должен быть пустым → выберите NoContent(204), чтобы
//     не посылать лишнее тело;
//
// Status: 200 Ok
// Code:   - (успешным ответам machine‑code не требуется)
//
// Пример:
//
//	article, _ := svc.Get(id)
//	httpx.Ok(w, r, article)
func Ok(w http.ResponseWriter, r *http.Request, data interface{}) {
	JSON(w, r, http.StatusOK, data)
}

/* 201 */

// Created - 201 CREATED
//
// Ресурс успешно создан. Передайте canonical‑URL в location
// (например, "/users/123"). Если location пустой - заголовок
// не ставится.
//
// Status: 201 Created
func Created(w http.ResponseWriter, r *http.Request, location string, data any) {
	if location != "" {
		w.Header().Set("Location", location)
	}
	JSON(w, r, http.StatusCreated, data)
}

/* 202 */

// Accepted - 202 ACCEPTED
//
// Запрос принят на асинхронную обработку. Верните DTO со статусом
// или Job‑ID, по которому клиент сможет опросить результат.
//
// Status: 202 Accepted
func Accepted(w http.ResponseWriter, r *http.Request, data any) {
	JSON(w, r, http.StatusAccepted, data)
}

/* 203 */

// NonAuthoritative - 203 NON_AUTHORITATIVE_INFORMATION
//
// Содержимое получено из третьего источника (прокси, кеш).
// Используется редко.
//
// Status: 203 Non‑Authoritative Information
func NonAuthoritative(w http.ResponseWriter, r *http.Request, data any) {
	JSON(w, r, http.StatusNonAuthoritativeInfo, data)
}

/* 204 */

// NoContent - 204 NO_CONTENT
//
// Успешно, но тело не требуется (DELETE, PUT‑idempotent, healthz).
// Не возвращаем Envelope, чтобы не нарушать спецификацию.
//
// Status: 204 No Content
func NoContent(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNoContent)
}

/* 205 */

// ResetContent - 205 RESET_CONTENT
//
// Клиент должен сбросить форму/UI. Аналогично 204 - без тела.
//
// Status: 205 Reset Content
func ResetContent(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusResetContent)
}

/* 206 */

// PartialContent - 206 PARTIAL_CONTENT
//
// Ответ на запрос с Range‑заголовком (скачивание куска файла).
//
// Status: 206 Partial Content
func PartialContent(w http.ResponseWriter, r *http.Request, data any) {
	JSON(w, r, http.StatusPartialContent, data)
}
