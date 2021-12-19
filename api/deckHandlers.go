package api

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/AbhilashJN/cards/database"
	db "github.com/AbhilashJN/cards/database"
	"github.com/AbhilashJN/cards/deck"
	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter"
	"go.mongodb.org/mongo-driver/mongo"
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

func HandleCreateDeck(r *http.Request, ps httprouter.Params, dc database.DeckCRUDer, ctx context.Context) (CreateDeckResponseBody, int, error) {
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

	err = dc.InsertDeck(ctx, deckItem)
	if err != nil {
		log.Println("Error occurred while inserting document into db.", err)
		return responseBody, http.StatusInternalServerError, ApiError{Message: "Internal Server Error"}
	}

	responseBody.DeckId = deckId
	responseBody.Shuffled = reqBody.Shuffle
	responseBody.Remaining = len(cards)
	return responseBody, http.StatusCreated, nil
}

func HandleGetDeck(r *http.Request, ps httprouter.Params, dc database.DeckCRUDer, ctx context.Context) (GetDeckResponseBody, int, error) {
	var (
		responseBody GetDeckResponseBody
		resultDeck   db.DeckModel
	)
	reqUUID := ps.ByName("uuid")
	resultDeck, err := dc.FindDeckByUUID(ctx, reqUUID)
	if err == mongo.ErrNoDocuments {
		return responseBody, http.StatusNotFound, ApiError{Message: "Deck with this id does not exist"}
	} else if err != nil {
		log.Println("Error occurred while searching for document in db.", err)
		return responseBody, http.StatusInternalServerError, ApiError{Message: "Internal Server Error"}
	}

	responseBody.DeckId = resultDeck.UUID
	responseBody.Shuffled = resultDeck.Shuffled
	responseBody.Remaining = len(resultDeck.Cards)
	responseBody.Cards = resultDeck.Cards.ToDeckJSON()
	return responseBody, http.StatusOK, nil
}
