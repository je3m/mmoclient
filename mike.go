package main

func (state *CharacterState) mikeLoop() error {
	return state.fightGameLoop("gingerbread", "cooked_beef", 150)
}
