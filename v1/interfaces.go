package httpx

// Envelope - универсальный шаблон HTTP-ответа, как для успеха, так и для ошибки.
type Envelope struct {
	Success bool        `json:"success"`            // true/false
	Data    any         `json:"data,omitempty"`     // полезная нагрузка (если success)
	Error   *ErrorBlock `json:"error,omitempty"`    // описание ошибки (если !success)
	TraceID string      `json:"trace_id,omitempty"` // X-Request-ID, сквозной идентификатор
}

// ErrorBlock - структура поля "error" в теле ответа.
type ErrorBlock struct {
	Code    string `json:"code"`              // код ошибки: VALIDATION, INTERNAL, etc
	Message string `json:"message"`           // человекочитаемое сообщение
	Details any    `json:"details,omitempty"` // map[string]string или любая структура
}
