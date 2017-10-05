package basic

import cd "github.com/nickaubert/dominionator/cards"

func DefCopper() cd.Card {
	var c cd.Card
	c.Name = "Copper"
	c.Cost = 0
	c.Coins = 1
	c.CTypes.Treasure = true
	return c
}

func DefSilver() cd.Card {
	var c cd.Card
	c.Name = "Silver"
	c.Cost = 3
	c.Coins = 2
	c.CTypes.Treasure = true
	return c
}

func DefGold() cd.Card {
	var c cd.Card
	c.Name = "Gold"
	c.Cost = 6
	c.Coins = 3
	c.CTypes.Treasure = true
	return c
}

func DefEstate() cd.Card {
	var c cd.Card
	c.Name = "Estate"
	c.Cost = 2
	c.VP = 1
	c.CTypes.Victory = true
	return c
}

func DefDuchy() cd.Card {
	var c cd.Card
	c.Name = "Duchy"
	c.Cost = 5
	c.VP = 3
	c.CTypes.Victory = true
	return c
}

func DefProvince() cd.Card {
	var c cd.Card
	c.Name = "Province"
	c.Cost = 8
	c.VP = 6
	c.CTypes.Victory = true
	return c
}

func DefCurse() cd.Card {
	var c cd.Card
	c.Name = "Curse"
	c.Cost = 0
	c.VP = -1
	c.CTypes.Curse = true
	return c
}

func DefVillage() cd.Card {
	var c cd.Card
	c.Name = "Village"
	c.Cost = 3
	c.CTypes.Action = true
	c.Effects.DrawCard = 1
	c.Effects.ExtraActions = 2
	return c
}

func DefSmithy() cd.Card {
	var c cd.Card
	c.Name = "Smithy"
	c.Cost = 4
	c.CTypes.Action = true
	c.Effects.DrawCard = 3
	return c
}

func DefFestival() cd.Card {
	var c cd.Card
	c.Name = "Festival"
	c.Cost = 5
	c.CTypes.Action = true
	c.Effects.ExtraActions = 2
	c.Effects.ExtraBuys = 1
	c.Effects.ExtraCoins = 2
	return c
}

func DefLaboratory() cd.Card {
	var c cd.Card
	c.Name = "Laboratory"
	c.Cost = 5
	c.CTypes.Action = true
	c.Effects.DrawCard = 2
	c.Effects.ExtraActions = 1
	return c
}

func DefMarket() cd.Card {
	var c cd.Card
	c.Name = "Market"
	c.Cost = 5
	c.CTypes.Action = true
	c.Effects.DrawCard = 1
	c.Effects.ExtraActions = 1
	c.Effects.ExtraBuys = 1
	c.Effects.ExtraCoins = 1
	return c
}

func DefCellar() cd.Card {
	var c cd.Card
	c.Name = "Cellar"
	c.Cost = 2
	c.CTypes.Action = true
	c.Effects.ExtraActions = 1
	// actually nonUsable should be a decision point
	sq := c.Effects.Sequence
	sq = append(sq, cd.Seq{Seq: []string{"getHandType", "nonUsable", "unusables"}})
	sq = append(sq, cd.Seq{Seq: []string{"removeFromHands", "unusables"}})
	sq = append(sq, cd.Seq{Seq: []string{"countCards", "unusables", "cardCount"}})
	sq = append(sq, cd.Seq{Seq: []string{"placeDiscards", "unusables"}})
	sq = append(sq, cd.Seq{Seq: []string{"drawDeck", "cardCount", "newCards"}})
	sq = append(sq, cd.Seq{Seq: []string{"placeHand", "newCards"}})
	c.Effects.Sequence = sq
	/*
		c.Effects.Sequence = append(c.Effects.Sequence, cd.Sequence{SetVal: cd.SeqVar{Name: "cellar", Type: "nonUsable"}})
		c.Effects.Sequence = append(c.Effects.Sequence, cd.Sequence{GetHandType: "cellar"})
		c.Effects.Sequence = append(c.Effects.Sequence, cd.Sequence{RemoveFromHand: "cellar"})
		c.Effects.Sequence = append(c.Effects.Sequence, cd.Sequence{CountCards: "cellar"})
		c.Effects.Sequence = append(c.Effects.Sequence, cd.Sequence{PlaceDiscards: "cellar"})
		c.Effects.Sequence = append(c.Effects.Sequence, cd.Sequence{ClearSet: "cellar"})
		c.Effects.Sequence = append(c.Effects.Sequence, cd.Sequence{DrawDeck: "cellar"})
		c.Effects.Sequence = append(c.Effects.Sequence, cd.Sequence{PlaceHands: "cellar"})
	*/
	return c
}

/*
func DefChapel() cd.Card {
	var c cd.Card
	c.Name = "Chapel"
	c.Cost = 2
	c.CTypes.Action = true
	c.Effects.Sequence = append(c.Effects.Sequence, cd.Sequence{SetVal: cd.SeqVar{Name: "trash", Val: 4}})
	c.Effects.Sequence = append(c.Effects.Sequence, cd.Sequence{TrashMax: "trash"})
	return c
}

func DefCellar() cd.Card {
	var c cd.Card
	c.Name = "Cellar"
	c.Cost = 2
	c.CTypes.Action = true
	c.Effects.ExtraActions = 1
	// actually nonUsable should be a decision point
	c.Effects.Sequence = append(c.Effects.Sequence, cd.Sequence{SetVal: cd.SeqVar{Name: "cellar", Type: "nonUsable"}})
	c.Effects.Sequence = append(c.Effects.Sequence, cd.Sequence{GetHandType: "cellar"})
	c.Effects.Sequence = append(c.Effects.Sequence, cd.Sequence{RemoveFromHand: "cellar"})
	c.Effects.Sequence = append(c.Effects.Sequence, cd.Sequence{CountCards: "cellar"})
	c.Effects.Sequence = append(c.Effects.Sequence, cd.Sequence{PlaceDiscards: "cellar"})
	c.Effects.Sequence = append(c.Effects.Sequence, cd.Sequence{ClearSet: "cellar"})
	c.Effects.Sequence = append(c.Effects.Sequence, cd.Sequence{DrawDeck: "cellar"})
	c.Effects.Sequence = append(c.Effects.Sequence, cd.Sequence{PlaceHands: "cellar"})
	return c
}

func DefMoat() cd.Card {
	var c cd.Card
	c.Name = "Moat"
	c.Cost = 2
	c.CTypes.Action = true
	c.CTypes.Reaction = true
	c.Effects.DrawCard = 2
	c.Reactions.Defend = true
	return c
}

func DefHarbinger() cd.Card {
	var c cd.Card
	c.Name = "Harbinger"
	c.Cost = 3
	c.CTypes.Action = true
	c.Effects.DrawCard = 1
	c.Effects.ExtraActions = 1
	c.Effects.Sequence = append(c.Effects.Sequence, cd.Sequence{SetVal: cd.SeqVar{Name: "retrieve", Val: 1}})
	c.Effects.Sequence = append(c.Effects.Sequence, cd.Sequence{RetrieveDiscard: "retrieve"})
	c.Effects.Sequence = append(c.Effects.Sequence, cd.Sequence{PlaceDeck: "retrieve"})
	return c
}

func DefBureaucrat() cd.Card {
	var c cd.Card
	c.Name = "Bureaucrat"
	c.Cost = 4
	c.CTypes.Action = true
	c.CTypes.Attack = true
	c.Effects.Sequence = append(c.Effects.Sequence, cd.Sequence{SetVal: cd.SeqVar{Name: "silver", Card: DefSilver()}})
	c.Effects.Sequence = append(c.Effects.Sequence, cd.Sequence{GetSupplyCard: "silver"})
	c.Effects.Sequence = append(c.Effects.Sequence, cd.Sequence{PlaceDeck: "silver"})
	c.Attacks.Sequence = append(c.Attacks.Sequence, cd.Sequence{SetVal: cd.SeqVar{Name: "putdeck", Type: "victory", Val: 1}})
	c.Attacks.Sequence = append(c.Attacks.Sequence, cd.Sequence{GetHandType: "putdeck"})
	c.Attacks.Sequence = append(c.Attacks.Sequence, cd.Sequence{PlaceDeck: "putdeck"})
	// or reveals a hand with no victory cards...
	return c
}

func DefVassal() cd.Card {
	var c cd.Card
	c.Name = "Vassal"
	c.Cost = 3
	c.CTypes.Action = true
	c.Effects.ExtraCoins = 2
	c.Effects.Sequence = append(c.Effects.Sequence, cd.Sequence{SetVal: cd.SeqVar{Name: "vassal", Val: 1, Type: "action"}})
	c.Effects.Sequence = append(c.Effects.Sequence, cd.Sequence{DrawDeck: "vassal"})
	c.Effects.Sequence = append(c.Effects.Sequence, cd.Sequence{DiscardNonMatch: "vassal"})
	c.Effects.Sequence = append(c.Effects.Sequence, cd.Sequence{PlayAction: "vassal"})
	return c
}

func DefWorkshop() cd.Card {
	var c cd.Card
	c.Name = "Workshop"
	c.Cost = 3
	c.CTypes.Action = true
	c.Effects.Sequence = append(c.Effects.Sequence, cd.Sequence{SetVal: cd.SeqVar{Name: "workshop", Val: 4, Type: "any"}})
	c.Effects.Sequence = append(c.Effects.Sequence, cd.Sequence{ClearSet: "workshop"})
	c.Effects.Sequence = append(c.Effects.Sequence, cd.Sequence{GainCard: "workshop"})
	c.Effects.Sequence = append(c.Effects.Sequence, cd.Sequence{PlaceDiscards: "workshop"})
	return c
}

func DefRemodel() cd.Card {
	var c cd.Card
	c.Name = "Remodel"
	c.Cost = 4
	c.CTypes.Action = true
	c.Effects.Sequence = append(c.Effects.Sequence, cd.Sequence{SetVal: cd.SeqVar{Name: "remodel", Type: "any", Val: 1}})
	c.Effects.Sequence = append(c.Effects.Sequence, cd.Sequence{GetHandType: "remodel"})
	c.Effects.Sequence = append(c.Effects.Sequence, cd.Sequence{SetVal: cd.SeqVar{Name: "remodel", Type: "any", Val: 2}})
	c.Effects.Sequence = append(c.Effects.Sequence, cd.Sequence{AddCost: "remodel"})
	c.Effects.Sequence = append(c.Effects.Sequence, cd.Sequence{TrashCards: "remodel"})
	c.Effects.Sequence = append(c.Effects.Sequence, cd.Sequence{ClearSet: "remodel"})
	c.Effects.Sequence = append(c.Effects.Sequence, cd.Sequence{GainCard: "remodel"})
	c.Effects.Sequence = append(c.Effects.Sequence, cd.Sequence{PlaceDiscards: "remodel"})
	return c
}

func DefMine() cd.Card {
	var c cd.Card
	c.Name = "Mine"
	c.Cost = 5
	c.CTypes.Action = true
	c.Effects.Sequence = append(c.Effects.Sequence, cd.Sequence{SetVal: cd.SeqVar{Name: "mine", Type: "treasure", Val: 1}})
	c.Effects.Sequence = append(c.Effects.Sequence, cd.Sequence{GetHandType: "mine"})
	c.Effects.Sequence = append(c.Effects.Sequence, cd.Sequence{SetVal: cd.SeqVar{Name: "mine", Type: "treasure", Val: 3}})
	c.Effects.Sequence = append(c.Effects.Sequence, cd.Sequence{AddCost: "mine"})
	c.Effects.Sequence = append(c.Effects.Sequence, cd.Sequence{TrashCards: "mine"})
	c.Effects.Sequence = append(c.Effects.Sequence, cd.Sequence{ClearSet: "mine"})
	c.Effects.Sequence = append(c.Effects.Sequence, cd.Sequence{GainCard: "mine"})
	c.Effects.Sequence = append(c.Effects.Sequence, cd.Sequence{PlaceHands: "mine"})
	return c
}

func DefGardens() cd.Card {
	var c cd.Card
	c.Name = "Gardens"
	c.Cost = 4
	c.CTypes.Victory = true
	c.Victories.CardsPerPoint = 10
	return c
}

func DefMilitia() cd.Card {
	var c cd.Card
	c.Name = "Militia"
	c.Cost = 4
	c.CTypes.Action = true
	c.CTypes.Attack = true
	c.Effects.ExtraCoins = 2
	c.Attacks.DiscardTo = 3
	return c
}

func DefMoneylender() cd.Card {
	var c cd.Card
	c.Name = "Moneylender"
	c.Cost = 4
	c.CTypes.Action = true
	c.Effects.Sequence = append(c.Effects.Sequence, cd.Sequence{SetVal: cd.SeqVar{Name: "moneylender", Card: DefCopper(), Val: 3}})
	c.Effects.Sequence = append(c.Effects.Sequence, cd.Sequence{GetHandMatch: "moneylender"})
	c.Effects.Sequence = append(c.Effects.Sequence, cd.Sequence{AddXCoins: "moneylender"})
	c.Effects.Sequence = append(c.Effects.Sequence, cd.Sequence{TrashCards: "moneylender"})
	return c
}

func DefWoodcutter() cd.Card {
	var c cd.Card
	c.Name = "Woodcutter"
	c.Cost = 3
	c.CTypes.Action = true
	c.Effects.ExtraBuys = 1
	c.Effects.ExtraCoins = 2
	return c
}

func DefWitch() cd.Card {
	var c cd.Card
	c.Name = "Witch"
	c.Cost = 5
	c.CTypes.Action = true
	c.CTypes.Attack = true
	c.Effects.DrawCard = 2
	c.Attacks.GainCurse = 1
	return c
}
*/

func InitialDeck() []cd.Card {
	var d []cd.Card
	for i := 0; i < 7; i++ {
		c := DefCopper()
		d = append(d, c)
	}
	for i := 0; i < 3; i++ {
		c := DefEstate()
		d = append(d, c)
	}
	return d
}

//func InitializeSupply(pl int) cd.Supply {
//
//	var s cd.Supply
//
//	/* coin cards */
//	s.Piles = append(s.Piles, cd.SupplyPile{Card: DefCopper(), Count: 60 - (pl * 7)})
//	s.Piles = append(s.Piles, cd.SupplyPile{Card: DefSilver(), Count: 40})
//	s.Piles = append(s.Piles, cd.SupplyPile{Card: DefGold(), Count: 30})
//
//	/* victory cards */
//	vc := 12
//	pc := vc
//	switch pl {
//	case 2:
//		vc = 8
//		pc = 8
//	case 5:
//		pc = 15
//	case 6:
//		pc = 18
//	}
//	s.Piles = append(s.Piles, cd.SupplyPile{Card: DefEstate(), Count: vc})
//	s.Piles = append(s.Piles, cd.SupplyPile{Card: DefDuchy(), Count: vc})
//	s.Piles = append(s.Piles, cd.SupplyPile{Card: DefProvince(), Count: pc})
//
//	/* curses */
//	s.Piles = append(s.Piles, cd.SupplyPile{Card: DefCurse(), Count: 10 * (pl - 1)})
//
//	/* kingdom */
//	s.Piles = append(s.Piles, cd.SupplyPile{Card: DefCellar(), Count: 10})
//	s.Piles = append(s.Piles, cd.SupplyPile{Card: DefChapel(), Count: 10})
//	s.Piles = append(s.Piles, cd.SupplyPile{Card: DefMoat(), Count: 10})
//	s.Piles = append(s.Piles, cd.SupplyPile{Card: DefVillage(), Count: 10})
//	s.Piles = append(s.Piles, cd.SupplyPile{Card: DefWoodcutter(), Count: 10})
//	s.Piles = append(s.Piles, cd.SupplyPile{Card: DefHarbinger(), Count: 10})
//	s.Piles = append(s.Piles, cd.SupplyPile{Card: DefSmithy(), Count: 10})
//	s.Piles = append(s.Piles, cd.SupplyPile{Card: DefMilitia(), Count: 10})
//	s.Piles = append(s.Piles, cd.SupplyPile{Card: DefGardens(), Count: 10})
//	s.Piles = append(s.Piles, cd.SupplyPile{Card: DefFestival(), Count: 10})
//	s.Piles = append(s.Piles, cd.SupplyPile{Card: DefLaboratory(), Count: 10})
//	s.Piles = append(s.Piles, cd.SupplyPile{Card: DefMarket(), Count: 10})
//	s.Piles = append(s.Piles, cd.SupplyPile{Card: DefWitch(), Count: 10})
//
//	return s
//}
