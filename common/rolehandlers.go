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
// On RoleCollection, return a RoleCollection
//------------------------------------------------------------------------------

func RoleCollectionHandler(s *System) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {	
		// Set response header for JSON
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		// Get a user from context (put there by Auth middleware)
		vUsrKey, ok := r.Context().Value("usrkey").(int);
		if !ok {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("HTTP %v - Role Collection Context Error", http.StatusInternalServerError)))
			return
		}
		// Open a db, defer close
		db, err := sql.Open("mysql", s.SqlString())
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("HTTP %v - RoleList DB Open Error: %v", http.StatusInternalServerError, err.Error())))
			return
		}
		defer db.Close()
		// Fetch scourge rows
		var vRoleCollection = new(RoleCollection)
		var vQuery strings.Builder
		vQuery.WriteString("select App_Key, App_Code, App_Name, App_Desc, App_Lock ")
		vQuery.WriteString("from application ")
		vQuery.WriteString("order by App_Lock, App_Name ")
		// Execute query into rows
		vRows, err := db.Query(vQuery.String(), vUsrKey)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("HTTP %v - RoleList DB Query Error: %v", http.StatusInternalServerError, err.Error())))
			return
		}
		// Iterate rows into Role
		defer vRows.Close()
		for vRows.Next() {
			var vRole Role
			err = vRows.Scan(&vRoleCollection.Account.Email, &vRoleCollection.Account.Fullname,
				&vRole.Key, &vRole.Code, &vRole.Name, &vRole.Desc, &vRole.Lock)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(fmt.Sprintf("HTTP %v - RoleList DB RowScan Error: %v", http.StatusInternalServerError, err.Error())))
				return
			}
			vRoleCollection.List = append(vRoleCollection.List, vRole)
		}
		// get any error encountered during iteration
		err = vRows.Err()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("HTTP %v - RoleList DB Rows Error: %v", http.StatusInternalServerError, err.Error())))
			return
		}
		// Send results as json
		if err := json.NewEncoder(w).Encode(vRoleCollection); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Your call was politely taken, but the JSON Encoder failed to format a reply!"))
		}
	}
}