package main

func (state *CharacterState) lilyLoop() error {
	//var err error
	//state.dumpAtBank()
	//state.withdrawItemAtBank("life_amulet", 20)
	//state.moveToJewelryCraftingStation()
	//state.recycleItem("life_amulet", 20)
	//for state.JewelrycraftingLevel < 10 {
	//	err = state.dumpAtBank()
	//	if err != nil {
	//		return err
	//	}
	//	err = state.goCraftItem("life_amulet", 5)
	//	if err != nil {
	//		return err
	//	}
	//}
	//lvl10Items := []string{"fire_and_earth_amulet", "air_and_water_amulet", "iron_ring"}
	//
	//for _, item := range lvl10Items {
	//	err = state.dumpAtBank()
	//	if err != nil {
	//		return err
	//	}
	//	err = state.goCraftItem(item, 1)
	//	if err != nil {
	//		return err
	//	}
	//	if state.JewelrycraftingLevel >= 15 {
	//		break
	//	}
	//}
	//
	//lvl15Items := []string{"air_ring", "earth_ring", "fire_ring", "life_ring", "water_ring"}
	//for _, item := range lvl15Items {
	//	err = state.dumpAtBank()
	//	if err != nil {
	//		return err
	//	}
	//	err = state.goCraftItem(item, 1)
	//	if err != nil {
	//		return err
	//	}
	//	if state.JewelrycraftingLevel >= 20 {
	//		break
	//	}
	//}

	return state.fightGameLoop("red_slime", "cooked_shrimp", 150)
	//for {
	//	state.dumpAtBank()
	//
	//	err := state.moveToSunflower()
	//	if err != nil {
	//		return err
	//	}
	//
	//	err = state.gatherUntil("sunflower", state.InventoryMaxItems)
	//	if err != nil {
	//		return err
	//	}
	//
	//	err = state.moveToAlchemy()
	//	if err != nil {
	//		return err
	//	}
	//
	//	err = state.craftUntil("small_health_potion", state.InventoryMaxItems/3)
	//	if err != nil {
	//		return err
	//	}
	//}
}
