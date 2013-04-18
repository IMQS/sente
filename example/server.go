package main

import (
	"io"
	"log"
	"net/http"
)

func Hello(w http.ResponseWriter, req *http.Request) {
	io.WriteString(w, "hello\n")
}

func StartServer() {
	http.HandleFunc("/hello", Hello)
	err := http.ListenAndServe(":2200", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
