package api

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/AbhilashJN/cards/database"
	"github.com/AbhilashJN/cards/deck"
	"github.com/google/go-cmp/cmp"
	"github.com/julienschmidt/httprouter"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type mockDeckCRUDOperator struct {
	mockInsertDeckFn     func(context.Context, database.DeckModel) error
	mockFindDeckByUUID   func(ctx context.Context, uuid string) (database.DeckModel, error)
	mockUpdateDeckByUUID func(context.Context, string, bson.D) error
}

func (d *mockDeckCRUDOperator) InsertDeck(ctx context.Context, deckItem database.DeckModel) error {
	return d.mockInsertDeckFn(ctx, deckItem)
}

func (d *mockDeckCRUDOperator) FindDeckByUUID(ctx context.Context, uuid string) (database.DeckModel, error) {
	return d.mockFindDeckByUUID(ctx, uuid)
}

func (d *mockDeckCRUDOperator) UpdateDeckByUUID(ctx context.Context, uuid string, updateQuery bson.D) error {
	return d.mockUpdateDeckByUUID(ctx, uuid, updateQuery)
}

type HandleCreateDeckTest struct {
	shuffle          bool
	customDeck       bool
	wantedCards      []string
	expectedNumCards int
}

var mdc mockDeckCRUDOperator

func TestHandleCreateDeck(t *testing.T) {
	mockParams := httprouter.Params{}
	mockCtx := context.TODO()
	mdc.mockInsertDeckFn = func(ctx context.Context, d database.DeckModel) error {
		return nil
	}

	tests := []HandleCreateDeckTest{
		{shuffle: false, customDeck: false, wantedCards: []string{}, expectedNumCards: 52},
		{shuffle: true, customDeck: false, wantedCards: []string{}, expectedNumCards: 52},
		{shuffle: false, customDeck: true, wantedCards: []string{"AS", "QS", "2H", "7D", "4C"}, expectedNumCards: 5},
		{shuffle: true, customDeck: true, wantedCards: []string{"AS", "QS", "2H", "7D", "4C"}, expectedNumCards: 5},
	}

	for _, test := range tests {
		mockBody, _ := json.Marshal(CreateDeckRequestBody{Shuffle: test.shuffle, CustomDeck: test.customDeck, WantedCards: test.wantedCards})
		req := httptest.NewRequest("POST", "/deck", bytes.NewReader(mockBody))
		responseBody, responseCode, err := HandleCreateDeck(req, mockParams, &mdc, mockCtx)
		if err != nil {
			t.Errorf("Failed for success case: expected err to be %v, got %v", nil, err)
		}
		if len(responseBody.DeckId) == 0 {
			t.Errorf("Failed for success case: expected deck id to be valid uuid, got %v", responseBody.DeckId)
		}
		if responseBody.Shuffled != test.shuffle {
			t.Errorf("Failed for success case: expected shuffled to be %t, got %t", test.shuffle, responseBody.Shuffled)
		}
		if responseBody.Remaining != test.expectedNumCards {
			t.Errorf("Failed for success case: expected remaining cards to be %d, got %d", test.expectedNumCards, responseBody.Remaining)
		}
		if responseCode != http.StatusCreated {
			t.Errorf("Failed for success case: expected status code to be %d, got %d", http.StatusCreated, responseCode)
		}
	}

}

func TestHandleCreateDeckNoWantedCardsErr(t *testing.T) {
	mockParams := httprouter.Params{}
	mockCtx := context.TODO()
	mdc.mockInsertDeckFn = func(ctx context.Context, d database.DeckModel) error {
		return nil
	}
	mockBody, _ := json.Marshal(CreateDeckRequestBody{Shuffle: false, CustomDeck: true, WantedCards: []string{}})
	req := httptest.NewRequest("POST", "/deck", bytes.NewReader(mockBody))
	expectedErr := ApiError{Message: "List of wanted cards must be provided for custom deck"}
	_, responseCode, err := HandleCreateDeck(req, mockParams, &mdc, mockCtx)
	if !cmp.Equal(err, expectedErr) {
		t.Errorf("Failed for no wanted cards given error case: expected error to be %v, got %v", expectedErr, err)
	}
	if responseCode != http.StatusBadRequest {
		t.Errorf("Failed for no wanted cards given error case: expected response code to be %d, got %d", http.StatusBadRequest, responseCode)
	}
}

func TestHandleCreateDeckDbError(t *testing.T) {
	mockParams := httprouter.Params{}
	mockCtx := context.TODO()
	mdc := mockDeckCRUDOperator{}
	mdc.mockInsertDeckFn = func(ctx context.Context, d database.DeckModel) error {
		return errors.New("Test error 123")
	}

	mockBody, _ := json.Marshal(CreateDeckRequestBody{Shuffle: false, CustomDeck: false, WantedCards: []string{}})
	req := httptest.NewRequest("POST", "/deck", bytes.NewReader(mockBody))
	expectedErr := ApiError{Message: "Internal Server Error"}
	_, responseCode, err := HandleCreateDeck(req, mockParams, &mdc, mockCtx)
	if !cmp.Equal(err, expectedErr) {
		t.Errorf("Failed for db write error case: expected error to be %v, got %v", expectedErr, err)
	}
	if responseCode != http.StatusInternalServerError {
		t.Errorf("Failed for db write error case: expected response code to be %d, got %d", http.StatusInternalServerError, responseCode)
	}
}

func TestHandleCreateDeckBadRequest(t *testing.T) {
	mockParams := httprouter.Params{}
	mockCtx := context.TODO()
	mdc := mockDeckCRUDOperator{}
	mdc.mockInsertDeckFn = func(ctx context.Context, d database.DeckModel) error {
		return errors.New("Test error 123")
	}

	mockBody, _ := json.Marshal(CreateDeckRequestBody{Shuffle: false, CustomDeck: false, WantedCards: []string{}})
	req := httptest.NewRequest("POST", "/deck", bytes.NewReader(mockBody[:len(mockBody)-2]))
	expectedErr := ApiError{Message: "Request body is malformed"}
	_, responseCode, err := HandleCreateDeck(req, mockParams, &mdc, mockCtx)
	if !cmp.Equal(err, expectedErr) {
		t.Errorf("Failed for bad request case: expected error to be %v, got %v", expectedErr, err)
	}
	if responseCode != http.StatusBadRequest {
		t.Errorf("Failed for bad request case: expected response code to be %d, got %d", http.StatusBadRequest, responseCode)
	}
}

func TestHandleGetDeck(t *testing.T) {
	mockParams := httprouter.Params{{Key: "uuid", Value: "test-uuid-123"}}
	mockCtx := context.TODO()
	mockResult := database.DeckModel{
		UUID:     "test-uuid-123",
		Shuffled: false,
		Cards:    deck.Deck{{Value: deck.Ace, Suit: deck.Spades}, {Value: deck.Three, Suit: deck.Clubs}},
	}
	mdc := mockDeckCRUDOperator{}
	mdc.mockFindDeckByUUID = func(ctx context.Context, uuid string) (database.DeckModel, error) {
		return mockResult, nil
	}

	req := httptest.NewRequest("GET", "/deck/test-uuid-123", bytes.NewReader([]byte{}))
	expectedResponse := GetDeckResponseBody{
		DeckId:    "test-uuid-123",
		Shuffled:  false,
		Remaining: 2,
		Cards:     deck.Deck{{Value: deck.Ace, Suit: deck.Spades}, {Value: deck.Three, Suit: deck.Clubs}}.ToDeckJSON(),
	}
	response, responseCode, err := HandleGetDeck(req, mockParams, &mdc, mockCtx)
	if !cmp.Equal(response, expectedResponse) {
		t.Errorf("Failed for success case: expected respnonse to be %v, got %v", expectedResponse, response)
	}
	if responseCode != http.StatusOK {
		t.Errorf("Failed for success case: expected response code to be %d, got %d", http.StatusOK, responseCode)
	}
	if err != nil {
		t.Errorf("Failed for success case: expected error to be %v, got %v", nil, err)
	}

}

func TestHandleGetDeckNotFound(t *testing.T) {
	mockParams := httprouter.Params{{Key: "uuid", Value: "test-uuid-123"}}
	mockCtx := context.TODO()
	mdc := mockDeckCRUDOperator{}
	mdc.mockFindDeckByUUID = func(ctx context.Context, uuid string) (database.DeckModel, error) {
		return database.DeckModel{}, mongo.ErrNoDocuments
	}

	req := httptest.NewRequest("GET", "/deck/test-uuid-123", bytes.NewReader([]byte{}))
	expectedErr := ApiError{Message: "Deck with this id does not exist"}
	_, responseCode, err := HandleGetDeck(req, mockParams, &mdc, mockCtx)
	if !cmp.Equal(err, expectedErr) {
		t.Errorf("Failed for deck not found case: expected error to be %v, got %v", expectedErr, err)
	}
	if responseCode != http.StatusNotFound {
		t.Errorf("Failed for deck not found case: expected response code to be %d, got %d", http.StatusNotFound, responseCode)
	}
}

func TestHandleGetDeckDbError(t *testing.T) {
	mockParams := httprouter.Params{{Key: "uuid", Value: "test-uuid-123"}}
	mockCtx := context.TODO()
	mdc := mockDeckCRUDOperator{}
	mdc.mockFindDeckByUUID = func(ctx context.Context, uuid string) (database.DeckModel, error) {
		return database.DeckModel{}, errors.New("test error 456")
	}

	req := httptest.NewRequest("GET", "/deck/test-uuid-123", bytes.NewReader([]byte{}))
	expectedErr := ApiError{Message: "Internal Server Error"}
	_, responseCode, err := HandleGetDeck(req, mockParams, &mdc, mockCtx)
	if !cmp.Equal(err, expectedErr) {
		t.Errorf("Failed for db read error case: expected error to be %v, got %v", expectedErr, err)
	}
	if responseCode != http.StatusInternalServerError {
		t.Errorf("Failed for db read error case: expected response code to be %d, got %d", http.StatusInternalServerError, responseCode)
	}
}

func TestHandleDrawCards(t *testing.T) {
	mockParams := httprouter.Params{{Key: "uuid", Value: "test-uuid-123"}}
	mockCtx := context.TODO()
	mockResult := database.DeckModel{
		UUID:     "test-uuid-123",
		Shuffled: false,
		Cards: deck.Deck{
			{Value: deck.Ace, Suit: deck.Spades},
			{Value: deck.Three, Suit: deck.Clubs},
			{Value: deck.Nine, Suit: deck.Hearts},
		},
	}
	mdc := mockDeckCRUDOperator{}
	mdc.mockFindDeckByUUID = func(ctx context.Context, uuid string) (database.DeckModel, error) {
		return mockResult, nil
	}
	mdc.mockUpdateDeckByUUID = func(ctx context.Context, uuid string, filterQuery bson.D) error {
		return nil
	}
	mockBody, _ := json.Marshal(DrawCardsRequestBody{NumberOfCards: 2})
	req := httptest.NewRequest("PATCH", "/deck/test-uuid-123", bytes.NewReader(mockBody))
	expectedResponse := DrawCardsResponseBody{
		Cards: deck.Deck{{Value: deck.Ace, Suit: deck.Spades}, {Value: deck.Three, Suit: deck.Clubs}}.ToDeckJSON(),
	}
	response, responseCode, err := HandleDrawCards(req, mockParams, &mdc, mockCtx)
	if !cmp.Equal(response, expectedResponse) {
		t.Errorf("Failed for success case: expected error to be %v, got %v", expectedResponse, response)
	}
	if responseCode != http.StatusOK {
		t.Errorf("Failed for success case: expected response code to be %d, got %d", http.StatusOK, responseCode)
	}
	if err != nil {
		t.Errorf("Failed for success case: expected error to be %v, got %v", nil, err)
	}
}

func TestHandleDrawCardsSizeExceededError(t *testing.T) {
	mockParams := httprouter.Params{{Key: "uuid", Value: "test-uuid-1234"}}
	mockCtx := context.TODO()
	mdc := mockDeckCRUDOperator{}
	mockResult := database.DeckModel{
		UUID:     "test-uuid-123",
		Shuffled: false,
		Cards: deck.Deck{
			{Value: deck.Ace, Suit: deck.Spades},
			{Value: deck.Three, Suit: deck.Clubs},
			{Value: deck.Nine, Suit: deck.Hearts},
		},
	}
	mdc.mockFindDeckByUUID = func(ctx context.Context, uuid string) (database.DeckModel, error) {
		return mockResult, nil
	}
	mdc.mockUpdateDeckByUUID = func(ctx context.Context, uuid string, filterQuery bson.D) error {
		return errors.New("test update error")
	}
	mockBody, _ := json.Marshal(DrawCardsRequestBody{NumberOfCards: 4})
	req := httptest.NewRequest("PATCH", "/deck/test-uuid-123", bytes.NewReader(mockBody))
	expectedErr := ApiError{Message: "Requested number of cards is greater than the cards remaining in the deck"}
	_, responseCode, err := HandleDrawCards(req, mockParams, &mdc, mockCtx)
	if !cmp.Equal(err, expectedErr) {
		t.Errorf("Failed for size exceeded error case: expected error to be %v, got %v", expectedErr, err)
	}
	if responseCode != http.StatusBadRequest {
		t.Errorf("Failed for size exceeded error case: expected response code to be %d, got %d", http.StatusBadRequest, responseCode)
	}
}

func TestHandleDrawCardsNumberNotGivenError(t *testing.T) {
	mockParams := httprouter.Params{{Key: "uuid", Value: "test-uuid-1234"}}
	mockCtx := context.TODO()
	mdc := mockDeckCRUDOperator{}
	mockResult := database.DeckModel{
		UUID:     "test-uuid-123",
		Shuffled: false,
		Cards: deck.Deck{
			{Value: deck.Ace, Suit: deck.Spades},
			{Value: deck.Three, Suit: deck.Clubs},
			{Value: deck.Nine, Suit: deck.Hearts},
		},
	}
	mdc.mockFindDeckByUUID = func(ctx context.Context, uuid string) (database.DeckModel, error) {
		return mockResult, nil
	}
	mdc.mockUpdateDeckByUUID = func(ctx context.Context, uuid string, filterQuery bson.D) error {
		return errors.New("test update error")
	}
	mockBody, _ := json.Marshal(struct{}{})
	req := httptest.NewRequest("PATCH", "/deck/test-uuid-123", bytes.NewReader(mockBody))
	expectedErr := ApiError{Message: "Number of cards must be specified and be greater than 0"}
	_, responseCode, err := HandleDrawCards(req, mockParams, &mdc, mockCtx)
	if !cmp.Equal(err, expectedErr) {
		t.Errorf("Failed for size exceeded error case: expected error to be %v, got %v", expectedErr, err)
	}
	if responseCode != http.StatusBadRequest {
		t.Errorf("Failed for size exceeded error case: expected response code to be %d, got %d", http.StatusBadRequest, responseCode)
	}
}

func TestHandleDrawCardsDeckNotFound(t *testing.T) {
	mockParams := httprouter.Params{{Key: "uuid", Value: "test-uuid-123"}}
	mockCtx := context.TODO()
	mdc := mockDeckCRUDOperator{}
	mdc.mockFindDeckByUUID = func(ctx context.Context, uuid string) (database.DeckModel, error) {
		return database.DeckModel{}, mongo.ErrNoDocuments
	}
	mdc.mockUpdateDeckByUUID = func(ctx context.Context, uuid string, filterQuery bson.D) error {
		return nil
	}
	mockBody, _ := json.Marshal(DrawCardsRequestBody{NumberOfCards: 2})
	req := httptest.NewRequest("PATCH", "/deck/test-uuid-123", bytes.NewReader(mockBody))
	expectedErr := ApiError{Message: "Deck with this id does not exist"}
	_, responseCode, err := HandleDrawCards(req, mockParams, &mdc, mockCtx)
	if !cmp.Equal(err, expectedErr) {
		t.Errorf("Failed for deck not found case: expected error to be %v, got %v", expectedErr, err)
	}
	if responseCode != http.StatusNotFound {
		t.Errorf("Failed for deck not found case: expected response code to be %d, got %d", http.StatusNotFound, responseCode)
	}
}

func TestHandleDrawCardsDbReadError(t *testing.T) {
	mockParams := httprouter.Params{{Key: "uuid", Value: "test-uuid-123"}}
	mockCtx := context.TODO()
	mdc := mockDeckCRUDOperator{}
	mdc.mockFindDeckByUUID = func(ctx context.Context, uuid string) (database.DeckModel, error) {
		return database.DeckModel{}, errors.New("test error 456")
	}
	mdc.mockUpdateDeckByUUID = func(ctx context.Context, uuid string, filterQuery bson.D) error {
		return nil
	}
	mockBody, _ := json.Marshal(DrawCardsRequestBody{NumberOfCards: 2})
	req := httptest.NewRequest("PATCH", "/deck/test-uuid-123", bytes.NewReader(mockBody))
	expectedErr := ApiError{Message: "Internal Server Error"}
	_, responseCode, err := HandleDrawCards(req, mockParams, &mdc, mockCtx)
	if !cmp.Equal(err, expectedErr) {
		t.Errorf("Failed for db read error case: expected error to be %v, got %v", expectedErr, err)
	}
	if responseCode != http.StatusInternalServerError {
		t.Errorf("Failed for db read error case: expected response code to be %d, got %d", http.StatusInternalServerError, responseCode)
	}
}

func TestHandleDrawCardsDbUpdateError(t *testing.T) {
	mockParams := httprouter.Params{{Key: "uuid", Value: "test-uuid-123"}}
	mockCtx := context.TODO()
	mdc := mockDeckCRUDOperator{}
	mockResult := database.DeckModel{
		UUID:     "test-uuid-123",
		Shuffled: false,
		Cards: deck.Deck{
			{Value: deck.Ace, Suit: deck.Spades},
			{Value: deck.Three, Suit: deck.Clubs},
			{Value: deck.Nine, Suit: deck.Hearts},
		},
	}
	mdc.mockFindDeckByUUID = func(ctx context.Context, uuid string) (database.DeckModel, error) {
		return mockResult, nil
	}
	mdc.mockUpdateDeckByUUID = func(ctx context.Context, uuid string, filterQuery bson.D) error {
		return errors.New("test update error")
	}
	mockBody, _ := json.Marshal(DrawCardsRequestBody{NumberOfCards: 2})
	req := httptest.NewRequest("PATCH", "/deck/test-uuid-123", bytes.NewReader(mockBody))
	expectedErr := ApiError{Message: "Internal Server Error"}
	_, responseCode, err := HandleDrawCards(req, mockParams, &mdc, mockCtx)
	if !cmp.Equal(err, expectedErr) {
		t.Errorf("Failed for db update error case: expected error to be %v, got %v", expectedErr, err)
	}
	if responseCode != http.StatusInternalServerError {
		t.Errorf("Failed for db update error case: expected response code to be %d, got %d", http.StatusInternalServerError, responseCode)
	}
}
