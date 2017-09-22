package players

import (
	"fmt"
)

import cards "github.com/nickaubert/dominionator/cards"
import basic "github.com/nickaubert/dominionator/basic"

type Player struct {
	Deck        cards.Deck
	Hand        cards.Hand
	DiscardPile cards.DiscardPile
	Name        string
}

type Playgroup struct {
	Players    []Player
	PlayerTurn int
	Supply     cards.Supply
	Trash      cards.Trash
}

func InitializePlaygroup(s int) Playgroup {
	var pg Playgroup
	for i := 0; i < s; i++ {
		var pl Player
		pl.Deck = InitialDeck()
		pl.Name = fmt.Sprintf("Player%2d", i)
		pg.Players = append(pg.Players, pl)
	}
	pg.PlayerTurn = 0
	pg.Supply = InitializeSupply(pg)
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

func InitializeSupply(pg Playgroup) cards.Supply {

	var s cards.Supply
	var sp cards.SupplyPile

	sp.Card = basic.DefCopper()
	sp.Count = 10
	s.Piles = append(s.Piles, sp)

	sp.Card = basic.DefSilver()
	sp.Count = 10
	s.Piles = append(s.Piles, sp)

	sp.Card = basic.DefGold()
	sp.Count = 10
	s.Piles = append(s.Piles, sp)

	sp.Card = basic.DefEstate()
	sp.Count = 10
	s.Piles = append(s.Piles, sp)

	sp.Card = basic.DefDuchy()
	sp.Count = 10
	s.Piles = append(s.Piles, sp)

	sp.Card = basic.DefProvince()
	sp.Count = 10
	s.Piles = append(s.Piles, sp)

	sp.Card = basic.DefCurse()
	sp.Count = 10
	s.Piles = append(s.Piles, sp)

	return s
}

/*
func InitializeSupplyPile( c cards.Card, n int ) cards.SupplyPile {
    var sp cards.SupplyPile
    // p := make( c, n ) // ???
    for i := 0; i < n; i++ {
        sp.Card
    }
    return sp
}
*/

func PlayTurn(pg Playgroup) Playgroup {

	fmt.Printf("%s's turn\n", pg.Players[pg.PlayerTurn].Name)

	pg = ActionPhase(pg)
	pg = BuyPhase(pg)
	pg = CleanupPhase(pg)

	pg.PlayerTurn++ // advance play to next turn
	if pg.PlayerTurn >= len(pg.Players) {
		pg.PlayerTurn = 0
	}

	return pg
}

func ActionPhase(pg Playgroup) Playgroup {
	fmt.Printf("\t%s's turn ActionPhase\n", pg.Players[pg.PlayerTurn].Name)
	return pg
}

func BuyPhase(pg Playgroup) Playgroup {
	fmt.Printf("\t%s's turn BuyPhase\n", pg.Players[pg.PlayerTurn].Name)
	return pg
}

func CleanupPhase(pg Playgroup) Playgroup {
	fmt.Printf("\t%s's turn CleanupPhase\n", pg.Players[pg.PlayerTurn].Name)
	return pg
}

func Draw(p Player, d int) Player {
	for i := 0; i < d; i++ {
		c, z := p.Deck.Cards[0], p.Deck.Cards[1:]
		p.Deck.Cards = z
		p.Hand.Cards = append(p.Hand.Cards, c)
	}
	return p
}
