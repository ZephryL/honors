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
// On ProgressContinentHandler, return a ProgressContinentCollection
// This is the big one [#3] - every date, every continent, daily count
//------------------------------------------------------------------------------
func ProgressContinentHandler(s *common.System) func(http.ResponseWriter, *http.Request) {
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
		var conKey int
		var conName string
		var conCode string
		var prgContracted int
		var prgDied int
		// Set response header for JSON
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		// Open a db, defer close
		db, err := sql.Open("mysql", s.SqlString())
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("HTTP %v - ProgressContinentHandler DB Open Error: %v", http.StatusInternalServerError, err.Error())))
			return
		}
		// defer close a db
		defer db.Close()
		var vQuery strings.Builder
		vQuery.WriteString("select prg.Scg_Key, prg.Prg_Date, con.Con_Key, ")
		vQuery.WriteString("	con.Con_Name, con.Con_Iso2, ")
		vQuery.WriteString("	sum(prg.Prg_Contracted) as Prg_Contracted, ")
		vQuery.WriteString("	sum(prg.Prg_Died) as Prg_Died ")
		vQuery.WriteString("from progress prg, country cnt, continent con ")
		vQuery.WriteString("where prg.Cnt_Key = cnt.Cnt_Key ")
		vQuery.WriteString("and   cnt.Con_Key = con.Con_Key ")
		vQuery.WriteString("and   prg.Scg_Key = ? ")
		vQuery.WriteString("group by prg.Scg_Key, prg.Prg_Date, con.Con_Key, ")
		vQuery.WriteString("	con.Con_Name, con.Con_Iso2 ")
		vQuery.WriteString("order by prg.Scg_Key, prg.Prg_Date, con.Con_Key, ")
		vQuery.WriteString("	con.Con_Name, con.Con_Iso2 ")
		// Execute
		vRows, err := db.Query(vQuery.String(), scgKey)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("HTTP %v - ProgressContinentCollection DB Query Error: %v", http.StatusInternalServerError, err.Error())))
			return
		}
		// Prepare ScourgeProgress query, create structs, maps
		var vProgressContinentCollection = new(ProgressContinentCollection)
		vProgressContinentCollection.Scourge_Metric.Lo_Date = time.Unix(1<<63-62135596801, 999999999)
		vProgressContinentCollection.Scg_Key = scgKey
		var conMap = make(map[int]ScourgeTotal)
		// Iterate rows into ScourgeProgressCollection ***in a level break by date***
		defer vRows.Close()
		var vHaveRows = vRows.Next()
		for vHaveRows == true {
			err = vRows.Scan(&scgKey, &prgDate, &conKey, &conName, &conCode, &prgContracted, &prgDied)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(fmt.Sprintf("HTTP %v - ProgressContinentCollection DB RowScan Error: %v", http.StatusInternalServerError, err.Error())))
				return
			}
			// Get max/min date
			if prgDate.After(vProgressContinentCollection.Scourge_Metric.Hi_Date) {
				vProgressContinentCollection.Scourge_Metric.Hi_Date = prgDate
			}
			if prgDate.Before(vProgressContinentCollection.Scourge_Metric.Lo_Date) {
				vProgressContinentCollection.Scourge_Metric.Lo_Date = prgDate
			}
			// Create a Date/Continent collection for the current date
			var vProgressContinent ProgressContinent
			vProgressContinent.Prg_Date = prgDate
			// now process and break on date (i.e. get all continents for current date)
			for vHaveRows == true && prgDate == vProgressContinent.Prg_Date {
				var vScourgeContinent ScourgeContinent
				vScourgeContinent.Con_Key = conKey
				vScourgeContinent.Con_Name = conName
				vScourgeContinent.Con_Code = conCode
				vScourgeContinent.Prg_Contracted = prgContracted
				vScourgeContinent.Prg_Died = prgDied
				// get map value by continent key
				if vContinentTotals, vOk := conMap[conKey]; vOk {
					vContinentTotals.cases += prgContracted
					vContinentTotals.died += prgDied
					conMap[conKey] = vContinentTotals
				} else {
					vContinentTotals.cases = prgContracted
					vContinentTotals.died = prgDied
					conMap[conKey] = vContinentTotals
				}
				// Accumulate continent
				vScourgeContinent.Prg_Contracted_Acc = conMap[conKey].cases
				vScourgeContinent.Prg_Died_Acc = conMap[conKey].died
				// Accumulate totals
				vProgressContinentCollection.Scourge_Metric.Case_Total += prgContracted
				vProgressContinentCollection.Scourge_Metric.Died_Total += prgDied
				// Append
				vProgressContinent.List = append(vProgressContinent.List, vScourgeContinent)
				// Replenishment Read - test for end-of-row first
				vHaveRows = vRows.Next()
				if vHaveRows {
					err = vRows.Scan(&scgKey, &prgDate, &conKey, &conName, &conCode, &prgContracted, &prgDied)
				}
			}
			vProgressContinentCollection.List = append(vProgressContinentCollection.List, vProgressContinent)
		}
		// Test any error encountered during iteration
		err = vRows.Err()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("HTTP %v - ProgressCountryCollection DB Rows Error: %v", http.StatusInternalServerError, err.Error())))
			return
		}
		// Get Continent Totals
		for _, vScourgeTotal := range conMap {
			if vScourgeTotal.cases > vProgressContinentCollection.Scourge_Metric.Case_Max {
				vProgressContinentCollection.Scourge_Metric.Case_Max = vScourgeTotal.cases
			}
			if vScourgeTotal.died > vProgressContinentCollection.Scourge_Metric.Died_Max {
				vProgressContinentCollection.Scourge_Metric.Died_Max = vScourgeTotal.died
			}
		}
		// Send results as json
		if err := json.NewEncoder(w).Encode(vProgressContinentCollection); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Your call was politely taken, but the JSON Encoder failed to format a reply!"))
		}
	}
}