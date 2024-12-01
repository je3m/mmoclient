package main

func (state *CharacterState) chadLoop() error {
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
}
