package main

import (
	"encoding/json"
	"errors"
	"os"
	"reflect"
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

func (state *CharacterState) equipItems(itemCodes []string) error {
	for _, code := range itemCodes {
		itemInfo, err := state.getItemInfo(code)
		if err != nil {
			return err
		}
		err = state.equipItem(code, itemInfo.Data.Type, 1)
		if err != nil {
			return err
		}
	}
	return nil
}

type CharacterEquipment struct {
	WeaponSlot           string `slot:"weapon"`
	ShieldSlot           string `slot:"shield"`
	HelmetSlot           string `slot:"helmet"`
	BodyArmorSlot        string `slot:"body_armor"`
	LegArmorSlot         string `slot:"leg_armor"`
	BootsSlot            string `slot:"boots"`
	Ring1Slot            string `slot:"ring1"`
	Ring2Slot            string `slot:"ring2"`
	AmuletSlot           string `slot:"amulet"`
	Artifact1Slot        string `slot:"artifact1"`
	Artifact2Slot        string `slot:"artifact2"`
	Artifact3Slot        string `slot:"artifact3"`
	Utility1Slot         string `slot:"utility1"`
	Utility1SlotQuantity int
	Utility2Slot         string `slot:"utility2"`
	Utility2SlotQuantity int
}

var SLOTS = []string{
	"weapon",
	"shield",
	"helmet",
	"body_armor",
	"leg_armor",
	"boots",
	"ring1",
	"ring2",
	"amulet",
	"artifact1",
	"artifact2",
	"artifact3",
	"utility1",
	"utility2",
}

func (state *CharacterState) goEquipLoadout(equipment CharacterEquipment) error {
	err := state.moveToBank()
	if err != nil {
		return err
	}

	// get rid of what we have equipped
	for _, slot := range SLOTS {
		// failure is okay if it wasn't equipped
		_ = state.unequip(slot)
	}

	err = state.dumpAtBank()
	if err != nil {
		return err
	}

	equipType := reflect.TypeOf(equipment)
	equipValue := reflect.ValueOf(equipment)

	for i := 0; i < equipType.NumField(); i++ {
		field := equipType.Field(i)

		if field.Type.Kind() != reflect.String {
			continue
		}

		itemValue := equipValue.Field(i)
		itemCode := itemValue.Interface().(string)
		slotName := field.Tag.Get("slot")

		qty := 1
		qtyField := field.Name + "Quantity"
		qtyValue := itemValue.FieldByName(qtyField).Interface().(int)

		if qtyValue != 0 {
			qty = qtyValue
		}

		err = state.withdrawItemAtBank(itemCode, qty)
		if err != nil {
			state.Logger.Warn("item not in bank!", "item", itemCode, "error", err)
			return err
		}

		err = state.equipItem(itemCode, slotName, qty)
		if err != nil {
			return err
		}
	}

	return nil
}
