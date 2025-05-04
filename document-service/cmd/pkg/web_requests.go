package pkg

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type ResponseJSON struct {
	Data interface{} `json:"data"`
}

type ErrorResponse struct {
	Status  int    `json:"-"`
	Code    string `json:"code"`
	Message string `json:"message"`
}

func Success(w http.ResponseWriter, status int, data interface{}) {
	RespondWithJSON(w, status, ResponseJSON{Data: data})
}

func Error(w http.ResponseWriter, status int, format string, args ...interface{}) {
	err := ErrorResponse{
		Code:    strings.ReplaceAll(strings.ToLower(http.StatusText(status)), " ", "_"),
		Message: fmt.Sprintf(format, args...),
		Status:  status,
	}

	RespondWithJSON(w, status, err)
}

func RespondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
