package controller

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/go-rod/rod"
	"github.com/invinity/linkedin-profile-grabber/cache"
	"github.com/invinity/linkedin-profile-grabber/linkedin"
)

type Controller struct {
	linkedinInst *linkedin.LinkedInBrowser
	cache        *cache.Cache
	lock         sync.Mutex
}

func NewController(browser *rod.Browser, cache *cache.Cache) *Controller {
	return &Controller{linkedinInst: linkedin.NewBrowser(browser), cache: cache, lock: sync.Mutex{}}
}

func (r *Controller) GetLinkedInProfile(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", req.Header.Get("Origin"))
	w.Header().Set("Access-Control-Allow-Methods", "GET,OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type,Accept")
	w.Header().Set("Cache-Control", fmt.Sprintf("max-age=%.0f", (60*time.Minute).Seconds()))
	profile, err := r.getProfile()
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	} else {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		json.NewEncoder(w).Encode(profile)
	}
}

func (r *Controller) getProfile() (*linkedin.LinkedInProfile, error) {
	var storedProfile *linkedin.LinkedInProfile
	err := r.cache.Get("myprofile", &storedProfile)
	if err != nil {
		log.Println("error during profile fetch from bucket", err)
	}
	var age time.Duration
	if storedProfile != nil {
		age = time.Since(storedProfile.GeneratedAt)
	}

	if storedProfile == nil || age >= 4*time.Hour {
		log.Println("Stored profile data is too old or empty, attempting to retrieve fresh data.")
		storedProfile, err = r.retrieveProfile()
		if err != nil {
			log.Println("error during linked in profile retrieval", err)
			if storedProfile != nil {
				log.Println("stored profile was present, just returning that for now")
				return storedProfile, nil
			}
			return nil, err
		}
		log.Println("storing profile for caching")
		err = r.cache.Put("myprofile", storedProfile)
		if err != nil {
			return nil, err
		}
	} else {
		log.Println("using cached profile copy")
	}
	return storedProfile, nil
}

func (r *Controller) retrieveProfile() (*linkedin.LinkedInProfile, error) {
	r.lock.Lock()
	defer r.lock.Unlock()
	email, password := os.Getenv("LINKEDIN_EMAIL"), os.Getenv("LINKEDIN_PASSWORD")
	if email != "" && password != "" {
		return r.linkedinInst.RetrieveProfileViaLogin(email, password)
	} else {
		return r.linkedinInst.RetrieveProfileViaSearch("matthew", "pitts", "mattpitts")
	}
}
