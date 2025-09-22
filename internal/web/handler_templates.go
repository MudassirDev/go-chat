package web

import "net/http"

func (c *apiConfig) handleLoginTemplate(w http.ResponseWriter, r *http.Request) {
	data := struct {
		AuthType bool
	}{
		AuthType: true,
	}

	c.templates.ExecuteTemplate(w, "auth", data)
}

func (c *apiConfig) handleRegisterTemplate(w http.ResponseWriter, r *http.Request) {
	c.templates.ExecuteTemplate(w, "auth", nil)
}
