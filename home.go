package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/zephryl/zephry/common"
)

// On index, return a StandardResponse funky message
func IndexHandler(s *common.System) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")

		// Create a successful greeting, send success response
		var vWelcome common.Welcome;
		vWelcome.Status = 0;
		vWelcome.Greeting = "Zephry v0.1.0 is up and running. The following REST endpoints are available";
		vWelcome.Date = time.Now();
		vWelcome.Endpoints = append(vWelcome.Endpoints, common.Endpoint{ Url: "/", Description: "Status and Greeting", ContentType: "JSON"});
		vWelcome.Endpoints = append(vWelcome.Endpoints, common.Endpoint{ Url: "/world/region", Description: "All World Bank Regions", ContentType: "JSON"});
		vWelcome.Endpoints = append(vWelcome.Endpoints, common.Endpoint{ Url: "/scourge/", Description: "All Scourges", ContentType: "JSON"});
		vWelcome.Endpoints = append(vWelcome.Endpoints, common.Endpoint{ Url: "/scourge/{scg-key}/", Description: "Specific Scourge", ContentType: "JSON"});
		vWelcome.Endpoints = append(vWelcome.Endpoints, common.Endpoint{ Url: "/scourge/{scg-key}/numbers", Description: "Key Metrics, including cases and deaths", ContentType: "JSON"});
		vWelcome.Endpoints = append(vWelcome.Endpoints, common.Endpoint{ Url: "/scourge/{scg-key}/progress", Description: "Daily progression totals", ContentType: "JSON"});
		vWelcome.Endpoints = append(vWelcome.Endpoints, common.Endpoint{ Url: "/scourge/{scg-key}/progress/continent", Description: "Progression by date by continent", ContentType: "JSON"});
		vWelcome.Endpoints = append(vWelcome.Endpoints, common.Endpoint{ Url: "/scourge/{scg-key}/progress/region", Description: "Progression by date by region", ContentType: "JSON"});
		vWelcome.Endpoints = append(vWelcome.Endpoints, common.Endpoint{ Url: "/scourge/{scg-key}/progress/country", Description: "Progression by date by country", ContentType: "JSON"});
		vWelcome.Endpoints = append(vWelcome.Endpoints, common.Endpoint{ Url: "/scourge/{scg-key}/country/progress", Description: "Progression by all countries by date", ContentType: "JSON"});
		vWelcome.Endpoints = append(vWelcome.Endpoints, common.Endpoint{ Url: "/scourge/{scg-key}/country/{cnt-key}/progress", Description: "Progression by specific country by date", ContentType: "JSON"});
		vWelcome.Endpoints = append(vWelcome.Endpoints, common.Endpoint{ Url: "/scourge/{scg-key}/country/deathmill", Description: "All countries deaths per million population", ContentType: "JSON"});
		vWelcome.Endpoints = append(vWelcome.Endpoints, common.Endpoint{ Url: "/scourge/{scg-key}/region/stream", Description: "Progression by day by stacked region", ContentType: "JSON"});

		if err := json.NewEncoder(w).Encode(vWelcome); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Your call was politely taken, but the JSON Encoder failed to format a reply!"))
		}
	}
}

