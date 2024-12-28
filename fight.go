package main

import (
	"encoding/json"
	"math"
	"os"
)

func (state *CharacterState) fight() error {
	_, err := state.performActionAndWait("fight", []byte{})
	if err != nil {
		return err
	}
	return nil
}

func (state *CharacterState) rest() error {
	_, err := state.performActionAndWait("rest", []byte{})
	return err
}

func (state *CharacterState) findWorthyEnemy() string {
	maxLevel := state.Level - 1
	mostWorthy := ""
	highestLevel := 0

	for _, monster := range MonsterDB.Data {
		if monster.Level < maxLevel && monster.Level > highestLevel {
			mostWorthy = monster.Code
			highestLevel = monster.Level
		}
	}

	state.Logger.Debug(mostWorthy + " deemed worthy")

	return mostWorthy
}

func (state *CharacterState) goFightEnemy(enemyName string, healingItem string, healAmount int) error {
	location, err := state.getMonsterLocation(enemyName)
	if err != nil {
		return err
	}

	err = state.moveToLocation(location)
	if err != nil {
		return err
	}

	err = state.fightUntilLowInventory(healingItem, healAmount)
	if err != nil {
		return err
	}

	return nil
}

// heal to full
func (state *CharacterState) healToFull(healingItem string, amountHeal int) error {
	numHave := state.getItemInventoryQty(healingItem)
	hpToHeal := state.MaxHp - state.Hp
	numNeeded := int(math.Ceil(float64(hpToHeal) / float64(amountHeal)))

	if numNeeded <= 0 {
		return nil
	}

	err := state.useItem(healingItem, min(numHave, numNeeded))
	if err != nil {
		return err
	}

	return nil
}

func (state *CharacterState) goFightEnemyRest(enemyName string) error {
	location, err := state.getMonsterLocation(enemyName)
	if err != nil {
		return err
	}

	jsonData, err := json.Marshal(location)
	if err != nil {
		state.Logger.Error("Error marshalling request body", "error", err)
		os.Exit(1)
	}

	state.Logger.Debug("moving", "location", location)
	_, err = state.performActionAndWait("move", jsonData)
	if err != nil {
		return err
	}
	for {
		if state.getInventoryFull() {
			// no point in fighting bc we get no loot
			return nil
		}
		err = state.fight()
		if err != nil {
			return err
		}
		err = state.rest()
		if err != nil {
			return err
		}
	}
}

// fight until out of hp and healing item
func (state *CharacterState) fightUntilLowInventory(healingItem string, amountHeal int) error {
	numHealItem := state.getItemInventoryQty(healingItem)
	fightCount := 0
	state.Logger.Info("fight_forever",
		"healing_item", healingItem)

	for numHealItem > 0 {
		err := state.healToFull(healingItem, amountHeal)

		if err != nil {
			state.Logger.Error("Error healing", "error", err)
			return err
		}
		if state.getInventoryFull() {
			// no point in fighting bc we get no loot
			return nil
		}

		err = state.fight()
		if err != nil {
			return err
		}

		fightCount++
		numHealItem = state.getItemInventoryQty(healingItem)

		state.Logger.Debug("progress made",
			"action", "fighting",
			"fights_won", fightCount,
			"remaining", numHealItem)
	}
	return nil
}

func (state *CharacterState) fightGameLoop(monsterToFight string, food string, healAmount int) error {
	err := state.rest()
	if err != nil {
		return err
	}

	for {
		err = state.unequip("small_health_potion")
		if err != nil {
			return err
		}
		err = state.dumpAtBank()
		if err != nil {
			return err
		}

		err = state.withdrawItemAtBank("small_health_potion", 100)
		if err == nil {
			err = state.equipItem("small_health_potion", "utility1", 100)
		}

		if err != nil {
			state.Logger.Warn("Could not get/equip potions. Maybe we'll just die")
		}

		err = state.withdrawItemAtBank(food, 50)
		if err != nil {
			state.Logger.Warn("Could not withdraw food. We'll just have to rest")
		}

		if state.getItemInventoryQty(food) > 0 {
			err := state.goFightEnemy(monsterToFight, food, healAmount)
			if err != nil {
				return err
			}
		} else {
			err := state.goFightEnemyRest(monsterToFight)
			if err != nil {
				return err
			}
		}
	}
}
