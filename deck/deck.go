package deck

import (
	"math/rand"
)

type Deck []Card

type DeckJSON []CardJSON

type NewDeckOpts struct {
	Shuffle         bool
	CustomDeck      bool
	CustomDeckCards []string
}

type ErrDrawCardsSizeExceeded struct {
}

func (e ErrDrawCardsSizeExceeded) Error() string {
	return "requested number of cards is greater than the cards remaining in the deck"
}

func (d Deck) ToDeckJSON() DeckJSON {
	deckJSON := make(DeckJSON, len(d))
	for i, card := range d {
		deckJSON[i] = card.ToCardJSON()
	}
	return deckJSON
}

func (d Deck) Shuffle() {
	rand.Shuffle(len(d), func(i, j int) {
		d[i], d[j] = d[j], d[i]
	})
}

func DrawCards(d Deck, n int) (Deck, Deck, error) {
	size := len(d)
	if n > size {
		return nil, nil, ErrDrawCardsSizeExceeded{}
	}
	draw, remaining := d[:n], d[n:]
	return draw, remaining, nil
}

func defaultDeckGenerator() Deck {
	newDeck := make(Deck, 52)
	idx := 0
	for i := Spades; i <= Hearts; i++ {
		for j := Ace; j <= King; j++ {
			newDeck[idx] = Card{Suit: i, Value: j}
			idx++
		}
	}
	return newDeck
}

func customDeckGenerator(cardCodesList []string) Deck {
	numCards := len(cardCodesList)
	newDeck := make(Deck, numCards)
	for i, cardCode := range cardCodesList {
		value, suit := DecodeValueAndSuit(cardCode)
		newDeck[i] = Card{Suit: suit, Value: value}
	}
	return newDeck
}

func New(opts *NewDeckOpts) Deck {
	var deck Deck
	if opts.CustomDeck {
		deck = customDeckGenerator(opts.CustomDeckCards)
	} else {
		deck = defaultDeckGenerator()
	}

	if opts.Shuffle {
		deck.Shuffle()
	}

	return deck
}
