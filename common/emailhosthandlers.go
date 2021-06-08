package common

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

//------------------------------------------------------------------------------
// On EmailHostCreate, insert an EmailHost, return the keyed object
//------------------------------------------------------------------------------
func EmailHostCreateHandler(s *System) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		// Set response header for JSON
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		
		// Decode an EmailHost from request body (Basically get the key)
		d := json.NewDecoder(r.Body);
		d.DisallowUnknownFields();
		var vEmailHost = new(EmailHost);
		err := d.Decode(&vEmailHost)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError);
			w.Write([]byte(fmt.Sprintf("HTTP %v -Decode EmailHost Error: %v", http.StatusInternalServerError, err.Error())));
			return;
		}

		// Call the inserter
		if err := CreateEmailHost(s, vEmailHost); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("HTTP %v - EmailHost Create DB Error: %v", http.StatusInternalServerError, err.Error())))
			return
		}

		// Send EmailHost as json
		if err := json.NewEncoder(w).Encode(vEmailHost); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Your call was politely taken, but the JSON Encoder failed to format a reply!"))
		}
	}
}

//------------------------------------------------------------------------------
// On EmailHostCollection, return a EmailHostCollection
//------------------------------------------------------------------------------
func EmailHostCollectionHandler(s *System) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {	
		// Set response header for JSON
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		var vEmailHostCollection = new(EmailHostCollection);

		// Call the list reader
		if err := ReadEmailHostList(s, vEmailHostCollection); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("HTTP %v - EmailHost Read List DB Error: %v", http.StatusInternalServerError, err.Error())))
			return
		}

		// Send results as json
		if err := json.NewEncoder(w).Encode(vEmailHostCollection); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Your call was politely taken, but the JSON Encoder failed to format a reply!"))
		}
	}
}

//------------------------------------------------------------------------------
// On EmailHostRead, return a specific EmailHost
//------------------------------------------------------------------------------

func EmailHostReadHandler(s *System) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {	
		// Set response header for JSON
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")

		// Get the EmailHost key from http Request
		params := mux.Vars(r)
		var vEmailHost = new(EmailHost);
		vEhsKey, err := strconv.Atoi(params["ehs-key"])
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("HTTP %v - EmailHost url '%v' does not contain an identifier: %v", http.StatusInternalServerError, params, err.Error())))
			return
		}
		vEmailHost.Ehs_Key = vEhsKey;

		// Call the standard reader
		if err := ReadEmailHost(s, vEmailHost); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("HTTP %v - EmailHost Read DB Error: %v", http.StatusInternalServerError, err.Error())))
			return
		}

		// Send results as json
		if err := json.NewEncoder(w).Encode(vEmailHost); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Your call was politely taken, but the JSON Encoder failed to format a reply!"))
		}
	}
}

//------------------------------------------------------------------------------
// On EmailHostDefault, return the default EmailHost
//------------------------------------------------------------------------------
func EmailHostDefaultHandler(s *System) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {	
		// Set response header for JSON
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")

		var vEmailHost = new(EmailHost);

		// Call the default reader
		if err := ReadEmailHostDefault(s, vEmailHost); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("HTTP %v - EmailHost Read Default DB Error: %v", http.StatusInternalServerError, err.Error())))
			return
		}

		// Send results as json
		if err := json.NewEncoder(w).Encode(vEmailHost); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Your call was politely taken, but the JSON Encoder failed to format a reply!"))
		}
	}
}
