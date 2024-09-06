package servermodels

import (
	"html/template"
	"log"
)

type ServerConfig struct {
	ListenAddr    string
	Logger        *log.Logger
	TemplatePaths []string
}

type ServerOptions struct {
	Logger     *log.Logger
	TpmlEngine *template.Template
}

type Response struct {
	Data  interface{} `json:"data"`
	Error error       `json:"error"`
	Code  int         `json:"code"`
}
