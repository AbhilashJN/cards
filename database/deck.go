package database

import (
	"context"

	"github.com/AbhilashJN/cards/deck"
	"go.mongodb.org/mongo-driver/bson"
)

type DeckModel struct {
	UUID     string    `bson:"uuid"`
	Cards    deck.Deck `bson:"cards"`
	Shuffled bool      `bson:"shuffled"`
}

type DeckCRUDer interface {
	InsertDeck(context.Context, DeckModel) error
	FindDeckByUUID(context.Context, string) (DeckModel, error)
}

type DeckCRUDOperator struct {
	Collection MongoCollection
}

func (d *DeckCRUDOperator) InsertDeck(ctx context.Context, deckItem DeckModel) error {
	_, err := d.Collection.InsertOne(ctx, deckItem)
	return err
}

func (d *DeckCRUDOperator) FindDeckByUUID(ctx context.Context, uuid string) (DeckModel, error) {
	var resultDeck DeckModel
	filterByUUID := bson.D{{Key: "uuid", Value: uuid}}
	err := d.Collection.FindOne(ctx, filterByUUID).Decode(&resultDeck)
	return resultDeck, err
}
