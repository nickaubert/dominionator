package main

import (
	"fmt"
)

import players "github.com/nickaubert/dominionator/players"

func main() {

	fmt.Println("Dominion!")
	fmt.Println()

	pg := players.InitializePlaygroup(3)
	for n := range pg.Players {
		players.ShuffleDeck(&pg.Players[n])
		players.Draw(&pg.Players[n], 5)
	}

	fmt.Println("starting supply:")
	for _, p := range pg.Supply.Piles {
		fmt.Println("pile", p.Count, p.Card.Name)
	}
	fmt.Println()

	turnCount := 0
	for {

		turnCount++
		endGame := players.PlayTurn(&pg)

		if endGame == true {
			break
		}

		if turnCount > 200 {
			break
		}

	}

	fmt.Println(turnCount, "turns")
	fmt.Println("ending supply:")
	for _, p := range pg.Supply.Piles {
		fmt.Println("pile", p.Count, p.Card.Name)
	}
	fmt.Println()

	players.CheckScores(pg)

}
