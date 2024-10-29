package api

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Response struct {
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`
}

const (
	StatusOK    = "OK"
	StatusError = "Error"
)

func ResOK() Response {
	return Response{
		Status: StatusOK,
	}
}

func ResError(msg string) Response {
	return Response{
		Status: StatusError,
		Error:  msg,
	}
}

func WriteHTML(w http.ResponseWriter, statusCode int, v string) error {
	const op = "api.wtireHTML"

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(statusCode)
	_, err := w.Write([]byte(v))
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func WriteText(w http.ResponseWriter, statusCode int, v string) error {
	const op = "api.wtireText"

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(statusCode)
	_, err := w.Write([]byte(v))
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func WriteJSON(w http.ResponseWriter, statusCode int, v any) error {
	const op = "api.wtireJSON"

	data, err := json.Marshal(v)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return fmt.Errorf("%s: %w", op, err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_, err = w.Write(data)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func WriteError(w http.ResponseWriter, statusCode int, msg string) error {
	const op = "api.wtireError"

	data, err := json.Marshal(ResError(msg))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return fmt.Errorf("%s: %w", op, err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_, err = w.Write(data)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}
