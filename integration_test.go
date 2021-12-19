package main

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/AbhilashJN/cards/api"
	"github.com/AbhilashJN/cards/database"
	"github.com/AbhilashJN/cards/deck"
	"github.com/google/go-cmp/cmp"
	"github.com/julienschmidt/httprouter"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func TestCreateDeckIntegration(t *testing.T) {
	client, _ := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	client.Connect(ctx)
	defer client.Disconnect(ctx)

	dbClient := client.Database("cardsdb_test")
	s := &server{
		router:   httprouter.New(),
		dbClient: dbClient,
	}
	s.initRouter()

	mockBody, _ := json.Marshal(api.CreateDeckRequestBody{Shuffle: false, CustomDeck: false, WantedCards: []string{}})
	req := httptest.NewRequest("POST", "/deck", bytes.NewReader(mockBody))
	response := httptest.NewRecorder()
	s.ServeHTTP(response, req)
	var respBody api.CreateDeckResponseBody
	json.NewDecoder(response.Body).Decode(&respBody)

	if response.Code != http.StatusCreated {
		t.Errorf("Failed create default deck integration test: expected response code %d, got %d", http.StatusCreated, response.Code)
	}
	if len(respBody.DeckId) == 0 {
		t.Errorf("Failed create default deck integration test: expected deck uuid to be valid, got %v", respBody.DeckId)
	}
	if respBody.Shuffled != false {
		t.Errorf("Failed create default deck integration test: expected deck shuffled to be %t, got %t", false, respBody.Shuffled)
	}
	if respBody.Remaining != 52 {
		t.Errorf("Failed create default deck integration test: expected deck size to be %d, got %d", 52, respBody.Remaining)
	}
	s.dbClient.Collection("decks").Drop(ctx)

}

func TestCreateDeckIntegrationErrorCase(t *testing.T) {
	client, _ := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	client.Connect(ctx)
	defer client.Disconnect(ctx)
	dbClient := client.Database("cardsdb_test")
	s := &server{
		router:   httprouter.New(),
		dbClient: dbClient,
	}
	s.initRouter()

	mockBody, _ := json.Marshal(api.CreateDeckRequestBody{Shuffle: false, CustomDeck: true, WantedCards: []string{}})
	req := httptest.NewRequest("POST", "/deck", bytes.NewReader(mockBody))
	response := httptest.NewRecorder()
	s.ServeHTTP(response, req)
	var respBody genericResponseMessage
	json.NewDecoder(response.Body).Decode(&respBody)
	expectedMessage := "List of wanted cards must be provided for custom deck"
	if response.Code != http.StatusBadRequest {
		t.Errorf("Failed create default deck integration test error case: expected response code %d, got %d", http.StatusBadRequest, response.Code)
	}
	if respBody.Message != expectedMessage {
		t.Errorf("Failed create default deck integration test error case: expected response message %s, got %s", expectedMessage, respBody.Message)
	}
	s.dbClient.Collection("decks").Drop(ctx)
}

func TestGetDeckIntegration(t *testing.T) {
	client, _ := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	client.Connect(ctx)
	defer client.Disconnect(ctx)

	dbClient := client.Database("cardsdb_test")
	s := &server{
		router:   httprouter.New(),
		dbClient: dbClient,
	}
	s.initRouter()
	mockDeck, _ := deck.New(&deck.NewDeckOpts{Shuffle: false, CustomDeck: false})
	deckItem := database.DeckModel{
		UUID:     "test-uuid-12345",
		Shuffled: false,
		Cards:    mockDeck,
	}
	s.dbClient.Collection("decks").InsertOne(ctx, deckItem)

	req := httptest.NewRequest("GET", "/deck/test-uuid-12345", bytes.NewReader([]byte{}))
	response := httptest.NewRecorder()
	s.ServeHTTP(response, req)
	var respBody api.GetDeckResponseBody
	json.NewDecoder(response.Body).Decode(&respBody)

	if response.Code != http.StatusOK {
		t.Errorf("Failed get deck integration test: expected response code %d, got %d", http.StatusOK, response.Code)
	}
	if respBody.DeckId != "test-uuid-12345" {
		t.Errorf("Failed get deck integration test: expected deck uuid to be valid, got %v", respBody.DeckId)
	}
	if respBody.Shuffled != false {
		t.Errorf("Failed get deck integration test: expected deck shuffled to be %t, got %t", false, respBody.Shuffled)
	}
	if respBody.Remaining != 52 {
		t.Errorf("Failed get deck integration test: expected deck size to be %d, got %d", 52, respBody.Remaining)
	}
	if len(respBody.Cards) != 52 {
		t.Errorf("Failed get deck integration test: expected %d cards, got %d", 52, len(respBody.Cards))
	}
	s.dbClient.Collection("decks").Drop(ctx)
}

func TestGetDeckIntegrationErrorCase(t *testing.T) {
	client, _ := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	client.Connect(ctx)
	defer client.Disconnect(ctx)

	dbClient := client.Database("cardsdb_test")
	s := &server{
		router:   httprouter.New(),
		dbClient: dbClient,
	}
	s.initRouter()
	mockDeck, _ := deck.New(&deck.NewDeckOpts{Shuffle: false, CustomDeck: false})
	deckItem := database.DeckModel{
		UUID:     "test-uuid-12345",
		Shuffled: false,
		Cards:    mockDeck,
	}
	s.dbClient.Collection("decks").InsertOne(ctx, deckItem)

	req := httptest.NewRequest("GET", "/deck/test-uuid-9876", bytes.NewReader([]byte{}))
	response := httptest.NewRecorder()
	s.ServeHTTP(response, req)
	var respBody genericResponseMessage
	json.NewDecoder(response.Body).Decode(&respBody)
	expectedMessage := "Deck with this id does not exist"
	if response.Code != http.StatusNotFound {
		t.Errorf("Failed get deck integration test error case: expected response code %d, got %d", http.StatusNotFound, response.Code)
	}
	if respBody.Message != expectedMessage {
		t.Errorf("Failed get deck integration test error case: expected response message %s, got %s", expectedMessage, respBody.Message)

	}
	s.dbClient.Collection("decks").Drop(ctx)
}

func TestDrawCardsIntegration(t *testing.T) {
	client, _ := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	client.Connect(ctx)
	defer client.Disconnect(ctx)

	dbClient := client.Database("cardsdb_test")
	s := &server{
		router:   httprouter.New(),
		dbClient: dbClient,
	}
	s.initRouter()
	mockDeck, _ := deck.New(&deck.NewDeckOpts{Shuffle: false, CustomDeck: false})
	deckItem := database.DeckModel{
		UUID:     "test-uuid-12345",
		Shuffled: false,
		Cards:    mockDeck,
	}
	s.dbClient.Collection("decks").InsertOne(ctx, deckItem)
	mockBody, _ := json.Marshal(api.DrawCardsRequestBody{NumberOfCards: 2})
	req := httptest.NewRequest("PATCH", "/deck/test-uuid-12345", bytes.NewReader(mockBody))
	response := httptest.NewRecorder()
	s.ServeHTTP(response, req)
	var respBody api.DrawCardsResponseBody
	json.NewDecoder(response.Body).Decode(&respBody)

	expectedCards := deck.Deck{
		{Value: deck.Ace, Suit: deck.Spades},
		{Value: deck.Two, Suit: deck.Spades},
	}.ToDeckJSON()

	if response.Code != http.StatusOK {
		t.Errorf("Failed draw cards integration test error case: expected response code %d, got %d", http.StatusOK, response.Code)
	}
	if !cmp.Equal(respBody.Cards, expectedCards) {
		t.Errorf("Failed draw cards integration test error case: expected cards %v, got %v", expectedCards, respBody.Cards)
	}
	s.dbClient.Collection("decks").Drop(ctx)
}

func TestDrawCardsIntegrationErrorCase(t *testing.T) {
	client, _ := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	client.Connect(ctx)
	defer client.Disconnect(ctx)

	dbClient := client.Database("cardsdb_test")
	s := &server{
		router:   httprouter.New(),
		dbClient: dbClient,
	}
	s.initRouter()
	mockDeck, _ := deck.New(&deck.NewDeckOpts{Shuffle: false, CustomDeck: false})
	deckItem := database.DeckModel{
		UUID:     "test-uuid-12345",
		Shuffled: false,
		Cards:    mockDeck,
	}
	s.dbClient.Collection("decks").InsertOne(ctx, deckItem)
	mockBody, _ := json.Marshal(api.DrawCardsRequestBody{NumberOfCards: 2})
	req := httptest.NewRequest("PATCH", "/deck/test-uuid-87654", bytes.NewReader(mockBody))
	response := httptest.NewRecorder()
	s.ServeHTTP(response, req)
	var respBody genericResponseMessage
	json.NewDecoder(response.Body).Decode(&respBody)
	expectedMessage := "Deck with this id does not exist"
	if response.Code != http.StatusNotFound {
		t.Errorf("Failed get deck integration test error case: expected response code %d, got %d", http.StatusNotFound, response.Code)
	}
	if respBody.Message != expectedMessage {
		t.Errorf("Failed get deck integration test error case: expected response message %s, got %s", expectedMessage, respBody.Message)

	}
	s.dbClient.Collection("decks").Drop(ctx)
}
