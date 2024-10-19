package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/invinity/linkedin-profile-grabber/cache"
	"github.com/invinity/linkedin-profile-grabber/controller"
	"github.com/invinity/linkedin-profile-grabber/linkedin"
	"github.com/invinity/linkedin-profile-grabber/routes"
	"github.com/kofalt/go-memoize"
)

func main() {
	browser := createBrowser()
	defer browser.MustClose()
	cache, err := createCache()
	if err != nil {
		log.Fatal(err)
	}
	defer cache.Close()
	retriever := createRetriever(&cache, linkedin.NewBrowser(browser))
	router := routes.AppRoutes(retriever)
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

func createBrowser() *rod.Browser {
	timeout, _ := time.ParseDuration("180s")
	path, found := launcher.LookPath()
	if !found {
		log.Fatal("Did not find chrome in go-rod standard locations")
	}
	log.Printf("Using detected chrome path: %s\n", path)
	return rod.New().ControlURL(launcher.New().Leakless(false).NoSandbox(true).Headless(true).Bin(path).MustLaunch()).Timeout(timeout).Trace(true).MustConnect()
}

func createCache() (cache.Cache, error) {
	ctx := context.Background()
	return cache.NewGoogleStorageCache(&ctx, "linkedin-profile-grabber")
}

func createMemozier() *memoize.Memoizer {
	cacheTime := os.Getenv("CACHE_TIME")
	if cacheTime == "" {
		cacheTime = "4h"
	}
	cacheDuration, err := time.ParseDuration(cacheTime)
	if err != nil {
		log.Fatal("unable to parse cache time value: " + cacheTime)
	}
	return memoize.NewMemoizer(cacheDuration, 5*time.Minute)
}

type RealLinkedInProfileRetriever struct {
	browser *linkedin.LinkedInBrowser
}

func (r RealLinkedInProfileRetriever) Get() (*linkedin.LinkedInProfile, error) {
	email, password := os.Getenv("LINKEDIN_EMAIL"), os.Getenv("LINKEDIN_PASSWORD")
	if email != "" && password != "" {
		return r.browser.RetrieveProfileViaLogin(email, password)
	} else {
		return r.browser.RetrieveProfileViaSearch("matthew", "pitts", "mattpitts")
	}
}

func createRetriever(cache *cache.Cache, browser *linkedin.LinkedInBrowser) controller.LinkedinProfileRetriever {
	return controller.NewCacheHandlingRetriever(*cache, &RealLinkedInProfileRetriever{browser: browser})
}
