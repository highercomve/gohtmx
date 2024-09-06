package server

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"os/signal"
	"reflect"
	"sync/atomic"
	"time"

	"github.com/highercomve/gohtmx/modules/endpoints"
	"github.com/highercomve/gohtmx/modules/server/servermodels"
)

type key int

const (
	requestIDKey key = 0
)

var (
	healthy int32
)

var functions template.FuncMap = template.FuncMap{
	"notNil": notNil,
}

const enableRequestLog = false

func Serve(conf *servermodels.ServerConfig) {
	conf.Logger.Println("Server is starting...")

	t := template.New("").Funcs(functions)
	tmpl, err := t.ParseFiles(conf.TemplatePaths...)
	if err != nil {
		conf.Logger.Fatalf("Can't load templates: %v\n", err)
	}
	opts := &servermodels.ServerOptions{
		Logger:     conf.Logger,
		TpmlEngine: tmpl,
	}
	router := http.NewServeMux()
	router.Handle("/", endpoints.HandleIndex(opts))
	router.Handle("/networks", endpoints.GetNetworksList(opts))
	router.Handle(
		"/static/",
		http.StripPrefix(
			"/static/",
			http.FileServer(http.Dir("./static")),
		),
	)

	nextRequestID := func() string {
		return fmt.Sprintf("%d", time.Now().UnixNano())
	}

	server := &http.Server{
		Addr:         conf.ListenAddr,
		Handler:      tracing(nextRequestID)(logging(conf.Logger)(router)),
		ErrorLog:     conf.Logger,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	done := make(chan bool)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	go func() {
		<-quit
		conf.Logger.Println("Server is shutting down...")
		atomic.StoreInt32(&healthy, 0)

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		server.SetKeepAlivesEnabled(false)
		if err := server.Shutdown(ctx); err != nil {
			conf.Logger.Fatalf("Could not gracefully shutdown the server: %v\n", err)
		}
		close(done)
	}()

	conf.Logger.Println("Server is ready to handle requests at", conf.ListenAddr)

	atomic.StoreInt32(&healthy, 1)

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		conf.Logger.Fatalf("Could not listen on %s: %v\n", conf.ListenAddr, err)
	}

	<-done
	conf.Logger.Println("Server stopped")
}

func notNil(i interface{}) bool {
	// First, check if the interface itself is nil
	if i == nil {
		return false
	}

	// Then use reflect package to check if the underlying value is nil
	v := reflect.ValueOf(i)
	return !(v.Kind() == reflect.Ptr && v.IsNil())
}

func healthz() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if atomic.LoadInt32(&healthy) == 1 {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		w.WriteHeader(http.StatusServiceUnavailable)
	})
}

func logging(logger *log.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if enableRequestLog {
					requestID, ok := r.Context().Value(requestIDKey).(string)
					if !ok {
						requestID = "unknown"
					}
					logger.Println(requestID, r.Method, r.URL.Path, r.RemoteAddr, r.UserAgent())
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}

func tracing(nextRequestID func() string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestID := r.Header.Get("X-Request-Id")
			if requestID == "" {
				requestID = nextRequestID()
			}
			ctx := context.WithValue(r.Context(), requestIDKey, requestID)
			w.Header().Set("X-Request-Id", requestID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
