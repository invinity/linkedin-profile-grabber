package controller

import (
	"encoding/json"
	"net/http"

	"github.com/invinity/linkedin-profile-grabber/linkedin"
)

func GetLinkedInProfile(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	linkedin := linkedin.New()
	profile := linkedin.RetrieveProfile()

	json.NewEncoder(w).Encode(profile)
}
