package main

import (
	"log"
	"net/http"
	"os"
	"os/exec"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/invinity/linkedin-profile-grabber/routes"
)

func main() {
	path, found := launcher.LookPath()
	if !found {
		log.Println("Did not find chrome in go-rod standard locations")
		log.Println("Looking for chrome in apt")
		cmd := exec.Command("/bin/sh", "-c", "dpkg -L chromium")
		_, err := cmd.Output()
		if err != nil {
			log.Fatal("Unable to find chrome: ", err)
		} else {
			path, err = exec.LookPath("chrome")
			if err != nil {
				log.Fatal("We can't seem to find chrome")
			}
		}
	}
	log.Printf("Using detected chrome path: %s\n", path)
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
