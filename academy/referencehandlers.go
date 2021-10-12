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
// On ReferenceCreate, insert an Reference, return the keyed object
//------------------------------------------------------------------------------
func ReferenceCreateHandler(s *common.System) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		// Set response header for JSON
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		
		// Decode a Reference from request body
		d := json.NewDecoder(r.Body);
		d.DisallowUnknownFields();
		var vReference = new(Reference);
		err := d.Decode(&vReference)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError);
			w.Write([]byte(fmt.Sprintf("Decode Reference Error: %v", fmt.Errorf("%v", err))));
			return;
		}

		// Call the inserter
		if err := CreateReference(s, vReference); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("Reference Create DB Error: %v", err.Error())))
			return
		}

		// Send Reference as json
		if err := json.NewEncoder(w).Encode(vReference); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Your call was politely taken, but the JSON Encoder failed to format a reply!"))
		}
	}
}

//------------------------------------------------------------------------------
// On ReferenceCollection, return a ReferenceCollection
//------------------------------------------------------------------------------
func ReferenceListHandler(s *common.System) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// Set response header for JSON
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		
		// Call the standard list reader
		var vReferenceList = new(ReferenceList)
		if err := ReadReferenceList(s, vReferenceList); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("Read Reference List DB Error: %v", err.Error())))
			return
		}

		// Send results as json
		if err := json.NewEncoder(w).Encode(vReferenceList); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Your call was politely taken, but the JSON Encoder failed to format a reply!"))
		}
	}
}

//------------------------------------------------------------------------------
// On ReferenceRead, return a specific Reference
//------------------------------------------------------------------------------
func ReferenceReadHandler(s *common.System) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {	
		// Set response header for JSON
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")

		// Get the reference key from http Request
		params := mux.Vars(r)
		var vReference = new(Reference);
		vRefKey, err := strconv.Atoi(params["ref-key"])
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("Reference url '%v' does not contain a reference identifier: %v", r.URL, err.Error())))
			return
		}
		vReference.Ref_Key = vRefKey;

		// Call the standard reader
		if err := ReadReference(s, vReference); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("Reference Read DB Error: %v", err.Error())))
			return
		}

		// Send results as json
		if err := json.NewEncoder(w).Encode(vReference); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Your call was politely taken, but the JSON Encoder failed to format a reply!"))
		}
	}
}

//------------------------------------------------------------------------------
// On ReferenceUpdate, update a Reference, return the keyed object
//------------------------------------------------------------------------------
func ReferenceUpdateHandler(s *common.System) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		// Set response header for JSON
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		
		// Get the Reference Key from http Request variables
		params := mux.Vars(r)
		vRefKey, err := strconv.Atoi(params["ref-key"])
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("Unit url '%v' does not contain an estate identifier: %v", r.URL, err.Error())))
			return
		}

		// Get the reference fields from request body
		d := json.NewDecoder(r.Body);
		d.DisallowUnknownFields();
		var vReference = new(Reference);
		err = d.Decode(&vReference)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError);
			w.Write([]byte(fmt.Sprintf("Decode Reference Error: %v", err.Error())));
			return;
		}

		// Call the updater
		vReference.Ref_Key = vRefKey;
		if err := UpdateReference(s, vReference); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("Reference Create DB Error: %v", err.Error())))
			return
		}

		// Send Reference as json
		if err := json.NewEncoder(w).Encode(vReference); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Your call was politely taken, but the JSON Encoder failed to format a reply!"))
		}
	}
}

//------------------------------------------------------------------------------
// On ReferenceDelete, delete a Reference, return the key
//------------------------------------------------------------------------------
func ReferenceDeleteHandler(s *common.System) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		// Set response header for JSON
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		
		// Get the Reference Key from http Request variables
		params := mux.Vars(r)
		vRefKey, err := strconv.Atoi(params["ref-key"])
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("Url '%v' does not contain a reference identifier: %v", r.URL, err.Error())))
			return
		}

		// Call the deleter
		var vReference = new(Reference);
		vReference.Ref_Key = vRefKey;
		if err := DeleteReference(s, vReference); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("Reference Delete DB Error: %v", err.Error())))
			return
		}

		// Send Reference as json
		if err := json.NewEncoder(w).Encode(vReference); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Your call was politely taken, but the JSON Encoder failed to format a reply!"))
		}
	}
}

//------------------------------------------------------------------------------
// Return a list of matching References
//------------------------------------------------------------------------------
func ReferenceFindHandler(s *common.System) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// Set response header for JSON
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")

		// Get the estate key from http Request
		vFindpart := r.URL.Query().Get("findpart");
		var vReferenceList = new(ReferenceList)
		if vFindpart == "" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(fmt.Sprintf("Reference url '%v' does not contain a search argument", r.URL.Query())))
			return
		}
		vReferenceList.Find_Part = vFindpart;

		// Call the list reader
		if err := ReadReferenceList(s, vReferenceList); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("HTTP %v - Reference Read List DB Error: %v", http.StatusInternalServerError, err.Error())))
			return
		}

		// Send results as json
		if err := json.NewEncoder(w).Encode(vReferenceList); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Your call was politely taken, but the JSON Encoder failed to format a reply!"))
		}
	}
}
