package routes

import (
	"net/http"

	"github.com/go-rod/rod"
	"github.com/gorilla/mux"
	"github.com/invinity/linkedin-profile-grabber/controller"
)

func AppRoutes(browser *rod.Browser) *mux.Router {
	var router = mux.NewRouter()
	router = mux.NewRouter().StrictSlash(true)

	//Other Routes
	contInst := controller.New(browser)
	router.HandleFunc("/api/linkedin/profile", contInst.GetLinkedInProfile).Methods(http.MethodGet)

	return router
}
