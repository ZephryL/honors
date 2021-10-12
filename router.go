package main

import (
	"github.com/gorilla/mux"
	"github.com/zephryl/honors/common"
	"github.com/zephryl/honors/academy"
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
	// institution routes
	router.HandleFunc("/institutions", common.Cors(s, common.Auth(s, academy.InstitutionListHandler(s)))).Methods("OPTIONS", "GET");
	router.HandleFunc("/institutions", common.Cors(s, common.Auth(s, academy.InstitutionCreateHandler(s)))).Methods("OPTIONS", "POST")
	router.HandleFunc("/institutions/{ins-key}", common.Cors(s, common.Auth(s, academy.InstitutionReadHandler(s)))).Methods("OPTIONS", "GET")
	router.HandleFunc("/institutions/{ins-key}", common.Cors(s, common.Auth(s, academy.InstitutionUpdateHandler(s)))).Methods("OPTIONS", "PUT")
	router.HandleFunc("/institutions/{ins-key}", common.Cors(s, common.Auth(s, academy.InstitutionDeleteHandler(s)))).Methods("OPTIONS", "DELETE")
	// project routes
	router.HandleFunc("/projects", common.Cors(s, common.Auth(s, academy.ProjectListHandler(s)))).Methods("OPTIONS", "GET");
	router.HandleFunc("/projects", common.Cors(s, common.Auth(s, academy.ProjectCreateHandler(s)))).Methods("OPTIONS", "POST")
	router.HandleFunc("/projects/{prj-key}", common.Cors(s, common.Auth(s, academy.ProjectReadHandler(s)))).Methods("OPTIONS", "GET")
	router.HandleFunc("/projects/{prj-key}", common.Cors(s, common.Auth(s, academy.ProjectUpdateHandler(s)))).Methods("OPTIONS", "PUT")
	router.HandleFunc("/projects/{prj-key}", common.Cors(s, common.Auth(s, academy.ProjectDeleteHandler(s)))).Methods("OPTIONS", "DELETE")
	// reference routes
	router.HandleFunc("/references", common.Cors(s, common.Auth(s, academy.ReferenceListHandler(s)))).Methods("OPTIONS", "GET");
	router.HandleFunc("/references", common.Cors(s, common.Auth(s, academy.ReferenceCreateHandler(s)))).Methods("OPTIONS", "POST")
	router.HandleFunc("/references/{ref-key}", common.Cors(s, common.Auth(s, academy.ReferenceReadHandler(s)))).Methods("OPTIONS", "GET")
	router.HandleFunc("/references/{ref-key}", common.Cors(s, common.Auth(s, academy.ReferenceUpdateHandler(s)))).Methods("OPTIONS", "PUT")
	router.HandleFunc("/references/{ref-key}", common.Cors(s, common.Auth(s, academy.ReferenceDeleteHandler(s)))).Methods("OPTIONS", "DELETE")
	// projref routes
	router.HandleFunc("/projects/{prj-key}/references", common.Cors(s, common.Auth(s, academy.ProjRefListHandler(s)))).Methods("OPTIONS", "GET");
	router.HandleFunc("/projects/{prj-key}/references/{ref-key}", common.Cors(s, common.Auth(s, academy.ProjRefCreateHandler(s)))).Methods("OPTIONS", "POST");
	// return a list of all routes - keep this last?
	router.HandleFunc("/routes", common.Cors(s, common.RouteWalker(s, router))).Methods("OPTIONS", "GET")
	// return a gorilla router
	return router
}
