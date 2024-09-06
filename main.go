package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/highercomve/gohtmx/modules/server"
	"github.com/highercomve/gohtmx/modules/server/servermodels"
)

var (
	listenAddr string
)

func main() {
	flag.StringVar(&listenAddr, "listen-addr", ":80", "server listen address")
	flag.Parse()

	listenAddr = strings.Trim(listenAddr, " ")
	logger := log.New(os.Stdout, "gohtmx: ", log.LstdFlags)

	paths := []string{
		filepath.Join("templates", "layout.html"),
		filepath.Join("templates", "index.html"),
		filepath.Join("templates", "error.html"),
	}

	conf := &servermodels.ServerConfig{
		ListenAddr:    listenAddr,
		Logger:        logger,
		TemplatePaths: paths,
	}

	server.Serve(conf)
}
