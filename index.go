package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/zephryl/honors/common"
)

// On index, return a StandardResponse funky message
func IndexHandler(s *common.System) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")

		// Create a successful greeting, send success response
		var vWelcome common.Welcome
		vWelcome.Status = 0
		vWelcome.Greeting = "Zephry Academy v0.1.0 is up and running. Navigate to /routes for a list of endpoints."
		vWelcome.Date = time.Now()

		if err := json.NewEncoder(w).Encode(vWelcome); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Your call was politely taken, but the JSON Encoder failed to format a reply!"))
		}
	}
}
