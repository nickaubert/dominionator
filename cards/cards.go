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
	CountDiscard    int    // max cards to discard
	CountTrash      int    // max cards to trash
	RetrieveDiscard int    // max cards to pull from discard
	DrawDeck        int    // max cards to draw from deck
	PlayAction      int    // play action this many times
	GainMax         int    // gain card up to this cost
	AddXCoins       int    // add "X" * int coins
	GetSupplyCard   Card   // get this card from supply
	MayTrash        Card   // may trash this card from hand
	GetHandType     string // get any of this card type from hand
	DiscardNonMatch string // discard non-matching cards from set
	DrawCount       bool   // Draw "X" cards
	PlaceDeck       bool   // Place cards onto deck
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
