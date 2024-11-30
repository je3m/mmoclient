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
	"AshTreeForest":         MoveRequest{-1, 0},
	"CopperMine":            MoveRequest{2, 0},
	"IronMine":              MoveRequest{1, 7},
	"Bank":                  MoveRequest{4, 1},
	"Spruce":                MoveRequest{2, 6},
	"Sunflower":             MoveRequest{2, 2},
	"Shrimp":                MoveRequest{5, 2},
	"Gudgeon":               MoveRequest{5, 2},
	"Cooking":               MoveRequest{1, 1},
	"WeaponCraftingStation": MoveRequest{2, 1},
	"Alchemy":               MoveRequest{2, 3},
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
