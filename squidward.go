package main

import "os"

func squidwardLoop(currentCharacter *CharacterState) error {
	for {
		dumpAtBank(currentCharacter)

		err := moveToIronMine(currentCharacter)
		if err != nil {
			os.Exit(1)
		}

		err = gatherUntil(currentCharacter, "iron_ore", 100)
		if err != nil {
			return err
		}
	}
}
