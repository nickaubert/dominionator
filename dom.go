package main

import (
	"fmt"
	// "strings"
)

import cards "github.com/nickaubert/dominionator/cards"
import basic "github.com/nickaubert/dominionator/basic"

func main() {

	fmt.Println("Dominion!")

	/*
	   card1 := basic.DefVillage()

	   fmt.Println("cost:", card1.Cost)
	   fmt.Println("coins:", card1.Coins)
	   fmt.Println("vp:", card1.VP)
	   fmt.Println("action:", card1.CTypes.Action)

	   if card1.CTypes.Action == true {
	       fmt.Println("it is an action!")
	   }

	   if card1.CTypes.Curse == true {
	       fmt.Println("it is a curse!")
	   }
	*/

	startDeck := InitialDeck()
	playerDeck := cards.ShuffleDeck(startDeck) // would be more efficient to pass pointer?

	for n, c := range playerDeck.Cards {
		fmt.Println("card", n+1, c.Name)
	}

}

func InitialDeck() cards.Deck {
	var d cards.Deck
	s := make([]cards.Card, 0)
	for i := 0; i < 7; i++ {
		c := basic.DefCopper()
		s = append(s, c)
	}
	for i := 0; i < 3; i++ {
		c := basic.DefEstate()
		s = append(s, c)
	}
	d.Cards = s
	return d
}
