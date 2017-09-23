package players

import (
	"fmt"
)

import cards "github.com/nickaubert/dominionator/cards"
import basic "github.com/nickaubert/dominionator/basic"

type Player struct {
	Deck        cards.Deck
	Hand        cards.Hand
	InPlay      cards.InPlay
	DiscardPile cards.DiscardPile
	Name        string
}

type Playgroup struct {
	Players    []Player
	PlayerTurn int
	ThisTurn   ThisTurn
	Supply     cards.Supply
	Trash      cards.Trash
}

type ThisTurn struct {
	Actions int
	Coins   int
	Buys    int
}

func InitializePlaygroup(s int) Playgroup {
	var pg Playgroup
	for i := 0; i < s; i++ {
		var pl Player
		pl.Deck = InitialDeck()
		pl.Name = fmt.Sprintf("Player%d", i)
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

	/* coin cards */
	sp.Card = basic.DefCopper()
	sp.Count = 60 - (len(pg.Players) * 7)
	s.Piles = append(s.Piles, sp)

	sp.Card = basic.DefSilver()
	sp.Count = 40
	s.Piles = append(s.Piles, sp)

	sp.Card = basic.DefGold()
	sp.Count = 30
	s.Piles = append(s.Piles, sp)

	/* victory cards */
	sp.Count = 12
	if len(pg.Players) == 2 {
		sp.Count = 8
	}
	sp.Card = basic.DefEstate()
	s.Piles = append(s.Piles, sp)
	sp.Card = basic.DefDuchy()
	s.Piles = append(s.Piles, sp)

	if len(pg.Players) == 5 {
		sp.Count = 15
	}
	if len(pg.Players) == 6 {
		sp.Count = 18
	}
	sp.Card = basic.DefProvince()
	s.Piles = append(s.Piles, sp)

	/* curses */
	sp.Card = basic.DefCurse()
	sp.Count = 10 * (len(pg.Players) - 1)
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

func PlayTurn(pg *Playgroup) {

	fmt.Printf("%s's turn\n", pg.Players[pg.PlayerTurn].Name)

	pg.ThisTurn.Actions = 1
	pg.ThisTurn.Buys = 1
	pg.ThisTurn.Coins = 0

	/*
	   fmt.Println("\t\thand:")
	   for _, c := range pg.Players[pg.PlayerTurn].Hand.Cards {
	       fmt.Println("\t\t", c.Name)
	   }
	*/

	ActionPhase(pg)
	BuyPhase(pg)
	CleanupPhase(pg)

	pg.PlayerTurn++ // advance play to next turn
	if pg.PlayerTurn >= len(pg.Players) {
		pg.PlayerTurn = 0
	}

}

func ActionPhase(pg *Playgroup) {
	fmt.Println("\t ActionPhase")
	ac := getActionCards(pg.Players[pg.PlayerTurn].Hand)
	for _, c := range ac {
		fmt.Println("\t\t action card", c.Name)
	}
}

func BuyPhase(pg *Playgroup) {

	p := pg.Players[pg.PlayerTurn]
	fmt.Println("\t BuyPhase")

	var tc []cards.Card
	tc, p.Hand.Cards = getTreasureCards(p.Hand)

	// decision whether to put each card into play will go here
	decide := true

	for _, c := range tc {
		if decide == true {
			fmt.Println("\t\t play", c.Name)
			p.InPlay.Cards = append(p.InPlay.Cards, c)
		} else {
			p.Hand.Cards = append(p.Hand.Cards, c)
		}
	}

	for _, c := range tc {
		pg.ThisTurn.Coins += c.Coins
	}

	fmt.Println("\t\t", pg.ThisTurn.Coins, "coins to spend")

	for i := 0; i < pg.ThisTurn.Buys; i++ {
		// decision which card to buy here
		c := SelectCardBuy(pg.ThisTurn.Coins, pg.Supply)
		fmt.Println("\t\t select buy", c.Name)
	}

	pg.Players[pg.PlayerTurn] = p
}

func CleanupPhase(pg *Playgroup) {
	fmt.Println("\t CleanupPhase")
	// pg.ThisTurn.Actions = 0
	// pg.ThisTurn.Buys = 0
	// pg.ThisTurn.Coins = 0
	p := pg.Players[pg.PlayerTurn]
	p.DiscardPile.Cards = append(p.DiscardPile.Cards, p.InPlay.Cards...)
	p.DiscardPile.Cards = append(p.DiscardPile.Cards, p.Hand.Cards...)
	p.Hand.Cards = p.Hand.Cards[:0]
	p.InPlay.Cards = p.InPlay.Cards[:0]
	Draw(&p, 5)
	pg.Players[pg.PlayerTurn] = p
}

func Draw(p *Player, d int) {
	for i := 0; i < d; i++ {
		c, z := p.Deck.Cards[0], p.Deck.Cards[1:]
		p.Deck.Cards = z
		p.Hand.Cards = append(p.Hand.Cards, c)
		if len(p.Deck.Cards) == 0 {
			p.Deck.Cards = p.DiscardPile.Cards
			p.DiscardPile.Cards = p.DiscardPile.Cards[:0]
			cards.ShuffleDeck(&p.Deck)
		}
	}
}

func getActionCards(hand cards.Hand) []cards.Card {
	var ac []cards.Card
	for _, c := range hand.Cards {
		if c.CTypes.Action == true {
			ac = append(ac, c)
		}
	}
	return ac
}

func getTreasureCards(hand cards.Hand) ([]cards.Card, []cards.Card) {
	var tc []cards.Card // treasure cards
	var oc []cards.Card // other cards
	for _, c := range hand.Cards {
		if c.CTypes.Treasure == true {
			tc = append(tc, c)
		} else {
			oc = append(oc, c)
		}
	}
	return tc, oc
}

func SelectCardBuy(o int, s cards.Supply) cards.Card {
	// this is not the best heuristic
	var bestCard cards.Card
	for _, c := range s.Piles {
		if c.Count == 0 {
			continue
		}
		if c.Card.Cost > o {
			continue
		}
		if c.Card.Cost > bestCard.Cost {
			bestCard = c.Card
		}
	}
	return bestCard
}
