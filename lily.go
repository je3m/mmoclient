package main

func (state *CharacterState) lilyLoop() error {

	for {
		state.dumpAtBank()

		err := state.withdrawItemAtBank("cooked_gudgeon", 30)
		if err != nil {
			return err
		}

		err = fightWorthyEnemy(state, "cooked_gudgeon", 75)
		if err != nil {
			return err
		}
	}
}
