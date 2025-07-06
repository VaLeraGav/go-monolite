package respond

import (
	"encoding/json"
	"io"
	"net/http"
)

const (
	StatusSuccess = "success"
	StatusError   = "error"
)

type Response struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
	Errors  any    `json:"errors,omitempty"`
}

type ErrorResponse struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
	Errors  any    `json:"errors,omitempty"`
}

type SuccessResponse struct {
	Status  string `json:"status"`
	Message string `json:"message" swaggertype:"string"`
	Data    any    `json:"data,omitempty"`
}

func ErrorHandler(w http.ResponseWriter, r *http.Request, code int, err any, mess ...any) {
	var errorsOut any
	switch v := err.(type) {
	case error:
		errorsOut = v.Error()
	default:
		errorsOut = v
	}

	var userMess string
	if len(mess) > 0 && mess[0] != nil {
		extraErrors := mess[0]
		switch v := extraErrors.(type) {
		case string:
			userMess = v
		case error:
			userMess = v.Error()
		default:
			userMess = ""
		}
	}

	Respond(w, r, code, ErrorResponse{
		Status:  StatusError,
		Message: userMess,
		Errors:  errorsOut,
	})
}

func SuccessHandler(w http.ResponseWriter, r *http.Request, code int, message string, data ...any) {
	if len(data) > 0 && data[0] != nil {
		Respond(w, r, code, SuccessResponse{
			Status:  StatusSuccess,
			Message: message,
			Data:    data[0],
		})
	} else {
		Respond(w, r, code, SuccessResponse{
			Status:  StatusSuccess,
			Message: message,
			Data:    nil,
		})
	}
}

func Respond(w http.ResponseWriter, _ *http.Request, code int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	if data == nil {
		return
	}

	if err := json.NewEncoder(w).Encode(data); err != nil {
		fallback := map[string]string{
			"status":  StatusError,
			"message": "failed to encode response",
		}
		b, _ := json.Marshal(fallback)
		http.Error(w, string(b), http.StatusInternalServerError)
	}
}

func ParseBody(w http.ResponseWriter, r *http.Request) []byte {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		ErrorHandler(w, r, http.StatusBadRequest, "ошибка чтения запроса")
		return nil
	}
	defer r.Body.Close()

	return body
}
