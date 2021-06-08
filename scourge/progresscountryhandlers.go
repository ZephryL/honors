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
// On ProgressCountryHandler, return a ProgressCountryCollection
// This is the big one [#1] - every date, every country, daily count
//------------------------------------------------------------------------------
func ProgressCountryHandler(s *common.System) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get url vars from http Request
		params := mux.Vars(r)
		scgKey, err := strconv.Atoi(params["scg-key"])
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("HTTP %v - Scourge url '%v' not numeric: %v", http.StatusInternalServerError, params, err.Error())))
			return
		}
		// We're doing a level break, cherry-picking fields from the row scan depending on which
		// level of break we're at. These are redundant vars to cater for positional scan and save keys
		var prgDate time.Time
		var cntKey int
		var cntName string
		var cntCode string
		var prgContracted int
		var prgDied int
		// Set response header for JSON
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		// Open a db, defer close
		db, err := sql.Open("mysql", s.SqlString())
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("HTTP %v - ProgressCountryHandler DB Open Error: %v", http.StatusInternalServerError, err.Error())))
			return
		}
		// defer close a db
		defer db.Close()
		var vQuery strings.Builder
		vQuery.WriteString("select prg.Scg_Key, prg.Prg_Date, cnt.Cnt_Key, ")
		vQuery.WriteString("	cnt.Cnt_Name, cnt.Cnt_Code, ")
		vQuery.WriteString("	sum(prg.Prg_Contracted) as Prg_Contracted, ")
		vQuery.WriteString("	sum(prg.Prg_Died) as Prg_Died ")
		vQuery.WriteString("from progress prg, country cnt ")
		vQuery.WriteString("where prg.Cnt_Key = cnt.Cnt_Key ")
		vQuery.WriteString("and   prg.Scg_Key = ? ")
		vQuery.WriteString("group by prg.Scg_Key, prg.Prg_Date, cnt.Cnt_Key, ")
		vQuery.WriteString("	cnt.Cnt_Name, cnt.Cnt_Code ")
		vQuery.WriteString("order by prg.Scg_Key, prg.Prg_Date, cnt.Cnt_Key, ")
		vQuery.WriteString("	cnt.Cnt_Name, cnt.Cnt_Code ")
		// Execute
		vRows, err := db.Query(vQuery.String(), scgKey)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("HTTP %v - ProgressCountryCollection DB Query Error: %v", http.StatusInternalServerError, err.Error())))
			return
		}
		// Prepare ScourgeProgress query, create structs, maps
		var vProgressCountryCollection = new(ProgressCountryCollection)
		vProgressCountryCollection.Scourge_Metric.Lo_Date = time.Unix(1<<63-62135596801, 999999999)
		vProgressCountryCollection.Scg_Key = scgKey
		var conMap = make(map[int]ScourgeTotal)
		// Iterate rows into ScourgeProgressCollection ***in a level break by date***
		defer vRows.Close()
		var vHaveRows = vRows.Next()
		for vHaveRows == true {
			err = vRows.Scan(&scgKey, &prgDate, &cntKey, &cntName, &cntCode, &prgContracted, &prgDied)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(fmt.Sprintf("HTTP %v - ProgressCountryCollection DB RowScan Error: %v", http.StatusInternalServerError, err.Error())))
				return
			}
			// Get max/min date
			if prgDate.After(vProgressCountryCollection.Scourge_Metric.Hi_Date) {
				vProgressCountryCollection.Scourge_Metric.Hi_Date = prgDate
			}
			if prgDate.Before(vProgressCountryCollection.Scourge_Metric.Lo_Date) {
				vProgressCountryCollection.Scourge_Metric.Lo_Date = prgDate
			}
			// Create a Date/Country collection for the current date
			var vProgressCountry ProgressCountry
			vProgressCountry.Prg_Date = prgDate
			// now process and break on date (i.e. get all continents for current date)
			for vHaveRows == true && prgDate == vProgressCountry.Prg_Date {
				var vScourgeCountry ScourgeCountry
				vScourgeCountry.Cnt_Key = cntKey
				vScourgeCountry.Cnt_Name = cntName
				vScourgeCountry.Cnt_Code = cntCode
				vScourgeCountry.Prg_Contracted = prgContracted
				vScourgeCountry.Prg_Died = prgDied
				// get map value by continent key
				if vCountryTotals, vOk := conMap[cntKey]; vOk {
					vCountryTotals.cases += prgContracted
					vCountryTotals.died += prgDied
					conMap[cntKey] = vCountryTotals
				} else {
					vCountryTotals.cases = prgContracted
					vCountryTotals.died = prgDied
					conMap[cntKey] = vCountryTotals
				}
				// Accumulate continent
				vScourgeCountry.Prg_Contracted_Acc = conMap[cntKey].cases
				vScourgeCountry.Prg_Died_Acc = conMap[cntKey].died
				// Accumulate totals
				vProgressCountryCollection.Scourge_Metric.Case_Total += prgContracted
				vProgressCountryCollection.Scourge_Metric.Died_Total += prgDied
				// Append
				vProgressCountry.List = append(vProgressCountry.List, vScourgeCountry)
				// Replenishment Read - test for end-of-row first
				vHaveRows = vRows.Next()
				if vHaveRows {
					err = vRows.Scan(&scgKey, &prgDate, &cntKey, &cntName, &cntCode, &prgContracted, &prgDied)
				}
			}
			vProgressCountryCollection.List = append(vProgressCountryCollection.List, vProgressCountry)
		}
		// Test any error encountered during iteration
		err = vRows.Err()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("HTTP %v - ProgressCountryCollection DB Rows Error: %v", http.StatusInternalServerError, err.Error())))
			return
		}
		// Get Country Totals
		for _, vScourgeTotal := range conMap {
			if vScourgeTotal.cases > vProgressCountryCollection.Scourge_Metric.Case_Max {
				vProgressCountryCollection.Scourge_Metric.Case_Max = vScourgeTotal.cases
			}
			if vScourgeTotal.died > vProgressCountryCollection.Scourge_Metric.Died_Max {
				vProgressCountryCollection.Scourge_Metric.Died_Max = vScourgeTotal.died
			}
		}
		// Send results as json
		if err := json.NewEncoder(w).Encode(vProgressCountryCollection); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Your call was politely taken, but the JSON Encoder failed to format a reply!"))
		}
	}
}