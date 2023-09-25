package main

import (
	"log"
	"net/http"
)


func handlePing(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("pong"))
}

func main() {
	http.HandleFunc("/ping", handlePing)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
