package main

import (
	"log"
	"net/http"
	"os"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/invinity/linkedin-profile-grabber/routes"
)

func main() {
	path, _ := launcher.LookPath()
	browser := rod.New().ControlURL(launcher.New().Leakless(false).NoSandbox(true).Bin(path).MustLaunch()).Trace(true).MustConnect()
	defer browser.MustClose()
	router := routes.AppRoutes(browser)
	http.Handle("/api", router)

	// Determine port for HTTP service.
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("defaulting to port %s", port)
	}

	log.Println("Listening...")
	log.Fatal(http.ListenAndServe(":"+port, router))
}
