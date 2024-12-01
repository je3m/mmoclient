package main

func mikeLoop(currentCharacter *CharacterState) error {
	for {
		dumpAtBank(currentCharacter)

		err := withdrawItemAtBank(currentCharacter, "shrimp", 100)
		if err != nil {
			return err
		}

		err = moveToCooking(currentCharacter)
		if err != nil {
			currentCharacter.Logger.Error("Failed to move to kitchen", "error", err)
			return err
		}
		err = craftUntil(currentCharacter, "cooked_shrimp", 100)
		if err != nil {
			return err
		}
	}
}
