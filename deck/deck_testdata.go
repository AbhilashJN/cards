package deck

import "math/rand"

func getDefaultDeck() Deck {
	return Deck{
		{Value: Ace, Suit: Spades},
		{Value: Two, Suit: Spades},
		{Value: Three, Suit: Spades},
		{Value: Four, Suit: Spades},
		{Value: Five, Suit: Spades},
		{Value: Six, Suit: Spades},
		{Value: Seven, Suit: Spades},
		{Value: Eight, Suit: Spades},
		{Value: Nine, Suit: Spades},
		{Value: Ten, Suit: Spades},
		{Value: Jack, Suit: Spades},
		{Value: Queen, Suit: Spades},
		{Value: King, Suit: Spades},
		{Value: Ace, Suit: Diamonds},
		{Value: Two, Suit: Diamonds},
		{Value: Three, Suit: Diamonds},
		{Value: Four, Suit: Diamonds},
		{Value: Five, Suit: Diamonds},
		{Value: Six, Suit: Diamonds},
		{Value: Seven, Suit: Diamonds},
		{Value: Eight, Suit: Diamonds},
		{Value: Nine, Suit: Diamonds},
		{Value: Ten, Suit: Diamonds},
		{Value: Jack, Suit: Diamonds},
		{Value: Queen, Suit: Diamonds},
		{Value: King, Suit: Diamonds},
		{Value: Ace, Suit: Clubs},
		{Value: Two, Suit: Clubs},
		{Value: Three, Suit: Clubs},
		{Value: Four, Suit: Clubs},
		{Value: Five, Suit: Clubs},
		{Value: Six, Suit: Clubs},
		{Value: Seven, Suit: Clubs},
		{Value: Eight, Suit: Clubs},
		{Value: Nine, Suit: Clubs},
		{Value: Ten, Suit: Clubs},
		{Value: Jack, Suit: Clubs},
		{Value: Queen, Suit: Clubs},
		{Value: King, Suit: Clubs},
		{Value: Ace, Suit: Hearts},
		{Value: Two, Suit: Hearts},
		{Value: Three, Suit: Hearts},
		{Value: Four, Suit: Hearts},
		{Value: Five, Suit: Hearts},
		{Value: Six, Suit: Hearts},
		{Value: Seven, Suit: Hearts},
		{Value: Eight, Suit: Hearts},
		{Value: Nine, Suit: Hearts},
		{Value: Ten, Suit: Hearts},
		{Value: Jack, Suit: Hearts},
		{Value: Queen, Suit: Hearts},
		{Value: King, Suit: Hearts},
	}
}

func getShuffledDefaultDeck() Deck {
	rand.Seed(123)
	d := getDefaultDeck()
	d.Shuffle()
	rand.Seed(1)
	return d
}
