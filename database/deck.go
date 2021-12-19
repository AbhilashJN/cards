package database

import (
	"context"

	"github.com/AbhilashJN/cards/deck"
)

type DeckModel struct {
	UUID     string    `bson:"uuid"`
	Cards    deck.Deck `bson:"cards"`
	Shuffled bool      `bson:"shuffled"`
}

type DeckCRUDer interface {
	InsertDeck(context.Context, DeckModel) error
}

type DeckCRUDOperator struct {
	Collection MongoCollection
}

func (d *DeckCRUDOperator) InsertDeck(ctx context.Context, deckItem DeckModel) error {
	_, err := d.Collection.InsertOne(ctx, deckItem)
	return err
}
