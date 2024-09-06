package endpoints

import (
	"encoding/json"
	"net/http"
	"strings"
)

const (
	ContentTypeJson = "application/json"
)

var availableFormats []string = []string{"json", "html", "css"}

func getFormat(c string) string {
	format := c
	for _, f := range availableFormats {
		if strings.Contains(c, f) {
			format = f
		}
	}
	return format
}

type Response struct {
	Data  interface{} `json:"data"`
	Error error       `json:"error"`
	Code  int         `json:"code"`
}

func render(w http.ResponseWriter, r *http.Request, opts *servermodels.ServerOptions, name string, res *Response) {
	var err error
	if res.Error != nil {
		opts.Logger.Println(res.Error.Error())
	}
	contentType := r.Header.Get("Accept")
	format := getFormat(contentType)

	switch format {
	case "json":
		enc := json.NewEncoder(w)
		w.Header().Set("Content-Type", ContentTypeJson)
		if res.Code > 0 {
			w.WriteHeader(res.Code)
		}
		if res.Error != nil {
			data := map[string]string{
				"error": res.Error.Error(),
			}
			err = enc.Encode(data)
		} else {
			err = enc.Encode(res.Data)
		}
		if err != nil {
			opts.Logger.Println(err.Error())
			w.WriteHeader(500)
			w.Write([]byte(err.Error()))
		}
	default:
		err := opts.TpmlEngine.ExecuteTemplate(w, name, res)
		if err != nil {
			opts.Logger.Println(err.Error())
			w.WriteHeader(500)
			w.Write([]byte(err.Error()))
		}
	}
}
