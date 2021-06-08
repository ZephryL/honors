package scourge

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/zephryl/zephry/common"
)

//------------------------------------------------------------------------------
// On ScourgeDateCollection, return a ScourgeDateCollection
//------------------------------------------------------------------------------
func ProgressHandler(s *common.System) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get url vars from http Request
		params := mux.Vars(r)
		scgKey, err := strconv.Atoi(params["scg-key"])
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("HTTP %v - Scourge url '%v' not numeric: %v", http.StatusInternalServerError, params, err.Error())))
			return
		}
		// Set response header for JSON
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		// Open a db, defer close
		db, err := sql.Open("mysql", s.SqlString())
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("HTTP %v - ScourgeIndex DB Open Error: %v", http.StatusInternalServerError, err.Error())))
			return
		}
		defer db.Close()
		// Prepare ScourgeDate query
		var vQuery strings.Builder
		vQuery.WriteString("select Prg_Date, ")
		vQuery.WriteString("	sum(prg.Prg_Contracted) as Prg_Contracted, ")
		vQuery.WriteString("	sum(prg.Prg_Died) as Prg_Died ")
		vQuery.WriteString("from progress prg ")
		vQuery.WriteString("where prg.Scg_Key = ? ")
		vQuery.WriteString("group by prg.Scg_Key, Prg_Date ")
		vQuery.WriteString("order by prg.Scg_Key, Prg_Date;")
		// Execute
		vRows, err := db.Query(vQuery.String(), scgKey)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("HTTP %v - ScourgeDateCollection DB Query Error: %v", http.StatusInternalServerError, err.Error())))
			return
		}
		// Iterate rows into ScourgeDateCollection
		defer vRows.Close()
		var vScourgeDateCollection = new(ScourgeDateCollection)
		vScourgeDateCollection.Scourge_Metric.Lo_Date = time.Unix(1<<63-62135596801, 999999999)
		vScourgeDateCollection.Scg_Key = scgKey
		var vScourgeTotal = new(ScourgeTotal)
		for vRows.Next() {
			var vScourgeDate ScourgeDate
			err = vRows.Scan(
				&vScourgeDate.Prg_Date, &vScourgeDate.Prg_Contracted, &vScourgeDate.Prg_Died)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(fmt.Sprintf("HTTP %v - ScourgeDateCollection DB RowScan Error: %v", http.StatusInternalServerError, err.Error())))
				return
			}
			// Get max/min date
			if vScourgeDate.Prg_Date.After(vScourgeDateCollection.Scourge_Metric.Hi_Date) {
				vScourgeDateCollection.Scourge_Metric.Hi_Date = vScourgeDate.Prg_Date
			}
			if vScourgeDate.Prg_Date.Before(vScourgeDateCollection.Scourge_Metric.Lo_Date) {
				vScourgeDateCollection.Scourge_Metric.Lo_Date = vScourgeDate.Prg_Date
			}
			// Accumulate daily
			// Accumulate
			vScourgeTotal.cases += vScourgeDate.Prg_Contracted
			vScourgeTotal.died += vScourgeDate.Prg_Died
			vScourgeDate.Prg_Contracted_Acc = vScourgeTotal.cases
			vScourgeDate.Prg_Died_Acc = vScourgeTotal.died
			// Accumulate totals
			vScourgeDateCollection.Scourge_Metric.Case_Total += vScourgeDate.Prg_Contracted
			vScourgeDateCollection.Scourge_Metric.Died_Total += vScourgeDate.Prg_Died
			// Set max
			if vScourgeDate.Prg_Contracted > vScourgeDateCollection.Scourge_Metric.Case_Max {
				vScourgeDateCollection.Scourge_Metric.Case_Max = vScourgeDate.Prg_Contracted
			}
			if vScourgeDate.Prg_Died > vScourgeDateCollection.Scourge_Metric.Died_Max {
				vScourgeDateCollection.Scourge_Metric.Died_Max = vScourgeDate.Prg_Died
			}
			vScourgeDateCollection.List = append(vScourgeDateCollection.List, vScourgeDate)
		}
		// Test any error encountered during iteration
		err = vRows.Err()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("HTTP %v - ScourgeDateCollection DB Rows Error: %v", http.StatusInternalServerError, err.Error())))
			return
		}
		// Send results as json
		if err := json.NewEncoder(w).Encode(vScourgeDateCollection); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Your call was politely taken, but the JSON Encoder failed to format a reply!"))
		}
	}
}