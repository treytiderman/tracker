package main

import (
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"time"
)

func handler_base_url(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/trackers", http.StatusSeeOther)
}

func handler_time(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, time.Now().Format(time.UnixDate))
}

func http_server_start() error {
	mux := http.NewServeMux()

	// Routes
	mux.Handle("/", mw_logger(http.HandlerFunc(handler_base_url)))
	mux.Handle("/time", mw_logger(http.HandlerFunc(handler_time)))
	mux.Handle("/public/", mw_logger(http.StripPrefix("/public/", http.FileServer(http.Dir("../public")))))

	handle_routes_ui(mux)
	handle_routes_page(mux)
	handle_routes_api_htmx(mux)

	// Get HTTP Port
	port := os.Getenv("HTTP_PORT")
	if port == "" {
		port = "8000"
	}
	port = fmt.Sprintf(":%s", port)

	// Start Web Server
	slog.Info("http server started", "url", "http://localhost"+port)
	err := http.ListenAndServe(port, mux)

	return err
}

// Middleware

func mw_logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/favicon.ico" {
			r.ParseForm()
			slog.Debug("http request",
				"method", r.Method,
				"url", r.URL.Path,
				"query", r.URL.RawQuery,
				"form", r.Form,
				"ip", r.RemoteAddr,
			)
		}
		next.ServeHTTP(w, r)
	})
}

func mw_auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token_is_valid := auth_token_valid(w, r)
		if !token_is_valid {
			w.Write([]byte("token is not valid"))
			return
		}
		next.ServeHTTP(w, r)
	})
}

func mw_read_only(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		READ_ONLY := os.Getenv("READ_ONLY")
		if READ_ONLY == "true" {
			w.Write([]byte("not allowed in READ_ONLY mode"))
			return
		}
		next.ServeHTTP(w, r)
	})
}
