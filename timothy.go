package main

func (state *CharacterState) timothyLoop() error {
	for {
		state.dumpAtBank()

		err := state.moveToTrout()
		if err != nil {
			return err
		}

		err = state.gatherUntil("trout", 100)
		if err != nil {
			return err
		}
	}
}
