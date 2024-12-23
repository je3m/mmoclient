package main

func (state *CharacterState) farmSpruce() error {
	for {
		err := state.dumpAtBank()
		if err != nil {
			return err
		}
		err = state.moveToSpruce()
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
}

func (state *CharacterState) chadLoop() error {
	return state.fightGameLoop("gingerbread", "cooked_shrimp", 150)
}
