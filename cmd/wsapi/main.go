package main

import (
	"log"
	"net/http"

	"github.com/sflewis2970/go-wsapi-app/routes"
)

func main() {
	// Initialize logging
	log.SetFlags(log.Ldate | log.Lshortfile)

	// Setup routes
	routes.SetupRoutes()

	// Start Server
	log.Print("Websocket server is ready...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
