package controller

import (
	"encoding/json"
	"net/http"

	"github.com/go-rod/rod"
	"github.com/invinity/linkedin-profile-grabber/linkedin"
)

type Controller struct {
	linkedinInst *linkedin.LinkedIn
}

func New(browser *rod.Browser) *Controller {
	return &Controller{linkedinInst: linkedin.New(browser)}
}

func (r *Controller) GetLinkedInProfile(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	profile := r.linkedinInst.RetrieveProfile()

	json.NewEncoder(w).Encode(profile)
}
