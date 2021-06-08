package scourge

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
	"sort"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/zephryl/zephry/common"
)

//------------------------------------------------------------------------------
// On ScourgeRegionDateHandler, return a RegionProgressCollection
// This is a big one - all countries, every date, daily count
//------------------------------------------------------------------------------
func RegionStreamHandler(s *common.System) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get url vars from http Request
		params := mux.Vars(r)
		scgKey, err := strconv.Atoi(params["scg-key"])
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("HTTP %v - RegionStreamHandler key '%v' not numeric: %v", http.StatusInternalServerError, params, err.Error())))
			return
		}
	
		// Open a db, defer close
		db, err := sql.Open("mysql", s.SqlString())
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("HTTP %v - RegionStreamHandler DB Open Error: %v", http.StatusInternalServerError, err.Error())))
			return
		}
		// defer close a db
		defer db.Close()
	
		// Get an initialized grid of series > points (all scourge dates, each date with all regions)
		vChartMap, err := getGrid(db, s, scgKey);
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("HTTP %v - RegionStreamHandler Map Grid error: %v", http.StatusInternalServerError, err.Error())))
			return
		}
	
		// Scan db values into vars, find map entries using date and reg
		var prgDate time.Time
		var regKey int
		var prgDied int
		// Set response header for JSON
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		var vQuery strings.Builder
		vQuery.WriteString("select prg.Prg_Date, reg.Reg_Key, ")
		vQuery.WriteString("	   sum(prg.Prg_Died) as Prg_Died ")
		vQuery.WriteString("from progress prg, country cnt, region reg ")
		vQuery.WriteString("where prg.Cnt_Key = cnt.Cnt_Key ")
		vQuery.WriteString("and   cnt.Reg_Key = reg.Reg_Key ")
		vQuery.WriteString("and   prg.Scg_Key = ? ")
		vQuery.WriteString("group by prg.Prg_Date, reg.Reg_Key ")
		vQuery.WriteString("order by prg.Prg_Date desc, reg.Reg_Key ")
		// Execute
		vRows, err := db.Query(vQuery.String(), scgKey)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("HTTP %v - RegionStreamHandler DB Query Error: %v", http.StatusInternalServerError, err.Error())))
			return
		}
		// Iterate: each row is a date/region coordinate in the map
		defer vRows.Close()
		for vRows.Next() {
			err = vRows.Scan(&prgDate, &regKey, &prgDied)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(fmt.Sprintf("HTTP %v - RegionStreamHandler DB RowScan Error: %v", http.StatusInternalServerError, err.Error())))
				return
			}
			if vDateEntry, ok := vChartMap.DateMap[prgDate]; ok {
				if vPointEntry, ok := vDateEntry[regKey]; ok {
					vPointEntry.Value = prgDied;
					vDateEntry[regKey] = vPointEntry;
				}
			}
		}
		// Test any error encountered during iteration
		err = vRows.Err()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("HTTP %v - RegionStreamHandler DB Rows Error: %v", http.StatusInternalServerError, err.Error())))
			return
		}
		// Convert map to formal struct
		var vHiValue int;
		var vChartStream = new(ChartStream)
		vChartStream.Stream_Metric.Lo_Date = time.Unix(1<<63-62135596801, 999999999)
		vChartStream.Title = "Daily Deaths by Region"
		for date, chartseries := range vChartMap.DateMap {
			var vChartSeries ChartSeries;
			vChartSeries.Date = date;
			vHiValue = 0;
			for key, chartpoint := range chartseries {
				var vChartPoint ChartPoint
				vChartPoint.Key = key;
				vChartPoint.Code = chartpoint.Code;
				vChartPoint.Name = chartpoint.Name;
				vChartPoint.Value = chartpoint.Value;
				vHiValue += vChartPoint.Value;
				vChartSeries.List = append(vChartSeries.List, vChartPoint)
			}
			// StreamMetrics
			if date.After(vChartStream.Stream_Metric.Hi_Date) {
				vChartStream.Stream_Metric.Hi_Date = date
			}
			if date.Before(vChartStream.Stream_Metric.Lo_Date) {
				vChartStream.Stream_Metric.Lo_Date = date
			}
			if vHiValue > vChartStream.Stream_Metric.Max {
				vChartStream.Stream_Metric.Max = vHiValue;
				vChartStream.Stream_Metric.Max_Date = date;
			}
			vChartStream.Stream_Metric.Total += vHiValue;
			// Sort Regions
			sort.Slice(vChartSeries.List, func(i, j int) bool {return vChartSeries.List[i].Key < vChartSeries.List[j].Key})
			vChartStream.List = append(vChartStream.List, vChartSeries)
		}
		// Sort Dates
		sort.Slice(vChartStream.List, func(i, j int) bool {return vChartStream.List[i].Date.Before(vChartStream.List[j].Date);})
	
		// Send results as json
		if err := json.NewEncoder(w).Encode(vChartStream); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Your call was politely taken, but the JSON Encoder failed to format a reply!"))
		}
	}
}

//------------------------------------------------------------------------------
// Return a fully initialized grid of all possible values
//------------------------------------------------------------------------------
func getGrid(db *sql.DB, s *common.System, scgKey int) (ChartMap, error) {
	//return type
	var vChartMap ChartMap;

	// Get all regions
	var vQuery strings.Builder
	var Reglist []ChartPoint
	vQuery.WriteString("select Reg_Key, Reg_Code, Reg_Name from region")
	vRows, err := db.Query(vQuery.String())
	if err != nil { return vChartMap, err; }
	defer vRows.Close()
	for vRows.Next() {
		var vChartPoint ChartPoint
		err = vRows.Scan(
			&vChartPoint.Key, &vChartPoint.Code, &vChartPoint.Name)
			if err != nil { return vChartMap, err; }
			Reglist = append(Reglist, vChartPoint)
	}
	err = vRows.Err()
	if err != nil { return vChartMap, err; }

	// Get all dates, create a map of maps
	var loDate time.Time;
	var hiDate time.Time;
	vQuery.Reset()
	vQuery.WriteString("select min(Prg_Date) as LoDate, max(Prg_Date) as HiDate ")
	vQuery.WriteString("from progress ")
	vQuery.WriteString("where Scg_Key = ? ")
	err = db.QueryRow(vQuery.String(), scgKey).Scan(&loDate, &hiDate)
	if err != nil { return vChartMap, err; }
	vChartMap.DateMap = make(map[time.Time]map[int]ChartPoint)
	for date := loDate; date.After(hiDate) == false; date = date.AddDate(0, 0, 1) {
		vPointMap, ok := vChartMap.DateMap[date]
		if !ok {
			vPointMap = make(map[int]ChartPoint)
			vChartMap.DateMap[date] = vPointMap
		}
		for _, reg := range Reglist {
			vPointMap[reg.Key] = reg;
		}
	}
	return vChartMap, nil;
}
