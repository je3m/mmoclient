package main

func (state *CharacterState) lilyLoop() error {

	for {
		state.dumpAtBank()

		err := state.withdrawItemAtBank("cooked_gudgeon", 30)
		if err != nil {
			return err
		}

		err = state.moveToChicken()
		if err != nil {
			return err
		}
		err = state.fightUntilLowInventory("cooked_gudgeon", 75)
		if err != nil {
			return err
		}
	}
}
