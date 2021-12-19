package api

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	db "github.com/AbhilashJN/cards/database"
	"github.com/AbhilashJN/cards/deck"
	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter"
)

type CreateDeckRequestBody struct {
	Shuffle     bool
	CustomDeck  bool
	WantedCards []string
}

type CreateDeckResponseBody struct {
	DeckId    string `json:"deck_id"`
	Shuffled  bool   `json:"shuffled"`
	Remaining int    `json:"remaining"`
}

type GetDeckResponseBody struct {
	DeckId    string        `json:"deckId"`
	Shuffled  bool          `json:"shuffled"`
	Remaining int           `json:"remaining"`
	Cards     deck.DeckJSON `json:"cards"`
}

type DrawCardsRequestBody struct {
	DeckUUID      string
	NumberOfCards int
}

type DrawCardsResponseBody struct {
	Cards []deck.CardJSON `json:"cards"`
}

type ApiError struct {
	Message string
}

func (e ApiError) Error() string {
	return e.Message
}

func HandleCreateDeck(r *http.Request, ps httprouter.Params, collection db.MongoCollection, ctx context.Context) (CreateDeckResponseBody, int, error) {
	var (
		reqBody      CreateDeckRequestBody
		responseBody CreateDeckResponseBody
	)
	err := json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		log.Println("Error parsing request body", err)
		return responseBody, http.StatusBadRequest, ApiError{Message: "Request body is malformed"}
	}

	deckId := uuid.NewString()
	cards := deck.New(&deck.NewDeckOpts{
		Shuffle:         reqBody.Shuffle,
		CustomDeck:      reqBody.CustomDeck,
		CustomDeckCards: reqBody.WantedCards,
	})
	deckItem := db.DeckModel{UUID: deckId, Cards: cards}

	_, err = collection.InsertOne(ctx, deckItem)
	if err != nil {
		log.Println("Error occurred while inserting document into db.", err)
		return responseBody, http.StatusInternalServerError, ApiError{Message: "Internal Server Error"}
	}

	responseBody.DeckId = deckId
	responseBody.Shuffled = reqBody.Shuffle
	responseBody.Remaining = len(cards)
	return responseBody, http.StatusCreated, nil
}