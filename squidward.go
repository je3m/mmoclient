package main

func (state *CharacterState) squidwardLoop() error {
	return state.fightGameLoop("blue_slime", "apple", 50)
	//return state.gatherGameLoop("coal")
}
