package scourge

import "time"

// Stream totals
type StreamMetric struct {
	Hi_Date  time.Time `json:"hidate"`
	Lo_Date  time.Time `json:"lodate"`
	Max_Date time.Time `json:"maxdate"`
	Max      int       `json:"max"`
	Total    int       `json:"total"`
}

type ChartPoint struct {
	Key   int    `json:"key"`
	Code  string `json:"code"`
	Name  string `json:"name"`
	Value int    `json:"value"`
}

type ChartSeries struct {
	Date time.Time    `json:"date"`
	List []ChartPoint `json:"list"`
}

type ChartStream struct {
	Title         string        `json:"title"`
	Stream_Metric StreamMetric  `json:"metrics"`
	List          []ChartSeries `json:"list"`
}

type ChartMap struct {
	DateMap map[time.Time]map[int]ChartPoint
}
