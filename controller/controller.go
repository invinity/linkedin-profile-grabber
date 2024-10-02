package controller

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/go-rod/rod"
	"github.com/invinity/linkedin-profile-grabber/linkedin"
	"github.com/kofalt/go-memoize"
)

type Controller struct {
	linkedinInst *linkedin.LinkedIn
	cache        *memoize.Memoizer
}

func New(browser *rod.Browser) *Controller {
	return &Controller{linkedinInst: linkedin.New(browser), cache: memoize.NewMemoizer(1*time.Minute, 1*time.Minute)}
}

func (r *Controller) GetLinkedInProfile(w http.ResponseWriter, req *http.Request) {
	profile, err := r.retrieveProfile()
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	} else {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		json.NewEncoder(w).Encode(profile)
	}
}

func (r *Controller) retrieveProfile() (*linkedin.LinkedInProfile, error) {
	profile, err, cached := memoize.Call(r.cache, "profile", r.linkedinInst.RetrieveProfile)
	if err != nil {
		return nil, err
	}
	if cached {
		log.Println("Using cached Profile")
	}
	return profile, nil
}
