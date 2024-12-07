package main

func (state *CharacterState) lilyLoop() error {
	for {
		state.dumpAtBank()

		err := state.moveToSunflower()
		if err != nil {
			return err
		}

		err = state.gatherUntil("sunflower", state.InventoryMaxItems)
		if err != nil {
			return err
		}

		err = state.moveToAlchemy()
		if err != nil {
			return err
		}

		err = state.craftUntil("small_health_potion", state.InventoryMaxItems/3)
		if err != nil {
			return err
		}
	}
}
