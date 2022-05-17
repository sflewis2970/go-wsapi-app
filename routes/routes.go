package routes

import (
	"log"
	"net/http"

	"github.com/sflewis2970/go-wsapi-app/cmd/wsapi/controllers"
)

func SetupRoutes() {
	// Display log message
	log.Print("Setting up websocket routes")

	// Setup routes
	http.HandleFunc("/api", controllers.APIEndPoint)
}
