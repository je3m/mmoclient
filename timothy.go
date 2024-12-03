package main

func (state *CharacterState) timothyLoop() error {
	for {
		state.dumpAtBank()

		//err := state.withdrawItemAtBank("apple", 30)
		//if err != nil {
		//	return err
		//}
		err := state.goFightEnemyRest("chicken", "apple", 50)

		if err != nil {
			state.Logger.Error("something went wrong: restarting")
		}
	}
}
