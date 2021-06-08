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
// On ScourgeCountryDateHandler, return a CountryProgressCollection
// This is a not-so-big one - a specific country, every date, daily count
//------------------------------------------------------------------------------
func CountryProgressHandler(s *common.System) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get url vars from http Request
		params := mux.Vars(r)
		scgKey, err := strconv.Atoi(params["scg-key"])
		cntKey, err := strconv.Atoi(params["cnt-key"])
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("HTTP %v - Scourge url '%v' not numeric: %v", http.StatusInternalServerError, params, err.Error())))
			return
		}
		// We're doing a level break, cherry-picking fields from the row scan depending on which
		// level of break we're at. These are redundant vars to cater for positional scan and save keys
		var cntName string
		var cntCode string
		var prgDate time.Time
		var prgContracted int
		var prgDied int
		// Set response header for JSON
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		// Open a db, defer close
		db, err := sql.Open("mysql", s.SqlString())
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("HTTP %v - ScourgeCountryDateHandler DB Open Error: %v", http.StatusInternalServerError, err.Error())))
			return
		}
		// defer close a db
		defer db.Close()
		var vQuery strings.Builder
		vQuery.WriteString("select cnt.Cnt_Key, cnt.Cnt_Name, cnt.Cnt_Isoa2, prg.Prg_Date, ")
		vQuery.WriteString("	   sum(prg.Prg_Contracted) as Prg_Contracted, ")
		vQuery.WriteString("	   sum(prg.Prg_Died) as Prg_Died ")
		vQuery.WriteString("from progress prg, country cnt ")
		vQuery.WriteString("where prg.Cnt_Key = cnt.Cnt_Key ")
		vQuery.WriteString("and   prg.Scg_Key = ? ")
		vQuery.WriteString("and   prg.Cnt_Key = ? ")
		vQuery.WriteString("group by cnt.Cnt_Key, cnt.Cnt_Name, cnt.Cnt_Isoa2, prg.Prg_Date ")
		vQuery.WriteString("order by cnt.Cnt_Key, cnt.Cnt_Name, cnt.Cnt_Isoa2, prg.Prg_Date ")
		// Execute
		vRows, err := db.Query(vQuery.String(), scgKey, cntKey)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("HTTP %v - CountryProgressCollection DB Query Error: %v", http.StatusInternalServerError, err.Error())))
			return
		}
		// Prepare ScourgeProgress query, create structs, maps
		var vCountryProgressCollection = new(CountryProgressCollection)
		vCountryProgressCollection.Scourge_Metric.Lo_Date = time.Unix(1<<63-62135596801, 999999999)
		vCountryProgressCollection.Scg_Key = scgKey
		var vScourgeTotal = new(ScourgeTotal)
		// Iterate rows into vCountryProgressCollection ***in a level break by country***
		defer vRows.Close()
		var vHaveRows = vRows.Next()
		for vHaveRows == true {
			err = vRows.Scan(&cntKey, &cntName, &cntCode, &prgDate, &prgContracted, &prgDied)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(fmt.Sprintf("HTTP %v - CountryProgressCollection DB RowScan Error: %v", http.StatusInternalServerError, err.Error())))
				return
			}
			// Create a Date/Country collection for the current date
			var vScourgeCountryDate CountryProgress
			vScourgeCountryDate.Cnt_Key = cntKey
			vScourgeCountryDate.Cnt_Name = cntName
			vScourgeCountryDate.Cnt_Code = cntCode
			// now process and break on country key (i.e. get all dates for the country - and in fact there is only one)
			for vHaveRows == true && cntKey == vScourgeCountryDate.Cnt_Key {
				var vScourgeDate ScourgeDate
				vScourgeDate.Prg_Date = prgDate
				vScourgeDate.Prg_Contracted = prgContracted
				vScourgeDate.Prg_Died = prgDied
				// Accumulate
				vScourgeTotal.cases += prgContracted
				vScourgeTotal.died += prgDied
				vScourgeDate.Prg_Contracted_Acc += vScourgeTotal.cases
				vScourgeDate.Prg_Died_Acc += vScourgeTotal.died
				// Min / Max
				if vScourgeDate.Prg_Contracted > vCountryProgressCollection.Scourge_Metric.Case_Max {
					vCountryProgressCollection.Scourge_Metric.Case_Max = vScourgeDate.Prg_Contracted
				}
				if vScourgeDate.Prg_Died > vCountryProgressCollection.Scourge_Metric.Died_Max {
					vCountryProgressCollection.Scourge_Metric.Died_Max = vScourgeDate.Prg_Died
				}
				// Get max/min date
				if prgDate.After(vCountryProgressCollection.Scourge_Metric.Hi_Date) {
					vCountryProgressCollection.Scourge_Metric.Hi_Date = prgDate
				}
				if prgDate.Before(vCountryProgressCollection.Scourge_Metric.Lo_Date) {
					vCountryProgressCollection.Scourge_Metric.Lo_Date = prgDate
				}
				// Append
				vScourgeCountryDate.List = append(vScourgeCountryDate.List, vScourgeDate)
				// Replenishment Read - test for end-of-row first
				vHaveRows = vRows.Next()
				if vHaveRows {
					err = vRows.Scan(&cntKey, &cntName, &cntCode, &prgDate, &prgContracted, &prgDied)
				}
			}
			// Accumulate totals
			vCountryProgressCollection.Scourge_Metric.Case_Total = vScourgeTotal.cases
			vCountryProgressCollection.Scourge_Metric.Died_Total = vScourgeTotal.died
			vCountryProgressCollection.List = append(vCountryProgressCollection.List, vScourgeCountryDate)
		}
		// Test any error encountered during iteration
		err = vRows.Err()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("HTTP %v - CountryProgressCollection DB Rows Error: %v", http.StatusInternalServerError, err.Error())))
			return
		}
		// Send results as json
		if err := json.NewEncoder(w).Encode(vCountryProgressCollection); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Your call was politely taken, but the JSON Encoder failed to format a reply!"))
		}
	}
}