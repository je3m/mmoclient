package main

func lilyLoop(currentCharacter *CharacterState) error {
	for {
		dumpAtBank(currentCharacter)

		err := withdrawItemAtBank(currentCharacter, "cooked_gudgeon", 30)
		if err != nil {
			return err
		}

		err = moveToChicken(currentCharacter)
		if err != nil {
			return err
		}
		err = fight(currentCharacter)
		err = craftUntil(currentCharacter, "small_health_potion", 30)
		if err != nil {
			return err
		}
	}
}
