package players

import (
	"fmt"
	"math/rand"
	"time"
)

import cards "github.com/nickaubert/dominionator/cards"
import basic "github.com/nickaubert/dominionator/basic"

type Player struct {
	Deck    []cards.Card
	Hand    []cards.Card
	InPlay  []cards.Card
	Discard []cards.Card
	Name    string
}

type Playgroup struct {
	Players    []Player
	PlayerTurn int
	ThisTurn   ThisTurn
	Supply     cards.Supply
	Trash      []cards.Card
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

func InitialDeck() []cards.Card {
	var d []cards.Card
	for i := 0; i < 7; i++ {
		c := basic.DefCopper()
		d = append(d, c)
	}
	for i := 0; i < 3; i++ {
		c := basic.DefEstate()
		d = append(d, c)
	}
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

	/* kingdom */
	sp.Card = basic.DefVillage()
	sp.Count = 10
	s.Piles = append(s.Piles, sp)

	sp.Card = basic.DefWoodcutter()
	sp.Count = 10
	s.Piles = append(s.Piles, sp)

	sp.Card = basic.DefSmithy()
	sp.Count = 10
	s.Piles = append(s.Piles, sp)

	sp.Card = basic.DefFestival()
	sp.Count = 10
	s.Piles = append(s.Piles, sp)

	sp.Card = basic.DefLaboratory()
	sp.Count = 10
	s.Piles = append(s.Piles, sp)

	sp.Card = basic.DefMarket()
	sp.Count = 10
	s.Piles = append(s.Piles, sp)

	return s
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
	/*
		// ac := getActionCards(pg.Players[pg.PlayerTurn].Hand)
		for _, c := range ac {
			fmt.Println("\t\t action card", c.Name)
		}
	*/
	// decision which action cards to play will go here
	for pg.ThisTurn.Actions > 0 {
		nt, _ := findActionTerminals(p.Hand)
		if len(nt) == 0 {
			break
		}
		c := highestCostCard(nt)
		playAction(pg, c)
		fmt.Println("\t\t played", c.Name)
		fmt.Println("\t\t\t actions", pg.ThisTurn.Actions)
		fmt.Println("\t\t\t buys", pg.ThisTurn.Buys)
		fmt.Println("\t\t\t coins", pg.ThisTurn.Coins)
	}
	_, tm := findActionTerminals(p.Hand)
	if len(tm) > 0 {
		c := highestCostCard(tm)
		playAction(pg, c)
		fmt.Println("\t\t played terminal", c.Name)
		fmt.Println("\t\t\t actions", pg.ThisTurn.Actions)
		fmt.Println("\t\t\t buys", pg.ThisTurn.Buys)
		fmt.Println("\t\t\t coins", pg.ThisTurn.Coins)
	}
}

func BuyPhase(pg *Playgroup) {

	p := pg.Players[pg.PlayerTurn]
	fmt.Println("\t BuyPhase")

	var tc []cards.Card
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

	fmt.Println("\t\t", pg.ThisTurn.Coins, "coins to spend")

	for i := 0; i < pg.ThisTurn.Buys; i++ {
		// decision which card to buy here
		c := SelectCardBuy(pg.ThisTurn.Coins, pg.Supply)
		fmt.Println("\t\t select buy", c.Name)
		buyCard(&p, &pg.Supply, c)
	}

	pg.Players[pg.PlayerTurn] = p
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
func getActionCards(hand []cards.Card) []cards.Card {
	var ac []cards.Card
	for _, c := range hand {
		if c.CTypes.Action == true {
			ac = append(ac, c)
		}
	}
	return ac
}
*/

func getTreasureCards(hand []cards.Card) ([]cards.Card, []cards.Card) {
	var tc []cards.Card // treasure cards
	var oc []cards.Card // other cards
	for _, c := range hand {
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

func buyCard(p *Player, s *cards.Supply, c cards.Card) {
	// assuming pile size is not zero
	for n, pl := range s.Piles {
		if pl.Card.Name == c.Name {
			s.Piles[n].Count--
			p.Discard = append(p.Discard, c)
		}
	}
}

func CheckEnd(s cards.Supply) bool {
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

func countVictoryPoints(d []cards.Card) int {
	vp := 0
	for _, c := range d {
		vp += c.VP
	}
	return vp
}

func ShuffleDeck(p *Player) {
	rand.Seed(time.Now().UnixNano())
	d := p.Deck
	var n []cards.Card
	for _ = range d {
		r := rand.Intn(len(d))
		n = append(n, d[r])
		d = append(d[:r], d[r+1:]...)
	}
	p.Deck = n
}

func findActionTerminals(ac []cards.Card) ([]cards.Card, []cards.Card) {
	var nt []cards.Card
	var tm []cards.Card
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

func highestCostCard(d []cards.Card) cards.Card {
	o := -1
	var h cards.Card
	for _, c := range d {
		if c.Cost > o {
			o = c.Cost
			h = c
		}
	}
	return h
}

func playAction(pg *Playgroup, c cards.Card) {
	pg.ThisTurn.Actions--
	pg.Players[pg.PlayerTurn].InPlay = append(pg.Players[pg.PlayerTurn].InPlay, c)
	h := pg.Players[pg.PlayerTurn].Hand
	for i, cd := range h {
		if cd.Name == c.Name {
			pg.Players[pg.PlayerTurn].Hand = append(h[:i], h[i+1:]...)
			break
		}
	}
	resolveEffects(pg, c)
}

func resolveEffects(pg *Playgroup, c cards.Card) {
	pg.ThisTurn.Actions += c.Effects.ExtraActions
	pg.ThisTurn.Buys += c.Effects.ExtraBuys
	pg.ThisTurn.Coins += c.Effects.ExtraCoins
	Draw(&pg.Players[pg.PlayerTurn], c.Effects.DrawCard)
}

func showHand(h []cards.Card) {
	fmt.Print("\t\t hand: ")
	for _, c := range h {
		fmt.Print(c.Name, ", ")
	}
	fmt.Print("\n")
}
