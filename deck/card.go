package deck

import (
	"fmt"
	"strconv"
	"strings"
)

type CardValue int8
type CardSuit int8

//go:generate stringer -type=CardValue,CardSuit -linecomment

const (
	Spades CardSuit = iota
	Diamonds
	Clubs
	Hearts
)

const (
	Ace   CardValue = iota + 1 // Ace
	Two                        // 2
	Three                      // 3
	Four                       // 4
	Five                       // 5
	Six                        // 6
	Seven                      // 7
	Eight                      // 8
	Nine                       // 9
	Ten                        // 10
	Jack                       // Jack
	Queen                      // Queen
	King                       // King
)

type Card struct {
	Value CardValue `json:"value" bson:"value"`
	Suit  CardSuit  `json:"suit" bson:"suit"`
}

type CardJSON struct {
	Value string `json:"value"`
	Suit  string `json:"suit"`
	Code  string `json:"code"`
}

type ErrInvalidCardCode struct {
	CardCode string
}

func (e ErrInvalidCardCode) Error() string {
	return fmt.Sprintf("Card code %s is invalid", e.CardCode)
}

func (c Card) ToCardJSON() CardJSON {
	valueString := strings.ToUpper(c.Value.String())
	suitString := strings.ToUpper(c.Suit.String())
	code := []byte{valueString[0], suitString[0]}

	return CardJSON{
		Value: valueString,
		Suit:  suitString,
		Code:  string(code),
	}
}

func DecodeValueAndSuit(code string) (CardValue, CardSuit, error) {
	var (
		value    CardValue
		suit     CardSuit
		parseErr error
	)
	valueCode := string(code[0])
	suitCode := string(code[1])

	valueInt, err := strconv.Atoi(valueCode)
	if err != nil {
		switch valueCode {
		case "A":
			value = Ace
		case "J":
			value = Jack
		case "Q":
			value = Queen
		case "K":
			value = King
		default:
			parseErr = ErrInvalidCardCode{CardCode: code}
		}
	} else {
		value = CardValue(valueInt)
	}

	switch suitCode {
	case "S":
		suit = Spades
	case "D":
		suit = Diamonds
	case "C":
		suit = Clubs
	case "H":
		suit = Hearts
	default:
		parseErr = ErrInvalidCardCode{CardCode: code}
	}
	return value, suit, parseErr
}
