package controllers

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/sflewis2970/go-wsapi-app/api"
)

// Upgrader struct used to defining the websocket buffer sizes
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type ClientRequest struct {
	Request  string `json:"request"`
	Category string `json:"category"`
	Limit    string `json:"limit"`
	Word     string `json:"word"`
}

type ServerResponse struct {
	Request  string `json:"request"`
	Quote    string `json:"quote"`
	Author   string `json:"author"`
	Category string `json:"category"`
}

func processQuoteRequest() {
}

func processDictionaryRequest() {
}

func socketReader(wsConn *websocket.Conn) {
	for {
		// Read message from client
		var clientRequest ClientRequest
		readErr := wsConn.ReadJSON(&clientRequest)
		if readErr != nil {
			log.Print("Error reading from socket: ", readErr)
			return
		}

		log.Printf("%s received %s request from client\n", wsConn.RemoteAddr(), clientRequest.Request)

		// Prcocess API request
		var quoteResponses []api.QuoteResponse
		var quoteErr error
		var dictResponse *http.Response
		var dictErr error
		serverResponses := make([]ServerResponse, 0)

		switch clientRequest.Request {
		case "Quote":
			quoteResponses, quoteErr = api.QuoteRequest(clientRequest.Category, clientRequest.Limit)
			if quoteErr != nil {
				log.Print("Error processing quote request: ", quoteErr)
				return
			}

			for _, quoteResponse := range quoteResponses {
				serverResponse := ServerResponse{
					Request:  clientRequest.Request,
					Quote:    quoteResponse.Quote,
					Author:   quoteResponse.Author,
					Category: quoteResponse.Category,
				}

				serverResponses = append(serverResponses, serverResponse)
			}

		case "Dictionary":
			dictResponse, dictErr = api.DictionaryRequest(clientRequest.Word)
			if dictErr != nil {
				log.Print("Error processing quote request: ", quoteErr)
				return
			}

			if dictResponse == nil {
				log.Print("Parsing error")
				return
			}

		default:
			log.Print("Unhandled request: ", clientRequest.Request)
		}

		// Send response to client
		// serverResponse.Request = clientRequest.Request
		writeErr := wsConn.WriteJSON(serverResponses)
		if writeErr != nil {
			log.Print("Error writing to socket: ", writeErr)
			return
		}

		log.Printf("%s sent client request response to client\n", wsConn.RemoteAddr())
	}
}

func APIEndPoint(w http.ResponseWriter, r *http.Request) {
	// Check Origin
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	// Upgrade connection
	wsConn, socketErr := upgrader.Upgrade(w, r, nil)
	if socketErr != nil {
		log.Print("Error upgrading socket: ", socketErr)
		return
	}
	defer wsConn.Close()

	// At this point the server is connected to the client
	log.Print("Client connected")

	// Function for waiting for messages from client
	socketReader(wsConn)

	// Display a log message when the socketReader has returned
	log.Print("sockerRead has return...")
}
