package main

import (
	"io"
	"io/fs"
	"log/slog"
	"net/http"
	"speedy/frontend"
	"strconv"
	"time"
)

// backend for simple speedtest (up/down) and ping
// for the web frontend

func main() {
	mux := http.NewServeMux()

	sub, err := fs.Sub(frontend.Assets, "dist")
	if err != nil {
		slog.Error("failed to create sub filesystem", "error", err)
		return
	}

	mux.HandleFunc("/", http.FileServerFS(sub).ServeHTTP)

	mux.HandleFunc("/download", func(w http.ResponseWriter, r *http.Request) {
		size, err := strconv.Atoi(r.URL.Query().Get("size"))
		if err != nil {
			http.Error(w, "invalid size", http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/octet-stream")
		w.Header().Set("Content-Length", strconv.Itoa(size))
		w.WriteHeader(http.StatusOK)

		data := make([]byte, size)
		_, err = w.Write(data)
		if err != nil {
			http.Error(w, "error writing body", http.StatusInternalServerError)
			return
		}
	})

	mux.HandleFunc("/upload", func(w http.ResponseWriter, r *http.Request) {
		size, err := strconv.Atoi(r.URL.Query().Get("size"))
		if err != nil {
			http.Error(w, "invalid size", http.StatusBadRequest)
			return
		}

		// read everything, but don't store it
		// to ensure the client sends everything
		_, err = io.CopyN(io.Discard, r.Body, int64(size))
		if err != nil {
			http.Error(w, "error reading body", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	})

	mux.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		unixMillis := time.Now().UnixMilli()

		w.Write([]byte(strconv.FormatInt(unixMillis, 10)))
	})

	http.ListenAndServe(":5689", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		mux.ServeHTTP(w, r)
	}))
}
