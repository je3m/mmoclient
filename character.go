package main

import (
	"log/slog"
	"time"
)

type CharacterState struct {
	Logger               *slog.Logger
	GameLoop             func() error
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

func (state *CharacterState) String() string {
	return state.Name
}

// I think go makes you do all this shit manually :c
func (state *CharacterState) updateState(resp *ActionResponse) {
	//state.Inventory = nil
	//for _, slot := range resp.Data.Character.Inventory {
	//	state.Inventory = append(state.Inventory, slot)
	//}
	//state.InventoryMaxItems = resp.Data.Character.InventoryMaxItems
	//state.Hp = resp.Data.Character.Hp
	//state.MaxHp = resp.Data.Character.MaxHp
	//
	newLevel := resp.Data.Character.Level
	if state.Level < newLevel {
		state.Logger.Info("Level up!", "new level", newLevel)
	}

	logger := state.Logger
	*state = resp.Data.Character
	state.Logger = logger

	//state.Level = newLevel
}
