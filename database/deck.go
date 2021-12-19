package database

import (
	"github.com/AbhilashJN/cards/deck"
)

type DeckModel struct {
	UUID     string    `bson:"uuid"`
	Cards    deck.Deck `bson:"cards"`
	Shuffled bool      `bson:"shuffled"`
}
