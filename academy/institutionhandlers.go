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
// On InstitutionCreate, insert an Institution, return the keyed object
//------------------------------------------------------------------------------
func InstitutionCreateHandler(s *common.System) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		// Set response header for JSON
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		
		// Decode a Institution from request body
		d := json.NewDecoder(r.Body);
		d.DisallowUnknownFields();
		var vInstitution = new(Institution);
		err := d.Decode(&vInstitution)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError);
			w.Write([]byte(fmt.Sprintf("Decode Institution Error: %v", fmt.Errorf("%v", err))));
			return;
		}

		// Call the inserter
		if err := CreateInstitution(s, vInstitution); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("Institution Create DB Error: %v", err.Error())))
			return
		}

		// Send Institution as json
		if err := json.NewEncoder(w).Encode(vInstitution); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Your call was politely taken, but the JSON Encoder failed to format a reply!"))
		}
	}
}

//------------------------------------------------------------------------------
// On InstitutionCollection, return a InstitutionCollection
//------------------------------------------------------------------------------
func InstitutionListHandler(s *common.System) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// Set response header for JSON
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		
		// Call the standard list reader
		var vInstitutionList = new(InstitutionList)
		if err := ReadInstitutionList(s, vInstitutionList); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("Read Institution List DB Error: %v", err.Error())))
			return
		}

		// Send results as json
		if err := json.NewEncoder(w).Encode(vInstitutionList); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Your call was politely taken, but the JSON Encoder failed to format a reply!"))
		}
	}
}

//------------------------------------------------------------------------------
// On InstitutionRead, return a specific Institution
//------------------------------------------------------------------------------
func InstitutionReadHandler(s *common.System) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {	
		// Set response header for JSON
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")

		// Get the institution key from http Request
		params := mux.Vars(r)
		var vInstitution = new(Institution);
		vInsKey, err := strconv.Atoi(params["ins-key"])
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("Institution url '%v' does not contain a institution identifier: %v", r.URL, err.Error())))
			return
		}
		vInstitution.Ins_Key = vInsKey;

		// Call the standard reader
		if err := ReadInstitution(s, vInstitution); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("Institution Read DB Error: %v", err.Error())))
			return
		}

		// Send results as json
		if err := json.NewEncoder(w).Encode(vInstitution); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Your call was politely taken, but the JSON Encoder failed to format a reply!"))
		}
	}
}

//------------------------------------------------------------------------------
// On InstitutionUpdate, update a Institution, return the keyed object
//------------------------------------------------------------------------------
func InstitutionUpdateHandler(s *common.System) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		// Set response header for JSON
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		
		// Get the Institution Key from http Request variables
		params := mux.Vars(r)
		vInsKey, err := strconv.Atoi(params["ins-key"])
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("Unit url '%v' does not contain an estate identifier: %v", r.URL, err.Error())))
			return
		}

		// Get the institution fields from request body
		d := json.NewDecoder(r.Body);
		d.DisallowUnknownFields();
		var vInstitution = new(Institution);
		err = d.Decode(&vInstitution)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError);
			w.Write([]byte(fmt.Sprintf("Decode Institution Error: %v", err.Error())));
			return;
		}

		// Call the updater
		vInstitution.Ins_Key = vInsKey;
		if err := UpdateInstitution(s, vInstitution); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("Institution Create DB Error: %v", err.Error())))
			return
		}

		// Send Institution as json
		if err := json.NewEncoder(w).Encode(vInstitution); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Your call was politely taken, but the JSON Encoder failed to format a reply!"))
		}
	}
}

//------------------------------------------------------------------------------
// On InstitutionDelete, delete a Institution, return the key
//------------------------------------------------------------------------------
func InstitutionDeleteHandler(s *common.System) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		// Set response header for JSON
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		
		// Get the Institution Key from http Request variables
		params := mux.Vars(r)
		vInsKey, err := strconv.Atoi(params["ins-key"])
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("Url '%v' does not contain a institution identifier: %v", r.URL, err.Error())))
			return
		}

		// Call the deleter
		var vInstitution = new(Institution);
		vInstitution.Ins_Key = vInsKey;
		if err := DeleteInstitution(s, vInstitution); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("Institution Delete DB Error: %v", err.Error())))
			return
		}

		// Send Institution as json
		if err := json.NewEncoder(w).Encode(vInstitution); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Your call was politely taken, but the JSON Encoder failed to format a reply!"))
		}
	}
}
