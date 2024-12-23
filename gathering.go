package main

import (
	"encoding/json"
	"os"
)

// Perform gathering action until inventory contains at least <quantity> of item or inventory is full
func (state *CharacterState) gatherUntil(item string, quantity int) error {
	numberRemaining := 1

	state.Logger.Info("gathering_until",
		"quantity", quantity,
		"item", item)

	for numberRemaining > 0 {
		if state.getInventoryFull() {
			state.Logger.Warn("Inventory full. gathering done\n")
			break
		}

		resp, err := state.performActionAndWait("gathering", []byte{})
		if err != nil {
			state.Logger.Error("Error making request", "error", err)
			return err
		}
		numberHas := resp.Data.Character.getItemInventoryQty(item)
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

// Go to location and gather <quantity> of <item> or exit on full inventory
func (state *CharacterState) goGather(item string, quantity int) error {
	location, err := state.getResourceLocation(item)
	if err != nil {
		return err
	}

	jsonData, err := json.Marshal(location)
	if err != nil {
		state.Logger.Error("Error marshalling request body", "error", err)
		os.Exit(1)
	}

	state.Logger.Debug("moving", "location", location)
	_, err = state.performActionAndWait("move", jsonData)
	if err != nil {
		return err
	}
	return state.gatherUntil(item, quantity)
}

func (state *CharacterState) gatherGameLoop(item string) error {
	for {
		err := state.dumpAtBank()
		if err != nil {
			return err
		}
		err = state.goGather(item, state.InventoryMaxItems)
		if err != nil {
			return err
		}
	}
}
