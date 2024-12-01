package main

func (state *CharacterState) mikeLoop() error {
	for {
		state.dumpAtBank()

		err := withdrawItemAtBank(currentCharacter, "ash_wood", 100)
		if err != nil {
			return err
		}

		err = moveWoodCraftStation(currentCharacter)
		if err != nil {
			currentCharacter.Logger.Error("Failed to move to woodcraft", "error", err)
			return err
		}
		err = craftUntil(currentCharacter, "ash_plank", 100/8)
		if err != nil {
			return err
		}
	}
}
