package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/invinity/linkedin-profile-grabber/linkedin"
)

type Controller struct {
	retriever LinkedinProfileRetriever
	lock      sync.Mutex
}

func NewController(retriever LinkedinProfileRetriever) *Controller {
	return &Controller{retriever: retriever, lock: sync.Mutex{}}
}

func (r *Controller) GetLinkedInProfile(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", req.Header.Get("Origin"))
	w.Header().Set("Access-Control-Allow-Methods", "GET,OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type,Accept")
	w.Header().Set("Cache-Control", fmt.Sprintf("max-age=%.0f", (60*time.Minute).Seconds()))
	profile, err := r.retrieveProfile()
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	} else {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		json.NewEncoder(w).Encode(profile)
	}
}

func (r *Controller) retrieveProfile() (*linkedin.LinkedInProfile, error) {
	r.lock.Lock()
	defer r.lock.Unlock()
	return r.retriever.Get()
}
