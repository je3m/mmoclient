package main

func (state *CharacterState) lilyLoop() error {

	for {
		state.dumpAtBank()

		err := state.withdrawItemAtBank("cooked_shrimp", 30)
		if err != nil {
			return err
		}
		err = state.fightWorthyEnemy("cooked_shrimp", 150)

		if err != nil {
			return err
		}
	}
}
