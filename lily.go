package main

func lilyLoop(state *CharacterState) error {

	for {
		dumpAtBank(state)

		err := withdrawItemAtBank(state, "cooked_gudgeon", 30)
		if err != nil {
			return err
		}

		err = moveToChicken(state)
		if err != nil {
			return err
		}
		err = fightUntilLowInventory(state, "cooked_gudgeon", 75)
		if err != nil {
			return err
		}
	}
}
