package alpaca

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

func createRequest(method, url string, body []byte) *http.Request {
	apiKey := os.Getenv("key")
	secretKey := os.Getenv("secret")

	if apiKey == "" || secretKey == "" {
		log.Printf("Api key or secret key not found in environment")
	}

	req, err := http.NewRequest(method, url, bytes.NewBuffer(body))
	if err != nil {
		log.Printf("Failed to create request: %v", err)
	}

	req.Header.Add("accept", "application/json")
	req.Header.Add("APCA-API-KEY-ID", apiKey)
	req.Header.Add("APCA-API-SECRET-KEY", secretKey)
	req.Header.Add("Content-Type", "application/json")

	return req
}

func sendRequest(req *http.Request) ([]byte, error) {
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %v", err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read body: %v", err)
	}

	return body, nil
}
