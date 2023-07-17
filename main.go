package main

import (
	"fmt"
	"log"
	"net/http"
)

func RootHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, World! Welcome Go App!")
}

func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "ok")
}

func main() {
	log.Println("Starting server on :8080")
	http.HandleFunc("/", RootHandler)
	http.HandleFunc("/health_check", HealthCheckHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
