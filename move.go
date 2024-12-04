package main

import (
	"encoding/json"
	"errors"
	"os"
)

type MoveRequest struct {
	X int `json:"x"`
	Y int `json:"y"`
}

var locations = map[string]MoveRequest{
	"AshTreeForest":         {-1, 0},
	"CopperMine":            {2, 0},
	"IronMine":              {1, 7},
	"Bank":                  {4, 1},
	"Spruce":                {2, 6},
	"Sunflower":             {2, 2},
	"Shrimp":                {5, 2},
	"Gudgeon":               {5, 2},
	"Cooking":               {1, 1},
	"WeaponCraftingStation": {2, 1},
	"Alchemy":               {2, 3},
	"Chicken":               {0, 1},
	"WoodCraftingStation":   {-2, -3},
	"MiningStation":         {1, 5},
	"Trout":                 {7, 12},
	"GearCraftingStation":   {3, 1},
}

func (state *CharacterState) move(location string) error {
	destination := locations[location]
	// Marshal the request body to JSON
	requestBody, err := json.Marshal(destination)
	if err != nil {
		state.Logger.Error("Error marshalling request body:", "error", err)
		os.Exit(1)
	}
	_, err = state.performActionAndWait("move", requestBody)

	if err != nil {
		var responseCodeError ResponseCodeError
		if errors.As(err, &responseCodeError) {
			if responseCodeError.code == CodeCharacterAlreadyMap {
				// we are already here so it's fine
				return nil
			}
		}
		state.Logger.Error("Failed to move", "location", location, "x", destination.X, "y", destination.Y)
		return err
	}
	return nil
}

func (state *CharacterState) moveToWeaponCraftingStation() error {
	return state.move("WeaponCraftingStation")
}

func (state *CharacterState) moveToWoodcraftStation() error {
	return state.move("WoodCraftingStation")
}

func (state *CharacterState) moveToAshTreeForest() error {
	return state.move("AshTreeForest")
}

func (state *CharacterState) moveToCopperMine() error {
	return state.move("CopperMine")
}

func (state *CharacterState) moveToIronMine() error {
	return state.move("IronMine")
}

func (state *CharacterState) moveToBank() error {
	return state.move("Bank")
}

func (state *CharacterState) moveToSpruce() error {
	return state.move("Spruce")
}

func (state *CharacterState) moveToSunflower() error {
	return state.move("Sunflower")
}

func (state *CharacterState) moveToShrimp() error {
	return state.move("Shrimp")
}

func (state *CharacterState) moveToGearCrafting() error {
	return state.move("GearCraftingStation")
}

func (state *CharacterState) moveToMiningStation() error {
	return state.move("MiningStation")
}
func (state *CharacterState) moveToGudgeon() error {
	return state.move("Gudgeon")
}

func (state *CharacterState) moveToTrout() error {
	return state.move("Trout")
}

func (state *CharacterState) moveToCooking() error {
	return state.move("Cooking")
}

func (state *CharacterState) moveToAlchemy() error {
	return state.move("Alchemy")
}

func (state *CharacterState) moveToChicken() error {
	return state.move("Chicken")
}
