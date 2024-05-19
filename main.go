package main

import (
	"errors"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

func main() {
	handler := func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "html/index.html")
	}
	getHandler := func(w http.ResponseWriter, r *http.Request) {
		fileName := r.URL.Query().Get("name")

		if !(strings.Contains(fileName, "html")) && !(strings.Contains(fileName, "uploads")) {
			w.Header().Set("WWW-Authenticate", `Basic realm="Restricted`)
			http.Error(w, "Forbidden", 403)
			return
		}

		_, err := os.Stat(fileName)
		if fileName == "" || errors.Is(err, os.ErrNotExist) {
			http.Error(w, "File not found", 404)
			return
		}

		http.ServeFile(w, r, fileName)
	}
	uploadHandler := func(w http.ResponseWriter, r *http.Request) {
		r.ParseMultipartForm(32 << 20)
		file, handler, err := r.FormFile("file")
		if err != nil {
			log.Println(err)
			return
		}
		f, err := os.OpenFile(handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			log.Println(err)
			return
		}
		io.Copy(f, file)
		defer file.Close()
		log.Printf("Uploaded file: %+v", handler.Filename)
	}
	deleteHandler := func(w http.ResponseWriter, r *http.Request) {
		fileName := r.URL.Query().Get("name")

		if !(strings.Contains(fileName, "html")) && !(strings.Contains(fileName, "uploads")) {
			w.Header().Set("WWW-Authenticate", `Basic realm="Restricted`)
			http.Error(w, "Forbidden", 403)
			return
		}

		_, err := os.Stat(fileName)
		if fileName == "" || errors.Is(err, os.ErrNotExist) {
			http.Error(w, "File not found", 404)
			return
		}

		e := os.Remove(fileName)
		if e != nil {
			log.Println(e)
			return
		}
	}

	http.HandleFunc("GET /{$}", handler)
	http.HandleFunc("GET /file", getHandler)
	http.HandleFunc("DELETE /file", deleteHandler)
	http.HandleFunc("POST /upload", uploadHandler)
	// TODO
	http.HandleFunc("PUT /upload", handler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
