package endpoints

import (
	"net/http"
)

func HandleIndex(opts *ServerOptions) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		default:
			getIndex(w, r, opts)
			return
		}
	})
}

func getIndex(w http.ResponseWriter, r *http.Request, opts *models.ServerOptions) {
	render(w, r, opts, "index.html", &data)
}
