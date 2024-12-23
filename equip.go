package main

import (
	"encoding/json"
	"errors"
	"os"
)

func (state *CharacterState) unequip(slot string) error {
	type UnequipRequest struct {
		Slot string `json:"slot"`
	}

	jsonData, err := json.Marshal(UnequipRequest{slot})
	if err != nil {
		state.Logger.Error("Error marshalling request body", "error", err)
		os.Exit(1)
	}

	_, err = state.performActionAndWait("unequip", jsonData)

	var respError ResponseCodeError

	if errors.As(err, &respError) {
		// CodeInvalidPayload means that it was not equipped
		if respError.code != CodeInvalidPayload {
			return err
		}
	}
	return nil
}

func (state *CharacterState) equipItem(code string, slot string, qty int) error {
	type EquipItemRequest struct {
		Code     string `json:"code"`
		Slot     string `json:"slot"`
		Quantity int    `json:"quantity"`
	}
	jsonData, err := json.Marshal(EquipItemRequest{code, slot, qty})
	if err != nil {
		state.Logger.Error("Error marshalling request body", "error", err)
		os.Exit(1)
	}
	_, err = state.performActionAndWait("equip", jsonData)
	if err != nil {
		state.Logger.Error("Error equipping item", "error", err)
		return err
	}
	return nil
}
