package academy

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/zephryl/honors/common"
)

//------------------------------------------------------------------------------
// On ProjectCreate, insert an Project, return the keyed object
//------------------------------------------------------------------------------
func ProjectCreateHandler(s *common.System) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		// Set response header for JSON
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		
		// Decode a Project from request body
		d := json.NewDecoder(r.Body);
		d.DisallowUnknownFields();
		var vProject = new(Project);
		err := d.Decode(&vProject)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError);
			w.Write([]byte(fmt.Sprintf("HTTP %v -Decode Project Error: %v", http.StatusInternalServerError, err.Error())));
			return;
		}

		// Call the inserter
		if err := CreateProject(s, vProject); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("HTTP %v - Project Create DB Error: %v", http.StatusInternalServerError, err.Error())))
			return
		}

		// Send Project as json
		if err := json.NewEncoder(w).Encode(vProject); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Your call was politely taken, but the JSON Encoder failed to format a reply!"))
		}
	}
}

//------------------------------------------------------------------------------
// On ProjectCollection, return a ProjectCollection
//------------------------------------------------------------------------------
func ProjectListHandler(s *common.System) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// Set response header for JSON
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		
		// Call the standard list reader
		var vProjectList = new(ProjectList)
		if err := ReadProjectList(s, vProjectList); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("HTTP %v -  Read Project List DB Error: %v", http.StatusInternalServerError, err.Error())))
			return
		}

		// Send results as json
		if err := json.NewEncoder(w).Encode(vProjectList); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Your call was politely taken, but the JSON Encoder failed to format a reply!"))
		}
	}
}

//------------------------------------------------------------------------------
// On ProjectRead, return a specific Project
//------------------------------------------------------------------------------
func ProjectReadHandler(s *common.System) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {	
		// Set response header for JSON
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")

		// Get the project key from http Request
		params := mux.Vars(r)
		var vProject = new(Project);
		vPrjKey, err := strconv.Atoi(params["prj-key"])
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("HTTP %v - Project url '%v' does not contain a project identifier: %v", http.StatusInternalServerError, params, err.Error())))
			return
		}
		vProject.Prj_Key = vPrjKey;

		// Call the standard reader
		if err := ReadProject(s, vProject); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("HTTP %v - Project Read DB Error: %v", http.StatusInternalServerError, err.Error())))
			return
		}

		// Send results as json
		if err := json.NewEncoder(w).Encode(vProject); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Your call was politely taken, but the JSON Encoder failed to format a reply!"))
		}
	}
}

//------------------------------------------------------------------------------
// On ProjectUpdate, update a Project, return the keyed object
//------------------------------------------------------------------------------
func ProjectUpdateHandler(s *common.System) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		// Set response header for JSON
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		

		// Get the Project Key from http Request variables
		params := mux.Vars(r)
		vPrjKey, err := strconv.Atoi(params["prj-key"])
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("Unit url '%v' does not contain an estate identifier: %v", params, err.Error())))
			return
		}

		// Get the project fields from request body
		d := json.NewDecoder(r.Body);
		d.DisallowUnknownFields();
		var vProject = new(Project);
		err = d.Decode(&vProject)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError);
			w.Write([]byte(fmt.Sprintf("Decode Project Error: %v", err.Error())));
			return;
		}

		// Call the updater
		vProject.Prj_Key = vPrjKey;
		if err := UpdateProject(s, vProject); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("Project Update DB Error: %v", err.Error())))
			return
		}

		// Send Project as json
		if err := json.NewEncoder(w).Encode(vProject); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Your call was politely taken, but the JSON Encoder failed to format a reply!"))
		}
	}
}

//------------------------------------------------------------------------------
// On ProjectDelete, delete a Project, return the key
//------------------------------------------------------------------------------
func ProjectDeleteHandler(s *common.System) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		// Set response header for JSON
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		
		// Get the Project Key from http Request variables
		params := mux.Vars(r)
		vPrjKey, err := strconv.Atoi(params["prj-key"])
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("HTTP %v - url '%v' does not contain a project identifier: %v", http.StatusInternalServerError, params, err.Error())))
			return
		}

		// Call the deleter
		var vProject = new(Project);
		vProject.Prj_Key = vPrjKey;
		if err := DeleteProject(s, vProject); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("HTTP %v - Project Delete DB Error: %v", http.StatusInternalServerError, err.Error())))
			return
		}

		// Send Project as json
		if err := json.NewEncoder(w).Encode(vProject); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Your call was politely taken, but the JSON Encoder failed to format a reply!"))
		}
	}
}
