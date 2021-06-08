package main

import (
	"github.com/gorilla/mux"
	"github.com/zephryl/zephry/common"
	"github.com/zephryl/zephry/world"
	"github.com/zephryl/zephry/scourge"
	"github.com/zephryl/zephry/estate"
)

func MuxRouter(s *common.System) *mux.Router {
	// get a gorilla router
	router := mux.NewRouter().StrictSlash(true)
	// Home routes
	router.HandleFunc("/", common.SetHeaders(s, IndexHandler(s))).Methods("GET")
	router.HandleFunc("/registers", common.SetHeaders(s, common.RegisterCreateHandler(s))).Methods("OPTIONS", "PUT")
	router.HandleFunc("/verify", common.SetHeaders(s, common.RegisterVerifyHandler(s))).Methods("OPTIONS", "POST")
	router.HandleFunc("/registers/{reg-key}", common.SetHeaders(s, common.RegisterReadHandler(s))).Methods("OPTIONS", "GET")
	router.HandleFunc("/login", common.SetHeaders(s, common.LoginHandler(s))).Methods("OPTIONS", "POST")
	router.HandleFunc("/application", common.SetHeaders(s, common.ApplicationCollectionHandler(s))).Methods("OPTIONS", "GET")
	router.HandleFunc("/role", common.SetHeaders(s, common.Auth(s, common.RoleCollectionHandler(s)))).Methods("OPTIONS", "GET")
	router.HandleFunc("/emailhosts", common.SetHeaders(s, common.Auth(s, common.EmailHostCollectionHandler(s)))).Methods("OPTIONS", "GET")
	router.HandleFunc("/emailhosts", common.SetHeaders(s, common.Auth(s, common.EmailHostCreateHandler(s)))).Methods("OPTIONS", "PUT")
	router.HandleFunc("/emailhosts/{ehs-key}", common.SetHeaders(s, common.Auth(s, common.EmailHostReadHandler(s)))).Methods("OPTIONS", "GET")
	// Estate Routes
	router.HandleFunc("/estate", common.SetHeaders(s, common.Auth(s, estate.EstateCollectionHandler(s)))).Methods("OPTIONS", "GET")
	// World routes
	router.HandleFunc("/world/region", common.SetHeaders(s, world.RegionCollectionHandler(s))).Methods("OPTIONS", "GET")
	// Scourge routes
	router.HandleFunc("/scourge/", common.SetHeaders(s, scourge.ScourgeCollectionHandler(s))).Methods("OPTIONS", "GET")
	router.HandleFunc("/scourge/{scg-key}/", common.SetHeaders(s, scourge.ScourgeHandler(s))).Methods("OPTIONS", "GET")
	router.HandleFunc("/scourge/{scg-key}/numbers", common.SetHeaders(s, scourge.ScourgeNumbersHandler(s))).Methods("OPTIONS", "GET")
	router.HandleFunc("/scourge/{scg-key}/progress", common.SetHeaders(s, scourge.ProgressHandler(s))).Methods("OPTIONS", "GET")
	router.HandleFunc("/scourge/{scg-key}/progress/continent", common.SetHeaders(s, scourge.ProgressContinentHandler(s))).Methods("OPTIONS", "GET")
	router.HandleFunc("/scourge/{scg-key}/progress/region", common.SetHeaders(s, scourge.ProgressRegionHandler(s))).Methods("OPTIONS", "GET")
	router.HandleFunc("/scourge/{scg-key}/progress/country", common.SetHeaders(s, scourge.ProgressCountryHandler(s))).Methods("OPTIONS", "GET")
	router.HandleFunc("/scourge/{scg-key}/country/deathmill", common.SetHeaders(s, scourge.CountryDeathsMillionHandler(s))).Methods("OPTIONS", "GET")
	router.HandleFunc("/scourge/{scg-key}/country/progress", common.SetHeaders(s, scourge.CountriesProgressHandler(s))).Methods("OPTIONS", "GET")
	router.HandleFunc("/scourge/{scg-key}/country/{cnt-key}/progress", common.SetHeaders(s, scourge.CountryProgressHandler(s))).Methods("OPTIONS", "GET")
	router.HandleFunc("/scourge/{scg-key}/region/stream", common.SetHeaders(s, scourge.RegionStreamHandler(s))).Methods("OPTIONS", "GET")
	// return a list of all routes - keep this last?
	router.HandleFunc("/routes", common.SetHeaders(s, common.RouteWalker(s, router))).Methods("OPTIONS", "GET")
	// return a gorilla router
	return router
}
