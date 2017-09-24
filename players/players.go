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
			showStatus(pg, c)
			continue
		}
		if len(tm) > 0 {
			c := highestCostCard(tm)
			playActionCard(pg, c)
			showStatus(pg, c)
			continue
		}
		break
	}

}

func BuyPhase(pg *Playgroup) {

	p := &pg.Players[pg.PlayerTurn]
	fmt.Println("\t BuyPhase")

	for pg.ThisTurn.Buys > 0 {

		// actual decisions about which cards to play will go here
		tc := findTreasureCards(p.Hand)
		playTreasureCards(pg, tc)

		fmt.Println("\t\t", pg.ThisTurn.Coins, "coins to spend")
		c := SelectCardBuy(pg.ThisTurn.Coins, pg.Supply)
		buyCard(p, &pg.Supply, c)
		pg.ThisTurn.Buys--

	}

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

func findVictoryCards(h []cd.Card) []cd.Card {
	// finds victory cards that arent also treasure or action cards
	var vc []cd.Card
	for _, c := range h {
		if c.CTypes.Victory == false {
			continue
		}
		if c.CTypes.Treasure == true {
			continue
		}
		if c.CTypes.Action == true {
			continue
		}
		vc = append(vc, c)
	}
	return vc
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

func lowestCostCard(d []cd.Card) cd.Card {
	o := 100
	var l cd.Card
	for _, c := range d {
		if c.Cost < o {
			o = c.Cost
			l = c
		}
	}
	return l
}

func resolveEffects(pg *Playgroup, c cd.Card) {
	pg.ThisTurn.Actions += c.Effects.ExtraActions
	pg.ThisTurn.Buys += c.Effects.ExtraBuys
	pg.ThisTurn.Coins += c.Effects.ExtraCoins
	Draw(&pg.Players[pg.PlayerTurn], c.Effects.DrawCard)
}

func resolveAttacks(pg *Playgroup, c cd.Card) {
	for i := range pg.Players {
		if i == pg.PlayerTurn {
			continue
		}
		applyAttack(&pg.Players[i], c)
	}
}

func applyAttack(p *Player, c cd.Card) {
	fmt.Println("\t\tAttacking", p.Name)
	if c.Attacks.DiscardTo > 0 {
		discardTo(p, 3)
	}
}

func showHand(h []cd.Card) {
	fmt.Print("\t\t hand: ")
	for _, c := range h {
		fmt.Print(c.Name, ", ")
	}
	fmt.Print("\n")
}

func showStatus(pg *Playgroup, c cd.Card) {
	fmt.Println("\t\t\t actions", pg.ThisTurn.Actions)
	fmt.Println("\t\t\t buys", pg.ThisTurn.Buys)
	fmt.Println("\t\t\t coins", pg.ThisTurn.Coins)
}

func playActionCard(pg *Playgroup, c cd.Card) {
	fmt.Println("\t\t playing", c.Name)
	pg.ThisTurn.Actions--
	playCard(&pg.Players[pg.PlayerTurn], c)
	resolveEffects(pg, c)
	if c.CTypes.Attack == true {
		resolveAttacks(pg, c)
	}
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

func discardCard(p *Player, c cd.Card) {
	p.Discard = append(p.Discard, c)
	for i, ch := range p.Hand {
		if ch.Name == c.Name {
			p.Hand = append(p.Hand[:i], p.Hand[i+1:]...)
			break
		}
	}
}

func selectDiscardOwn(p *Player) cd.Card {
	// discard victory cards first, then select lowest value
	vc := findVictoryCards(p.Hand)
	for _, c := range vc {
		return c
	}
	c := lowestCostCard(p.Hand)
	return c
}

func discardTo(p *Player, m int) {
	for len(p.Hand) > m {
		d := selectDiscardOwn(p)
		discardCard(p, d)
	}
}
