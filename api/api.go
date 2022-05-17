package api

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
)

const (
	QuoteURL          string = "https://quotes-by-api-ninjas.p.rapidapi.com/v1/quotes"
	DictionaryURL     string = "https://quotes-by-api-ninjas.p.rapidapi.com/v1/dictionary"
	RapidAPIHostKey   string = "X-RapidAPI-Host"
	RapidAPIHostValue string = "quotes-by-api-ninjas.p.rapidapi.com"
	RapidAPIKey       string = "X-RapidAPI-Key"
	RapidAPIValue     string = "1f8720c0c7msh43fe783209a6813p1833b2jsnc2300c30b9a9"
)

type QuoteResponse struct {
	Quote    string `json:"quote"`
	Author   string `json:"author"`
	Category string `json:"category"`
}

func QuoteRequest(categoryStr string, limitStr string) ([]QuoteResponse, error) {
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

	request, requestErr := http.NewRequest("GET", url, nil)
	if requestErr != nil {
		return nil, requestErr
	}

	request.Header.Add(RapidAPIHostKey, RapidAPIHostValue)
	request.Header.Add(RapidAPIKey, RapidAPIValue)

	response, responseErr := http.DefaultClient.Do(request)
	if responseErr != nil {
		return nil, requestErr
	}
	defer response.Body.Close()

	body, readErr := ioutil.ReadAll(response.Body)
	if readErr != nil {
		return nil, readErr
	}

	responses := make([]QuoteResponse, 0)
	unmarshalErr := json.Unmarshal(body, &responses)
	if unmarshalErr != nil {
		return nil, unmarshalErr
	}

	return responses, nil
}

func DictionaryRequest(wordStr string) (*http.Response, error) {
	// Build URL string with
	if len(wordStr) == 0 {
		return nil, errors.New("word si required")
	}

	url := DictionaryURL + wordStr

	request, requestErr := http.NewRequest("GET", url, nil)
	if requestErr != nil {
		return nil, requestErr
	}

	request.Header.Add(RapidAPIHostKey, RapidAPIHostValue)
	request.Header.Add(RapidAPIKey, RapidAPIValue)

	response, responseErr := http.DefaultClient.Do(request)
	if responseErr != nil {
		return nil, requestErr
	}
	defer response.Body.Close()

	body, _ := ioutil.ReadAll(response.Body)

	log.Print(response)
	log.Print(string(body))

	return response, nil
}
