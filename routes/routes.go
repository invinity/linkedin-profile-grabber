package routes

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/invinity/linkedin-profile-grabber/controller"
)

func AppRoutes(retriever controller.LinkedinProfileRetriever) *mux.Router {
	var router = mux.NewRouter()
	router = mux.NewRouter().StrictSlash(true)

	//Other Routes
	contInst := controller.NewController(retriever)
	router.HandleFunc("/api/linkedin/profile", contInst.GetLinkedInProfile).Methods(http.MethodGet)
	router.HandleFunc("/api/linkedin/profile", ProvideOptions).Methods(http.MethodOptions)

	return router
}

func ProvideOptions(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", req.Header.Get("Origin"))
	w.Header().Set("Access-Control-Allow-Methods", "GET,OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type,Accept")
}
