package web

import (
	"encoding/json"
	"log"
	"net/http"
)

func respondWithError(w http.ResponseWriter, err error, msg string, code int) {
	log.Printf("error: %v", err)

	data := struct {
		Message string `json:"message"`
	}{
		Message: msg,
	}

	respondWithJson(w, code, data)
}

func respondWithJson(w http.ResponseWriter, code int, payload any) {
	data, err := json.Marshal(payload)
	if err != nil {
		log.Printf("failed to encode response: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("failed to send data"))
		return
	}

	w.WriteHeader(code)
	w.Write(data)
}
