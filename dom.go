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

	for _, p := range pg.Supply.Piles {
		fmt.Println("pile", p.Count, p.Card.Name)
	}
	fmt.Println()

	stopit := 0 // break for testing
	for {

		players.PlayTurn(&pg)
		fmt.Println()

		stopit++
		if stopit > 10 {
			break
		}
	}

}
