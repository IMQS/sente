package main

import (
	"database/sql"
	_ "github.com/bmizerany/pq"
	"io"
	"log"
	"net/http"
)

var MainDBCon *sql.DB

func Hello(w http.ResponseWriter, req *http.Request) {
	io.WriteString(w, "hello\n")
}

func Get(w http.ResponseWriter, req *http.Request) {
	cookie, err := req.Cookie("ses")
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		io.WriteString(w, "Unauthorised")
		return
	}

	row := MainDBCon.QueryRow("SELECT \"Cookie\" FROM \"Session\" WHERE \"Cookie\" = $1", cookie.Value)
	var cookieVal string
	err = row.Scan(&cookieVal)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		io.WriteString(w, "Unauthorised")
		return
	} else {
		io.WriteString(w, "Authorised : "+cookie.Value)
	}
}

func StartServer() {
	var err error
	MainDBCon, err = sql.Open("postgres", "user=imqs password=1mq5p@55w0rd dbname=main sslmode=disable")
	if err != nil {
		log.Fatal("Open DB :", err)
	}
	defer MainDBCon.Close()

	http.HandleFunc("/hello", Hello)
	http.HandleFunc("/get", Get)
	var h http.Handler
	err = http.ListenAndServe(":2200", h)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
