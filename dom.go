package main

import (
	"fmt"
)

import pl "github.com/nickaubert/dominionator/players"

func main() {

	fmt.Println("Dominion!")
	fmt.Println()

	pg := pl.InitializePlaygroup(3)
	for n := range pg.Players {
		pl.ShuffleDeck(&pg.Players[n])
		pl.Draw(&pg.Players[n], 5)
	}

	fmt.Println("starting supply:")
	for _, p := range pg.Supply.Piles {
		fmt.Println("pile", p.Count, p.Card.Name)
	}
	fmt.Println()

	turnCount := 0
	for {

		turnCount++
		endGame := pl.PlayTurn(&pg)

		if endGame == true {
			break
		}

		if turnCount > 200 {
			fmt.Println("Interrupted game at turn 201")
			break
		}

	}

	fmt.Println(turnCount, "turns")
	fmt.Println("ending supply:")
	for _, p := range pg.Supply.Piles {
		fmt.Println("pile", p.Count, p.Card.Name)
	}
	fmt.Println()

	pl.CheckScores(pg)

}
