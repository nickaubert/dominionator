package players

import (
	"fmt"
	"math/rand"
	"time"
)

import cd "github.com/nickaubert/dominionator/cards"
import bs "github.com/nickaubert/dominionator/basic"

type Player struct {
	Deck    []cd.Card
	Hand    []cd.Card
	InPlay  []cd.Card
	Discard []cd.Card
	Name    string
}

type Playgroup struct {
	Players    []Player
	PlayerTurn int
	ThisTurn   ThisTurn
	Supply     cd.Supply
	Trash      []cd.Card
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
		// pl.Deck = InitialDeck()
		pl.Deck = bs.InitialDeck()
		pl.Name = fmt.Sprintf("Player%d", i)
		pg.Players = append(pg.Players, pl)
	}
	pg.PlayerTurn = 0
	pg.Supply = bs.InitializeSupply(s)
	return pg
}

func PlayTurn(pg *Playgroup) bool {

	fmt.Printf("%s's turn\n", pg.Players[pg.PlayerTurn].Name)

	pg.ThisTurn.Actions = 1
	pg.ThisTurn.Buys = 1
	pg.ThisTurn.Coins = 0

	showHand(pg.Players[pg.PlayerTurn].Hand)
	ActionPhase(pg)
	showHand(pg.Players[pg.PlayerTurn].Hand)

	BuyPhase(pg)
	CleanupPhase(pg)

	pg.PlayerTurn++ // advance play to next turn
	if pg.PlayerTurn >= len(pg.Players) {
		pg.PlayerTurn = 0
	}

	endGame := CheckEnd(pg.Supply)
	return endGame

}

func ActionPhase(pg *Playgroup) {

	p := &pg.Players[pg.PlayerTurn]
	fmt.Println("\t ActionPhase")

	// decision which action cards to play will go here

	// play action cards as long as we can starting with non-terminals
	for pg.ThisTurn.Actions > 0 {
		nt, tm := findActionCards(p.Hand)
		if len(nt) > 0 {
			c := highestCostCard(nt)
			playActionCard(pg, c)
			showAction(pg, c)
			continue
		}
		if len(tm) > 0 {
			c := highestCostCard(tm)
			playActionCard(pg, c)
			showAction(pg, c)
			continue
		}
		break
	}

}

func BuyPhase(pg *Playgroup) {

	p := &pg.Players[pg.PlayerTurn]
	fmt.Println("\t BuyPhase")

	// newer loop will look like action loop
	for pg.ThisTurn.Buys > 0 {

		tc := findTreasureCards(p.Hand)
		// actual decisions about which cards to play will go here
		playTreasureCards(pg, tc)

		fmt.Println("\t\t", pg.ThisTurn.Coins, "coins to spend")
		c := SelectCardBuy(pg.ThisTurn.Coins, pg.Supply)
		buyCard(p, &pg.Supply, c)
		pg.ThisTurn.Buys--

		/*
			    for i := 0; i < pg.ThisTurn.Buys; i++ {
				    // decision which card to buy here
				    fmt.Println("\t\t select buy", c.Name)
				    buyCard(p, &pg.Supply, c)
			    }
		*/

	}

	/*
		var tc []cd.Card
		tc, p.Hand = getTreasureCards(p.Hand)

		// decision whether to put each card into play will go here
		decide := true
		for _, c := range tc {
			if decide == true {
				// fmt.Println("\t\t play", c.Name)
				p.InPlay = append(p.InPlay, c)
			} else {
				p.Hand = append(p.Hand, c)
			}
		}

		for _, c := range tc {
			pg.ThisTurn.Coins += c.Coins
		}
	*/

}

func CleanupPhase(pg *Playgroup) {
	fmt.Println("\t CleanupPhase")
	p := pg.Players[pg.PlayerTurn]
	p.Discard = append(p.Discard, p.InPlay...)
	p.Discard = append(p.Discard, p.Hand...)
	p.Hand = p.Hand[:0]
	p.InPlay = p.InPlay[:0]
	Draw(&p, 5)
	pg.Players[pg.PlayerTurn] = p
}

func Draw(p *Player, d int) {
	for i := 0; i < d; i++ {
		c, z := p.Deck[0], p.Deck[1:]
		p.Deck = z
		p.Hand = append(p.Hand, c)
		if len(p.Deck) == 0 {
			p.Deck = p.Discard
			p.Discard = p.Discard[:0]
			ShuffleDeck(p)
		}
	}
}

/*
func getTreasureCards(hand []cd.Card) ([]cd.Card, []cd.Card) {
	var tc []cd.Card // treasure cards
	var oc []cd.Card // other cards
	for _, c := range hand {
		if c.CTypes.Treasure == true {
			tc = append(tc, c)
		} else {
			oc = append(oc, c)
		}
	}
	return tc, oc
}
*/

func SelectCardBuy(o int, s cd.Supply) cd.Card {
	// this is not the best heuristic
	var bestCard cd.Card
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

func buyCard(p *Player, s *cd.Supply, c cd.Card) {
	// assuming pile size is not zero
	for n, pl := range s.Piles {
		if pl.Card.Name == c.Name {
			s.Piles[n].Count--
			p.Discard = append(p.Discard, c)
		}
	}
}

func CheckEnd(s cd.Supply) bool {
	emptyPiles := 0
	for _, pl := range s.Piles {
		if pl.Card.Name == "Province" {
			if pl.Count == 0 {
				return true
			}
		}
		if pl.Count == 0 {
			emptyPiles++
		}
	}
	if emptyPiles >= 3 {
		return true
	}
	return false
}

func CheckScores(pg Playgroup) {
	for _, p := range pg.Players {
		p.Deck = append(p.Deck, p.Hand...)
		p.Deck = append(p.Deck, p.Discard...)
		p.Deck = append(p.Deck, p.InPlay...)
		vp := countVictoryPoints(p.Deck)
		fmt.Println(p.Name, vp, "points")
	}
}

func countVictoryPoints(d []cd.Card) int {
	vp := 0
	for _, c := range d {
		vp += c.VP
	}
	return vp
}

func ShuffleDeck(p *Player) {
	rand.Seed(time.Now().UnixNano())
	d := p.Deck
	var n []cd.Card
	for _ = range d {
		r := rand.Intn(len(d))
		n = append(n, d[r])
		d = append(d[:r], d[r+1:]...)
	}
	p.Deck = n
}

func findTreasureCards(h []cd.Card) []cd.Card {
	var tc []cd.Card
	for _, c := range h {
		if c.CTypes.Treasure == false {
			continue
		}
		tc = append(tc, c)
	}
	return tc
}

func findActionCards(ac []cd.Card) ([]cd.Card, []cd.Card) {
	var nt []cd.Card
	var tm []cd.Card
	for _, c := range ac {
		if c.CTypes.Action == false {
			continue
		}
		if c.Effects.ExtraActions > 0 {
			nt = append(nt, c)
		} else {
			tm = append(tm, c)
		}
	}
	return nt, tm
}

func highestCostCard(d []cd.Card) cd.Card {
	o := -1
	var h cd.Card
	for _, c := range d {
		if c.Cost > o {
			o = c.Cost
			h = c
		}
	}
	return h
}

func resolveEffects(pg *Playgroup, c cd.Card) {
	pg.ThisTurn.Actions += c.Effects.ExtraActions
	pg.ThisTurn.Buys += c.Effects.ExtraBuys
	pg.ThisTurn.Coins += c.Effects.ExtraCoins
	Draw(&pg.Players[pg.PlayerTurn], c.Effects.DrawCard)
}

func showHand(h []cd.Card) {
	fmt.Print("\t\t hand: ")
	for _, c := range h {
		fmt.Print(c.Name, ", ")
	}
	fmt.Print("\n")
}

func showAction(pg *Playgroup, c cd.Card) {
	fmt.Println("\t\t played", c.Name)
	fmt.Println("\t\t\t actions", pg.ThisTurn.Actions)
	fmt.Println("\t\t\t buys", pg.ThisTurn.Buys)
	fmt.Println("\t\t\t coins", pg.ThisTurn.Coins)
}

func playActionCard(pg *Playgroup, c cd.Card) {
	pg.ThisTurn.Actions--
	playCard(&pg.Players[pg.PlayerTurn], c)
	resolveEffects(pg, c)
}

func playTreasureCards(pg *Playgroup, tc []cd.Card) {

	for _, c := range tc {
		playCard(&pg.Players[pg.PlayerTurn], c)
		pg.ThisTurn.Coins += c.Coins
	}

}

func playCard(p *Player, c cd.Card) {
	p.InPlay = append(p.InPlay, c)
	for i, ch := range p.Hand {
		if ch.Name == c.Name {
			p.Hand = append(p.Hand[:i], p.Hand[i+1:]...)
			break
		}
	}
}
