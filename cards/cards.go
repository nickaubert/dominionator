package cards

type Card struct {
	Name      string
	Cost      int
	VP        int
	Coins     int
	CTypes    CType
	Effects   Effect
	Attacks   Attack
	Reactions Reaction
	Victories Victory
}

type CType struct {
	Action   bool
	Treasure bool
	Victory  bool
	Curse    bool
	Attack   bool
	Duration bool // Seaside
	Reaction bool
}

type Effect struct {
	DrawCard     int
	ExtraActions int
	ExtraBuys    int
	ExtraCoins   int
	Sequence     []Seq
	SeqVal       map[string]int
	/* http://wiki.dominionstrategy.com/index.php/Gameplay
	   Discard (from hand or from deck) to Discard
	   Gain (to hand or to deck) from Supply
	   Trash (from hand or from deck) to Trash
	*/
}

/*
   New logic: map of sequence slice
        need to split setup sequence and sequence range loop?
        most cards use card name as key to sequence map
            sq := c.Effects.Sequence[c.Name]
        Library:
            sq := c.Effects.Sequence[c.Name]
            c.Effects.SeqVal["wantCards"] = 7
            c.Effects.SeqVal["drawDeckMax"] = 1
            c.Effects.SeqVal["getCardTypeMax"] = 0
            "Library":          "checkHandCount", "wantCards", "getnonactioncard"
            "getnonactioncard": "drawDeck", "drawDeckMax", "newCard"
            "getnonactioncard": "getCardType", "newCard", "action", "actionCard", "getCardTypeMax"
            "getnonactioncard": "removeCards", "actionCard", "newCard"
            "getnonactioncard": "placeDiscards", "actionCard"
            "getnonactioncard": "placeHand", "newCard"

*/

type Seq struct {
	Seq []string
}

type Attack struct {
	DiscardTo int
	GainCurse int
	Sequence  []Seq
	SeqVal    map[string]int
}

type Reaction struct {
	Defend bool
}

type Victory struct {
	CardsPerPoint int
}

type SupplyPile struct {
	Card  Card
	Count int
}

type Supply struct {
	Piles []SupplyPile
}
