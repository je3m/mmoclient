package main

func (state *CharacterState) mikeLoop() error {
	state.goEquipLoadout(
	{
		"adventurer_helmet",
			"adventurer_vest",
			"adventurer_pants",
			"iron_dagger",

	})
	return state.fightGameLoop("gingerbread", "cooked_beef", 150)
}
