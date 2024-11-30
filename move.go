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
}

func move(state *CharacterState, location string) error {
	destination := locations[location]
	// Marshal the request body to JSON
	requestBody, err := json.Marshal(destination)
	if err != nil {
		state.Logger.Error("Error marshalling request body:", "error", err)
		os.Exit(1)
	}

	_, err = performActionAndWait(state, "move", requestBody)

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

func moveWeaponCraftStation(state *CharacterState) error {
	return move(state, "WeaponCraftingStation")
}
func moveToAshTreeForest(state *CharacterState) error {
	return move(state, "AshTreeForest")
}

func moveToCopperMine(state *CharacterState) error {
	return move(state, "CopperMine")
}

func moveToIronMine(state *CharacterState) error {
	return move(state, "IronMine")
}

func moveToBank(state *CharacterState) error {
	return move(state, "Bank")
}

func moveToSpruce(state *CharacterState) error {
	return move(state, "Spruce")
}

func moveToSunflower(state *CharacterState) error {
	return move(state, "Sunflower")
}

func moveToShrimp(state *CharacterState) error {
	return move(state, "Shrimp")
}

func moveToGudgeon(state *CharacterState) error {
	return move(state, "Gudgeon")
}

func moveToCooking(state *CharacterState) error {
	return move(state, "Cooking")
}

func moveToAlchemy(state *CharacterState) error {
	return move(state, "Alchemy")
}

func moveToChicken(state *CharacterState) error {
	return move(state, "Chicken")
}
