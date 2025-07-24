package httpx

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/middleware"
)

// Error формирует структурированный ответ с ошибкой (status 4xx, 5xx)
func Error(w http.ResponseWriter, r *http.Request, status int, code, message string, details interface{}) {
	traceID := middleware.GetReqID(r.Context())

	resp := Envelope{
		Success: false,
		Error: &ErrorBlock{
			Code:    code,
			Message: message,
			Details: details,
		},
		TraceID: traceID,
	}

	writeJSON(w, status, resp)
}

// JSON возвращает успешный ответ с заданным HTTP-статусом.
//
// Пример:
//
//	httpx.JSON(w, r, http.StatusCreated, myObject)
func JSON(w http.ResponseWriter, r *http.Request, status int, data interface{}) {
	traceID := middleware.GetReqID(r.Context())

	resp := Envelope{
		Success: true,
		Data:    data,
		TraceID: traceID,
	}

	writeJSON(w, status, resp)
}

// Вспомогательная функция для отправки JSON ответов
func writeJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}
