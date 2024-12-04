package main

func (state *CharacterState) timothyLoop() error {
	for {
		state.dumpAtBank()

		err := state.goFightEnemyRest("yellow_slime")

		if err != nil {
			return err
		}
		return nil
	}
}
