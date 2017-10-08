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
	fmt.Println("deck:", len(pg.Players[pg.PlayerTurn].Deck.Cards), "discard:", len(pg.Players[pg.PlayerTurn].Discard.Cards))

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

	validatePlayerCards(pg)

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
		c := SelectCardBuy(pg.ThisTurn.Coins, "any", pg.Supply)
		if c.Name == "" {
			break
		}
		fmt.Println("\t\t buying", c.Name)
		c = gainCard(&pg.Supply, c.Name)
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
	if len(nc) > 0 {
		p.Hand.Cards = append(p.Hand.Cards, nc...)
	}
}

func Draw(p *Player, d int) []cd.Card {
	var nc []cd.Card
	if len(p.Deck.Cards) == 0 {
		if len(p.Discard.Cards) == 0 {
			fmt.Println("WARNING: Not enough cards in deck to draw!")
			return nc
		}
	}
	for i := 0; i < d; i++ {
		// fmt.Println("deck size", len(p.Deck.Cards))
		c, z := p.Deck.Cards[0], p.Deck.Cards[1:]
		p.Deck.Cards = z
		nc = append(nc, c)
		if len(p.Deck.Cards) == 0 {
			p.Deck.Cards = p.Discard.Cards
			p.Discard.Cards = p.Discard.Cards[:0]
			ShuffleDeck(p)
		}
		if len(p.Deck.Cards) == 0 {
			fmt.Println("Out of cards!")
			break
		}
	}
	if len(nc) < d {
		fmt.Println("WARNING: Not enough cards in deck to draw!")
	}
	return nc
}

func SelectCardBuy(o int, t string, s cd.Supply) cd.Card {
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
		if matchType(pl.Card, t) != true {
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
		if matchType(pl.Card, t) != true {
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

/*
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
*/

func gainCard(s *cd.Supply, n string) cd.Card {
	for i, pl := range s.Piles {
		if pl.Card.Name == n {
			if s.Piles[i].Count == 0 {
				return cd.Card{} // empty card
			}
			s.Piles[i].Count--
			return pl.Card
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
		if matchType(c, t) {
			cs = append(cs, c)
		}
	}
	return cs
}

func matchType(c cd.Card, t string) bool {
	switch t {
	case "any":
		return true
	case "action":
		if c.CTypes.Action == true {
			return true
		}
	case "reaction":
		if c.CTypes.Reaction == true {
			return true
		}
	case "treasure":
		if c.CTypes.Treasure == true {
			return true
		}
	case "victory":
		if c.CTypes.Victory == true {
			return true
		}
	case "curse":
		if c.CTypes.Curse == true {
			return true
		}
	case "silver": // not a type?
		if c.Name == "Silver" {
			return true
		}
	case "actionExtra":
		if c.CTypes.Action == true {
			if c.Effects.ExtraActions > 0 {
				return true
			}
		}
	case "nonUsable":
		if c.CTypes.Action == false {
			if c.CTypes.Treasure == false {
				return true
			}
		}
	}
	return false
}

func findCards(h []cd.Card, n string, m int) []cd.Card {
	var fc []cd.Card
	f := 0
	for _, c := range h {
		if c.Name == n {
			fc = append(fc, c)
			f++
			if m > 0 {
				if f >= m {
					break
				}
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
		if len(nc) > 0 {
			p.Hand.Cards = append(p.Hand.Cards, nc...)
		}
	}
	if c.CTypes.Attack == true {
		resolveAttacks(pg, c)
	}
	// resolveSequence(pg, p, c.Effects.Sequence)
	// fmt.Println("discardsize:", len(p.Discard.Cards))
	// resolveSequence(pg, p, c.Effects.Sequence, c.Effects.SeqVal)
	resolveSequence(pg, p, c, "effect")
}

func resolveSequence(pg *Playgroup, p *Player, c cd.Card, effectType string) {
	// resolveSequence(pg *Playgroup, p *Player, seq []cd.Seq, seqVal map[string]int)
	seq := c.Effects.Sequence
	seqVal := c.Effects.SeqVal
	if effectType == "attack" {
		seq = c.Attacks.Sequence
		seqVal = c.Attacks.SeqVal
	}
	seqCards := make(map[string][]cd.Card)
Sequence:
	for _, seq := range seq {
		op := seq.Seq[0]
		switch op {
		case "getHandType":
			cardType := seq.Seq[1]
			matchingCards := seq.Seq[2]
			fmt.Println("\t\t\t getHandType", cardType, matchingCards)
			seqCards[matchingCards] = findCardType(p.Hand.Cards, cardType)
			if seqVal["getHandTypeMax"] > 0 {
				if seqVal["getHandTypeMax"] < len(seqCards[matchingCards]) {
					seqCards[matchingCards] = seqCards[matchingCards][:seqVal["getHandTypeMax"]]
				}
			}
			fmt.Println("\t\t\t found", showQuick(seqCards[matchingCards]))
		case "getCardType":
			cardSet := seq.Seq[1]
			cardType := seq.Seq[2]
			matchingCards := seq.Seq[3]
			fmt.Println("\t\t\t getCardType", cardSet, cardType, matchingCards)
			seqCards[matchingCards] = findCardType(seqCards[cardSet], cardType)
			if seqVal["getCardTypeMax"] > 0 {
				if seqVal["getCardTypeMax"] < len(seqCards[matchingCards]) {
					seqCards[matchingCards] = seqCards[matchingCards][:seqVal["getCardTypeMax"]]
				}
			}
			fmt.Println("\t\t\t found", showQuick(seqCards[matchingCards]))
		case "removeFromHands":
			removeThese := seq.Seq[1]
			fmt.Println("\t\t\t removeFromHands", removeThese, len(seqCards[removeThese]))
			removeFromHands(p, seqCards[removeThese])
		case "removeCards":
			removeThese := seq.Seq[1]
			fromThese := seq.Seq[2]
			fmt.Println("\t\t\t removeCards", showQuick(seqCards[removeThese]), "from set", showQuick(seqCards[fromThese]))
			seqCards[fromThese] = removeCards(seqCards[removeThese], seqCards[fromThese])
		case "countCards":
			countThese := seq.Seq[1]
			counted := seq.Seq[2]
			seqVal[counted] = len(seqCards[countThese])
			fmt.Println("\t\t\t countCards", countThese, counted, seqVal[counted])
		case "placeDiscards":
			discards := seq.Seq[1]
			fmt.Println("\t\t\t placeDiscards", discards, len(seqCards[discards]))
			if len(seqCards[discards]) < 1 {
				continue
			}
			discardCards(p, seqCards[discards])
		case "drawDeck":
			drawMax := seq.Seq[1]
			drewCards := seq.Seq[2]
			fmt.Println("\t\t\t drawDeck", drawMax, drewCards, seqVal[drawMax])
			nc := Draw(p, seqVal[drawMax])
			if len(nc) > 0 {
				seqCards[drewCards] = nc
			}
			fmt.Println("\t\t\t drew cards", showQuick(seqCards[drewCards]))
		case "placeHand":
			newCards := seq.Seq[1]
			fmt.Println("\t\t\t placeHand", newCards, len(seqCards[newCards]))
			p.Hand.Cards = append(p.Hand.Cards, seqCards[newCards]...)
		case "placeTrash":
			trashCards := seq.Seq[1]
			fmt.Println("\t\t\t placeTrash", trashCards, len(seqCards[trashCards]))
			pg.Trash.Cards = append(pg.Trash.Cards, seqCards[trashCards]...)
		case "GainCardName":
			wantCardName := seq.Seq[1]
			newCard := seq.Seq[2]
			fmt.Println("\t\t\t GainCardName", wantCardName, newCard)
			c := gainCard(&pg.Supply, wantCardName)
			if c.Name == "" {
				fmt.Println("\t\t\t nothing to gain!")
				continue
			}
			fmt.Println("\t\t\t gained", c.Name)
			seqCards[newCard] = append(seqCards[newCard], c)
		case "GainCardType":
			cardType := seq.Seq[1]
			newCard := seq.Seq[2]
			maxVal := seqVal["GainCardTypeMaxVal"]
			fmt.Println("\t\t\t GainCardType", cardType, newCard, maxVal)
			c := SelectCardBuy(maxVal, cardType, pg.Supply)
			if c.Name == "" {
				fmt.Println("\t\t\t nothing to gain!")
				continue
			}
			fmt.Println("\t\t\t gained", c.Name)
			seqCards[newCard] = append(seqCards[newCard], c)
		case "LoadDiscards":
			discards := seq.Seq[1]
			fmt.Println("\t\t\t LoadDiscards", discards)
			seqCards[discards] = p.Discard.Cards
			fmt.Println("\t\t\t LoadDiscards", discards, len(seqCards[discards]))
		case "findBestPlayable":
			cardSet := seq.Seq[1]
			bestCard := seq.Seq[2]
			fmt.Println("\t\t\t findBestPlayable", cardSet, bestCard)
			bc := bestPlayableCard(seqCards[cardSet])
			if bc.Name == "" {
				fmt.Println("\t\t\t nothing to play!")
				continue
			}
			seqCards[bestCard] = append(seqCards[bestCard], bc)
			fmt.Println("\t\t\t found", bc.Name)
		case "RetrieveDiscard":
			// TODO: remove multiple cards?
			removeCard := seq.Seq[1]
			fmt.Println("\t\t\t RetrieveDiscard", removeCard)
			if len(seqCards[removeCard]) < 1 {
				continue
			}
			removeFromDiscard(p, seqCards[removeCard][0])
		case "PlaceDeck":
			placeCards := seq.Seq[1]
			fmt.Println("\t\t\t PlaceDeck", placeCards)
			if len(seqCards[placeCards]) < 1 {
				continue
			}
			p.Deck.Cards = append(seqCards[placeCards], p.Deck.Cards...)
		case "PlayAction":
			playCards := seq.Seq[1]
			if seqVal["PlayActionTimes"] == 0 {
				seqVal["PlayActionTimes"] = 1
			}
			fmt.Println("\t\t\t PlayAction", showQuick(seqCards[playCards]))
			// decision point here whether to play any action
			for _, c := range seqCards[playCards] {
				fmt.Println("\t\t\t seq play action", c.Name)
				p.InPlay.Cards = append(p.InPlay.Cards, c)
				for j := 0; j < seqVal["PlayActionTimes"]; j++ {
					resolveEffects(pg, c)
				}
			}
		case "GetHandMatch":
			matchCard := seq.Seq[1]
			cardSet := seq.Seq[2]
			maxMatches := seqVal["GetHandMatchMax"]
			mc := findCards(p.Hand.Cards, matchCard, maxMatches)
			fmt.Println("\t\t\t GetHandMatch", matchCard, cardSet, maxMatches, "found", showQuick(mc))
			seqCards[cardSet] = mc
		case "AddXCoins":
			cardSet := seq.Seq[1]
			cardVal := seqVal["AddXCoinsVal"]
			pg.ThisTurn.Coins += (len(seqCards[cardSet]) * cardVal)
		case "getCost":
			cardSet := seq.Seq[1]
			costSum := seq.Seq[2]
			seqVal[costSum] = 0
			for _, c := range seqCards[cardSet] {
				seqVal[costSum] += c.Cost
			}
			fmt.Println("\t\t\t getCost", cardSet, showQuick(seqCards[cardSet]), costSum, seqVal[costSum])
		case "addVal":
			val := seq.Seq[1]
			oldval := seqVal[val] // debug only
			seqVal[val] += seqVal["addValVal"]
			fmt.Println("\t\t\t addVal", val, oldval, "+", seqVal["addValVal"], "=", seqVal[val])
		case "copyVal":
			srcVal := seq.Seq[1]
			dstVal := seq.Seq[2]
			seqVal[dstVal] = seqVal[srcVal]
			fmt.Println("\t\t\t copyVal", srcVal, seqVal[srcVal], dstVal, seqVal[dstVal])
		case "breakSet":
			needSet := seq.Seq[1]
			fmt.Println("\t\t\t breakSet", needSet, len(seqCards[needSet]))
			if len(seqCards[needSet]) < 1 {
				break Sequence
			}
		default:
			fmt.Println("ERROR: No operation", op)
		}
	}
}

/*
func resolveSequence(pg *Playgroup, p *Player, seq []cd.Sequence) {
	var cardSet []cd.Card
	// var gainCost int
	seqCards := make(map[string][]cd.Card)
	seqCardX := make(map[string]cd.Card)
	seqVal := make(map[string]int)
	seqType := make(map[string]string)
	for i, s := range seq {
		fmt.Println("\t\t\t Sequence", i)
		if s.SetVal.Name != "" {
			fmt.Println("\t\t\t\t SetVal", s.SetVal.Name)
			seqVal[s.SetVal.Name] = s.SetVal.Val
			seqType[s.SetVal.Name] = s.SetVal.Type
			seqCardX[s.SetVal.Name] = s.SetVal.Card
			seqCards[s.SetVal.Name] = append(seqCards[s.SetVal.Name], s.SetVal.Card)
		}
			if s.CountDiscard != "" {
				// decision point here
				vc := findCardType(p.Hand.Cards, "nonUsable")
				discardCards(p, vc)
				seqVal[s.CountDiscard] = len(vc)
				fmt.Println("\t\t\t\t CountDiscard", seqVal[s.CountDiscard])
			}
			if s.DrawCount != "" {
				nc := Draw(p, seqVal[s.DrawCount])
				p.Hand.Cards = append(p.Hand.Cards, nc...)
				fmt.Println("\t\t\t\t DrawCount", seqVal[s.DrawCount])
			}
		if s.TrashMax != "" {
			// decision point here
			cc := findCardType(p.Hand.Cards, "curse")
			for j, u := range cc {
				if j > seqVal[s.TrashMax] {
					break
				}
				removeFromHand(p, u)
				trashFromHand(p, pg, u)
				cardSet = append(cardSet, u)
			}
			fmt.Println("\t\t\t\t TrashMax", s.TrashMax, len(cardSet))
		}
		if s.RetrieveDiscard != "" {
			fmt.Println("\t\t\t\t RetrieveDiscard", seqVal[s.RetrieveDiscard])
			var cs []cd.Card
			for j := 0; j < seqVal[s.RetrieveDiscard]; j++ {
				// decision point here
				bc := bestPlayableCard(p.Discard.Cards)
				// must have found something
				if bc.Name != "" {
					removeFromDiscard(p, bc)
					cs = append(cs, bc)
				}
			}
			seqCards[s.RetrieveDiscard] = cs
		}
		if s.PlaceDeck != "" {
			fmt.Println("\t\t\t\t PlaceDeck", showQuick(seqCards[s.PlaceDeck]))
			addDeckTop(p, seqCards[s.PlaceDeck])
		}
		if s.GetSupplyCard != "" {
			fmt.Println("\t\t\t\t GetSupplyCard", s.GetSupplyCard)
			c := gainCard(&pg.Supply, seqCardX[s.GetSupplyCard])
			if c.Name == "" {
				fmt.Println("\t\t\t\t GetSupplyCard", s.GetSupplyCard, "no card to match")
				continue
			}
			var cs []cd.Card
			cs = append(cs, c)
			seqCards[s.GetSupplyCard] = cs
		}
			if s.GetHandTypeX != "" {
				// TODO: set count
				vc := findCardType(p.Hand.Cards, seqType[s.GetHandTypeX])
				var cs []cd.Card
				if len(vc) > 0 {
					cs = append(cs, vc[0])
				}
				fmt.Println("\t\t\t\t GetHandTypeX", s.GetHandTypeX, seqType[s.GetHandTypeX], "found", showQuick(cs))
				seqCards[s.GetHandTypeX] = cs
			}
		if s.GetHandType != "" {
			vc := findCardType(p.Hand.Cards, seqType[s.GetHandType])
			if len(vc) == 0 {
				fmt.Println("\t\t\t\t GetHandType", s.GetHandType, seqType[s.GetHandType], "no match")
			}
			maxCards := len(vc)
			if seqVal[s.GetHandType] > 0 {
				if seqVal[s.GetHandType] < maxCards {
					maxCards = seqVal[s.GetHandType]
				}
			}
			var cs []cd.Card
			cs = append(cs, vc[:maxCards]...)
			fmt.Println("\t\t\t\t GetHandType", s.GetHandType, seqType[s.GetHandType], "found", showQuick(cs))
			seqCards[s.GetHandType] = cs
		}
		if s.GetHandMatch != "" {
			// TODO: set count
			fmt.Println("\t\t\t\t GetHandMatch", s.GetHandMatch, seqCardX[s.GetHandMatch].Name)
			fmt.Println("\t\t\t\t GetHandMatch set1", showQuick(seqCards[s.GetHandMatch]))
			mc := findCards(p.Hand.Cards, seqCardX[s.GetHandMatch], 1)
			fmt.Println("\t\t\t\t GetHandMatch cs", showQuick(mc))
			var cs []cd.Card
			if len(mc) > 0 {
				cs = append(cs, mc[0])
			}
			seqCards[s.GetHandMatch] = cs
			fmt.Println("\t\t\t\t GetHandMatch set2", showQuick(seqCards[s.GetHandMatch]))
		}
		if s.DrawDeck != "" {
			fmt.Println("\t\t\t\t DrawDeck", s.DrawDeck, seqVal[s.DrawDeck])
			var cs []cd.Card
			cs = append(cs, Draw(p, seqVal[s.DrawDeck])...)
			seqCards[s.DrawDeck] = cs
			fmt.Println("\t\t\t\t drew", showQuick(cs))
		}
		if s.DiscardNonMatch != "" {
			fmt.Println("\t\t\t\t DiscardNonMatch", s.DiscardNonMatch, seqType[s.DiscardNonMatch])
			var oldSet Cards
			oldSet.Cards = seqCards[s.DiscardNonMatch]
			seqCards[s.DiscardNonMatch] = getCards(&oldSet, findCardType(oldSet.Cards, seqType[s.DiscardNonMatch]))
			p.Discard.Cards = append(p.Discard.Cards, oldSet.Cards...)
		}
		if s.PlayAction != "" {
			fmt.Println("\t\t\t\t PlayAction", s.PlayAction, seqVal[s.PlayAction])
			// possible decision point here whether to play any action
			for _, c := range seqCards[s.PlayAction] {
				fmt.Println("\t\t\t seq play action", c.Name)
				p.InPlay.Cards = append(p.InPlay.Cards, c)
				for j := 0; j < seqVal[s.PlayAction]; j++ {
					resolveEffects(pg, c)
				}
			}
		}
		if s.GainCard != "" {
			fmt.Println("\t\t\t\t GainCard", s.GainCard, seqVal[s.GainCard], seqType[s.GainCard])
			c := SelectCardBuy(seqVal[s.GainCard], seqType[s.GainCard], pg.Supply)
			if c.Name == "" {
				fmt.Println("\t\t\t\t GainCard no card to match")
				continue
			}
			fmt.Println("\t\t\t\t seqCards", showQuick(seqCards[s.GainCard]), len(seqCards[s.GainCard]), "gained", c.Name)
			seqCards[s.GainCard] = append(seqCards[s.GainCard], c)
		}
		if s.PlaceDiscards != "" {
			fmt.Println("\t\t\t\t PlaceDiscards", s.PlaceDiscards, showQuick(seqCards[s.PlaceDiscards]))
			if len(seqCards[s.PlaceDiscards]) > 0 {
				p.Discard.Cards = append(p.Discard.Cards, seqCards[s.PlaceDiscards]...)
			}
		}
			if s.PlaceHand != "" {
				fmt.Println("\t\t\t\t PlaceHand", s.PlaceHand, seqCardX[s.PlaceHand].Name)
				p.Hand.Cards = append(p.Hand.Cards, seqCardX[s.PlaceHand])
			}
		if s.PlaceHands != "" {
			fmt.Println("\t\t\t\t PlaceHands", s.PlaceHands, seqCardX[s.PlaceHands].Name)
			if len(seqCards[s.PlaceHands]) > 0 {
				p.Hand.Cards = append(p.Hand.Cards, seqCards[s.PlaceHands]...)
			}
		}
		if s.AddCost != "" {
			o := 0
			for _, c := range seqCards[s.AddCost] {
				fmt.Println("\t\t\t\t AddCost", s.AddCost, c.Name)
				o += c.Cost
			}
			seqVal[s.AddCost] += o
		}
		if s.TrashCards != "" {
			for _, c := range seqCards[s.TrashCards] {
				fmt.Println("\t\t\t\t TrashCards", c.Name)
				getCard(&p.Hand, c) // remove from hand
				pg.Trash.Cards = append(pg.Trash.Cards, c)
			}
		}
		if s.AddXCoins != "" {
			fmt.Println("\t\t\t\t AddXCoins", s.AddXCoins, seqVal[s.AddXCoins])
			pg.ThisTurn.Coins += (len(seqCards[s.AddXCoins]) * seqVal[s.AddXCoins])
		}
		if s.CountCards != "" {
			seqVal[s.CountCards] = len(seqCards[s.CountCards])
			fmt.Println("\t\t\t\t CountCards", s.CountCards, seqVal[s.CountCards])
		}
		if s.ClearSet != "" {
			fmt.Println("\t\t\t\t ClearSet", s.ClearSet)
			seqCards[s.ClearSet] = seqCards[s.ClearSet][:0]
		}
		if s.RemoveFromHand != "" {
			fmt.Println("\t\t\t\t RemoveFromHand", s.RemoveFromHand, showQuick(seqCards[s.RemoveFromHand]))
			getCards(&p.Hand, seqCards[s.RemoveFromHand]) // remove from hand
		}
	}
}
*/

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
		// resolveSequence(pg, p, c.Attacks.Sequence, c.Attacks.SeqVal)
		resolveSequence(pg, p, c, "attack")
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

func showQuick(cs []cd.Card) string {
	var disp string
	for _, c := range cs {
		disp += c.Name + ", "
	}
	return disp
}

func showStatus(pg *Playgroup) {
	fmt.Println("\t\t\t actions", pg.ThisTurn.Actions)
	fmt.Println("\t\t\t buys", pg.ThisTurn.Buys)
	fmt.Println("\t\t\t coins", pg.ThisTurn.Coins)
}

func discardCards(p *Player, cs []cd.Card) {
	fmt.Print("\t\t\t discarding ")
	for _, c := range cs {
		// removeFromHand(p, c)
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
		removeFromHand(p, d)
		discardCards(p, []cd.Card{d})
	}
}

func gainCurse(p *Player, s *cd.Supply, m int) {
	for i := 0; i < m; i++ {
		// c := gainCard(s, bs.DefCurse())
		c := gainCard(s, "Curse")
		if c.Name == "Curse" {
			p.Discard.Cards = append(p.Discard.Cards, c)
		}
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

func removeFromHands(p *Player, cs []cd.Card) {
	for _, c := range cs {
		removeFromHand(p, c)
	}
}

func removeCards(rs, cs []cd.Card) []cd.Card {
	for _, r := range rs {
		for i, c := range cs {
			if r.Name == c.Name {
				cs = append(cs[:i], cs[i+1:]...)
				break
			}
		}
	}
	return cs
}

func removeFromDiscard(p *Player, c cd.Card) {
	for i, h := range p.Discard.Cards {
		if h.Name == c.Name {
			p.Discard.Cards = append(p.Discard.Cards[:i], p.Discard.Cards[i+1:]...)
			return
		}
	}
	fmt.Println("ERROR: could not remove from discard:", c.Name)
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

func validatePlayerCards(pg *Playgroup) {
	for _, p := range pg.Players {
		validateCards(p.Hand.Cards, p.Name, "hand")
		validateCards(p.Discard.Cards, p.Name, "discard")
		validateCards(p.Deck.Cards, p.Name, "deck")
	}
}

func validateCards(cs []cd.Card, pname, cardset string) {
	for _, c := range cs {
		if c.Name == "" {
			fmt.Println("ERROR: empty card! %s %s", pname, cardset)
		}
	}
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
	s.Piles = append(s.Piles, cd.SupplyPile{Card: bs.DefHarbinger(), Count: 10})
	s.Piles = append(s.Piles, cd.SupplyPile{Card: bs.DefVillage(), Count: 10})
	s.Piles = append(s.Piles, cd.SupplyPile{Card: bs.DefWorkshop(), Count: 10})
	s.Piles = append(s.Piles, cd.SupplyPile{Card: bs.DefWoodcutter(), Count: 10})
	s.Piles = append(s.Piles, cd.SupplyPile{Card: bs.DefVassal(), Count: 10})
	s.Piles = append(s.Piles, cd.SupplyPile{Card: bs.DefBureaucrat(), Count: 10})
	s.Piles = append(s.Piles, cd.SupplyPile{Card: bs.DefRemodel(), Count: 10})
	s.Piles = append(s.Piles, cd.SupplyPile{Card: bs.DefMoneylender(), Count: 10})
	s.Piles = append(s.Piles, cd.SupplyPile{Card: bs.DefMilitia(), Count: 10})
	s.Piles = append(s.Piles, cd.SupplyPile{Card: bs.DefGardens(), Count: 10})
	s.Piles = append(s.Piles, cd.SupplyPile{Card: bs.DefSmithy(), Count: 10})
	s.Piles = append(s.Piles, cd.SupplyPile{Card: bs.DefMine(), Count: 10})
	s.Piles = append(s.Piles, cd.SupplyPile{Card: bs.DefFestival(), Count: 10})
	s.Piles = append(s.Piles, cd.SupplyPile{Card: bs.DefLaboratory(), Count: 10})
	s.Piles = append(s.Piles, cd.SupplyPile{Card: bs.DefMarket(), Count: 10})
	s.Piles = append(s.Piles, cd.SupplyPile{Card: bs.DefWitch(), Count: 10})

	return s
}

/*
func initializeRandomizer() []cd.Card {

	var rd []cd.Card

	// kingdom
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
*/
