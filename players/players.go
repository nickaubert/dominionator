package players

import (
	"fmt"
)

import cards "github.com/nickaubert/dominionator/cards"
// import basic "github.com/nickaubert/dominionator/basic"

type Player struct {
	Deck cards.Deck
	Name string
}

type Playgroup struct {
	Players     []Player
    PlayerTurn  int
}

func InitializePlaygroup(s int) Playgroup {
	var pg Playgroup
	for i := 0; i < s; i++ {
		var pl Player
		pl.Deck = cards.InitialDeck()
		pl.Name = fmt.Sprintf("Player%2d", s)
		// shuffle here?
		pg.Players = append(pg.Players, pl)
	}
    pg.PlayerTurn = 0
	return pg
}

