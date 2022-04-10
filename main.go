package main

import (
	"log"
	"movie_store/handlers"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", handlers.MainPageHandler)
	r.HandleFunc("/login", handlers.Login).Methods("GET")
	r.HandleFunc("/login", handlers.DoLogin).Methods("POST")
	r.HandleFunc("/images/{image_filename:.+}", handlers.ImageHandler)

	r.Handle("/{movie_id:\\d+}", handlers.AuthMiddleware(http.HandlerFunc(handlers.MoviePageHandler)))

	http.Handle("/", r)
	log.Fatal(http.ListenAndServe(":8888", nil)) // 121.199.60.139
}
