package main

func (state *CharacterState) timothyLoop() error {
	for {
		state.dumpAtBank()

		err := state.moveToShrimp()
		if err != nil {
			state.Logger.Warn("Failed to move to shrimp: %v\n", err)
			return err
		}
		err = state.gatherUntil("shrimp", 100)
		if err != nil {
			return err
		}
	}
}
