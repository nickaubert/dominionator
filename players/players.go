package players

import (
	"fmt"
	"math/rand"
	"time"
)

import cd "github.com/nickaubert/dominionator/cards"
import bs "github.com/nickaubert/dominionator/basic"

type Player struct {
	Deck    Cards
	Hand    Cards
	InPlay  Cards
	Discard Cards
	Name    string
}

type Cards struct {
	Cards []cd.Card
}

type Playgroup struct {
	Players    []Player
	PlayerTurn int
	ThisTurn   ThisTurn
	Supply     cd.Supply
	Trash      Cards
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
		pl.Deck.Cards = bs.InitialDeck()
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
	for pg.ThisTurn.Actions > 0 {
		ac := findCardType(p.Hand.Cards, "actionExtra")
		if len(ac) == 0 {
			ac = findCardType(p.Hand.Cards, "action")
		}
		if len(ac) == 0 {
			break
		}
		showCards("hand", p.Hand)
		c := getCard(&p.Hand, highestCostCard(ac))
		fmt.Println("\t\t playing", c.Name)
		pg.ThisTurn.Actions--
		p.InPlay.Cards = append(p.InPlay.Cards, c)
		resolveEffects(pg, c)
		showStatus(pg)
	}

}

func BuyPhase(pg *Playgroup) {

	p := &pg.Players[pg.PlayerTurn]
	fmt.Println("\t BuyPhase")
	showCards("hand", p.Hand)

	for pg.ThisTurn.Buys > 0 {

		// decision point about which cards to play will go here
		// usually all treasure cards will be activated on fist buy phase
		tc := getCards(&p.Hand, findCardType(p.Hand.Cards, "treasure"))
		for _, c := range tc {
			p.InPlay.Cards = append(p.InPlay.Cards, c)
			resolveEffects(pg, c)
			pg.ThisTurn.Coins += c.Coins
		}

		fmt.Println("\t\t", pg.ThisTurn.Coins, "coins to spend")
		c := SelectCardBuy(pg.ThisTurn.Coins, pg.Supply)
		if c.Name == "" {
			break
		}
		fmt.Println("\t\t buying", c.Name)
		c = gainCard(&pg.Supply, c)
		if c.Name == "" {
			panic(fmt.Sprintf("ERROR: missing gain card! %s %v", c, pg.Supply))
		}
		p.Discard.Cards = append(p.Discard.Cards, c)
		pg.ThisTurn.Coins -= c.Cost
		pg.ThisTurn.Buys--

	}

}

func CleanupPhase(pg *Playgroup) {
	fmt.Println("\t CleanupPhase")
	p := &pg.Players[pg.PlayerTurn]
	p.Discard.Cards = append(p.Discard.Cards, p.InPlay.Cards...)
	p.InPlay.Cards = p.InPlay.Cards[:0]
	p.Discard.Cards = append(p.Discard.Cards, p.Hand.Cards...)
	p.Hand.Cards = p.Hand.Cards[:0]
	nc := Draw(p, 5)
	p.Hand.Cards = append(p.Hand.Cards, nc...)
}

func Draw(p *Player, d int) []cd.Card {
	var nc []cd.Card
	for i := 0; i < d; i++ {
		c, z := p.Deck.Cards[0], p.Deck.Cards[1:]
		p.Deck.Cards = z
		nc = append(nc, c)
		if len(p.Deck.Cards) == 0 {
			p.Deck.Cards = p.Discard.Cards
			p.Discard.Cards = p.Discard.Cards[:0]
			ShuffleDeck(p)
		}
	}
	return nc
}

func SelectCardBuy(o int, s cd.Supply) cd.Card {
	// decision point here
	// this is not the best heuristic
	var highestCost int
	for _, pl := range s.Piles {
		if pl.Count == 0 {
			continue
		}
		if pl.Card.Cost > o {
			continue
		}
		if pl.Card.Cost > highestCost {
			highestCost = pl.Card.Cost
		}
	}
	var bestCards []cd.Card
	for _, pl := range s.Piles {
		if pl.Count == 0 {
			continue
		}
		if pl.Card.Cost != highestCost {
			continue
		}
		if pl.Card.CTypes.Curse == true {
			continue // decision point, yeah dont buy curses
		}
		bestCards = append(bestCards, pl.Card)
	}
	if len(bestCards) == 0 {
		return cd.Card{}
	}
	rand.Seed(time.Now().UnixNano())
	r := rand.Intn(len(bestCards))
	c := bestCards[r]
	return c
}

func gainCard(s *cd.Supply, c cd.Card) cd.Card {
	for n, pl := range s.Piles {
		if pl.Card.Name == c.Name {
			if s.Piles[n].Count == 0 {
				return cd.Card{}
			}
			s.Piles[n].Count--
			return c
		}
	}
	return cd.Card{}
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
		p.Deck.Cards = append(p.Deck.Cards, p.Hand.Cards...)
		p.Deck.Cards = append(p.Deck.Cards, p.Discard.Cards...)
		p.Deck.Cards = append(p.Deck.Cards, p.InPlay.Cards...)

		vp := countVictoryPoints(p.Deck.Cards)
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
	d := p.Deck.Cards
	var n []cd.Card
	for _ = range d {
		r := rand.Intn(len(d))
		n = append(n, d[r])
		d = append(d[:r], d[r+1:]...)
	}
	p.Deck.Cards = n
}

func findCardType(h []cd.Card, t string) []cd.Card {
	var cs []cd.Card
	for _, c := range h {
		switch t {
		case "action":
			if c.CTypes.Action == true {
				cs = append(cs, c)
			}
		case "reaction":
			if c.CTypes.Reaction == true {
				cs = append(cs, c)
			}
		case "treasure":
			if c.CTypes.Treasure == true {
				cs = append(cs, c)
			}
		case "victory":
			if c.CTypes.Victory == true {
				cs = append(cs, c)
			}
		case "curse":
			if c.CTypes.Curse == true {
				cs = append(cs, c)
			}
		case "actionExtra":
			if c.CTypes.Action == true {
				if c.Effects.ExtraActions > 0 {
					cs = append(cs, c)
				}
			}
		case "nonUsable":
			if c.CTypes.Action == false {
				if c.CTypes.Treasure == false {
					cs = append(cs, c)
				}
			}
		}
	}
	return cs
}

func findCards(h []cd.Card, mc cd.Card, m int) []cd.Card {
	var fc []cd.Card
	f := 0
	for _, c := range h {
		if c.Name == mc.Name {
			fc = append(fc, c)
			f++
			if f >= m {
				break
			}
		}
	}
	return fc
}

// Must either pass pointer to encapsulating object, or return modified slice
func getCard(stack *Cards, mc cd.Card) cd.Card {
	var fc cd.Card
	for i, c := range stack.Cards {
		if c.Name == mc.Name {
			stack.Cards = append(stack.Cards[:i], stack.Cards[i+1:]...)
			fc = mc
			break
		}
	}
	if fc.Name != mc.Name {
		panic(fmt.Sprintf("ERROR: missing card! %s %v", mc, stack))
	}
	return fc
}

func getCards(stack *Cards, set []cd.Card) []cd.Card {
	var fc []cd.Card
	for _, c := range set {
		fc = append(fc, getCard(stack, c))
	}
	return fc
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
	p := &pg.Players[pg.PlayerTurn]
	pg.ThisTurn.Actions += c.Effects.ExtraActions
	pg.ThisTurn.Buys += c.Effects.ExtraBuys
	pg.ThisTurn.Coins += c.Effects.ExtraCoins
	if c.Effects.DrawCard > 0 {
		nc := Draw(p, c.Effects.DrawCard)
		p.Hand.Cards = append(p.Hand.Cards, nc...)
	}
	if c.CTypes.Attack == true {
		resolveAttacks(pg, c)
	}
	resolveSequence(pg, p, c.Effects.Sequence)
}

func resolveSequence(pg *Playgroup, p *Player, seq []cd.Sequence) {
	var cardSet []cd.Card
	// countX := 0 // replace with len(cardSet)?
	for i, s := range seq {
		fmt.Println("\t\t\t Sequence", i)
		if s.CountDiscard > 0 {
			// decision point here
			vc := findCardType(p.Hand.Cards, "nonUsable")
			for j, v := range vc {
				if j > s.CountDiscard {
					break
				}
				discardCards(p, []cd.Card{v})
				cardSet = append(cardSet, v)
			}
			fmt.Println("\t\t\t\t CountDiscard", len(cardSet))
		}
		if s.DrawCount == true {
			nc := Draw(p, len(cardSet))
			p.Hand.Cards = append(p.Hand.Cards, nc...)
			fmt.Println("\t\t\t\t DrawCount", len(cardSet))
		}
		if s.CountTrash > 0 {
			// decision point here
			fmt.Println("\t\t\t\t CountTrash")
			cc := findCardType(p.Hand.Cards, "curse")
			for j, u := range cc {
				if j > s.CountTrash {
					break
				}
				removeFromHand(p, u)
				trashFromHand(p, pg, u)
				// countX++
				cardSet = append(cardSet, u)
			}
		}
		if s.RetrieveDiscard > 0 {
			fmt.Println("\t\t\t\t RetrieveDiscard")
			for j := 0; j < s.RetrieveDiscard; j++ {
				// decision point here
				bc := bestPlayableCard(p.Discard.Cards)
				// must have found something
				if bc.Name != "" {
					removeFromDiscard(p, bc)
					cardSet = append(cardSet, bc)
				}
			}
		}
		if s.PlaceDeck == true {
			fmt.Println("\t\t\t\t PlaceDeck")
			addDeckTop(p, cardSet)
		}
		if s.DrawDeck > 0 {
			fmt.Println("\t\t\t\t DrawDeck")
			cardSet = append(cardSet, Draw(p, s.DrawDeck)...)
		}
		if s.DiscardNonMatch != "" {
			fmt.Println("\t\t\t\t DiscardNonMatch", s.DiscardNonMatch)
			var oldSet Cards
			oldSet.Cards = cardSet
			cardSet = getCards(&oldSet, findCardType(oldSet.Cards, s.DiscardNonMatch))
			discardCards(p, oldSet.Cards)
		}
		if s.PlayAction > 0 {
			fmt.Println("\t\t\t\t PlayAction")
			// possible decision point here whether to play any action
			for _, c := range cardSet {
				fmt.Println("\t\t\t seq play action", c.Name)
				p.InPlay.Cards = append(p.InPlay.Cards, c)
				for j := 0; j < s.PlayAction; j++ {
					resolveEffects(pg, c)
				}
			}
		}
		if s.GainMax > 0 {
			c := SelectCardBuy(s.GainMax, pg.Supply)
			if c.Name == "" {
				fmt.Println("\t\t\t\t GainMax", s.GainMax, "no card to match")
				continue
			}
			c = gainCard(&pg.Supply, c)
			if c.Name == "" {
				panic(fmt.Sprintf("ERROR: missing gainmax card! %s %v", c, pg.Supply))
			}
			p.Discard.Cards = append(p.Discard.Cards, c)
			fmt.Println("\t\t\t\t GainMax", s.GainMax, c.Name)
		}
		if s.GetSupplyCard.Name != "" {
			fmt.Println("\t\t\t\t GetCard", s.GetSupplyCard.Name)
			c := gainCard(&pg.Supply, s.GetSupplyCard)
			if c.Name == "" {
				fmt.Println("\t\t\t\t GetSupplyCard", s.GetSupplyCard.Name, "no card to match")
				continue
			}
			cardSet = append(cardSet, c)
		}
		if s.GetHandType != "" {
			fmt.Println("\t\t\t\t GetHandType", s.GetHandType)
			vc := findCardType(p.Hand.Cards, s.GetHandType)
			if len(vc) > 0 {
				cardSet = append(cardSet, vc[0])
			}
		}
		if s.MayTrash.Name != "" {
			fmt.Println("\t\t\t\t MayTrash", s.MayTrash.Name)
			// decision point whether to trash here
			cs := findCards(p.Hand.Cards, s.MayTrash, 1)
			if len(cs) > 0 {
				c := getCard(&p.Hand, s.MayTrash)
				pg.Trash.Cards = append(pg.Trash.Cards, c)
				cardSet = append(cardSet, c)
			}
		}
		if s.AddXCoins > 0 {
			fmt.Println("\t\t\t\t AddXCoins", s.AddXCoins)
			pg.ThisTurn.Coins += (len(cardSet) * s.AddXCoins)
		}
	}
}

func resolveAttacks(pg *Playgroup, c cd.Card) {
	for i := range pg.Players {
		p := &pg.Players[i]
		if i == pg.PlayerTurn {
			continue
		}
		fmt.Println("\t\t Attacking", p.Name)
		defended := checkReactions(p)
		if defended == true {
			fmt.Println("\t\t defended!")
			continue
		}
		if c.Attacks.DiscardTo > 0 {
			discardTo(p, c.Attacks.DiscardTo)
		}
		if c.Attacks.GainCurse > 0 {
			gainCurse(p, &pg.Supply, c.Attacks.GainCurse)
		}
		resolveSequence(pg, p, c.Attacks.Sequence)
	}
	fmt.Println("\t\t finished attacks")
}

func showCards(s string, h Cards) {
	fmt.Print("\t\t", s, ": ")
	for _, c := range h.Cards {
		fmt.Print(c.Name, ", ")
	}
	fmt.Print("\n")
}

func showStatus(pg *Playgroup) {
	fmt.Println("\t\t\t actions", pg.ThisTurn.Actions)
	fmt.Println("\t\t\t buys", pg.ThisTurn.Buys)
	fmt.Println("\t\t\t coins", pg.ThisTurn.Coins)
}

func discardCards(p *Player, cs []cd.Card) {
	fmt.Print("\t\t\t discarding ")
	for _, c := range cs {
		removeFromHand(p, c)
		fmt.Print(c.Name, ", ")
		p.Discard.Cards = append(p.Discard.Cards, c)
	}
	fmt.Print("\n")
}

func selectDiscardOwn(p *Player) cd.Card {
	// discard victory cards first, then select lowest value
	vc := findCardType(p.Hand.Cards, "nonUsable")
	for _, c := range vc {
		return c
	}
	c := lowestCostCard(p.Hand.Cards)
	return c
}

func discardTo(p *Player, m int) {
	for len(p.Hand.Cards) > m {
		d := selectDiscardOwn(p)
		discardCards(p, []cd.Card{d})
	}
}

func gainCurse(p *Player, s *cd.Supply, m int) {
	for i := 0; i < m; i++ {
		c := gainCard(s, bs.DefCurse())
		p.Discard.Cards = append(p.Discard.Cards, c)
	}
}

func removeFromHand(p *Player, c cd.Card) {
	for i, h := range p.Hand.Cards {
		if h.Name == c.Name {
			p.Hand.Cards = append(p.Hand.Cards[:i], p.Hand.Cards[i+1:]...)
			break
		}
	}
}

func removeFromDiscard(p *Player, c cd.Card) {
	for i, h := range p.Discard.Cards {
		if h.Name == c.Name {
			p.Discard.Cards = append(p.Discard.Cards[:i], p.Discard.Cards[i+1:]...)
			break
		}
	}
}

func trashFromHand(p *Player, pg *Playgroup, c cd.Card) {
	fmt.Println("\t\t\t trashing", c.Name)
	removeFromHand(p, c)
	pg.Trash.Cards = append(pg.Trash.Cards, c)
}

func trashUpTo(p *Player, pg *Playgroup, t int) {
	for i := 0; i < t; i++ {
		// decision point here
		cc := findCardType(p.Hand.Cards, "curse")
		if len(cc) > 0 {
			removeFromHand(p, cc[0])
			trashFromHand(p, pg, cc[0])
		}
	}
}

func checkReactions(p *Player) bool {
	defended := false
	rc := findCardType(p.Hand.Cards, "reaction")
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
	p.Deck.Cards = append(cs, p.Deck.Cards...)
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
	s.Piles = append(s.Piles, cd.SupplyPile{Card: bs.DefSilver(), Count: 40})
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
	s.Piles = append(s.Piles, cd.SupplyPile{Card: bs.DefVillage(), Count: 10})
	s.Piles = append(s.Piles, cd.SupplyPile{Card: bs.DefWoodcutter(), Count: 10})
	s.Piles = append(s.Piles, cd.SupplyPile{Card: bs.DefHarbinger(), Count: 10})
	s.Piles = append(s.Piles, cd.SupplyPile{Card: bs.DefVassal(), Count: 10})
	s.Piles = append(s.Piles, cd.SupplyPile{Card: bs.DefWorkshop(), Count: 10})
	s.Piles = append(s.Piles, cd.SupplyPile{Card: bs.DefSmithy(), Count: 10})
	s.Piles = append(s.Piles, cd.SupplyPile{Card: bs.DefMilitia(), Count: 10})
	s.Piles = append(s.Piles, cd.SupplyPile{Card: bs.DefBureaucrat(), Count: 10})
	s.Piles = append(s.Piles, cd.SupplyPile{Card: bs.DefMoneylender(), Count: 10})
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
