package main

import (
	"go-server/accumulate"
	"go-server/bonus"
	"go-server/stack"
)

var (
	// Application state - global variables
	// Ideally each of these should be managed by a dedicated micro-service
	jackpot    *accumulate.JackpotType
	luckyPairs *stack.BytePairStack

	// TODO Remove from global context; in presence of Web UI implement a bonus game as a separate endpoint with unique url known only to the eligible player
	// Since there is no UI we have to keep track of all players who are entiled to play a bonus game as a part of the server state
	bonusGames *bonus.BonusRegistry
)

func Init() {
	jackpot = accumulate.New()
	luckyPairs = stack.New()
	bonusGames = bonus.New()

	luckyPairs.FillUp()
}
