package main

func (state *CharacterState) squidwardLoop() error {
	for {
		state.dumpAtBank()

		err := state.moveToIronMine()
		if err != nil {
			return err
		}

		err = state.gatherUntil("iron_ore", 100)
		if err != nil {
			return err
		}

		err = state.moveToMiningStation()
		if err != nil {
			return err
		}

		err = state.craftUntil("iron", 100/8)
		if err != nil {
			return err
		}

	}
}
