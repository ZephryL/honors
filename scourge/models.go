package scourge

import "time"

// Case totals
type ScourgeMetric struct {
	Hi_Date    time.Time `json:"hidate"`
	Lo_Date    time.Time `json:"lodate"`
	Case_Total int       `json:"casetotal"`
	Died_Total int       `json:"diedtotal"`
	Case_Max   int       `json:"casemax"`
	Died_Max   int       `json:"diedmax"`
}

// Scourge Accumulation Totals
type ScourgeTotal struct {
	cases int
	died  int
}

// Scourge, Collection and Detail
type Scourge struct {
	Scg_Key        int       `json:"key"`
	Scg_Name       string    `json:"name"`
	Scg_Origin     string    `json:"origin"`
	Scg_Date       time.Time `json:"date"`
	Prg_Contracted int       `json:"cases"`
	Prg_Died       int       `json:"died"`
}
type ScourgeCollection struct {
	List []Scourge `json:"list"`
}
type ScourgeDetail struct {
	Scg_Key         int       `json:"key"`
	Scg_Name        string    `json:"name"`
	Scg_Origin      string    `json:"origin"`
	Scg_Date        time.Time `json:"date"`
	Scg_Cause       string    `json:"cause"`
	Scg_Description string    `json:"desc"`
	Prg_Date        time.Time `json:"lastdate"`
	Prg_Contracted  int       `json:"cases"`
	Prg_Died        int       `json:"died"`
}

// ScourgeDate
type ScourgeDate struct {
	Prg_Date           time.Time `json:"dt"`
	Prg_Contracted     int       `json:"cs"`
	Prg_Contracted_Acc int       `json:"csa"`
	Prg_Died           int       `json:"dd"`
	Prg_Died_Acc       int       `json:"dda"`
}

// DateProgressCollection
type ScourgeDateCollection struct {
	Scg_Key        int           `json:"scg"`
	Scourge_Metric ScourgeMetric `json:"metrics"`
	List           []ScourgeDate `json:"list"`
}

// CountryProgressCollection
type CountryProgress struct {
	Cnt_Key  int           `json:"cnt"`
	Cnt_Name string        `json:"nm"`
	Cnt_Code string        `json:"id"`
	List     []ScourgeDate `json:"list"`
}
type CountryProgressCollection struct {
	Scg_Key        int               `json:"scg"`
	Scourge_Metric ScourgeMetric     `json:"metrics"`
	List           []CountryProgress `json:"list"`
}

// RegionProgressCollection
type RegionProgress struct {
	Reg_Key  int           `json:"reg"`
	Reg_Name string        `json:"nm"`
	Reg_Code string        `json:"id"`
	List     []ScourgeDate `json:"list"`
}
type RegionProgressCollection struct {
	Scg_Key        int              `json:"scg"`
	Scourge_Metric ScourgeMetric    `json:"metrics"`
	List           []RegionProgress `json:"list"`
}

// ProgressCountryCollection
type ScourgeCountry struct {
	Cnt_Key            int    `json:"cnt"`
	Cnt_Name           string `json:"nm"`
	Cnt_Code           string `json:"id"`
	Prg_Contracted     int    `json:"cs"`
	Prg_Contracted_Acc int    `json:"csa"`
	Prg_Died           int    `json:"dd"`
	Prg_Died_Acc       int    `json:"dda"`
}
type ProgressCountry struct {
	Prg_Date time.Time        `json:"dt"`
	List     []ScourgeCountry `json:"list"`
}
type ProgressCountryCollection struct {
	Scg_Key        int               `json:"scg"`
	Scourge_Metric ScourgeMetric     `json:"metrics"`
	List           []ProgressCountry `json:"list"`
}

// ProgressContinentCollection
type ScourgeContinent struct {
	Con_Key            int    `json:"con"`
	Con_Name           string `json:"nm"`
	Con_Code           string `json:"id"`
	Prg_Contracted     int    `json:"cs"`
	Prg_Contracted_Acc int    `json:"csa"`
	Prg_Died           int    `json:"dd"`
	Prg_Died_Acc       int    `json:"dda"`
}
type ProgressContinent struct {
	Prg_Date time.Time          `json:"dt"`
	List     []ScourgeContinent `json:"list"`
}
type ProgressContinentCollection struct {
	Scg_Key        int                 `json:"scg"`
	Scourge_Metric ScourgeMetric       `json:"metrics"`
	List           []ProgressContinent `json:"list"`
}

// ProgressRegionCollection
type ScourgeRegion struct {
	Reg_Key            int    `json:"key"`
	Reg_Name           string `json:"nm"`
	Reg_Code           string `json:"reg"`
	Prg_Contracted     int    `json:"cs"`
	Prg_Contracted_Acc int    `json:"csa"`
	Prg_Died           int    `json:"dd"`
	Prg_Died_Acc       int    `json:"dda"`
}
type ProgressRegion struct {
	Prg_Date time.Time       `json:"dt"`
	List     []ScourgeRegion `json:"list"`
}
type ProgressRegionCollection struct {
	Scg_Key        int              `json:"scg"`
	Scourge_Metric ScourgeMetric    `json:"metrics"`
	List           []ProgressRegion `json:"list"`
}
