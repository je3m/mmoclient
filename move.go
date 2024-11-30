package main

import (
	"encoding/json"
	"fmt"
)

func move(state *CharacterState, x int, y int) error {
	type MoveRequest struct {
		X int `json:"x"`
		Y int `json:"y"`
	}

	moveRequest := MoveRequest{
		X: x,
		Y: y,
	}
	// Marshal the request body to JSON
	requestBody, err := json.Marshal(moveRequest)
	if err != nil {
		fmt.Printf("Error marshalling request body: %v\n", err)
		return err
	}

	_, err = performActionAndWait(state, "move", requestBody)
	return err
}

func moveToAshTreeForest(state *CharacterState) error {
	err := move(state, -1, 0)
	if err != nil {
		fmt.Printf("Failed to move to ash tree forest: %v\n", err)
		return err
	}
	return nil
}

func moveToCopperMine(state *CharacterState) error {
	err := move(state, 2, 0)
	if err != nil {
		fmt.Printf("Failed to move to ash tree forest: %v\n", err)
		return err
	}
	return nil
}

func moveToIronMine(state *CharacterState) error {
	err := move(state, 1, 7)
	if err != nil {
		fmt.Printf("Failed to move to iron mine: %v\n", err)
		return err
	}
	return nil
}

func moveToBank(state *CharacterState) error {
	err := move(state, 4, 1)
	if err != nil {
		fmt.Printf("Failed to move to bank: %v\n", err)
		return err
	}
	return nil
}

func moveToSpruce(state *CharacterState) error {
	err := move(state, 2, 6)
	if err != nil {
		fmt.Printf("Failed to move to spruce: %v\n", err)
		return err
	}
	return nil
}

func moveToSunflower(state *CharacterState) error {
	err := move(state, 2, 2)
	if err != nil {
		fmt.Printf("Failed to move to bank: %v\n", err)
		return err
	}
	return nil
}

func moveToShrimp(state *CharacterState) error {
	err := move(state, 5, 2)
	if err != nil {
		fmt.Printf("Failed to move to shrimp: %v\n", err)
		return err
	}
	return nil
}

func moveToGudgeon(state *CharacterState) error {
	err := move(state, 4, 2)
	if err != nil {
		fmt.Printf("Failed to move to bank: %v\n", err)
		return err
	}
	return nil
}

func moveToCooking(state *CharacterState) error {
	err := move(state, 1, 1)
	if err != nil {
		fmt.Printf("Failed to move to kitchen: %v\n", err)
		return err
	}
	return nil
}
