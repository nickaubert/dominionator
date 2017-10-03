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
	Sequence     []Sequence
	/* http://wiki.dominionstrategy.com/index.php/Gameplay
	   Discard (from hand or from deck) to Discard
	   Gain (to hand or to deck) from Supply
	   Trash (from hand or from deck) to Trash
	*/
}

type Sequence struct {
	SetVal SeqVar // set value
	// CountDiscard    string // max cards to discard
	// DrawCount       string // Draw "X" cards
	TrashMax        string // max cards to trash
	RetrieveDiscard string // max cards to pull from discard
	PlaceDeck       string // Place cards onto deck
	// PlaceDiscard    string // place one card in discard
	PlaceDiscards   string // place cardSet in discard
	GetSupplyCard   string // get this card from supply
	GetHandTypeX    string // get one of this card type from hand
	GetHandType     string // get any of this card type from hand
	DrawDeck        string // max cards to draw from deck
	DiscardNonMatch string // discard non-matching cards from set
	PlayAction      string // play action this many times
	CountCards      string
	ClearSet        string
	GainCard        string
	AddCost         string
	TrashCards      string
	PlaceHand       string
	PlaceHands      string
	AddXCoins       string // add "X" * int coins
	GetHandMatch    string
	RemoveFromHand  string
}

type SeqVar struct {
	Name    string
	Val     int
	Card    Card
	Type    string
	NewName string
}

type Attack struct {
	DiscardTo int
	GainCurse int
	Sequence  []Sequence
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
