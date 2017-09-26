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
	pg.Supply = InitializeSupply(s)
	return pg
}

func PlayTurn(pg *Playgroup) bool {

	fmt.Printf("%s's turn\n", pg.Players[pg.PlayerTurn].Name)

	pg.ThisTurn.Actions = 1
	pg.ThisTurn.Buys = 1
	pg.ThisTurn.Coins = 0

	ActionPhase(pg)

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
	showCards("hand", pg.Players[pg.PlayerTurn].Hand)

	for pg.ThisTurn.Buys > 0 {

		// actual decisions about which cards to play will go here
		tc := findTreasureCards(p.Hand)
		playTreasureCards(pg, tc)

		fmt.Println("\t\t", pg.ThisTurn.Coins, "coins to spend")
		c := SelectCardBuy(pg.ThisTurn.Coins, pg.Supply)
		fmt.Println("\t\t buying", c.Name)
		gainCard(p, &pg.Supply, c)
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
	nc := Draw(&p, 5)
	AddHand(&p, nc)
	pg.Players[pg.PlayerTurn] = p
}

func Draw(p *Player, d int) []cd.Card {
	var nc []cd.Card
	for i := 0; i < d; i++ {
		c, z := p.Deck[0], p.Deck[1:]
		// fmt.Println("\t\t\t drawing", c.Name)
		p.Deck = z
		nc = append(nc, c)
		if len(p.Deck) == 0 {
			p.Deck = p.Discard
			p.Discard = p.Discard[:0]
			ShuffleDeck(p)
		}
	}
	return nc
}

func AddHand(p *Player, nc []cd.Card) {
	if len(nc) > 0 {
		p.Hand = append(p.Hand, nc...)
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

func gainCard(p *Player, s *cd.Supply, c cd.Card) {
	for n, pl := range s.Piles {
		if pl.Card.Name == c.Name {
			if s.Piles[n].Count == 0 {
				return
			}
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

		/*
		   fmt.Println(p.Name, "deck:")
		   for _, c := range p.Deck {
		       fmt.Print(c.Name, ", ")
		   }
		   fmt.Print("\n")
		*/

		vp := countVictoryPoints(p.Deck)
		fmt.Println(p.Name, vp, "points")
	}
}

func countVictoryPoints(d []cd.Card) int {
	vp := 0
	for _, c := range d {
		vp += c.VP
		if c.Victories.CardsPerPoint > 0 {
			vp += (len(d) / c.Victories.CardsPerPoint)
		}
	}
	return vp
}

// replace with ShuffleCards ?
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

/*
// cant range over pointer
func ShuffleCards(d *[]cd.Card) {
	rand.Seed(time.Now().UnixNano())
	var n []cd.Card
    ud := *d
	for _ = range ud {
		r := rand.Intn(len(ud))
		n = append(n, ud[r])
		ud = append(ud[:r], ud[r+1:]...)
	}
	d = &n
}
*/

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

func siftActionCards(cs []cd.Card) ([]cd.Card, []cd.Card) {
	var ac []cd.Card
	var na []cd.Card
	for _, c := range cs {
		if c.CTypes.Action == true {
			ac = append(ac, c)
			continue
		}
		na = append(na, c)
	}
	return ac, na
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

func findCurses(h []cd.Card) []cd.Card {
	var cc []cd.Card
	for _, c := range h {
		if c.CTypes.Curse == false {
			continue
		}
		cc = append(cc, c)
	}
	return cc
}

func findReactions(h []cd.Card) []cd.Card {
	var rc []cd.Card
	for _, c := range h {
		if c.CTypes.Reaction == false {
			continue
		}
		rc = append(rc, c)
	}
	return rc
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
	nc := Draw(&pg.Players[pg.PlayerTurn], c.Effects.DrawCard)
	AddHand(&pg.Players[pg.PlayerTurn], nc)
	if c.CTypes.Attack == true {
		resolveAttacks(pg, c)
	}
	resolveSequence(pg, c)
}

func resolveSequence(pg *Playgroup, c cd.Card) {
	var cardSet []cd.Card
	countX := 0 // replace with len(cardSet)?
	p := &pg.Players[pg.PlayerTurn]
	for i, s := range c.Effects.Sequence {
		fmt.Println("\t\t\t Sequence", i)
		if s.CountDiscard > 0 {
			// decision point here
			vc := findVictoryCards(p.Hand)
			for j, v := range vc {
				if j > s.CountDiscard {
					break
				}
				// discardCard(p, v)
				discardCards(p, []cd.Card{v})
				countX++
			}
		}
		if s.DrawCount == true {
			nc := Draw(p, countX)
			AddHand(p, nc)
		}
		if s.CountTrash > 0 {
			// decision point here
			cc := findCurses(p.Hand)
			for j, u := range cc {
				if j > s.CountTrash {
					break
				}
				removeFromHand(p, u)
				trashFromHand(p, pg, u)
				countX++
			}
		}
		if s.RetrieveDiscard > 0 {
			for j := 0; j < s.RetrieveDiscard; j++ {
				// decision point here
				bc := bestPlayableCard(p.Discard)
				// must have found something
				if bc.Name != "" {
					removeFromDiscard(p, bc)
					cardSet = append(cardSet, bc)
				}
			}
		}
		if s.PlaceDeck == true {
			addDeckTop(p, cardSet)
		}
		if s.DrawDeck > 0 {
			// fmt.Println("seq draw", s.DrawDeck)
			cardSet = append(cardSet, Draw(p, s.DrawDeck)...)
			// showCards("hand", p.Hand)
			// showCards("cardSet", cardSet)
		}
		if s.DiscardNonAction == true {
			ac, na := siftActionCards(cardSet)
			// fmt.Println("seq keep", len(ac), "discard", len(na))
			discardCards(p, na)
			cardSet = ac
		}
		if s.PlayAction > 0 {
			// possible decision point here whether to play any action
			for _, c := range cardSet {
				fmt.Println("\t\t\t seq play action", c.Name)
				p.InPlay = append(p.InPlay, c)
				for j := 0; j < s.PlayAction; j++ {
					resolveEffects(pg, c)
				}
			}
		}
	}
}

func resolveAttacks(pg *Playgroup, c cd.Card) {
	for i := range pg.Players {
		if i == pg.PlayerTurn {
			continue
		}
		fmt.Println("\t\t Attacking", pg.Players[i].Name)
		defended := checkReactions(&pg.Players[i])
		if defended == true {
			fmt.Println("\t\t defended!")
			continue
		}
		if c.Attacks.DiscardTo > 0 {
			discardTo(&pg.Players[i], c.Attacks.DiscardTo)
		}
		if c.Attacks.GainCurse > 0 {
			gainCurse(&pg.Players[i], &pg.Supply, c.Attacks.GainCurse)
		}
	}
}

func showCards(s string, h []cd.Card) {
	fmt.Print("\t\t", s, ": ")
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
	showCards("hand", pg.Players[pg.PlayerTurn].Hand)
	fmt.Println("\t\t playing", c.Name)
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
	removeFromHand(p, c)
	p.InPlay = append(p.InPlay, c)
}

func discardCards(p *Player, cs []cd.Card) {
	fmt.Print("\t\t\t discarding ")
	for _, c := range cs {
		removeFromHand(p, c)
		fmt.Print(c.Name, ", ")
		p.Discard = append(p.Discard, c)
	}
	fmt.Print("\n")
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
		// discardCards(p, d)
		discardCards(p, []cd.Card{d})
	}
}

func gainCurse(p *Player, s *cd.Supply, m int) {
	for i := 0; i < m; i++ {
		gainCard(p, s, bs.DefCurse())
	}
}

func removeFromHand(p *Player, c cd.Card) {
	for i, h := range p.Hand {
		if h.Name == c.Name {
			p.Hand = append(p.Hand[:i], p.Hand[i+1:]...)
			break
		}
	}
}

func removeFromDiscard(p *Player, c cd.Card) {
	for i, h := range p.Discard {
		if h.Name == c.Name {
			p.Discard = append(p.Discard[:i], p.Discard[i+1:]...)
			break
		}
	}
}

func trashFromHand(p *Player, pg *Playgroup, c cd.Card) {
	fmt.Println("\t\t\t trashing", c.Name)
	removeFromHand(p, c)
	pg.Trash = append(pg.Trash, c)
}

func trashUpTo(p *Player, pg *Playgroup, t int) {
	for i := 0; i < t; i++ {
		// decision point here
		cc := findCurses(p.Hand)
		if len(cc) > 0 {
			removeFromHand(p, cc[0])
			trashFromHand(p, pg, cc[0])
		}
	}
}

func checkReactions(p *Player) bool {
	defended := false
	rc := findReactions(p.Hand)
	for _, c := range rc {
		// decision point here
		if c.Reactions.Defend == true {
			defended = true
		}
	}
	return defended
}

func addDeckTop(p *Player, cs []cd.Card) {
	for _, c := range cs {
		fmt.Println("\t\t\t place on top of deck:", c.Name)
	}
	p.Deck = append(cs, p.Deck...)
}

func bestPlayableCard(cs []cd.Card) cd.Card {
	o := -1
	var b cd.Card
	var keepit bool
	for _, c := range cs {
		keepit = false
		if c.CTypes.Action == true {
			keepit = true
		}
		if c.CTypes.Treasure == true {
			keepit = true
		}
		if keepit == false {
			continue
		}
		if c.Cost > o {
			o = c.Cost
			b = c
		}
	}
	return b
}

func InitializeSupply(pl int) cd.Supply {

	var s cd.Supply

	/* coin cards */
	s.Piles = append(s.Piles, cd.SupplyPile{Card: bs.DefCopper(), Count: 60 - (pl * 7)})
	// s.Piles = append(s.Piles, cd.SupplyPile{Card: bs.DefSilver(), Count: 40})
	s.Piles = append(s.Piles, cd.SupplyPile{Card: bs.DefGold(), Count: 30})

	/* victory cards */
	vc := 12
	pc := vc
	switch pl {
	case 2:
		vc = 8
		pc = 8
	case 5:
		pc = 15
	case 6:
		pc = 18
	}
	s.Piles = append(s.Piles, cd.SupplyPile{Card: bs.DefEstate(), Count: vc})
	s.Piles = append(s.Piles, cd.SupplyPile{Card: bs.DefDuchy(), Count: vc})
	s.Piles = append(s.Piles, cd.SupplyPile{Card: bs.DefProvince(), Count: pc})

	/* curses */
	s.Piles = append(s.Piles, cd.SupplyPile{Card: bs.DefCurse(), Count: 10 * (pl - 1)})

	/* kingdom */
	s.Piles = append(s.Piles, cd.SupplyPile{Card: bs.DefCellar(), Count: 10})
	s.Piles = append(s.Piles, cd.SupplyPile{Card: bs.DefChapel(), Count: 10})
	s.Piles = append(s.Piles, cd.SupplyPile{Card: bs.DefMoat(), Count: 10})
	// s.Piles = append(s.Piles, cd.SupplyPile{Card: bs.DefVillage(), Count: 10})
	// s.Piles = append(s.Piles, cd.SupplyPile{Card: bs.DefWoodcutter(), Count: 10})
	// s.Piles = append(s.Piles, cd.SupplyPile{Card: bs.DefHarbinger(), Count: 10})
	s.Piles = append(s.Piles, cd.SupplyPile{Card: bs.DefVassal(), Count: 10})
	s.Piles = append(s.Piles, cd.SupplyPile{Card: bs.DefSmithy(), Count: 10})
	s.Piles = append(s.Piles, cd.SupplyPile{Card: bs.DefMilitia(), Count: 10})
	s.Piles = append(s.Piles, cd.SupplyPile{Card: bs.DefGardens(), Count: 10})
	s.Piles = append(s.Piles, cd.SupplyPile{Card: bs.DefFestival(), Count: 10})
	s.Piles = append(s.Piles, cd.SupplyPile{Card: bs.DefLaboratory(), Count: 10})
	s.Piles = append(s.Piles, cd.SupplyPile{Card: bs.DefMarket(), Count: 10})
	s.Piles = append(s.Piles, cd.SupplyPile{Card: bs.DefWitch(), Count: 10})

	return s
}

func initializeRandomizer() []cd.Card {

	var rd []cd.Card

	/* kingdom */
	rd = append(rd, bs.DefCellar())
	rd = append(rd, bs.DefChapel())
	rd = append(rd, bs.DefMoat())
	rd = append(rd, bs.DefVillage())
	rd = append(rd, bs.DefWoodcutter())
	rd = append(rd, bs.DefHarbinger())
	rd = append(rd, bs.DefSmithy())
	rd = append(rd, bs.DefMilitia())
	rd = append(rd, bs.DefGardens())
	rd = append(rd, bs.DefFestival())
	rd = append(rd, bs.DefLaboratory())
	rd = append(rd, bs.DefMarket())
	rd = append(rd, bs.DefWitch())

	return rd
}
