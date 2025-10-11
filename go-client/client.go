package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func MakeRequest(client *http.Client, req *http.Request, headers map[string]string) (*http.Response, error) {
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	req.Header.Set("User-Agent", "archi~/Downloads/vlipsy-gamer-smashes-computer-1w5ENmlV.mp4 vus-bigup-client/1.0")
	req.Header.Set("Origin", "http://localhost:1323")
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		defer resp.Body.Close()
		var errResp map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&errResp); err != nil {
			return nil, fmt.Errorf("request failed with status: %s", resp.Status)
		}
		return nil, fmt.Errorf("request failed: %v", errResp)
	}
	return resp, nil
}

func decodeResponse(resp *http.Response, target interface{}) error {
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("request failed with status: %s", resp.Status)
	}
	return json.NewDecoder(resp.Body).Decode(target)
}
