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
// On ProgressRegionHandler, return a ProgressRegionCollection
// This is the not-so-big one [#2] - every date, every region, daily count
//------------------------------------------------------------------------------
func ProgressRegionHandler(s *common.System) func(http.ResponseWriter, *http.Request) {
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
		var regKey int
		var regName string
		var regCode string
		var prgContracted int
		var prgDied int
		// Set response header for JSON
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		// Open a db, defer close
		db, err := sql.Open("mysql", s.SqlString())
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("HTTP %v - ProgressRegionHandler DB Open Error: %v", http.StatusInternalServerError, err.Error())))
			return
		}
		// defer close a db
		defer db.Close()
		var vQuery strings.Builder
		vQuery.WriteString("select prg.Scg_Key, prg.Prg_Date, reg.Reg_Key, ")
		vQuery.WriteString("	reg.Reg_Name, reg.Reg_Code, ")
		vQuery.WriteString("	sum(prg.Prg_Contracted) as Prg_Contracted, ")
		vQuery.WriteString("	sum(prg.Prg_Died) as Prg_Died ")
		vQuery.WriteString("from progress prg, country cnt, region reg ")
		vQuery.WriteString("where prg.Cnt_Key = cnt.Cnt_Key ")
		vQuery.WriteString("and   cnt.Reg_Key = reg.Reg_Key ")
		vQuery.WriteString("and   prg.Scg_Key = ? ")
		vQuery.WriteString("group by prg.Scg_Key, prg.Prg_Date, reg.Reg_Key, ")
		vQuery.WriteString("	reg.Reg_Name, reg.Reg_Code ")
		vQuery.WriteString("order by prg.Scg_Key, prg.Prg_Date, reg.Reg_Key, ")
		vQuery.WriteString("	reg.Reg_Name, reg.Reg_Code ")
		// Execute
		vRows, err := db.Query(vQuery.String(), scgKey)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("HTTP %v - ProgressRegionCollection DB Query Error: %v", http.StatusInternalServerError, err.Error())))
			return
		}
		// Prepare ScourgeProgress query, create structs, maps
		var vProgressRegionCollection = new(ProgressRegionCollection)
		vProgressRegionCollection.Scourge_Metric.Lo_Date = time.Unix(1<<63-62135596801, 999999999)
		vProgressRegionCollection.Scg_Key = scgKey
		var conMap = make(map[int]ScourgeTotal)
		// Iterate rows into ScourgeProgressCollection ***in a level break by date***
		defer vRows.Close()
		var vHaveRows = vRows.Next()
		for vHaveRows == true {
			err = vRows.Scan(&scgKey, &prgDate, &regKey, &regName, &regCode, &prgContracted, &prgDied)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(fmt.Sprintf("HTTP %v - ProgressRegionCollection DB RowScan Error: %v", http.StatusInternalServerError, err.Error())))
				return
			}
			// Get max/min date
			if prgDate.After(vProgressRegionCollection.Scourge_Metric.Hi_Date) {
				vProgressRegionCollection.Scourge_Metric.Hi_Date = prgDate
			}
			if prgDate.Before(vProgressRegionCollection.Scourge_Metric.Lo_Date) {
				vProgressRegionCollection.Scourge_Metric.Lo_Date = prgDate
			}
			// Create a Date/Region collection for the current date
			var vProgressRegion ProgressRegion
			vProgressRegion.Prg_Date = prgDate
			// now process and break on date (i.e. get all continents for current date)
			for vHaveRows == true && prgDate == vProgressRegion.Prg_Date {
				var vScourgeRegion ScourgeRegion
				vScourgeRegion.Reg_Key = regKey
				vScourgeRegion.Reg_Name = regName
				vScourgeRegion.Reg_Code = regCode
				vScourgeRegion.Prg_Contracted = prgContracted
				vScourgeRegion.Prg_Died = prgDied
				// get map value by continent key
				if vRegionTotals, vOk := conMap[regKey]; vOk {
					vRegionTotals.cases += prgContracted
					vRegionTotals.died += prgDied
					conMap[regKey] = vRegionTotals
				} else {
					vRegionTotals.cases = prgContracted
					vRegionTotals.died = prgDied
					conMap[regKey] = vRegionTotals
				}
				// Accumulate continent
				vScourgeRegion.Prg_Contracted_Acc = conMap[regKey].cases
				vScourgeRegion.Prg_Died_Acc = conMap[regKey].died
				// Accumulate totals
				vProgressRegionCollection.Scourge_Metric.Case_Total += prgContracted
				vProgressRegionCollection.Scourge_Metric.Died_Total += prgDied
				// Append
				vProgressRegion.List = append(vProgressRegion.List, vScourgeRegion)
				// Replenishment Read - test for end-of-row first
				vHaveRows = vRows.Next()
				if vHaveRows {
					err = vRows.Scan(&scgKey, &prgDate, &regKey, &regName, &regCode, &prgContracted, &prgDied)
				}
			}
			vProgressRegionCollection.List = append(vProgressRegionCollection.List, vProgressRegion)
		}
		// Test any error encountered during iteration
		err = vRows.Err()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("HTTP %v - ProgressCountryCollection DB Rows Error: %v", http.StatusInternalServerError, err.Error())))
			return
		}
		// Get Region Totals
		for _, vScourgeTotal := range conMap {
			if vScourgeTotal.cases > vProgressRegionCollection.Scourge_Metric.Case_Max {
				vProgressRegionCollection.Scourge_Metric.Case_Max = vScourgeTotal.cases
			}
			if vScourgeTotal.died > vProgressRegionCollection.Scourge_Metric.Died_Max {
				vProgressRegionCollection.Scourge_Metric.Died_Max = vScourgeTotal.died
			}
		}
		// Send results as json
		if err := json.NewEncoder(w).Encode(vProgressRegionCollection); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Your call was politely taken, but the JSON Encoder failed to format a reply!"))
		}
	}
}