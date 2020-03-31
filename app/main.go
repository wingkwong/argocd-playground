package main

import (
	"io"
	"log"
	"net/http"
)

func main() {

	http.HandleFunc("/", Handler)

	if err := http.ListenAndServe(":8888", nil); err != nil {
		log.Fatal(err)
	}
}

func Handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	io.WriteString(w, `{"status":"ok", "message": "hello-world"}`)
}
