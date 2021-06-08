package scourge

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/zephryl/zephry/common"
)

//------------------------------------------------------------------------------
// On ScourgeCollection, return a ScourgeCollection
//------------------------------------------------------------------------------
func ScourgeCollectionHandler(s *common.System) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {	
		// Set response header for JSON
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		// Open a db, defer close
		db, err := sql.Open("mysql", s.SqlString())
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("HTTP %v - ScourgeList DB Open Error: %v", http.StatusInternalServerError, err.Error())))
			return
		}
		defer db.Close()
		// Fetch scourge rows
		var vScourgeCollection = new(ScourgeCollection)
		var vQuery strings.Builder
		vQuery.WriteString("select scg.Scg_Key, scg.Scg_Name, scg.Scg_Origin, scg.Scg_Date, ")
		vQuery.WriteString("       ifnull(sum(prg.Prg_Contracted), 0) as Prg_Contracted, ")
		vQuery.WriteString("       ifnull(sum(prg.Prg_Died), 0) as Prg_Died ")
		vQuery.WriteString("from scourge scg left outer join progress prg ")
		vQuery.WriteString("			     on scg.Scg_Key = prg.Scg_key ")
		vQuery.WriteString("group by scg.Scg_Key, scg.Scg_Name, scg.Scg_Origin, scg.Scg_Date ")
		vQuery.WriteString("order by scg.Scg_Date desc")
		// Execute query into rows
		vRows, err := db.Query(vQuery.String())
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("HTTP %v - ScourgeList DB Query Error: %v", http.StatusInternalServerError, err.Error())))
			return
		}
		// Iterate rows into Scourge
		defer vRows.Close()
		for vRows.Next() {
			var vScourge Scourge
			err = vRows.Scan(
				&vScourge.Scg_Key, &vScourge.Scg_Name, &vScourge.Scg_Origin, &vScourge.Scg_Date, &vScourge.Prg_Contracted, &vScourge.Prg_Died)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(fmt.Sprintf("HTTP %v - ScourgeList DB RowScan Error: %v", http.StatusInternalServerError, err.Error())))
				return
			}
			vScourgeCollection.List = append(vScourgeCollection.List, vScourge)
		}
		// get any error encountered during iteration
		err = vRows.Err()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("HTTP %v - ScourgeList DB Rows Error: %v", http.StatusInternalServerError, err.Error())))
			return
		}
		// Send results as json
		if err := json.NewEncoder(w).Encode(vScourgeCollection); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Your call was politely taken, but the JSON Encoder failed to format a reply!"))
		}
	}
}

//------------------------------------------------------------------------------
// On Scourge, return a Scourge detail
//------------------------------------------------------------------------------
func ScourgeHandler(s *common.System) func(http.ResponseWriter, *http.Request) {
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
		// Fetch a sum row
		var vScourgeDetail = new(ScourgeDetail)
		var vQuery strings.Builder
		vQuery.WriteString("select scg.Scg_Key, scg.Scg_Name, scg.Scg_Origin, ")
		vQuery.WriteString("       scg.Scg_Date, scg.Scg_Cause, scg.Scg_Description, ")
		vQuery.WriteString("       max(prg.Prg_Date) as Prg_Date, ")
		vQuery.WriteString("       sum(prg.Prg_Contracted) as Prg_Contracted, ")
		vQuery.WriteString("       sum(prg.Prg_Died) as Prg_Died ")
		vQuery.WriteString("from scourge scg, progress prg ")
		vQuery.WriteString("where scg.Scg_Key = ? ")
		vQuery.WriteString("group by scg.Scg_Key, scg.Scg_Name")
		err = db.QueryRow(vQuery.String(), scgKey).Scan(
			&vScourgeDetail.Scg_Key, &vScourgeDetail.Scg_Name, &vScourgeDetail.Scg_Origin, &vScourgeDetail.Scg_Date, &vScourgeDetail.Scg_Cause,
			&vScourgeDetail.Scg_Description, &vScourgeDetail.Prg_Date, &vScourgeDetail.Prg_Contracted, &vScourgeDetail.Prg_Died)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("HTTP %v - ScourgeIndex DB QueryRow Error: %v", http.StatusInternalServerError, err.Error())))
			return
		}
		// Send results as json
		if err := json.NewEncoder(w).Encode(vScourgeDetail); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Your call was politely taken, but the JSON Encoder failed to format a reply!"))
		}
	}
}

//------------------------------------------------------------------------------
// On ScourgeNumbersHandler, return a NumberCards
//------------------------------------------------------------------------------
func ScourgeNumbersHandler(s *common.System) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get url vars from http Request
		params := mux.Vars(r)
		scgKey, err := strconv.Atoi(params["scg-key"])
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("HTTP %v - ScourgeNumbersHandler url '%v' not numeric: %v", http.StatusInternalServerError, params, err.Error())))
			return
		}
		// Set response header for JSON
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		// Open a db, defer close
		db, err := sql.Open("mysql", s.SqlString())
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("HTTP %v - ScourgeNumbersHandler DB Open Error: %v", http.StatusInternalServerError, err.Error())))
			return
		}
		var vCases, vDied, vDays, vDiedDay, vPopulation, vDiedPop float32;
		defer db.Close()
		// Fetch a sum row
		var vNumberCards = new(common.NumberCards)
		var vQuery strings.Builder
		vQuery.WriteString("select sum(Prg_Contracted) as cases, ")
		vQuery.WriteString("       sum(PRG_Died) as died, ")
		vQuery.WriteString("       datediff(max(Prg_Date), min(Prg_Date)) as days, ")
		vQuery.WriteString("       (select sum(Cnt_Population) from country) as population ")
		vQuery.WriteString("from progress ")
		vQuery.WriteString("where Scg_Key = ? ")
		err = db.QueryRow(vQuery.String(), scgKey).Scan(
			&vCases, &vDied, &vDays, &vPopulation)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("HTTP %v - ScourgeNumbersHandler DB QueryRow Error: %v", http.StatusInternalServerError, err.Error())))
			return
		}
	
		vNumberCards.List = append(vNumberCards.List, common.KeyValue{Key: "Cases", Value: vCases, Precision: 0});
		vNumberCards.List = append(vNumberCards.List, common.KeyValue{Key: "Deaths", Value: vDied, Precision: 0});
		vNumberCards.List = append(vNumberCards.List, common.KeyValue{Key: "Days", Value: vDays, Precision: 0});
		if vDays < 1 {
			vDiedDay = 0;
		} else {
			vDiedDay = float32(vDied) / float32(vDays);
		}
		vNumberCards.List = append(vNumberCards.List, common.KeyValue{Key: "Deaths/Day", Value: vDiedDay, Precision: 2});
		vNumberCards.List = append(vNumberCards.List, common.KeyValue{Key: "Population (Million)", Value: vPopulation / 1000000, Precision: 0});
		if vPopulation < 1 {
			vDiedPop = 0;
		} else {
			vDiedPop = float32(vDied) / float32(vPopulation / 1000000);
		}
		vNumberCards.List = append(vNumberCards.List, common.KeyValue{Key: "Deaths/Million", Value: vDiedPop, Precision: 2});
		// Send results as json
		if err := json.NewEncoder(w).Encode(vNumberCards); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Your call was politely taken, but the JSON Encoder failed to format a reply!"))
		}
	}
}