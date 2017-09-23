package main

import (
	"fmt"
	// "strings"
)

import cards "github.com/nickaubert/dominionator/cards"
import players "github.com/nickaubert/dominionator/players"

func main() {

	fmt.Println("Dominion!")
	fmt.Println()

	pg := players.InitializePlaygroup(3)
	for n := range pg.Players {
		cards.ShuffleDeck(&pg.Players[n].Deck)
		players.Draw(&pg.Players[n], 5)
	}

	fmt.Println("starting supply:")
	for _, p := range pg.Supply.Piles {
		fmt.Println("pile", p.Count, p.Card.Name)
	}
	fmt.Println()

	// stopit := 0 // break for testing
	turnCount := 0
	for {

		turnCount++
		endGame := players.PlayTurn(&pg)

		if endGame == true {
			break
		}

		/*
			stopit++
			if stopit > 50 {
				break
			}
		*/
	}

	fmt.Println(turnCount, "turns")
	fmt.Println("ending supply:")
	for _, p := range pg.Supply.Piles {
		fmt.Println("pile", p.Count, p.Card.Name)
	}
	fmt.Println()

	players.CheckScores(pg)

}
