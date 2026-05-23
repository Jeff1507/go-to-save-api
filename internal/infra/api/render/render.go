package render

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type response map[string]any

func Bind(r *http.Request, data any) error {
	if err := json.NewDecoder(r.Body).Decode(data); err != nil {
		return fmt.Errorf("decode body: %v", err)
	}

	return nil
}

func NoContent(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNoContent)
}

func writeJSON(w http.ResponseWriter, status int, res response) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(res)
}

func Success(w http.ResponseWriter, status int, data any) {
	writeJSON(w, status, response{
		"data": data,
	})
}

func ValidationError(w http.ResponseWriter, errs error) {
	writeJSON(w, http.StatusUnprocessableEntity, response{
		"errors": errs,
	})
}

func Error(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, response{
		"error": msg,
	})
}

func InternalServerError(w http.ResponseWriter) {
	Error(w, http.StatusInternalServerError, "Erro interno no servidor")
}

func InvalidRequestBody(w http.ResponseWriter) {
	Error(w, http.StatusBadRequest, "Corpo da requisição inválido")
}

func InvalidMultipartForm(w http.ResponseWriter) {
	Error(w, http.StatusBadRequest, "Formulário inválido")
}

func Unauthorized(w http.ResponseWriter) {
	Error(w, http.StatusUnauthorized, "Sem autorização")
}
