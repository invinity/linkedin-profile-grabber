package main

import (
	"log"
	"net/http"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/invinity/linkedin-profile-grabber/routes"
)

func main() {
	browser := rod.New().ControlURL(launcher.New().Leakless(false).MustLaunch()).Trace(true).MustConnect()
	defer browser.MustClose()
	router := routes.AppRoutes(browser)
	http.Handle("/api", router)

	log.Println("Listening...")
	log.Fatal(http.ListenAndServe(":8081", router))
}
