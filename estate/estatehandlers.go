package estate

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	_ "github.com/go-sql-driver/mysql"
	"github.com/zephryl/zephry/common"
)

//------------------------------------------------------------------------------
// On EstateCollection, return a EstateCollection
//------------------------------------------------------------------------------

func EstateCollectionHandler(s *common.System) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {	
		// Set response header for JSON
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		// Open a db, defer close
		db, err := sql.Open("mysql", s.SqlString())
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("HTTP %v - EstateList DB Open Error: %v", http.StatusInternalServerError, err.Error())))
			return
		}
		defer db.Close()
		// Fetch scourge rows
		var vEstateCollection = new(EstateCollection)
		var vQuery strings.Builder
		vQuery.WriteString("select reg.Reg_Key, reg.Reg_Name, reg.Reg_Code ")
		vQuery.WriteString("from region reg ")
		vQuery.WriteString("order by reg.Reg_Name")
		// Execute query into rows
		vRows, err := db.Query(vQuery.String())
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("HTTP %v - EstateList DB Query Error: %v", http.StatusInternalServerError, err.Error())))
			return
		}
		// Iterate rows into Estate
		defer vRows.Close()
		for vRows.Next() {
			var vEstate Estate
			err = vRows.Scan(
				&vEstate.Reg_Key, &vEstate.Reg_Name, &vEstate.Reg_Code)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(fmt.Sprintf("HTTP %v - EstateList DB RowScan Error: %v", http.StatusInternalServerError, err.Error())))
				return
			}
			vEstateCollection.List = append(vEstateCollection.List, vEstate)
		}
		// get any error encountered during iteration
		err = vRows.Err()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("HTTP %v - EstateList DB Rows Error: %v", http.StatusInternalServerError, err.Error())))
			return
		}
		// Send results as json
		if err := json.NewEncoder(w).Encode(vEstateCollection); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Your call was politely taken, but the JSON Encoder failed to format a reply!"))
		}
	}
}