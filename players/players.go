package players

import (
	"fmt"
	// "strings"
)

import cards "github.com/nickaubert/dominionator/cards"
import basic "github.com/nickaubert/dominionator/basic"

type Player struct {
	Deck cards.Deck
	Name string
}

type Playgroup struct {
	Players []Player
}

func InitializePlaygroup(s int) Playgroup {
	var pg Playgroup
	for i := 0; i < s; i++ {
		var pl Player
		pl.Deck = InitialDeck()
		pl.Name = fmt.Sprintf("Player%2d", s)
		// shuffle here?
		pg.Players = append(pg.Players, pl)
	}
	return pg
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
