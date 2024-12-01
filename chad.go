package main

func chadLoop(currentCharacter *CharacterState) error {
	for {
		dumpAtBank(currentCharacter)

		err := moveToSpruce(currentCharacter)
		if err != nil {
			return err
		}
		err = gatherUntil(currentCharacter, "spruce_wood", 100)
		if err != nil {
			return err
		}
	}
}
