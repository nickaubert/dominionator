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

func DefVillage() cd.Card {
	var c cd.Card
	c.Name = "Village"
	c.Cost = 3
	c.CTypes.Action = true
	c.Effects.DrawCard = 1
	c.Effects.ExtraActions = 2
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
	c.Effects.SeqVal = make(map[string]int)
	c.Effects.SeqVal["getHandTypeMax"] = 0
	sq := c.Effects.Sequence
	sq = append(sq, cd.Seq{Seq: []string{"getHandType", "nonUsable", "unusables", "getHandTypeMax"}})
	sq = append(sq, cd.Seq{Seq: []string{"breakSet", "unusables"}})
	sq = append(sq, cd.Seq{Seq: []string{"removeFromHands", "unusables"}})
	sq = append(sq, cd.Seq{Seq: []string{"countCards", "unusables", "cardCount"}})
	sq = append(sq, cd.Seq{Seq: []string{"placeDiscards", "unusables"}})
	sq = append(sq, cd.Seq{Seq: []string{"drawDeck", "cardCount", "newCards"}})
	sq = append(sq, cd.Seq{Seq: []string{"placeHand", "newCards"}})
	c.Effects.Sequence = sq
	return c
}

func DefChapel() cd.Card {
	var c cd.Card
	c.Name = "Chapel"
	c.Cost = 2
	c.CTypes.Action = true
	c.Effects.SeqVal = make(map[string]int)
	c.Effects.SeqVal["getHandTypeMax"] = 4
	sq := c.Effects.Sequence
	sq = append(sq, cd.Seq{Seq: []string{"getHandType", "curse", "trashThese", "getHandTypeMax"}})
	sq = append(sq, cd.Seq{Seq: []string{"removeFromHands", "trashThese"}})
	sq = append(sq, cd.Seq{Seq: []string{"placeTrash", "trashThese"}})
	c.Effects.Sequence = sq
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

func DefWorkshop() cd.Card {
	var c cd.Card
	c.Name = "Workshop"
	c.Cost = 3
	c.CTypes.Action = true
	c.Effects.SeqVal = make(map[string]int)
	c.Effects.SeqVal["GainCardTypeMaxVal"] = 4
	sq := c.Effects.Sequence
	sq = append(sq, cd.Seq{Seq: []string{"GainCardType", "any", "newCard", "GainCardTypeMaxVal"}})
	sq = append(sq, cd.Seq{Seq: []string{"placeDiscards", "newCard"}})
	c.Effects.Sequence = sq
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

func DefHarbinger() cd.Card {
	var c cd.Card
	c.Name = "Harbinger"
	c.Cost = 3
	c.CTypes.Action = true
	c.Effects.DrawCard = 1
	c.Effects.ExtraActions = 1
	c.Effects.SeqVal = make(map[string]int)
	sq := c.Effects.Sequence
	sq = append(sq, cd.Seq{Seq: []string{"LoadDiscards", "discards"}})
	sq = append(sq, cd.Seq{Seq: []string{"findBestPlayable", "discards", "bestCard"}})
	sq = append(sq, cd.Seq{Seq: []string{"RetrieveDiscard", "bestCard"}})
	sq = append(sq, cd.Seq{Seq: []string{"PlaceDeck", "bestCard"}})
	c.Effects.Sequence = sq
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

func DefVassal() cd.Card {
	var c cd.Card
	c.Name = "Vassal"
	c.Cost = 3
	c.CTypes.Action = true
	c.Effects.ExtraCoins = 2
	c.Effects.SeqVal = make(map[string]int)
	c.Effects.SeqVal["drawDeckMax"] = 1
	c.Effects.SeqVal["getCardTypeMax"] = 0
	sq := c.Effects.Sequence
	sq = append(sq, cd.Seq{Seq: []string{"drawDeck", "drawDeckMax", "newCard"}})
	sq = append(sq, cd.Seq{Seq: []string{"getCardType", "newCard", "action", "actionCard", "getCardTypeMax"}})
	sq = append(sq, cd.Seq{Seq: []string{"removeCards", "actionCard", "newCard"}})
	sq = append(sq, cd.Seq{Seq: []string{"placeDiscards", "newCard"}})
	sq = append(sq, cd.Seq{Seq: []string{"PlayAction", "actionCard"}})
	c.Effects.Sequence = sq
	return c
}

func DefBureaucrat() cd.Card {
	var c cd.Card
	c.Name = "Bureaucrat"
	c.Cost = 4
	c.CTypes.Action = true
	c.CTypes.Attack = true

	c.Effects.SeqVal = make(map[string]int)
	sq := c.Effects.Sequence
	sq = append(sq, cd.Seq{Seq: []string{"GainCardName", "Silver", "newSilver"}})
	sq = append(sq, cd.Seq{Seq: []string{"PlaceDeck", "newSilver"}})
	c.Effects.Sequence = sq

	c.Attacks.SeqVal = make(map[string]int)
	c.Attacks.SeqVal["getHandTypeMax"] = 1
	aq := c.Attacks.Sequence
	aq = append(aq, cd.Seq{Seq: []string{"getHandType", "victory", "putdeck", "getHandTypeMax"}})
	aq = append(aq, cd.Seq{Seq: []string{"removeFromHands", "putdeck"}})
	aq = append(aq, cd.Seq{Seq: []string{"PlaceDeck", "putdeck"}})
	// or reveals a hand with no victory cards...
	c.Attacks.Sequence = aq

	return c

}

func DefCouncilRoom() cd.Card {
	var c cd.Card
	c.Name = "CouncilRoom"
	c.Cost = 5
	c.Effects.DrawCard = 4
	c.Effects.ExtraBuys = 1
	c.CTypes.Action = true
	c.CTypes.Attack = false // but using "attacks" anyway

	// not an "attack" per se
	c.Attacks.SeqVal = make(map[string]int)
	c.Attacks.SeqVal["drawDeckMax"] = 1
	aq := c.Attacks.Sequence
	aq = append(aq, cd.Seq{Seq: []string{"drawDeck", "drawDeckMax", "newCard"}})
	aq = append(aq, cd.Seq{Seq: []string{"placeHand", "newCard"}})
	c.Attacks.Sequence = aq
	return c

}

func DefMoneylender() cd.Card {
	var c cd.Card
	c.Name = "Moneylender"
	c.Cost = 4
	c.CTypes.Action = true

	c.Effects.SeqVal = make(map[string]int)
	c.Effects.SeqVal["GetHandMatchMax"] = 1
	c.Effects.SeqVal["AddXCoinsVal"] = 3
	sq := c.Effects.Sequence
	sq = append(sq, cd.Seq{Seq: []string{"GetHandMatch", "Copper", "coppers"}})
	sq = append(sq, cd.Seq{Seq: []string{"breakSet", "coppers"}})
	sq = append(sq, cd.Seq{Seq: []string{"removeFromHands", "coppers"}})
	sq = append(sq, cd.Seq{Seq: []string{"AddXCoins", "coppers"}})
	sq = append(sq, cd.Seq{Seq: []string{"placeTrash", "coppers"}})
	c.Effects.Sequence = sq
	return c
}

func DefMine() cd.Card {
	var c cd.Card
	c.Name = "Mine"
	c.Cost = 5
	c.CTypes.Action = true

	c.Effects.SeqVal = make(map[string]int)
	c.Effects.SeqVal["getHandTypeMax"] = 1
	c.Effects.SeqVal["addValVal"] = 3
	sq := c.Effects.Sequence
	sq = append(sq, cd.Seq{Seq: []string{"getHandType", "treasure", "oldtreasure", "getHandTypeMax"}})
	sq = append(sq, cd.Seq{Seq: []string{"breakSet", "oldtreasure"}})
	sq = append(sq, cd.Seq{Seq: []string{"removeFromHands", "oldtreasure"}})
	sq = append(sq, cd.Seq{Seq: []string{"getCost", "oldtreasure", "treasurecost"}})
	sq = append(sq, cd.Seq{Seq: []string{"placeTrash", "oldtreasure"}})
	sq = append(sq, cd.Seq{Seq: []string{"addVal", "treasurecost"}})
	sq = append(sq, cd.Seq{Seq: []string{"copyVal", "treasurecost", "GainCardTypeMaxVal"}})
	sq = append(sq, cd.Seq{Seq: []string{"GainCardType", "treasure", "newtreasure", "GainCardTypeMaxVal"}})
	sq = append(sq, cd.Seq{Seq: []string{"placeHand", "newtreasure"}})
	c.Effects.Sequence = sq
	return c
}

func DefRemodel() cd.Card {
	var c cd.Card
	c.Name = "Remodel"
	c.Cost = 4
	c.CTypes.Action = true

	c.Effects.SeqVal = make(map[string]int)
	c.Effects.SeqVal["getHandTypeMax"] = 1
	c.Effects.SeqVal["addValVal"] = 2
	sq := c.Effects.Sequence
	sq = append(sq, cd.Seq{Seq: []string{"getHandType", "any", "oldcard", "getHandTypeMax"}})
	sq = append(sq, cd.Seq{Seq: []string{"breakSet", "oldcard"}})
	sq = append(sq, cd.Seq{Seq: []string{"removeFromHands", "oldcard"}})
	sq = append(sq, cd.Seq{Seq: []string{"getCost", "oldcard", "cardcost"}})
	sq = append(sq, cd.Seq{Seq: []string{"placeTrash", "oldcard"}})
	sq = append(sq, cd.Seq{Seq: []string{"addVal", "cardcost"}})
	sq = append(sq, cd.Seq{Seq: []string{"copyVal", "cardcost", "GainCardTypeMaxVal"}})
	sq = append(sq, cd.Seq{Seq: []string{"GainCardType", "any", "newcard", "GainCardTypeMaxVal"}})
	sq = append(sq, cd.Seq{Seq: []string{"placeDiscards", "newtreasure"}})
	c.Effects.Sequence = sq
	return c
}

func DefThroneRoom() cd.Card {
	var c cd.Card
	c.Name = "ThroneRoom"
	c.Cost = 4
	c.CTypes.Action = true

	c.Effects.SeqVal = make(map[string]int)
	c.Effects.SeqVal["getHandTypeMax"] = 1
	c.Effects.SeqVal["PlayActionTimes"] = 2
	sq := c.Effects.Sequence
	sq = append(sq, cd.Seq{Seq: []string{"getHandType", "action", "actioncard", "getHandTypeMax"}})
	sq = append(sq, cd.Seq{Seq: []string{"breakSet", "actioncard"}})
	sq = append(sq, cd.Seq{Seq: []string{"removeFromHands", "actioncard"}})
	sq = append(sq, cd.Seq{Seq: []string{"PlayAction", "actionCard"}})
	c.Effects.Sequence = sq
	return c
}

func DefArtisan() cd.Card {
	var c cd.Card
	c.Name = "Artisan"
	c.Cost = 6
	c.CTypes.Action = true
	c.Effects.SeqVal = make(map[string]int)
	c.Effects.SeqVal["GainCardTypeMaxVal"] = 5
	c.Effects.SeqVal["discardCount"] = 1
	sq := c.Effects.Sequence
	sq = append(sq, cd.Seq{Seq: []string{"GainCardType", "any", "newCard", "GainCardTypeMaxVal"}})
	sq = append(sq, cd.Seq{Seq: []string{"placeHand", "newCard"}})
	sq = append(sq, cd.Seq{Seq: []string{"selectDiscard", "putdeck", "discardCount"}})
	sq = append(sq, cd.Seq{Seq: []string{"removeFromHands", "putdeck"}})
	sq = append(sq, cd.Seq{Seq: []string{"PlaceDeck", "putdeck"}})
	c.Effects.Sequence = sq
	return c
}

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

func AvailableCards() []cd.Card {

	var rd []cd.Card

	rd = append(rd, DefCellar())
	rd = append(rd, DefChapel())
	rd = append(rd, DefMoat())
	rd = append(rd, DefHarbinger())
	rd = append(rd, DefVillage())
	rd = append(rd, DefWorkshop())
	rd = append(rd, DefWoodcutter())
	rd = append(rd, DefVassal())
	rd = append(rd, DefBureaucrat())
	rd = append(rd, DefRemodel())
	rd = append(rd, DefMoneylender())
	rd = append(rd, DefMilitia())
	rd = append(rd, DefGardens())
	rd = append(rd, DefSmithy())
	rd = append(rd, DefThroneRoom())
	rd = append(rd, DefCouncilRoom())
	rd = append(rd, DefMine())
	rd = append(rd, DefFestival())
	rd = append(rd, DefLaboratory())
	rd = append(rd, DefMarket())
	rd = append(rd, DefWitch())
	rd = append(rd, DefArtisan())

	return rd
}
