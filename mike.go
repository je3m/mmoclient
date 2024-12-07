package main

func (state *CharacterState) mikeLoop() error {
	//for {
	//	state.dumpAtBank()
	//	err := state.goCraftItem("iron", state.InventoryMaxItems/8)
	//	if err != nil {
	//		return err
	//	}
	//}
	return state.fightGameLoop("chicken", "apple", 50)
	//return nil
}
