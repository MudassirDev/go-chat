package main

import "net/http"

func handleLoginTemplate(w http.ResponseWriter, r *http.Request) {
	data := struct {
		AuthType bool
	}{
		AuthType: true,
	}

	templates.ExecuteTemplate(w, "auth", data)
}

func handleRegisterTemplate(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "auth", nil)
}
