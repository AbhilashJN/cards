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
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type CreateDeckRequestBody struct {
	Shuffle     bool     `json:"shuffle"`
	CustomDeck  bool     `json:"customDeck"`
	WantedCards []string `json:"wantedCards"`
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
	NumberOfCards int `json:"numberOfCards"`
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
	if reqBody.CustomDeck && len(reqBody.WantedCards) == 0 {
		return responseBody, http.StatusBadRequest, ApiError{Message: "List of wanted cards must be provided for custom deck"}
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

func HandleDrawCards(r *http.Request, ps httprouter.Params, dc database.DeckCRUDer, ctx context.Context) (DrawCardsResponseBody, int, error) {
	var (
		reqBody      DrawCardsRequestBody
		responseBody DrawCardsResponseBody
		resultDeck   db.DeckModel
	)
	reqUUID := ps.ByName("uuid")
	err := json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		log.Println("Error parsing request body", err)
		return responseBody, http.StatusBadRequest, ApiError{Message: "Request body is malformed"}
	}

	resultDeck, err = dc.FindDeckByUUID(ctx, reqUUID)
	if err == mongo.ErrNoDocuments {
		return responseBody, http.StatusNotFound, ApiError{Message: "Deck with this id does not exist"}
	} else if err != nil {
		log.Println("Error occurred while searching for document in db.", err)
		return responseBody, http.StatusInternalServerError, ApiError{Message: "Internal Server Error"}
	}

	drawnCards, remainingCards, err := deck.DrawCards(resultDeck.Cards, reqBody.NumberOfCards)
	if err != nil {
		return responseBody, http.StatusBadRequest, ApiError{Message: err.Error()}
	}
	updateQuery := bson.D{{Key: "$set", Value: bson.D{{Key: "cards", Value: remainingCards}}}}
	err = dc.UpdateDeckByUUID(ctx, reqUUID, updateQuery)
	if err != nil {
		log.Println("Error occurred while updating document in db.", err)
		return responseBody, http.StatusInternalServerError, ApiError{Message: "Internal Server Error"}
	}

	responseBody.Cards = drawnCards.ToDeckJSON()
	return responseBody, http.StatusOK, nil
}
