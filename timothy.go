package main

func timothyLoop(currentCharacter *CharacterState) error {
	for {
		dumpAtBank(currentCharacter)

		err := moveToShrimp(currentCharacter)
		if err != nil {
			currentCharacter.Logger.Warn("Failed to move to shrimp: %v\n", err)
			return err
		}
		err = gatherUntil(currentCharacter, "shrimp", 100)
		if err != nil {
			return err
		}
	}
}
