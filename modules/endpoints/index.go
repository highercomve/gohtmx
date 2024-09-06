package endpoints

import (
	"net/http"

	"github.com/highercomve/gohtmx/modules/server/servermodels"
)

func HandleIndex(opts *servermodels.ServerOptions) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		default:
			getIndex(w, r, opts)
			return
		}
	})
}

func getIndex(w http.ResponseWriter, r *http.Request, opts *servermodels.ServerOptions) {
	data := servermodels.Response{
		Code:  200,
		Error: nil,
		Data:  map[string]string{},
	}
	render(w, r, opts, "index.html", &data)
}
