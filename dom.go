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
	for i, p := range pg.Players {
		p.Deck = cards.ShuffleDeck(p.Deck)
		p = players.Draw(p, 5)
		pg.Players[i] = p
	}

	/*
		for _, p := range pg.Players {
	        fmt.Println(p.Name, "Hand")
	        for _, c := range p.Hand.Cards {
			    fmt.Printf("\t%s\n", c.Name)
	        }
	        fmt.Println(p.Name, "Deck")
	        for _, c := range p.Deck.Cards {
			    fmt.Printf("\t%s\n", c.Name)
	        }
	    }
	*/

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
