package main

func (state *CharacterState) mikeLoop() error {
	//for i := 0; i < 4; i++ {
	//	err := state.dumpAtBank()
	//	if err != nil {
	//		return err
	//	}
	//	err = state.goCraftItem("iron_boots", 1)
	//	if err != nil {
	//		return err
	//	}
	//	err = state.goCraftItem("iron_helm", 1)
	//	if err != nil {
	//		return err
	//	}
	//	err = state.goCraftItem("slime_shield", 1)
	//	if err != nil {
	//		return err
	//	}
	err := state.dumpAtBank()
	//err = state.goCraftItem("iron_legs_armor", 2)
	if err != nil {
		return err
	}
	//err = state.goCraftItem("iron_armor", 2)
	//if err != nil {
	//	return err
	//}

	//}
	for {
		err := state.dumpAtBank()
		if err != nil {
			break
		}
		err = state.goCraftItem("small_health_potion", state.InventoryMaxItems/3)
		if err != nil {
			break
		}
	}
	state.Logger.Info("moving to stage 2")
	for {
		err := state.dumpAtBank()
		if err != nil {
			return err
		}
		err = state.moveToSunflower()
		if err != nil {
			return err
		}

		err = state.gatherUntil("sunflower", state.InventoryMaxItems)
		if err != nil {
			return err
		}
		if state.AlchemyLevel >= 5 {
			err = state.goCraftItem("small_health_potion", state.InventoryMaxItems/3)
			if err != nil {
				return err
			}
		}
	}

	//return state.fightGameLoop("chicken", "apple", 50)
}
