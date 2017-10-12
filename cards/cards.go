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

/*
type Sequence struct {
	SetVal          SeqVar // set value
	TrashMax        string // max cards to trash
	RetrieveDiscard string // max cards to pull from discard
	PlaceDeck       string // Place cards onto deck
	PlaceDiscards   string // place cardSet in discard
	GetSupplyCard   string // get this card from supply
	// GetHandTypeX    string // get one of this card type from hand
	GetHandType string // get any of this card type from hand
	// DrawCount       string // Draw "X" cards
	DrawDeck        string // max cards to draw from deck
	DiscardNonMatch string // discard non-matching cards from set
	PlayAction      string // play action this many times
	CountCards      string
	ClearSet        string
	GainCard        string
	AddCost         string
	TrashCards      string
	// PlaceHand       string
	PlaceHands     string
	AddXCoins      string // add "X" * int coins
	GetHandMatch   string
	RemoveFromHand string
}
*/

/*
type SeqVar struct {
	Name    string
	Val     int
	Card    Card
	Type    string
	NewName string
}
*/

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
