package deck

import (
	"math/rand"
	"testing"

	"github.com/google/go-cmp/cmp"
)

type CustomDeckGeneratorTest struct {
	input          []string
	expectedOutput Deck
}

type ShuffleTest struct {
	input          Deck
	expectedOutput Deck
}

type DrawCardsTest struct {
	inputDeck             Deck
	n                     int
	expectedDrawnDeck     Deck
	expectedRemainingDeck Deck
	expectedError         error
}

type ToDeckJSONTest struct {
	inputDeck        Deck
	expectedDeckJSON DeckJSON
}

type NewDeckTest struct {
	inputOpts      NewDeckOpts
	expectedOutput Deck
}

func TestDefaultDeckGenerator(t *testing.T) {
	expectedDeck := getDefaultDeck()
	deck := defaultDeckGenerator()

	if !cmp.Equal(deck, expectedDeck) {
		t.Errorf("Incorrect default deck generated, got %v", deck)
	}
}

func TestCustomDeckGenerator(t *testing.T) {
	tests := []CustomDeckGeneratorTest{
		{[]string{},
			Deck{},
		},
		{[]string{"KC"},
			Deck{
				{Value: King, Suit: Clubs},
			},
		},
		{[]string{"8C", "5H", "QS", "AD", "7S", "3C"},
			Deck{
				{Value: Eight, Suit: Clubs},
				{Value: Five, Suit: Hearts},
				{Value: Queen, Suit: Spades},
				{Value: Ace, Suit: Diamonds},
				{Value: Seven, Suit: Spades},
				{Value: Three, Suit: Clubs},
			},
		},
	}

	for _, test := range tests {
		output := customDeckGenerator(test.input)
		if !cmp.Equal(output, test.expectedOutput) {
			t.Errorf("Failed for input %v: expected %v, got %v", test.input, test.expectedOutput, output)
		}
	}

}

func TestShuffle(t *testing.T) {
	tests := []ShuffleTest{
		{Deck{
			{Value: Eight, Suit: Clubs},
			{Value: Five, Suit: Hearts},
			{Value: Queen, Suit: Spades},
			{Value: Ace, Suit: Diamonds},
			{Value: Seven, Suit: Spades},
			{Value: Three, Suit: Clubs},
		},
			Deck{
				{Value: Queen, Suit: Spades},
				{Value: Three, Suit: Clubs},
				{Value: Five, Suit: Hearts},
				{Value: Seven, Suit: Spades},
				{Value: Eight, Suit: Clubs},
				{Value: Ace, Suit: Diamonds},
			},
		},
	}
	rand.Seed(123)
	for _, test := range tests {
		test.input.Shuffle()
		if !cmp.Equal(test.input, test.expectedOutput) {
			t.Errorf("Failed: expected %v, got %v", test.expectedOutput, test.input)
		}
	}
	rand.Seed(1)
}

func TestDrawCards(t *testing.T) {
	tests := []DrawCardsTest{
		{
			Deck{
				{Value: Queen, Suit: Spades},
				{Value: Three, Suit: Clubs},
				{Value: Five, Suit: Hearts},
				{Value: Seven, Suit: Spades},
				{Value: Eight, Suit: Clubs},
				{Value: Ace, Suit: Diamonds},
			},
			0,
			Deck{},
			Deck{
				{Value: Queen, Suit: Spades},
				{Value: Three, Suit: Clubs},
				{Value: Five, Suit: Hearts},
				{Value: Seven, Suit: Spades},
				{Value: Eight, Suit: Clubs},
				{Value: Ace, Suit: Diamonds},
			}, nil,
		},
		{
			Deck{
				{Value: Queen, Suit: Spades},
				{Value: Three, Suit: Clubs},
				{Value: Five, Suit: Hearts},
				{Value: Seven, Suit: Spades},
				{Value: Eight, Suit: Clubs},
				{Value: Ace, Suit: Diamonds},
			},
			2,
			Deck{
				{Value: Queen, Suit: Spades},
				{Value: Three, Suit: Clubs},
			},
			Deck{
				{Value: Five, Suit: Hearts},
				{Value: Seven, Suit: Spades},
				{Value: Eight, Suit: Clubs},
				{Value: Ace, Suit: Diamonds},
			}, nil,
		},
		{
			Deck{
				{Value: Queen, Suit: Spades},
				{Value: Three, Suit: Clubs},
				{Value: Five, Suit: Hearts},
				{Value: Seven, Suit: Spades},
				{Value: Eight, Suit: Clubs},
				{Value: Ace, Suit: Diamonds},
			},
			6,
			Deck{
				{Value: Queen, Suit: Spades},
				{Value: Three, Suit: Clubs},
				{Value: Five, Suit: Hearts},
				{Value: Seven, Suit: Spades},
				{Value: Eight, Suit: Clubs},
				{Value: Ace, Suit: Diamonds},
			},
			Deck{}, nil,
		},
		{
			Deck{
				{Value: Queen, Suit: Spades},
				{Value: Three, Suit: Clubs},
				{Value: Five, Suit: Hearts},
				{Value: Seven, Suit: Spades},
				{Value: Eight, Suit: Clubs},
				{Value: Ace, Suit: Diamonds},
			},
			7,
			nil,
			nil,
			ErrDrawCardsSizeExceeded{},
		},
	}

	for _, test := range tests {
		drawnDeck, remainingDeck, err := DrawCards(test.inputDeck, test.n)

		if !cmp.Equal(drawnDeck, test.expectedDrawnDeck) {
			t.Errorf("Failed for input %v %v: Expected drawn deck to be %v, got %v", test.inputDeck, test.n, test.expectedDrawnDeck, drawnDeck)
		}
		if !cmp.Equal(remainingDeck, test.expectedRemainingDeck) {
			t.Errorf("Failed for input %v %v: Expected remaining deck to be %v, got %v", test.inputDeck, test.n, test.expectedRemainingDeck, remainingDeck)
		}
		if err != test.expectedError {
			t.Errorf("Failed for input %v %v: Expected error to be %v, got %v", test.inputDeck, test.n, test.expectedError, err)
		}

	}

}

func TestToDeckJSON(t *testing.T) {
	tests := []ToDeckJSONTest{
		{
			Deck{},
			DeckJSON{},
		},
		{
			Deck{
				{Value: Queen, Suit: Spades},
				{Value: Three, Suit: Clubs},
				{Value: Five, Suit: Hearts},
			},
			DeckJSON{
				{Value: "QUEEN", Suit: "SPADES", Code: "QS"},
				{Value: "3", Suit: "CLUBS", Code: "3C"},
				{Value: "5", Suit: "HEARTS", Code: "5H"},
			},
		},
	}

	for _, test := range tests {
		output := test.inputDeck.ToDeckJSON()
		if !cmp.Equal(output, test.expectedDeckJSON) {
			t.Errorf("Failed for input %v: Expected deck json to be %v, got %v", test.inputDeck, test.expectedDeckJSON, output)
		}
	}
}

func TestNewDeck(t *testing.T) {
	tests := []NewDeckTest{
		{
			NewDeckOpts{Shuffle: false, CustomDeck: false, CustomDeckCards: []string{}},
			getDefaultDeck(),
		},
		{
			NewDeckOpts{Shuffle: true, CustomDeck: false, CustomDeckCards: []string{}},
			getShuffledDefaultDeck(),
		},
		{
			NewDeckOpts{Shuffle: false, CustomDeck: true, CustomDeckCards: []string{"AS", "QS", "2H", "7D", "4C"}},
			Deck{
				{Value: Ace, Suit: Spades},
				{Value: Queen, Suit: Spades},
				{Value: Two, Suit: Hearts},
				{Value: Seven, Suit: Diamonds},
				{Value: Four, Suit: Clubs},
			},
		},
		{
			NewDeckOpts{Shuffle: true, CustomDeck: true, CustomDeckCards: []string{"AS", "QS", "2H", "7D", "4C"}},
			Deck{
				{Value: Queen, Suit: Spades},
				{Value: Ace, Suit: Spades},
				{Value: Seven, Suit: Diamonds},
				{Value: Four, Suit: Clubs},
				{Value: Two, Suit: Hearts},
			},
		},
		{
			NewDeckOpts{Shuffle: false, CustomDeck: false, CustomDeckCards: []string{"AS", "QS", "2H", "7D", "4C"}},
			getDefaultDeck(),
		},
	}
	rand.Seed(123)
	for _, test := range tests {
		output := New(&test.inputOpts)
		if !cmp.Equal(output, test.expectedOutput) {
			t.Errorf("Failed for input %v: expected %v, got %v", test.inputOpts, test.expectedOutput, output)
		}
	}
	rand.Seed(1)
}
