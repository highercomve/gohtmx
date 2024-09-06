package endpoints

import (
	"fmt"
	"net/http"

	"github.com/highercomve/gohtmx/modules/nm"
	"github.com/highercomve/gohtmx/modules/nm/nmmodules"
	"github.com/highercomve/gohtmx/modules/server/servermodels"
)

func GetNetworksList(opts *servermodels.ServerOptions) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "POST":
			handleSaveNetwork(w, r, opts)
			return
		default:
			handleNetworkList(w, r, opts)
			return
		}
	})
}

func handleSaveNetwork(w http.ResponseWriter, r *http.Request, opts *servermodels.ServerOptions) {
	ssid := r.FormValue("ssid")
	password := r.FormValue("password")

	fmt.Printf("%s %s\n", ssid, password)
	networkmanager := nm.Init()

	// Call the function to connect using nmcli
	err := networkmanager.Save(ssid, password)
	if err != nil {
		opts.Logger.Println(err)
		data := servermodels.Response{
			Code:  http.StatusOK,
			Data:  nil,
			Error: err,
		}
		render(w, r, opts, "error.html", &data)
		return
	}

	data := servermodels.Response{
		Code: http.StatusOK,
		Data: nmmodules.WifiConn{
			ID:   ssid,
			SSID: ssid,
		},
		Error: nil,
	}
	opts.Logger.Printf("Connected successfully to %s\n", ssid)
	render(w, r, opts, "", &data)
}

func handleNetworkList(w http.ResponseWriter, r *http.Request, opts *servermodels.ServerOptions) {
	networkmanager := nm.Init()
	conns, err := networkmanager.List()
	if err != nil {
		opts.Logger.Println(err)
		data := servermodels.Response{
			Code:  http.StatusInternalServerError,
			Data:  nil,
			Error: err,
		}
		render(w, r, opts, "error.html", &data)
		return
	}

	data := servermodels.Response{
		Code:  http.StatusOK,
		Data:  conns,
		Error: nil,
	}
	render(w, r, opts, "network_list", &data)
}
