package web

import "net/http"

func (c *APIConfig) handlerLoginTemplate(w http.ResponseWriter, r *http.Request) {
	data := struct {
		AuthType bool
	}{
		AuthType: true,
	}

	c.Templates.ExecuteTemplate(w, "auth", data)
}

func (c *APIConfig) handlerRegisterTemplate(w http.ResponseWriter, r *http.Request) {
	c.Templates.ExecuteTemplate(w, "auth", nil)
}
