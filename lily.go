package main

func (state *CharacterState) lilyLoop() error {
	return state.fightGameLoop("frost_slime", "cooked_shrimp", 150)
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

func (state *CharacterState) grindTo20Jewel(err error) error {
	lvl15Items := []string{"air_ring", "earth_ring", "fire_ring", "life_ring", "water_ring"}
	for {
		for _, item := range lvl15Items {
			err = state.dumpAtBank()
			if err != nil {
				return err
			}
			err = state.goCraftItem(item, 1)
			if err != nil {
				break
			}
			if state.JewelrycraftingLevel >= 20 {
				break
			}
		}
		if err != nil {
			break
		}
		if state.JewelrycraftingLevel >= 20 {
			break
		}
	}
	return nil
}

func (state *CharacterState) grindTo15Jewel(err error) error {
	lvl10Items := []string{"fire_and_earth_amulet", "air_and_water_amulet", "iron_ring"}

	for {
		for _, item := range lvl10Items {
			err = state.dumpAtBank()
			if err != nil {
				return err
			}
			state.Logger.Info("going to craft", "item", item)
			err = state.goCraftItem(item, 1)
			if err != nil {
				break
			}
			if state.JewelrycraftingLevel >= 15 {
				break
			}
		}
		if err != nil {
			break
		}
		if state.JewelrycraftingLevel >= 15 {
			break
		}

	}
	return nil
}
