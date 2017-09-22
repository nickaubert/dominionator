package cards

// import "fmt"
import "math/rand"
import "time"

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
	DrawCard    int
	ExtraAction int
	ExtraBuy    int
	ExtraCoins  int
	/* http://wiki.dominionstrategy.com/index.php/Gameplay
	   Discard (from hand or from deck) to DiscardPile
	   Gain (to hand or to deck) from Supply
	   Trash (from hand or from deck) to Trash
	*/
}

type Card struct {
	Name    string
	Cost    int
	VP      int
	Coins   int
	CTypes  CType
	Effects Effect
}

type Deck struct {
	Cards []Card
}

type Hand struct {
	Cards []Card
}

type InPlay struct {
	Cards []Card
}

type DiscardPile struct {
	Cards []Card
}

type Trash struct {
	Cards []Card
}

/*
type WholeDeck struct {
	Deck        Deck
	Hand        Hand
	DiscardPile DiscardPile
}
*/

type SupplyPile struct {
	Card  Card
	Count int
}

type Supply struct {
	Piles []SupplyPile
}

func ShuffleDeck(d *Deck) {
	rand.Seed(time.Now().UTC().UnixNano())
	time.Sleep(time.Duration(rand.Int31n(10)) * time.Nanosecond) // more randomness?
	var n Deck
	for _ = range d.Cards {
		r := rand.Intn(len(d.Cards))
		n.Cards = append(n.Cards, d.Cards[r])
		d.Cards = append(d.Cards[:r], d.Cards[r+1:]...)
	}
	d.Cards = n.Cards
}
