package main

import (
	"ascii-art-web/handlers"
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", handlers.Home)
	mux.HandleFunc("/ascii-art", handlers.Output)
	mux.Handle("/static/", http.StripPrefix("/static", http.FileServer(http.Dir("./web/static/"))))
	log.Println("start server in http://localhost:1337")
	err := http.ListenAndServe(":1337", mux)
	if err != nil {
		log.Fatal(err)
	}
}
