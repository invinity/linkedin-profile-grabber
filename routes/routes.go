package routes

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/invinity/linkedin-profile-grabber/controller"
)

func AppRoutes() *mux.Router {
	var router = mux.NewRouter()
	router = mux.NewRouter().StrictSlash(true)

	//Other Routes
	router.HandleFunc("/api/linkedin/profile", controller.GetLinkedInProfile).Methods(http.MethodGet)

	return router
}
