package main

import (
	"fmt"
	"log"
	"os"
    "flag"
	"net/http"
	"github.com/zephryl/zephry/common"
)

func main() {
	// Get a pointer to a System struct, typically from CL Flags. Any error here, and blow out.
	// These are required: -user, -password, -schema and -port
	var vSystem = new(common.System); 
	if err := vSystem.GetFlags(); err != nil {
		fmt.Fprintf(os.Stderr, "Flag error: %v Flag defaults are:\n", err)
		fmt.Fprintf(os.Stderr, "Flag defaults are:\n")
		flag.PrintDefaults()
		log.Fatal("Zephry Rest Service has stopped.");
	}
	// Setup cookie security
	vSystem.SetCookie();
	// System middleware
	
	// Fetch a route-filled gorilla router, and serve
	router := MuxRouter(vSystem)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", vSystem.Port), (router))) // handlers.CORS(corsHeaders, corsMethods, corsOrigins, corsCreds)(router)
}
