package main

func (state *CharacterState) mikeLoop() error {

	for {
		err := state.dumpAtBank()
		if err != nil {
			return err
		}

		err = state.moveToCopperMine()
		if err != nil {
			return err
		}

		err = state.gatherUntil("copper_ore", 100)
		if err != nil {
			return err
		}

		err = state.moveToMiningStation()
		if err != nil {
			return err
		}

		err = state.craftUntil("copper", 100/8)
	}
}
