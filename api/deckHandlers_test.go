package api

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/julienschmidt/httprouter"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type mockMongoCollection struct {
	mockInsertOne func(context.Context, interface{},
		...*options.InsertOneOptions) (*mongo.InsertOneResult, error)
}

func (c mockMongoCollection) InsertOne(ctx context.Context, document interface{},
	opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error) {
	return c.mockInsertOne(ctx, document, opts...)
}

type HandleCreateDeckTest struct {
	shuffle          bool
	customDeck       bool
	wantedCards      []string
	expectedNumCards int
}

func TestHandleCreateDeck(t *testing.T) {
	collection := mockMongoCollection{}
	mockParams := httprouter.Params{}
	mockCtx := context.TODO()
	collection.mockInsertOne = func(ctx context.Context, document interface{},
		opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error) {
		return &mongo.InsertOneResult{InsertedID: 123}, nil
	}

	tests := []HandleCreateDeckTest{
		{shuffle: false, customDeck: false, wantedCards: []string{}, expectedNumCards: 52},
		{shuffle: true, customDeck: false, wantedCards: []string{}, expectedNumCards: 52},
		{shuffle: false, customDeck: true, wantedCards: []string{"AS", "QS", "2H", "7D", "4C"}, expectedNumCards: 5},
		{shuffle: true, customDeck: true, wantedCards: []string{"AS", "QS", "2H", "7D", "4C"}, expectedNumCards: 5},
		{shuffle: true, customDeck: true, wantedCards: []string{}, expectedNumCards: 0},
	}

	for _, test := range tests {
		mockBody, _ := json.Marshal(CreateDeckRequestBody{Shuffle: test.shuffle, CustomDeck: test.customDeck, WantedCards: test.wantedCards})
		req := httptest.NewRequest("POST", "http://www.test.com", bytes.NewReader(mockBody))
		responseBody, responseCode, err := HandleCreateDeck(req, mockParams, collection, mockCtx)
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

func TestHandleCreateDeckDbError(t *testing.T) {
	collection := mockMongoCollection{}
	mockParams := httprouter.Params{}
	mockCtx := context.TODO()
	collection.mockInsertOne = func(ctx context.Context, document interface{},
		opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error) {
		return &mongo.InsertOneResult{}, errors.New("test error")
	}
	mockBody, _ := json.Marshal(CreateDeckRequestBody{Shuffle: false, CustomDeck: false, WantedCards: []string{}})
	req := httptest.NewRequest("POST", "http://www.test.com", bytes.NewReader(mockBody))
	expectedErr := ApiError{Message: "Internal Server Error"}
	_, responseCode, err := HandleCreateDeck(req, mockParams, collection, mockCtx)
	if !cmp.Equal(err, expectedErr) {
		t.Errorf("Failed for error case: expected error to be %v, got %v", expectedErr, err)
	}
	if responseCode != http.StatusInternalServerError {
		t.Errorf("Failed for error case: expected response code to be %d, got %d", http.StatusInternalServerError, responseCode)
	}
}

func TestHandleCreateDeckBadRequest(t *testing.T) {
	collection := mockMongoCollection{}
	mockParams := httprouter.Params{}
	mockCtx := context.TODO()
	collection.mockInsertOne = func(ctx context.Context, document interface{},
		opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error) {
		return &mongo.InsertOneResult{InsertedID: 123}, errors.New("test error")
	}
	mockBody, _ := json.Marshal(CreateDeckRequestBody{Shuffle: false, CustomDeck: false, WantedCards: []string{}})
	req := httptest.NewRequest("POST", "http://www.test.com", bytes.NewReader(mockBody[:len(mockBody)-2]))
	expectedErr := ApiError{Message: "Request body is malformed"}
	_, responseCode, err := HandleCreateDeck(req, mockParams, collection, mockCtx)
	if !cmp.Equal(err, expectedErr) {
		t.Errorf("Failed for error case: expected error to be %v, got %v", expectedErr, err)
	}
	if responseCode != http.StatusBadRequest {
		t.Errorf("Failed for error case: expected response code to be %d, got %d", http.StatusBadRequest, responseCode)
	}
}
