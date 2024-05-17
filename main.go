package main

import (
	"io"
	"log"
	"net/http"
)

func main() {
	handler := func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			io.WriteString(w, "This is a get request...\n")
		case "POST":
			io.WriteString(w, "This is a post request...\n")
		case "PUT":
			io.WriteString(w, "This is a put request...\n")
		default:
			r.Header.Set("Allow", "GET, POST, PUT")
			http.Error(w, "Method Not Allowed", 405)
		}
	}
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
