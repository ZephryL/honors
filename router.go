package main

import (
	"github.com/gorilla/mux"
	"github.com/zephryl/honors/common"
)

func MuxRouter(s *common.System) *mux.Router {
	// get a gorilla router
	router := mux.NewRouter().StrictSlash(true)
	// common routes
	router.HandleFunc("/", common.Cors(s, IndexHandler(s))).Methods("GET")
	router.HandleFunc("/login", common.Cors(s, common.LoginHandler(s))).Methods("OPTIONS", "POST")
	router.HandleFunc("/logout", common.Cors(s, common.AuthOnly(s, common.LogoutHandler(s)))).Methods("OPTIONS", "POST")
	router.HandleFunc("/register", common.Cors(s, common.RegisterCreateHandler(s))).Methods("OPTIONS", "POST")
	router.HandleFunc("/register/verify", common.Cors(s, common.RegisterVerifyHandler(s))).Methods("OPTIONS", "POST")
	router.HandleFunc("/registers/{reg-key}", common.Cors(s, common.RegisterReadHandler(s))).Methods("OPTIONS", "GET")
	router.HandleFunc("/forgot", common.Cors(s, common.ForgotCreateHandler(s))).Methods("OPTIONS", "POST")
	router.HandleFunc("/forgot/verify", common.Cors(s, common.ForgotVerifyHandler(s))).Methods("OPTIONS", "POST")
	// return a list of all routes - keep this last?
	router.HandleFunc("/routes", common.Cors(s, common.RouteWalker(s, router))).Methods("OPTIONS", "GET")
	// return a gorilla router
	return router
}
