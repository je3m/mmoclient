package main

import (
	"encoding/json"
	"os"
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

// Perform cooking action until inventory contains at least <quantity> of item
func (state *CharacterState) craftUntil(item string, quantity int) error {
	numberRemaining := 1

	state.Logger.Info("craft_until",
		"quantity", quantity,
		"item", item)

	for numberRemaining > 0 {
		err := state.craftItem(item, quantity)

		if err != nil {
			state.Logger.Error("Error crafting item", "error", err)
			return err
		}
		numberHas := state.getItemInventoryQty(item)
		numberRemaining = quantity - numberHas

		state.Logger.Debug("progress made",
			"action", "gathering",
			"item", item,
			"have", numberHas,
			"need", quantity,
			"remaining", numberRemaining)
	}
	return nil
}

func (state *CharacterState) craftItem(code string, qty int) error {
	type CraftItemRequest struct {
		Code     string `json:"code"`
		Quantity int    `json:"quantity"`
	}
	jsonData, err := json.Marshal(CraftItemRequest{code, qty})
	if err != nil {
		state.Logger.Error("Error marshalling request body", "error", err)
		os.Exit(1)
	}
	_, err = state.performActionAndWait("crafting", jsonData)
	if err != nil {
		state.Logger.Error("Error making crafting item", "error", err)
		return err
	}
	return nil
}
