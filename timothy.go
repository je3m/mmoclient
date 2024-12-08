package main

func (state *CharacterState) timothyLoop() error {
	return state.fishForShrimpLoop()
}

func (state *CharacterState) fishForShrimpLoop() error {
	for {
		err := state.dumpAtBank()
		if err != nil {
			return err
		}
		err = state.moveToShrimp()
		if err != nil {
			return err
		}
		err = state.gatherUntil("shrimp", state.InventoryMaxItems)
		if err != nil {
			return err
		}
	}
}
