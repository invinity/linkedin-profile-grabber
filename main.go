package main

import (
	"log"
	"net/http"

	"github.com/invinity/linkedin-profile-grabber/routes"
)

func main() {
	router := routes.AppRoutes()
	http.Handle("/api", router)

	log.Println("Listening...")
	log.Fatal(http.ListenAndServe(":8081", router))
}
