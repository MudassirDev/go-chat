package main

import "net/http"

func handleLoginTemplate(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "auth", struct {
		AuthType bool
	}{
		AuthType: true,
	})
}

func handleRegisterTemplate(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "auth", nil)
}
