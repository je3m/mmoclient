package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

type BankContents struct {
	Data []struct {
		Code     string `json:"code"`
		Quantity int    `json:"quantity"`
	} `json:"data"`
	Total int `json:"total"`
	Page  int `json:"page"`
	Size  int `json:"size"`
	Pages int `json:"pages"`
}

func (state *CharacterState) getBankContents() (*BankContents, error) {
	response := new(BankContents)

	// Define the endpoint and token
	apiURL := "https://api.artifactsmmo.com/my/bank/items"

	// Create the HTTP request
	req, err := http.NewRequest("GET", apiURL, bytes.NewBuffer([]byte{}))
	if err != nil {
		state.Logger.Error("Error creating request", "error", err)
		return nil, err
	}

	// Set headers
	req.Header.Set("Accept", "application/json")
	//req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+ApiToken)

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		state.Logger.Error("Error making request", "error", err)
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	// Read and display the response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		state.Logger.Error("Error reading response body", "error", err)
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		state.Logger.Warn("Request failed with getBankContents",
			"response_code", ArtifactsResponseCode(resp.StatusCode),
			"val", int32(resp.StatusCode))
		return nil, ResponseCodeError{ArtifactsResponseCode(resp.StatusCode)}
	}

	err = json.Unmarshal(body, &response)

	if err != nil {
		state.Logger.Error("Error parsing response", "error", err)
		return nil, err
	}

	return response, err
}

func (state *CharacterState) getBankItemQty(itemCode string) (int, error) {

	items, err := state.getBankContents()
	if err != nil {
		return 0, err
	}

	for _, item := range items.Data {
		if item.Code == itemCode {
			return item.Quantity, nil
		}
	}

	return 0, nil
}

// go to bank and get at most qty items from bank and return actual amount withdrawn
func (state *CharacterState) goGetItemFromBank(itemCode string, qty int) (int, error) {

	err := state.moveToBank()
	if err != nil {
		return 0, err
	}

	qtyInBank, err := state.getBankItemQty(itemCode)
	if err != nil {
		return 0, err
	}

	if qtyInBank == 0 {
		return 0, nil
	}
	qtyToWithdraw := min(qtyInBank, qty)

	err = state.withdrawItemAtBank(itemCode, qtyToWithdraw)
	if err != nil {
		return 0, err
	}
	return qtyToWithdraw, nil
}
