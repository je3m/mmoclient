package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

// Map skill->location
var craftingSpots = map[string]string{
	"weaponcrafting":  "WeaponCraftingStation",
	"cooking":         "Cooking",
	"alchemy":         "Alchemy",
	"woodcutting":     "WoodCraftingStation",
	"mining":          "MiningStation",
	"gearcrafting":    "GearCraftingStation",
	"jewelrycrafting": "JewelryCraftingStation",
}

type ItemQueryResponse struct {
	Data struct {
		Name        string `json:"name"`
		Code        string `json:"code"`
		Level       int    `json:"level"`
		Type        string `json:"type"`
		Subtype     string `json:"subtype"`
		Description string `json:"description"`
		Effects     []struct {
			Name  string `json:"name"`
			Value int    `json:"value"`
		} `json:"effects"`
		Craft struct {
			Skill string `json:"skill"`
			Level int    `json:"level"`
			Items []struct {
				Code     string `json:"code"`
				Quantity int    `json:"quantity"`
			} `json:"items"`
			Quantity int `json:"quantity"`
		} `json:"craft"`
		Tradeable bool `json:"tradeable"`
	} `json:"data"`
}

func (state *CharacterState) getItemInfo(itemCode string) (*ItemQueryResponse, error) {
	response := new(ItemQueryResponse)

	// Define the endpoint and token
	apiURL := "https://api.artifactsmmo.com/items/" + itemCode

	// Create the HTTP request
	req, err := http.NewRequest("GET", apiURL, bytes.NewBuffer([]byte{}))
	if err != nil {
		state.Logger.Error("Error creating request: %v\n", err)
		return nil, err
	}

	// Set headers
	req.Header.Set("Accept", "application/json")
	//req.Header.Set("Content-Type", "application/json")

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		state.Logger.Error("Error making request: %v\n", err)
		return nil, err
	}
	defer resp.Body.Close()

	// Read and display the response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		state.Logger.Error("Error reading response body\n", err)
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		state.Logger.Warn("getItemInfo Request failed with status",
			"response_code", ArtifactsResponseCode(resp.StatusCode),
			"val", resp.StatusCode)
		return nil, ResponseCodeError{ArtifactsResponseCode(resp.StatusCode)}
	}

	err = json.Unmarshal(body, &response)

	if err != nil {
		state.Logger.Error("Error parsing response: %v\n", err)
		return nil, err
	}

	return response, nil
}

// go get ingredients and craft an item at it's workshop
func (state *CharacterState) goCraftItem(itemCode string, qty int) error {
	itemInfo, err := state.getItemInfo(itemCode)
	if err != nil {
		return err
	}

	err = state.moveToBank()
	if err != nil {
		return err
	}

	for _, item := range itemInfo.Data.Craft.Items {
		totalRequired := item.Quantity * qty
		err = state.withdrawItemAtBank(item.Code, totalRequired)
		if err != nil {
			return err
		}
	}

	err = state.move(craftingSpots[itemInfo.Data.Craft.Skill])
	if err != nil {
		return err
	}

	err = state.craftItem(itemCode, qty)
	if err != nil {
		return err
	}
	return nil
}

// try and get item from bank, craft if need more
func (state *CharacterState) goGetOrCraftItem(itemCode string, qty int) error {
	_, err := state.goGetItemFromBank(itemCode, qty)
	if err != nil {
		return err
	}

	amountInInventory := state.getItemInventoryQty(itemCode)

	if amountInInventory == qty {
		return nil
	}

	err = state.goCraftItem(itemCode, qty-amountInInventory)
	if err != nil {
		return err
	}
	return nil
}

// do whatever it takes to get an item crafted
func (state *CharacterState) goCraftItemAndDependencies(itemCode string, qty int) error {
	qtyInInventory := state.getItemInventoryQty(itemCode)

	if qtyInInventory >= qty {
		state.Logger.Debug("inventory already contains enough item. skipping...",
			"item", itemCode,
			"qty", qty,
			"qtyInInventory", qtyInInventory)
		return nil
	}

	itemInfo, err := state.getItemInfo(itemCode)
	if err != nil {
		return err
	}

	// make prerequisite items
	for _, ingredient := range itemInfo.Data.Craft.Items {
		err = state.goCraftItemAndDependencies(ingredient.Code, ingredient.Quantity*qty)
		if err != nil {
			return err
		}
	}

	err = state.goGetOrCraftItem(itemCode, qty)
	if err != nil {
		return err
	}
	return nil
}
