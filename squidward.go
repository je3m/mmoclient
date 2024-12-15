package main

func (state *CharacterState) squidwardLoop() error {
	return state.fightGameLoop("gingerbread", "apple", 50)
	//return state.gatherGameLoop("coal")
}
