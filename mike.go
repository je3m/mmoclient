package main

func (state *CharacterState) mikeLoop() error {
	for {
		state.dumpAtBank()

		err := state.withdrawItemAtBank("shrimp", 100)
		if err != nil {
			return err
		}

		err = state.moveToCooking()
		if err != nil {
			state.Logger.Error("Failed to move to kitchen", "error", err)
			return err
		}
		err = state.craftUntil("cooked_shrimp", 100)
		if err != nil {
			return err
		}
	}
}
