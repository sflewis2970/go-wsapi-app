package api

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/sflewis2970/go-wsapi-app/common"
)

const (
	RapidAPIHostKey string = "X-RapidAPI-Host"
	RapidAPIKey     string = "X-RapidAPI-Key"
	RapidAPIValue   string = "1f8720c0c7msh43fe783209a6813p1833b2jsnc2300c30b9a9"

	QuoteURL          string = "https://quotes-by-api-ninjas.p.rapidapi.com/v1/quotes"
	QuoteAPIHostValue string = "quotes-by-api-ninjas.p.rapidapi.com"

	DictionaryURL          string = "https://dictionary-by-api-ninjas.p.rapidapi.com/v1/dictionary"
	DictionaryAPIHostValue string = "dictionary-by-api-ninjas.p.rapidapi.com"
)

type QuoteResponse struct {
	Quote    string `json:"quote"`
	Author   string `json:"author"`
	Category string `json:"category"`
}

type DictionaryResponse struct {
	Definition string `json:"definition"`
	Word       string `json:"word"`
	Valid      bool   `json:"valid"`
}

func QuoteRequest(categoryStr string, limitStr string) (error, []QuoteResponse, string) {
	// Build URL string
	url := QuoteURL

	// Add optional category string
	if len(categoryStr) > 0 {
		url = url + "?category=" + categoryStr
	}

	// Add optional limit string
	if len(limitStr) > 0 {
		if len(categoryStr) > 0 {
			url = url + "&limit=" + limitStr
		} else {
			url = url + "?limit=" + limitStr
		}
	}

	// Create new http request
	request, requestErr := http.NewRequest("GET", url, nil)
	if requestErr != nil {
		return requestErr, nil, ""
	}

	// Setup request headers
	request.Header.Add(RapidAPIHostKey, QuoteAPIHostValue)
	request.Header.Add(RapidAPIKey, RapidAPIValue)

	// Get response from http request
	response, responseErr := http.DefaultClient.Do(request)
	if responseErr != nil {
		return requestErr, nil, ""
	}
	defer response.Body.Close()

	// Get timestamp right after receiving a valid request
	timestamp := common.GetFormattedTime()

	// Parse request body
	body, readErr := ioutil.ReadAll(response.Body)
	if readErr != nil {
		return readErr, nil, ""
	}

	// Parse response into JSON format
	responses := make([]QuoteResponse, 0)
	unmarshalErr := json.Unmarshal(body, &responses)
	if unmarshalErr != nil {
		return unmarshalErr, nil, ""
	}

	// Return a valid response (in JSON format) as well as a timestamp
	return nil, responses, timestamp
}

func DictionaryRequest(wordStr string) (error, DictionaryResponse, string) {
	// Build URL string with
	if len(wordStr) == 0 {
		return errors.New("word is required"), DictionaryResponse{}, ""
	}

	url := DictionaryURL + "?word=" + wordStr

	request, requestErr := http.NewRequest("GET", url, nil)
	if requestErr != nil {
		return requestErr, DictionaryResponse{}, ""
	}

	request.Header.Add(RapidAPIHostKey, DictionaryAPIHostValue)
	request.Header.Add(RapidAPIKey, RapidAPIValue)

	response, responseErr := http.DefaultClient.Do(request)
	if responseErr != nil {
		return requestErr, DictionaryResponse{}, ""
	}
	defer response.Body.Close()

	// Get timestamp right after receiving a valid request
	timestamp := common.GetFormattedTime()

	body, readErr := ioutil.ReadAll(response.Body)
	if readErr != nil {
		return readErr, DictionaryResponse{}, ""
	}

	dictResponse := DictionaryResponse{}

	unmarshalErr := json.Unmarshal(body, &dictResponse)
	if unmarshalErr != nil {
		return unmarshalErr, DictionaryResponse{}, ""
	}

	return nil, dictResponse, timestamp
}
