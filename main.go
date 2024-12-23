package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"time"
)

var MonsterDB *MonsterResponse

var ApiToken = ""

func (state *CharacterState) getItemInventoryQty(itemName string) int {
	inv := state.Inventory
	for _, item := range inv {
		if item.Code == itemName {
			return item.Quantity
		}
	}
	return 0
}

func (state *CharacterState) grindCrafting(itemToCraft string, ingredientName string, numIngredientPerCraft int) error {
	ingredientQty := state.getItemInventoryQty(ingredientName)

	for ingredientQty > numIngredientPerCraft {
		itemsToCraft := ingredientQty / numIngredientPerCraft

		err := state.craftItem(itemToCraft, itemsToCraft)
		if err != nil {
			return err
		}

		err = state.recycleItem(itemToCraft, itemsToCraft)
		if err != nil {
			return err
		}
		ingredientQty = state.getItemInventoryQty(ingredientName)
	}
	return nil
}

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
func (state *CharacterState) recycleItem(code string, qty int) error {
	type RecycleItemRequest struct {
		Code     string `json:"code"`
		Quantity int    `json:"quantity"`
	}
	jsonData, err := json.Marshal(RecycleItemRequest{code, qty})
	if err != nil {
		return err
	}

	_, err = state.performActionAndWait("recycling", jsonData)
	if err != nil {
		return err
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

func (state *CharacterState) useItem(code string, quantity int) error {
	type UseItemRequest struct {
		Code     string `json:"code"`
		Quantity int    `json:"quantity"`
	}
	jsonData, err := json.Marshal(UseItemRequest{code, quantity})
	if err != nil {
		state.Logger.Error("Error marshalling request body", "error", err)
		os.Exit(1)
	}
	_, err = state.performActionAndWait("use", jsonData)
	if err != nil {
		state.Logger.Error("Error using item", "error", err)
		return err
	}
	return nil
}

func (state *CharacterState) depositItemAtBank(code string, qty int) error {
	type DepositItemRequest struct {
		Code     string `json:"code"`
		Quantity int    `json:"quantity"`
	}
	jsonData, err := json.Marshal(DepositItemRequest{code, qty})
	if err != nil {
		state.Logger.Error("Error marshalling request body", "error", err)
		os.Exit(1)
	}
	_, err = state.performActionAndWait("bank/deposit", jsonData)
	if err != nil {
		state.Logger.Error("Error depositing item", "reason", err)
		return err
	}
	return nil
}
func (state *CharacterState) withdrawItemAtBank(code string, qty int) error {
	type WithdrawItemRequest struct {
		Code     string `json:"code"`
		Quantity int    `json:"quantity"`
	}
	jsonData, err := json.Marshal(WithdrawItemRequest{code, qty})
	if err != nil {
		state.Logger.Error("Error marshalling request body", "error", err)
		os.Exit(1)
	}
	_, err = state.performActionAndWait("bank/withdraw", jsonData)
	if err != nil {
		state.Logger.Warn("Error withdrawing item", "error", err)
		return err
	}
	return nil
}

func (state *CharacterState) dumpAtBank() error {
	err := state.moveToBank()
	if err != nil {
		return err
	}

	for _, item := range state.Inventory {
		if item.Quantity > 0 {
			err := state.depositItemAtBank(item.Code, item.Quantity)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (state *CharacterState) getInventoryFull() bool {
	inv := state.Inventory
	count := 0
	for _, item := range inv {
		count += item.Quantity
	}

	return state.InventoryMaxItems <= count
}

func setApiToken() {
	apiTok, err := os.ReadFile("token.txt")
	if err != nil {
		slog.Error("Failed to read API token", "error", err)
		os.Exit(1)
	}
	ApiToken = string(apiTok)
}
func (state *CharacterState) setupLogging() error {
	handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})

	logger := slog.New(handler)
	state.Logger = logger.WithGroup(state.Name)
	slog.SetDefault(logger)
	return nil
}

func setupMonsterDB() error {
	db, err := getMonsterDB()
	if err != nil {
		return err
	}
	MonsterDB = db
	return nil
}

func doGameLoop(state *CharacterState) error {
	switch state.Name {
	case "lily":
		return state.lilyLoop()
	case "timothy":
		return state.timothyLoop()
	case "chad":
		return state.chadLoop()
	case "squidward":
		return state.squidwardLoop()
	case "mike":
		return state.mikeLoop()
	default:
		return errors.New("unknown character: " + state.Name)
	}
}

func setupStates(stateRefs map[string]*CharacterState) error {
	states, err := getGameStatus()
	if err != nil {
		slog.Error("Failed to get game status", "error", err)
		os.Exit(1)
	}

	for i, state := range states {
		stateRefs[state.Name] = &states[i]

		err = states[i].setupLogging()
		if err != nil {
			println("Failed to setup logging", "error", err)
			os.Exit(1)
		}
	}
	return err
}

func makePidfile(characterName string) (string, error) {
	pid := os.Getpid()
	pidFile := characterName + ".pid"

	file, err := os.Create(pidFile)
	if err != nil {
		return "", err
	}

	defer func(file *os.File) {
		_ = file.Close()
	}(file)

	// Write the PID to the file
	_, err = file.WriteString(strconv.Itoa(pid))
	if err != nil {
		return "", err
	}
	return pidFile, nil
}

func main() {
	stateRefs := make(map[string]*CharacterState)
	characterName := os.Args[1]

	pidFile, err := makePidfile(characterName)
	if err != nil {
		fmt.Println("Error creating PID file:", err)
		return
	}

	defer func(name string) {
		_ = os.Remove(name)
	}(pidFile)

	setApiToken()

	err = setupStates(stateRefs)
	if err != nil {
		slog.Error("Failed to setup states", "error", err)
		return
	}

	state := stateRefs[characterName]
	if state == nil {
		slog.Error("Character not found", "characterName", characterName)
		return
	}

	err = setupMonsterDB()
	if err != nil {
		slog.Error("Failed to get Monster DB", "error", err)
		return
	}

	lastFail := time.Now().Add(time.Duration(-1) * time.Hour)

	for {
		err := doGameLoop(state)
		if err != nil {
			currentTime := time.Now()
			timeSinceLastFail := currentTime.Sub(lastFail)
			if timeSinceLastFail < time.Duration(5)*time.Minute {
				state.Logger.Error("Error in gameloop. killing program", "error", err)
				return
			} else {
				state.Logger.Error("Error in gameloop. rebooting character...", "error", err)
			}
			lastFail = time.Now()
		}
	}
}
