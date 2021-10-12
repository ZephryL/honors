package academy

import (
	"encoding/json"
	"fmt"
	"strconv"
	"net/http"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/zephryl/honors/common"
)

//------------------------------------------------------------------------------
// On ProjRefCreate, insert a ProjRef, return the keyed object
//------------------------------------------------------------------------------
func ProjRefCreateHandler(s *common.System) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		// Set response header for JSON
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		
		// Decode a ProjRef from request body
		d := json.NewDecoder(r.Body);
		d.DisallowUnknownFields();
		var vProjRef = new(ProjRef);
		err := d.Decode(&vProjRef)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError);
			w.Write([]byte(fmt.Sprintf("Decode ProjRef Error: %v", fmt.Errorf("%v", err))));
			return;
		}

		// Call the inserter
		if err := CreateProjRef(s, vProjRef); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("ProjRef Create DB Error: %v", err.Error())))
			return
		}

		// Send ProjRef as json
		if err := json.NewEncoder(w).Encode(vProjRef); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Your call was politely taken, but the JSON Encoder failed to format a reply!"))
		}
	}
}

//------------------------------------------------------------------------------
// On ProjRefCollection, return a ProjRefCollection
//------------------------------------------------------------------------------
func ProjRefListHandler(s *common.System) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// Set response header for JSON
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		
		// Get the project key from http Request
		params := mux.Vars(r)
		var vProjRefList = new(ProjRefList);
		vPrjKey, err := strconv.Atoi(params["prj-key"]);
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("HTTP %v - ProjRef url '%v' does not contain a project identifier: %v", http.StatusInternalServerError, params, err.Error())))
			return
		}
		vProjRefList.Prj_Key = vPrjKey;

		// Call the standard list reader
		if err := ReadProjRefList(s, vProjRefList); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("Read ProjRef List DB Error: %v", err.Error())))
			return
		}

		// Send results as json
		if err := json.NewEncoder(w).Encode(vProjRefList); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Your call was politely taken, but the JSON Encoder failed to format a reply!"))
		}
	}
}
