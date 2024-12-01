package main

import (
	"bytes"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"sync"
	"time"
)

type CharacterResponse struct {
	Data []CharacterState
}

type ActionResponse struct {
	Data struct {
		Cooldown struct {
			TotalSeconds     int       `json:"total_seconds"`
			RemainingSeconds int       `json:"remaining_seconds"`
			StartedAt        time.Time `json:"started_at"`
			Expiration       time.Time `json:"expiration"`
			Reason           string    `json:"reason"`
		} `json:"cooldown"`
		Details struct {
			Xp    int `json:"xp"`
			Items []struct {
				Code     string `json:"code"`
				Quantity int    `json:"quantity"`
			} `json:"items"`
		} `json:"details"`
		Character CharacterState `json:"character"`
	} `json:"data"`
}

type MonsterResponse struct {
	Data []struct {
		Name        string `json:"name"`
		Code        string `json:"code"`
		Level       int    `json:"level"`
		Hp          int    `json:"hp"`
		AttackFire  int    `json:"attack_fire"`
		AttackEarth int    `json:"attack_earth"`
		AttackWater int    `json:"attack_water"`
		AttackAir   int    `json:"attack_air"`
		ResFire     int    `json:"res_fire"`
		ResEarth    int    `json:"res_earth"`
		ResWater    int    `json:"res_water"`
		ResAir      int    `json:"res_air"`
		MinGold     int    `json:"min_gold"`
		MaxGold     int    `json:"max_gold"`
		Drops       []struct {
			Code        string `json:"code"`
			Rate        int    `json:"rate"`
			MinQuantity int    `json:"min_quantity"`
			MaxQuantity int    `json:"max_quantity"`
		} `json:"drops"`
	} `json:"data"`
	Total int `json:"total"`
	Page  int `json:"page"`
	Size  int `json:"size"`
	Pages int `json:"pages"`
}

var MonsterDB *MonsterResponse

var API_TOKEN = ""

type MapResponse struct {
	Data []struct {
		Name    string `json:"name"`
		Skin    string `json:"skin"`
		X       int    `json:"x"`
		Y       int    `json:"y"`
		Content struct {
			Type string `json:"type"`
			Code string `json:"code"`
		} `json:"content"`
	} `json:"data"`
	Total int `json:"total"`
	Page  int `json:"page"`
	Size  int `json:"size"`
	Pages int `json:"pages"`
}

func getMonsterLocation(state *CharacterState, monsterName string) (*MoveRequest, error) {
	response := new(MapResponse)
	retval := new(MoveRequest)

	// Define the endpoint and token
	apiURL := "https://api.artifactsmmo.com/maps"

	u, _ := url.Parse(apiURL)

	q := u.Query()
	q.Add("content_code", monsterName)
	q.Add("content_type", "monster")

	u.RawQuery = q.Encode()

	// Create the HTTP request
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, err
	}

	// Set headers
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+API_TOKEN)

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Read and display the response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, ResponseCodeError{ArtifactsResponseCode(resp.StatusCode)}
	}

	err = json.Unmarshal(body, &response)

	if err != nil {
		slog.Error("Error parsing response: %v\n", err)
		return nil, err
	}

	for _, spot := range response.Data {
		// TODO: go to closest
		retval.X = spot.X
		retval.Y = spot.Y
	}
	return retval, nil
}

// perform given action and block until cooldown is up
func (state *CharacterState) performActionAndWait(actionName string, actionData []byte) (*ActionResponse, error) {
	response := new(ActionResponse)

	// Define the endpoint and token
	apiURL := "https://api.artifactsmmo.com/my/" + state.Name + "/action/" + actionName

	// Create the HTTP request
	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(actionData))
	if err != nil {
		state.Logger.Error("Error creating request: %v\n", err)
		return nil, err
	}

	// Set headers
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+API_TOKEN)

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
		state.Logger.Warn("Request failed with status",
			"action_name", actionName,
			"response_code", ArtifactsResponseCode(resp.StatusCode))
		return nil, ResponseCodeError{ArtifactsResponseCode(resp.StatusCode)}
	}

	err = json.Unmarshal(body, &response)

	if err != nil {
		state.Logger.Error("Error parsing response: %v\n", err)
		return nil, err
	}

	state.updateInventory(response)

	cooldown := response.Data.Cooldown.RemainingSeconds
	state.Logger.Debug("Waiting finish action", "cooldown", cooldown, "action", actionName)
	time.Sleep(time.Duration(cooldown) * time.Second)

	return response, err
}

// query game for initial status of all characters
func getGameStatus() ([]CharacterState, error) {
	response := new(CharacterResponse)

	// Define the endpoint and token
	apiURL := "https://api.artifactsmmo.com/my/characters"

	// Create the HTTP request
	req, err := http.NewRequest("GET", apiURL, bytes.NewBuffer([]byte{}))
	if err != nil {
		return nil, err
	}

	// Set headers
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+API_TOKEN)

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Read and display the response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, ResponseCodeError{ArtifactsResponseCode(resp.StatusCode)}
	}

	err = json.Unmarshal(body, &response)

	if err != nil {
		slog.Error("Error parsing response: %v\n", err)
		return nil, err
	}

	return response.Data, nil
}

func (state *CharacterState) fight() {
// query game for initial status of all characters
func getMonsterDB() (*MonsterResponse, error) {
	response := new(MonsterResponse)

	// Define the endpoint and token
	apiURL := "https://api.artifactsmmo.com/monsters"

	// Create the HTTP request
	req, err := http.NewRequest("GET", apiURL, bytes.NewBuffer([]byte{}))
	if err != nil {
		return nil, err
	}

	// Set headers
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+API_TOKEN)

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Read and display the response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, ResponseCodeError{ArtifactsResponseCode(resp.StatusCode)}
	}

	err = json.Unmarshal(body, &response)

	if err != nil {
		slog.Error("Error parsing response: %v\n", err)
		return nil, err
	}

	return response, nil
}

func fight(state *CharacterState) {
	state.Logger.Debug("fighting", "start_hp", state.Hp)
	state.performActionAndWait("fight", []byte{})
	state.Logger.Debug("fighting", "end_hp", state.Hp)
}

func rest(state *CharacterState) {
	state.performActionAndWait("rest", []byte{})
}

func gathering(state *CharacterState) {
	state.performActionAndWait("gathering", []byte{})
}

func (state *CharacterState) getItemInventoryQty(itemName string) int {
	inv := state.Inventory
	for _, item := range inv {
		if item.Code == itemName {
			return item.Quantity
		}
	}
	return 0
}

// Perform gathering action until inventory contains at least <quantity> of item
func (state *CharacterState) gatherUntil(item string, quantity int) error {
	numberRemaining := 1

	state.Logger.Info("gathering_until",
		"quantity", quantity,
		"item", item)

	for numberRemaining > 0 {
		if state.getInventoryFull() {
			state.Logger.Warn("Inventory full. returning early\n")
			break
		}

		resp, err := state.performActionAndWait("gathering", []byte{})
		if err != nil {
			state.Logger.Error("Error making request", err)
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

// Perform cooking action until inventory contains at least <quantity> of item
func (state *CharacterState) craftUntil(item string, quantity int) error {
	numberRemaining := 1

	state.Logger.Info("craft_until",
		"quantity", quantity,
		"item", item)

	for numberRemaining > 0 {
		err := craftItem(state, item, quantity)

		if err != nil {
			state.Logger.Error("Error crafting item: %v\n", err)
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

func findWorthyEnemy(state *CharacterState) string {
	maxLevel := state.Level - 1
	mostWorthy := ""
	highestLevel := 0

	for _, monster := range MonsterDB.Data {
		if monster.Level < maxLevel && monster.Level > highestLevel {
			mostWorthy = monster.Code
			highestLevel = monster.Level
		}
	}

	state.Logger.Debug(mostWorthy + " deemed worthy")

	return mostWorthy
}

func fightWorthyEnemy(state *CharacterState, healing_item string, heal_amount int) error {
	//enemy := findWorthyEnemy(state)
	location, err := getMonsterLocation(state, "yellow_slime")
	if err != nil {
		return err
	}

	jsonData, err := json.Marshal(location)
	if err != nil {
		state.Logger.Error("Error marshalling request body: %v\n", err)
		os.Exit(1)
	}

	state.Logger.Debug("moving", "location", location)
	_, err = performActionAndWait(state, "move", jsonData)
	if err != nil {
		return err
	}
	err = fightUntilLowInventory(state, healing_item, heal_amount)
	if err != nil {
		return err
	}
	return nil
}

// heal as much as possible without waste
func (state *CharacterState) healEfficient(healing_item string, amount_heal int) error {
	numHave := state.getItemInventoryQty(healing_item)
	hpToHeal := state.MaxHp - state.Hp
	numNeeded := hpToHeal / amount_heal

	numToConsume := min(numNeeded, numHave)
	if hpToHeal > 0 { //TODO: hack
		state.Logger.Info("healing", "start_hp", state.Hp)

		err := useItem(state, healing_item, numToConsume+1) //TODO: +1 is bad hack to keep killing yellow slimes overnight
		if err != nil {
			return err
		}
		state.Logger.Info("healing", "end_hp", state.Hp)
	}

	return nil
}

// fight until out of hp and healing item
func (state *CharacterState) fightUntilLowInventory(healing_item string, amount_heal int) error {
	numHealItem := state.getItemInventoryQty(healing_item)
	fight_count := 0
	state.Logger.Info("fight_forever",
		"healing_item", healing_item)

	for numHealItem > 0 {
		err := state.healEfficient(healing_item, amount_heal)

		if err != nil {
			state.Logger.Error("Error healing", "error", err)
			return err
		}
		if state.getInventoryFull() {
			// no point in fighting bc we get no loot
			return nil
		}

		state.fight()
		fight_count++
		numHealItem = state.getItemInventoryQty(healing_item)

		state.Logger.Debug("progress made",
			"action", "fighting",
			"fights_won", fight_count,
			"remaining", numHealItem)
	}
	return nil
}

func (state *CharacterState) unequip(slot string) {
	type UnequipRequest struct {
		Slot string `json:"slot"`
	}

	jsonData, err := json.Marshal(UnequipRequest{slot})
	if err != nil {
		state.Logger.Error("Error marshalling request body: %v\n", err)
		os.Exit(1)
	}

	state.performActionAndWait("unequip", jsonData)
}

func (state *CharacterState) craftItem(code string) error {
	type CraftItemRequest struct {
		Code     string `json:"code"`
		Quantity int    `json:"quantity"`
	}
	jsonData, err := json.Marshal(CraftItemRequest{code, qty})
	if err != nil {
		state.Logger.Error("Error marshalling request body: %v\n", err)
		os.Exit(1)
	}
	_, err = state.performActionAndWait("crafting", jsonData)
	if err != nil {
		state.Logger.Error("Error making crafting item: %v\n", err)
		return err
	}
	return nil
}

func (state *CharacterState) equipItem(code string, slot string) error {
	type EquipItemRequest struct {
		Code string `json:"code"`
		Slot string `json:"slot"`
	}
	jsonData, err := json.Marshal(EquipItemRequest{code, slot})
	if err != nil {
		state.Logger.Error("Error marshalling request body: %v\n", err)
		os.Exit(1)
	}
	_, err = state.performActionAndWait("equip", jsonData)
	if err != nil {
		state.Logger.Error("Error equiping item: %v\n", err)
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
		state.Logger.Error("Error marshalling request body: %v\n", err)
		os.Exit(1)
	}
	_, err = state.performActionAndWait("use", jsonData)
	if err != nil {
		state.Logger.Error("Error using item", err)
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
		state.Logger.Error("Error marshalling request body", err)
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
		state.Logger.Error("Error marshalling request body: %v\n", err)
		os.Exit(1)
	}
	_, err = state.performActionAndWait("bank/withdraw", jsonData)
	if err != nil {
		state.Logger.Warn("Error withdrawing item: %v\n", err)
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
	api_tok, err := os.ReadFile("token.txt")
	if err != nil {
		slog.Error("Failed to read API token: %v\n", err)
		os.Exit(1)
	}
	API_TOKEN = string(api_tok)
}
func (state *CharacterState) setupLogging() error {
	handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})

	logger := slog.New(handler)
	state.Logger = logger.WithGroup(state.Name)

	return nil
}

func main() {
	stateRefs := make(map[string]*CharacterState)

	setApiToken()

	states, err := getGameStatus()
	if err != nil {
		slog.Error("Failed to get game status", "error", err)
		os.Exit(1)
	}

	db, err := getMonsterDB()
	if err != nil {
		slog.Error("Failed to get Monster DB", "error", err)
		os.Exit(1)
	}
	MonsterDB = db

	for i, state := range states {
		stateRefs[state.Name] = &states[i]
		err = states[i].setupLogging()
		if err != nil {
			println("Failed to setup logging: %v\n", err)
			os.Exit(1)
		}

	}
	go func() {
		chadState := stateRefs["chad"]
		err := chadState.chadLoop()
		if err != nil {
			chadState.Logger.Error("Failed to chad loop: %v\n", err)
		}
	}()
	go func() {
		squidwardState := stateRefs["squidward"]
		err := squidwardState.squidwardLoop()
		if err != nil {
			squidwardState.Logger.Error("Failed to squward loop: %v\n", err)
		}
	}()
	go func() {
		lilyState := stateRefs["lily"]
		err := lilyState.lilyLoop()
		if err != nil {
			lilyState.Logger.Error("Failed to lily loop: %v\n", err)
		}
	}()

	go func() {
		timothyState := stateRefs["timothy"]
		err := timothyState.timothyLoop()
		if err != nil {
			timothyState.Logger.Error("Failed to timothy loop: %v\n", err)
		}
	}()

	go func() {
		mikeState := stateRefs["mike"]
		err := mikeState.mikeLoop()
		if err != nil {
			mikeState.Logger.Error("Failed to mike loop: %v\n", err)
		}
	}()

	var wg = sync.WaitGroup{}
	wg.Add(1)
	defer wg.Wait()

}
