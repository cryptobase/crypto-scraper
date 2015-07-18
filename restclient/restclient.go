package restclient

import (
	"io/ioutil"
	"net/http"
)

func getContent(url string) ([]byte, error) {
	// Build the request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	return handleRequest(req)
}

func handleRequest(req http.Request) ([]byte, error) {
	// Send the request via a client
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	// Defer the closing of the body
	defer resp.Body.Close()

	// Read the content into a byte array
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}