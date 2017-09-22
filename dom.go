package main

import (
	"fmt"
	// "strings"
)

import cards "github.com/nickaubert/dominionator/cards"
import players "github.com/nickaubert/dominionator/players"

func main() {

	fmt.Println("Dominion!")

	pg := players.InitializePlaygroup(3)
	for _, p := range pg.Players {
		// fmt.Println(p.Name)
		p.Deck = cards.ShuffleDeck(p.Deck)
		/*
			for n, c := range p.Deck.Cards {
				fmt.Println("card", n+1, c.Name)
			}
		*/
	}
	// fmt.Println()

	// sp := players.InitializeSupply(pg)
	for _, p := range pg.Supply.Piles {
		fmt.Println("pile", p.Count, p.Card.Name)
	}
	fmt.Println()

	stopit := 0 // break for testing
	for {

		pg = players.PlayTurn(pg)

		stopit++
		if stopit > 10 {
			break
		}
	}

}
