package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"sync"
	"time"
)

type CharacterState struct {
	Name                 string    `json:"name"`
	Account              string    `json:"account"`
	Skin                 string    `json:"skin"`
	Level                int       `json:"level"`
	Xp                   int       `json:"xp"`
	MaxXp                int       `json:"max_xp"`
	Gold                 int       `json:"gold"`
	Speed                int       `json:"speed"`
	MiningLevel          int       `json:"mining_level"`
	MiningXp             int       `json:"mining_xp"`
	MiningMaxXp          int       `json:"mining_max_xp"`
	WoodcuttingLevel     int       `json:"woodcutting_level"`
	WoodcuttingXp        int       `json:"woodcutting_xp"`
	WoodcuttingMaxXp     int       `json:"woodcutting_max_xp"`
	FishingLevel         int       `json:"fishing_level"`
	FishingXp            int       `json:"fishing_xp"`
	FishingMaxXp         int       `json:"fishing_max_xp"`
	WeaponcraftingLevel  int       `json:"weaponcrafting_level"`
	WeaponcraftingXp     int       `json:"weaponcrafting_xp"`
	WeaponcraftingMaxXp  int       `json:"weaponcrafting_max_xp"`
	GearcraftingLevel    int       `json:"gearcrafting_level"`
	GearcraftingXp       int       `json:"gearcrafting_xp"`
	GearcraftingMaxXp    int       `json:"gearcrafting_max_xp"`
	JewelrycraftingLevel int       `json:"jewelrycrafting_level"`
	JewelrycraftingXp    int       `json:"jewelrycrafting_xp"`
	JewelrycraftingMaxXp int       `json:"jewelrycrafting_max_xp"`
	CookingLevel         int       `json:"cooking_level"`
	CookingXp            int       `json:"cooking_xp"`
	CookingMaxXp         int       `json:"cooking_max_xp"`
	AlchemyLevel         int       `json:"alchemy_level"`
	AlchemyXp            int       `json:"alchemy_xp"`
	AlchemyMaxXp         int       `json:"alchemy_max_xp"`
	Hp                   int       `json:"hp"`
	MaxHp                int       `json:"max_hp"`
	Haste                int       `json:"haste"`
	CriticalStrike       int       `json:"critical_strike"`
	Stamina              int       `json:"stamina"`
	AttackFire           int       `json:"attack_fire"`
	AttackEarth          int       `json:"attack_earth"`
	AttackWater          int       `json:"attack_water"`
	AttackAir            int       `json:"attack_air"`
	DmgFire              int       `json:"dmg_fire"`
	DmgEarth             int       `json:"dmg_earth"`
	DmgWater             int       `json:"dmg_water"`
	DmgAir               int       `json:"dmg_air"`
	ResFire              int       `json:"res_fire"`
	ResEarth             int       `json:"res_earth"`
	ResWater             int       `json:"res_water"`
	ResAir               int       `json:"res_air"`
	X                    int       `json:"x"`
	Y                    int       `json:"y"`
	Cooldown             int       `json:"cooldown"`
	CooldownExpiration   time.Time `json:"cooldown_expiration"`
	WeaponSlot           string    `json:"weapon_slot"`
	ShieldSlot           string    `json:"shield_slot"`
	HelmetSlot           string    `json:"helmet_slot"`
	BodyArmorSlot        string    `json:"body_armor_slot"`
	LegArmorSlot         string    `json:"leg_armor_slot"`
	BootsSlot            string    `json:"boots_slot"`
	Ring1Slot            string    `json:"ring1_slot"`
	Ring2Slot            string    `json:"ring2_slot"`
	AmuletSlot           string    `json:"amulet_slot"`
	Artifact1Slot        string    `json:"artifact1_slot"`
	Artifact2Slot        string    `json:"artifact2_slot"`
	Artifact3Slot        string    `json:"artifact3_slot"`
	Utility1Slot         string    `json:"utility1_slot"`
	Utility1SlotQuantity int       `json:"utility1_slot_quantity"`
	Utility2Slot         string    `json:"utility2_slot"`
	Utility2SlotQuantity int       `json:"utility2_slot_quantity"`
	Task                 string    `json:"task"`
	TaskType             string    `json:"task_type"`
	TaskProgress         int       `json:"task_progress"`
	TaskTotal            int       `json:"task_total"`
	InventoryMaxItems    int       `json:"inventory_max_items"`
	Inventory            []struct {
		Slot     int    `json:"slot"`
		Code     string `json:"code"`
		Quantity int    `json:"quantity"`
	} `json:"inventory"`
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

// I think go makes you do all this shit manually :c
func updateInventory(state *CharacterState, resp *ActionResponse) {
	state.Inventory = nil
	for _, slot := range resp.Data.Character.Inventory {
		state.Inventory = append(state.Inventory, slot)
	}
	state.InventoryMaxItems = resp.Data.Character.InventoryMaxItems

}

var API_TOKEN = ""

// perform given action and block until cooldown is up
func performActionAndWait(state *CharacterState, actionName string, actionData []byte) (*ActionResponse, error) {
	response := new(ActionResponse)

	// Define the endpoint and token
	apiURL := "https://api.artifactsmmo.com/my/" + state.Name + "/action/" + actionName

	// Create the HTTP request
	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(actionData))
	if err != nil {
		fmt.Printf("Error creating request: %v\n", err)
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
		fmt.Printf("Error making request: %v\n", err)
		return nil, err
	}
	defer resp.Body.Close()

	// Read and display the response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading response body: %v\n", err)
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Request failed with status: %s\n", ArtifactsResponseCode(resp.StatusCode))
		return nil, err
	}

	err = json.Unmarshal(body, &response)

	if err != nil {
		fmt.Printf("Error parsing response: %v\n", err)
		return nil, err
	}

	updateInventory(state, response)

	cooldown := response.Data.Cooldown.RemainingSeconds
	fmt.Printf("Waiting %v seconds to finish action %s\n", cooldown, actionName)
	time.Sleep(time.Duration(cooldown) * time.Second)

	return response, err
}

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

func fight(state *CharacterState) {
	performActionAndWait(state, "fight", []byte{})
}

func rest(state *CharacterState) {
	performActionAndWait(state, "rest", []byte{})
}

func gathering(state *CharacterState) {
	performActionAndWait(state, "gathering", []byte{})
}

func getItemInventoryQty(state *CharacterState, itemName string) int {
	inv := state.Inventory
	for _, item := range inv {
		if item.Code == itemName {
			return item.Quantity
		}
	}
	return 0
}

// Perform gathering action until inventory contains at least <quantity> of item
func gatherUntil(state *CharacterState, item string, quantity int) error {
	numberRemaining := 1

	for numberRemaining > 0 {
		if getInventoryFull(state) {
			fmt.Printf("Inventory full. returning early\n")
			break
		}

		resp, err := performActionAndWait(state, "gathering", []byte{})
		if err != nil {
			fmt.Printf("Error making request: %v\n", err)
			return err
		}
		numberHas := getItemInventoryQty(&resp.Data.Character, item)
		numberRemaining = quantity - numberHas

		fmt.Printf(
			"Have: %v\n"+
				"Need: %v\n"+
				"Remaining: %v\n", numberHas, quantity, numberRemaining)

	}
	return nil
}

// Perform cooking action until inventory contains at least <quantity> of item
func craftUntil(state *CharacterState, item string, quantity int) error {
	numberRemaining := 1

	for numberRemaining > 0 {
		err := craftItem(state, item)

		if err != nil {
			fmt.Printf("Error crafting item: %v\n", err)
			return err
		}
		numberHas := getItemInventoryQty(state, item)
		numberRemaining = quantity - numberHas

		fmt.Printf(
			"Have: %v\n"+
				"Need: %v\n"+
				"Remaining: %v\n", numberHas, quantity, numberRemaining)

	}
	return nil
}

func unequip(state *CharacterState, slot string) {
	type UnequipRequest struct {
		Slot string `json:"slot"`
	}

	jsonData, err := json.Marshal(UnequipRequest{slot})
	if err != nil {
		fmt.Printf("Error marshalling request body: %v\n", err)
		os.Exit(1)
	}

	performActionAndWait(state, "unequip", jsonData)
}

func moveWeaponCraftStation(state *CharacterState) {
	move(state, 2, 1)
}
func craftItem(state *CharacterState, code string) error {
	type CraftItemRequest struct {
		Code string `json:"code"`
	}
	jsonData, err := json.Marshal(CraftItemRequest{code})
	if err != nil {
		fmt.Printf("Error marshalling request body: %v\n", err)
		os.Exit(1)
	}
	_, err = performActionAndWait(state, "crafting", jsonData)
	if err != nil {
		fmt.Printf("Error making crafting item: %v\n", err)
		return err
	}
	return nil
}

func equipItem(state *CharacterState, code string, slot string) error {
	type EquipItemRequest struct {
		Code string `json:"code"`
		Slot string `json:"slot"`
	}
	jsonData, err := json.Marshal(EquipItemRequest{code, slot})
	if err != nil {
		fmt.Printf("Error marshalling request body: %v\n", err)
		os.Exit(1)
	}
	_, err = performActionAndWait(state, "equip", jsonData)
	if err != nil {
		fmt.Printf("Error equiping item: %v\n", err)
		return err
	}
	return nil
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

func depositItemAtBank(state *CharacterState, code string, qty int) error {
	type DepositItemRequest struct {
		Code     string `json:"code"`
		Quantity int    `json:"quantity"`
	}
	jsonData, err := json.Marshal(DepositItemRequest{code, qty})
	if err != nil {
		fmt.Printf("Error marshalling request body: %v\n", err)
		os.Exit(1)
	}
	_, err = performActionAndWait(state, "bank/deposit", jsonData)
	if err != nil {
		fmt.Printf("Error depositing item: %v\n", err)
		return err
	}
	return nil
}
func withdrawItemAtBank(state *CharacterState, code string, qty int) error {
	type WithdrawItemRequest struct {
		Code     string `json:"code"`
		Quantity int    `json:"quantity"`
	}
	jsonData, err := json.Marshal(WithdrawItemRequest{code, qty})
	if err != nil {
		fmt.Printf("Error marshalling request body: %v\n", err)
		os.Exit(1)
	}
	_, err = performActionAndWait(state, "bank/withdraw", jsonData)
	if err != nil {
		fmt.Printf("Error depositing item: %v\n", err)
		return err
	}
	return nil
}

func dumpAtBank(state *CharacterState) {
	moveToBank(state)

	for _, item := range state.Inventory {
		if item.Quantity > 0 {
			depositItemAtBank(state, item.Code, item.Quantity)
		}
	}
}

func getInventoryFull(state *CharacterState) bool {
	inv := state.Inventory
	count := 0
	for _, item := range inv {
		count += item.Quantity
	}

	return state.InventoryMaxItems <= count
}

func squidwardLoop(currentCharacter *CharacterState) error {
	for {
		dumpAtBank(currentCharacter)

		err := moveToIronMine(currentCharacter)
		if err != nil {
			os.Exit(1)
		}

		err = gatherUntil(currentCharacter, "iron_ore", 100)
		if err != nil {
			return err
		}
	}
}

func chadLoop(currentCharacter *CharacterState) error {
	for {
		dumpAtBank(currentCharacter)

		err := moveToSpruce(currentCharacter)
		if err != nil {
			return err
		}
		err = gatherUntil(currentCharacter, "spruce_wood", 100)
		if err != nil {
			return err
		}
	}
}

func lilyLoop(currentCharacter *CharacterState) error {
	for {
		dumpAtBank(currentCharacter)

		err := moveToSunflower(currentCharacter)
		if err != nil {
			fmt.Printf("Failed to move to sunflower: %v\n", err)
			return err
		}
		err = gatherUntil(currentCharacter, "sunflower", 100)
		if err != nil {
			return err
		}
	}
}

func timothyLoop(currentCharacter *CharacterState) error {
	for {
		dumpAtBank(currentCharacter)

		err := moveToShrimp(currentCharacter)
		if err != nil {
			fmt.Printf("Failed to move to shrimp: %v\n", err)
			return err
		}
		err = gatherUntil(currentCharacter, "shrimp", 100)
		if err != nil {
			return err
		}
	}
}

func mikeLoop(currentCharacter *CharacterState) error {
	for {
		dumpAtBank(currentCharacter)
		withdrawItemAtBank(currentCharacter, "gudgeon", 100)
		err := moveToCooking(currentCharacter)
		if err != nil {
			fmt.Printf("Failed to move to kitchen: %v\n", err)
			return err
		}
		err = craftUntil(currentCharacter, "cooked_gudgeon", 100)
		if err != nil {
			return err
		}
	}
}

func setApiToken() {
	api_tok, err := os.ReadFile("token.txt")
	if err != nil {
		fmt.Printf("Failed to read API token: %v\n", err)
		os.Exit(1)
	}
	API_TOKEN = string(api_tok)
}

func main() {
	setApiToken()

	chadState := CharacterState{Name: "chad"}
	squidwardState := CharacterState{Name: "squidward"}
	lilyState := CharacterState{Name: "lily"}
	timothyState := CharacterState{Name: "timothy"}
	mikeState := CharacterState{Name: "mike"}

	go func() {
		err := chadLoop(&chadState)
		if err != nil {
			fmt.Printf("Failed to chad loop: %v\n", err)
		}
	}()
	go func() {
		err := squidwardLoop(&squidwardState)
		if err != nil {
			fmt.Printf("Failed to squward loop: %v\n", err)
		}
	}()
	go func() {
		err := lilyLoop(&lilyState)
		if err != nil {
			fmt.Printf("Failed to lily loop: %v\n", err)
		}
	}()

	go func() {
		err := timothyLoop(&timothyState)
		if err != nil {
			fmt.Printf("Failed to timothy loop: %v\n", err)
		}
	}()
	go func() {
		err := mikeLoop(&mikeState)
		if err != nil {
			fmt.Printf("Failed to mike loop: %v\n", err)
		}
	}()
	var wg = sync.WaitGroup{}
	wg.Add(1)
	defer wg.Wait()

}
