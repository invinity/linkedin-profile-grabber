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
	"github.com/invinity/linkedin-profile-grabber/linkedin"
	"github.com/kofalt/go-memoize"
)

type Controller struct {
	linkedinInst *linkedin.LinkedInBrowser
	cache        *memoize.Memoizer
	lock         sync.Mutex
}

func New(browser *rod.Browser) *Controller {
	return &Controller{linkedinInst: linkedin.NewBrowser(browser), cache: memoize.NewMemoizer(60*time.Minute, 5*time.Minute), lock: sync.Mutex{}}
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
	profile, err, cached := memoize.Call(r.cache, "profile", r.retrieveProfile)
	if err != nil {
		return nil, err
	}
	if cached {
		log.Println("Using cached Profile")
	}
	return profile, nil
}

func (r *Controller) retrieveProfile() (*linkedin.LinkedInProfile, error) {
	r.lock.Lock()
	defer r.lock.Unlock()
	return r.linkedinInst.RetrieveProfile(os.Getenv("LINKEDIN_EMAIL"), os.Getenv("LINKEDIN_PASSWORD"))
}
