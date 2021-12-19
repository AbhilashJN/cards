package deck

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

type DecodeValueAndSuitTest struct {
	input         string
	expectedValue CardValue
	expectedSuit  CardSuit
	expectedErr   error
}

type CardValueStringerTest struct {
	input          CardValue
	expectedOutput string
}

type CardSuitStringerTest struct {
	input          CardSuit
	expectedOutput string
}

type ToCardJSONTest struct {
	card             Card
	expectedCardJSON CardJSON
}

func TestTableDecodeValueAndSuit(t *testing.T) {
	var tests = []DecodeValueAndSuitTest{
		{"AS", Ace, Spades, nil},
		{"3C", Three, Clubs, nil},
		{"8H", Eight, Hearts, nil},
		{"QD", Queen, Diamonds, nil},
		{"0M", CardValue(0), CardSuit(0), ErrInvalidCardCode{CardCode: "0M"}},
	}

	for _, test := range tests {
		value, suit, err := DecodeValueAndSuit(test.input)
		if value != test.expectedValue {
			t.Errorf("Failed for input '%s': Expected value to be %v, got %v", test.input, test.expectedValue, value)
		}
		if suit != test.expectedSuit {
			t.Errorf("Failed for input '%s': Expected suit to be %v, got %v", test.input, test.expectedSuit, suit)
		}
		if !cmp.Equal(err, test.expectedErr) {
			t.Errorf("Failed for input '%s': Expected error to be %v, got %v", test.input, test.expectedErr, err)

		}
	}
}

func TestCardValueStringer(t *testing.T) {
	tests := []CardValueStringerTest{
		{Ace, "Ace"},
		{Five, "5"},
		{Eight, "8"},
		{King, "King"},
	}

	for _, test := range tests {
		if output := test.input.String(); output != test.expectedOutput {
			t.Errorf("Failed for input %v: Expected card value string to be '%v', got '%v'", test.input, test.expectedOutput, output)
		}
	}
}

func TestCardSuitStringer(t *testing.T) {
	tests := []CardSuitStringerTest{
		{Spades, "Spades"},
		{Diamonds, "Diamonds"},
		{Clubs, "Clubs"},
		{Hearts, "Hearts"},
	}

	for _, test := range tests {
		if output := test.input.String(); output != test.expectedOutput {
			t.Errorf("Failed for input %v: Expected card suit string to be '%v', got '%v'", test.input, test.expectedOutput, output)
		}
	}
}

func TestTableToCardJSON(t *testing.T) {
	var tests = []ToCardJSONTest{
		{Card{Value: Six, Suit: Clubs}, CardJSON{Value: "6", Suit: "CLUBS", Code: "6C"}},
		{Card{Value: Jack, Suit: Hearts}, CardJSON{Value: "JACK", Suit: "HEARTS", Code: "JH"}},
		{Card{Value: Ace, Suit: Spades}, CardJSON{Value: "ACE", Suit: "SPADES", Code: "AS"}},
	}

	for _, test := range tests {
		output := test.card.ToCardJSON()
		if !cmp.Equal(output, test.expectedCardJSON) {
			t.Errorf("Failed for input %v: Expected card json to be %v, got %v", test.card, test.expectedCardJSON, output)
		}
	}
}
