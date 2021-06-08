package common

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	_ "github.com/go-sql-driver/mysql"
)

//------------------------------------------------------------------------------
// On ApplicationCollection, return a ApplicationCollection (No auth required)
//------------------------------------------------------------------------------

func ApplicationCollectionHandler(s *System) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {	
		// Set response header for JSON
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		// Open a db, defer close
		db, err := sql.Open("mysql", s.SqlString())
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("HTTP %v - ApplicationList DB Open Error: %v", http.StatusInternalServerError, err.Error())))
			return
		}
		defer db.Close()
		// Fetch scourge rows
		var vApplicationCollection = new(ApplicationCollection)
		var vQuery strings.Builder
		vQuery.WriteString("select App_Key, App_Code, App_Name, App_Desc, App_Lock ")
		vQuery.WriteString("from application ")
		vQuery.WriteString("order by App_Lock, App_Name ")
		// Execute query into rows
		vRows, err := db.Query(vQuery.String())
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("HTTP %v - ApplicationList DB Query Error: %v", http.StatusInternalServerError, err.Error())))
			return
		}
		// Iterate rows into Application
		defer vRows.Close()
		for vRows.Next() {
			var vApplication Application
			err = vRows.Scan(&vApplication.App_Key, &vApplication.App_Code, &vApplication.App_Name, &vApplication.App_Desc, &vApplication.App_Lock)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(fmt.Sprintf("HTTP %v - ApplicationList DB RowScan Error: %v", http.StatusInternalServerError, err.Error())))
				return
			}
			vApplicationCollection.List = append(vApplicationCollection.List, vApplication)
		}
		// get any error encountered during iteration
		err = vRows.Err()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("HTTP %v - ApplicationList DB Rows Error: %v", http.StatusInternalServerError, err.Error())))
			return
		}
		// Send results as json
		if err := json.NewEncoder(w).Encode(vApplicationCollection); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Your call was politely taken, but the JSON Encoder failed to format a reply!"))
		}
	}
}