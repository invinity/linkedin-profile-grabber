package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/invinity/linkedin-profile-grabber/routes"
)

func main() {
	timeout, _ := time.ParseDuration("180s")
	path, found := launcher.LookPath()
	if !found {
		log.Fatal("Did not find chrome in go-rod standard locations")
	}
	log.Printf("Using detected chrome path: %s\n", path)
	browser := rod.New().ControlURL(launcher.New().Leakless(false).NoSandbox(true).Headless(false).Bin(path).KeepUserDataDir().MustLaunch()).Timeout(timeout).Trace(true).MustConnect()
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
