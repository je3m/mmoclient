package main

func (state *CharacterState) farmSpruce() error {
	for {
		state.dumpAtBank()

		err := state.moveToSpruce()
		if err != nil {
			return err
		}
		err = state.gatherUntil("spruce_wood", 100)
		if err != nil {
			return err
		}

		err = state.moveToWoodcraftStation()
		if err != nil {
			return err
		}

		err = state.craftUntil("spruce_plank", 100/8)
		if err != nil {
			return err
		}
	}
	return nil
}

func (state *CharacterState) chadLoop() error {
	return state.fightGameLoop("red_slime", "cooked_shrimp", 150)
	//return state.farmSpruce()
}
