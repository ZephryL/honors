package world

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
// On RegionCollection, return a RegionCollection
//------------------------------------------------------------------------------
func RegionCollectionHandler(s *common.System) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {	
		// Set response header for JSON
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		// Open a db, defer close
		db, err := sql.Open("mysql", s.SqlString())
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("HTTP %v - RegionList DB Open Error: %v", http.StatusInternalServerError, err.Error())))
			return
		}
		defer db.Close()
		// Fetch scourge rows
		var vRegionCollection = new(RegionCollection)
		var vQuery strings.Builder
		vQuery.WriteString("select reg.Reg_Key, reg.Reg_Name, reg.Reg_Code ")
		vQuery.WriteString("from region reg ")
		vQuery.WriteString("order by reg.Reg_Name")
		// Execute query into rows
		vRows, err := db.Query(vQuery.String())
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("HTTP %v - RegionList DB Query Error: %v", http.StatusInternalServerError, err.Error())))
			return
		}
		// Iterate rows into Region
		defer vRows.Close()
		for vRows.Next() {
			var vRegion Region
			err = vRows.Scan(
				&vRegion.Reg_Key, &vRegion.Reg_Name, &vRegion.Reg_Code)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(fmt.Sprintf("HTTP %v - RegionList DB RowScan Error: %v", http.StatusInternalServerError, err.Error())))
				return
			}
			vRegionCollection.List = append(vRegionCollection.List, vRegion)
		}
		// get any error encountered during iteration
		err = vRows.Err()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("HTTP %v - RegionList DB Rows Error: %v", http.StatusInternalServerError, err.Error())))
			return
		}
		// Send results as json
		if err := json.NewEncoder(w).Encode(vRegionCollection); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Your call was politely taken, but the JSON Encoder failed to format a reply!"))
		}
	}
}


func RegionCollectionHandler2(w http.ResponseWriter, r *http.Request, s *common.System) {
	// Set response header for JSON
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	// Open a db, defer close
	db, err := sql.Open("mysql", s.SqlString())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("HTTP %v - RegionList DB Open Error: %v", http.StatusInternalServerError, err.Error())))
		return
	}
	defer db.Close()
	// Fetch scourge rows
	var vRegionCollection = new(RegionCollection)
	var vQuery strings.Builder
	vQuery.WriteString("select reg.Reg_Key, reg.Reg_Name, reg.Reg_Code ")
	vQuery.WriteString("from region reg ")
	vQuery.WriteString("order by reg.Reg_Name")
	// Execute query into rows
	vRows, err := db.Query(vQuery.String())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("HTTP %v - RegionList DB Query Error: %v", http.StatusInternalServerError, err.Error())))
		return
	}
	// Iterate rows into Region
	defer vRows.Close()
	for vRows.Next() {
		var vRegion Region
		err = vRows.Scan(
			&vRegion.Reg_Key, &vRegion.Reg_Name, &vRegion.Reg_Code)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("HTTP %v - RegionList DB RowScan Error: %v", http.StatusInternalServerError, err.Error())))
			return
		}
		vRegionCollection.List = append(vRegionCollection.List, vRegion)
	}
	// get any error encountered during iteration
	err = vRows.Err()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("HTTP %v - RegionList DB Rows Error: %v", http.StatusInternalServerError, err.Error())))
		return
	}
	// Send results as json
	if err := json.NewEncoder(w).Encode(vRegionCollection); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Your call was politely taken, but the JSON Encoder failed to format a reply!"))
	}
}
