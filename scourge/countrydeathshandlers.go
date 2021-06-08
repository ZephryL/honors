package scourge

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"strconv"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/zephryl/zephry/common"
)

//------------------------------------------------------------------------------
// On CountryDeathsMillionHandler, return a CountryDeathsMillion
//------------------------------------------------------------------------------
func CountryDeathsMillionHandler(s *common.System) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get url vars from http Request
		params := mux.Vars(r)
		scgKey, err := strconv.Atoi(params["scg-key"])
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("HTTP %v - CountryDeathsMillionHandler url '%v' not numeric: %v", http.StatusInternalServerError, params, err.Error())))
			return
		}
		// Set response header for JSON
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		// Open a db, defer close
		db, err := sql.Open("mysql", s.SqlString())
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("HTTP %v - CountryDeathsMillionHandler DB Open Error: %v", http.StatusInternalServerError, err.Error())))
			return
		}
		// defer close a db
		defer db.Close()
		var vQuery strings.Builder
		vQuery.WriteString("select cnt.Cnt_Isoa2, ")
		vQuery.WriteString("	   case when Cnt_Population = 0 then 0 else sum(prg.Prg_Died) / (Cnt_Population / 1000000) end as Cnt_DeathsMillPop ")
		vQuery.WriteString("from country cnt, progress prg ")
		vQuery.WriteString("where cnt.Cnt_key = prg.Cnt_Key ")
		vQuery.WriteString("and   prg.Scg_Key = ? ")
		vQuery.WriteString("group by cnt.Cnt_Isoa2")
		// Execute
		vRows, err := db.Query(vQuery.String(), scgKey)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("HTTP %v - CountryDeathsMillionHandler DB Query Error: %v", http.StatusInternalServerError, err.Error())))
			return
		}
		// Scan storage
		var vKeyValueList common.KeyValueList
		vKeyValueList.Loval = math.MaxFloat32;
		var vIsoa2 string;
		var vDeathMillPop float32;
		defer vRows.Close()
		for vRows.Next() {
			err = vRows.Scan(&vIsoa2, &vDeathMillPop)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(fmt.Sprintf("HTTP %v - CountryDeathsMillionHandler DB RowScan Error: %v", http.StatusInternalServerError, err.Error())))
				return
			}
			// Create a Date/Country collection for the current date
			vKeyValueList.Count++;
			if vDeathMillPop < vKeyValueList.Loval { vKeyValueList.Loval = vDeathMillPop }
			if vDeathMillPop > vKeyValueList.Hival { vKeyValueList.Hival = vDeathMillPop }
			vKeyValueList.List = append(vKeyValueList.List, common.KeyValue{Key: vIsoa2, Value: vDeathMillPop, Precision: 2});
		}
		// Test any error encountered during iteration
		err = vRows.Err()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("HTTP %v - CountryDeathsMillionHandler DB Rows Error: %v", http.StatusInternalServerError, err.Error())))
			return
		}
		// Send results as json
		if err := json.NewEncoder(w).Encode(vKeyValueList); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Your call was politely taken, but the JSON Encoder failed to format a reply!"))
		}
	}
}