package main

import (
	"fmt"
	"log"
	"net/http"

	"example.com/chirpy/handlers"
)

type homeHandler struct{}

func (homeHandler) ServeHTTP(http.ResponseWriter, *http.Request) {
	fmt.Println("Serving index.html")
	http.FileServer(http.Dir("index.html"))
}

func main() {
	mux := http.NewServeMux()

	mux.Handle("/app", http.StripPrefix("/app", http.FileServer(http.Dir("./"))))
	mux.Handle("/app/assets/", http.StripPrefix("/app/assets/", http.FileServer(http.Dir("./assets/"))))

	mux.HandleFunc("GET /healthz", handlers.GETHealthzHandler)

	server := http.Server{
		Handler: mux,
		Addr:    ":8080",
	}

	log.Fatal(server.ListenAndServe())
}
