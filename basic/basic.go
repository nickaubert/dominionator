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

func DefChapel() cd.Card {
	var c cd.Card
	c.Name = "Chapel"
	c.Cost = 2
	c.CTypes.Action = true
	var s1 cd.Sequence
	s1.CountTrash = 4
	c.Effects.Sequence = append(c.Effects.Sequence, s1)
	return c
}

func DefCellar() cd.Card {
	var c cd.Card
	c.Name = "Cellar"
	c.Cost = 2
	c.CTypes.Action = true
	c.Effects.ExtraActions = 1
	var s1 cd.Sequence
	s1.CountDiscard = 1000 // infinite?
	c.Effects.Sequence = append(c.Effects.Sequence, s1)
	var s2 cd.Sequence
	s2.DrawCount = true
	c.Effects.Sequence = append(c.Effects.Sequence, s2)
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

func DefHarbinger() cd.Card {
	var c cd.Card
	c.Name = "Harbinger"
	c.Cost = 3
	c.CTypes.Action = true
	c.Effects.DrawCard = 1
	c.Effects.ExtraActions = 1
	var s1 cd.Sequence
	s1.RetrieveDiscard = 1
	c.Effects.Sequence = append(c.Effects.Sequence, s1)
	var s2 cd.Sequence
	s2.PlaceDeck = true
	c.Effects.Sequence = append(c.Effects.Sequence, s2)
	return c
}

func DefVassal() cd.Card {
	var c cd.Card
	c.Name = "Vassal"
	c.Cost = 3
	c.CTypes.Action = true
	c.Effects.ExtraCoins = 2
	var s1 cd.Sequence
	s1.DrawDeck = 1
	c.Effects.Sequence = append(c.Effects.Sequence, s1)
	var s2 cd.Sequence
	s2.DiscardNonAction = true
	c.Effects.Sequence = append(c.Effects.Sequence, s2)
	// need decision point here whether to play action
	var s3 cd.Sequence
	s3.PlayAction = 1
	c.Effects.Sequence = append(c.Effects.Sequence, s3)
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

func DefGardens() cd.Card {
	var c cd.Card
	c.Name = "Gardens"
	c.Cost = 4
	c.CTypes.Victory = true
	c.Victories.CardsPerPoint = 10
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
