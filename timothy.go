package main

func (state *CharacterState) timothyLoop() error {
	for {
		state.dumpAtBank()

		err := state.goFightEnemyRest("chicken", "apple", 50)

		if err != nil {
			state.Logger.Error("something went wrong: restarting")
		}
	}
}
