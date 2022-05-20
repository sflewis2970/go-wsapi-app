package controllers

import (
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/websocket"
	"github.com/sflewis2970/go-wsapi-app/api"
)

// Upgrader struct used to defining the websocket buffer sizes
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type WSRequest struct {
	Request  string `json:"request"`
	Category string `json:"category"`
	Limit    string `json:"limit"`
	Word     string `json:"word"`
}

type WSQuoteResponse struct {
	Request   string `json:"request"`
	Timestamp string `json:"timestamp"`
	Quote     string `json:"quote"`
	Category  string `json:"category"`
	Author    string `json:"author"`
	Warning   string `json:"warning,omitempty"`
	Error     string `json:"error,omitempty"`
}

type WSDictionaryResponse struct {
	Request    string `json:"request"`
	Timestamp  string `json:"timestamp"`
	Definition string `json:"definition"`
	Word       string `json:"word"`
	Valid      bool   `json:"valid"`
	Error      string `json:"error,omitempty"`
}

func processQuoteRequest(wsRequest WSRequest) []WSQuoteResponse {
	var quoteResponses []api.QuoteResponse
	var quoteErr error

	// Send request to Public API
	var wsQuoteResponses []WSQuoteResponse
	var timestamp string
	quoteErr, quoteResponses, timestamp = api.QuoteRequest(wsRequest.Category, wsRequest.Limit)
	if quoteErr != nil {
		wsQuoteResponse := WSQuoteResponse{
			Request:   wsRequest.Request,
			Timestamp: timestamp,
			Error:     quoteErr.Error(),
		}

		wsQuoteResponses = append(wsQuoteResponses, wsQuoteResponse)
	} else {
		// Construct Websocket Response message
		if len(quoteResponses) == 0 {
			warningMsg := "no data returned, category value may be misspelled or not found"
			wsQuoteResponse := WSQuoteResponse{
				Request:   wsRequest.Request,
				Timestamp: timestamp,
				Warning:   warningMsg,
			}

			wsQuoteResponses = append(wsQuoteResponses, wsQuoteResponse)
		} else {
			for _, quoteResponse := range quoteResponses {
				wsQuoteResponse := WSQuoteResponse{
					Request:   wsRequest.Request,
					Timestamp: timestamp,
					Quote:     quoteResponse.Quote,
					Author:    quoteResponse.Author,
					Category:  quoteResponse.Category,
				}

				wsQuoteResponses = append(wsQuoteResponses, wsQuoteResponse)
			}
		}
	}

	return wsQuoteResponses
}

func processDictionaryRequest(wsRequest WSRequest) *WSDictionaryResponse {
	var dictResponse api.DictionaryResponse
	var dictErr error
	var timestamp string

	dictErr, dictResponse, timestamp = api.DictionaryRequest(wsRequest.Word)

	var wsDictResponse *WSDictionaryResponse
	if dictErr != nil {
		wsDictResponse = &WSDictionaryResponse{
			Request:   wsRequest.Request,
			Timestamp: timestamp,
			Error:     dictErr.Error(),
		}

	} else {
		wsDictResponse = &WSDictionaryResponse{
			Request:    wsRequest.Request,
			Timestamp:  timestamp,
			Definition: dictResponse.Definition,
			Word:       dictResponse.Word,
			Valid:      dictResponse.Valid,
		}
	}

	return wsDictResponse
}

func socketReader(wsConn *websocket.Conn) {
	for {
		// Read message from client
		var wsRequest WSRequest
		readErr := wsConn.ReadJSON(&wsRequest)
		if readErr != nil {
			log.Print("Error reading from socket: ", readErr)
			return
		}

		log.Printf("%s received %s request from client\n", wsConn.RemoteAddr(), wsRequest.Request)

		// Prcocess API request
		switch strings.ToLower(wsRequest.Request) {
		case "quote":
			wsQuoteResponses := make([]WSQuoteResponse, 0)
			wsQuoteResponses = processQuoteRequest(wsRequest)

			// Send response to client
			writeErr := wsConn.WriteJSON(wsQuoteResponses)
			if writeErr != nil {
				log.Print("Error writing to socket: ", writeErr)
				return
			}

		case "dictionary":
			wsDictResponse := processDictionaryRequest(wsRequest)

			// Send response to client
			writeErr := wsConn.WriteJSON(*wsDictResponse)
			if writeErr != nil {
				log.Print("Error writing to socket: ", writeErr)
				return
			}

		default:
			log.Print("Unhandled request: ", wsRequest.Request)
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
