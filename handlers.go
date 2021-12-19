package main

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/AbhilashJN/cards/api"
	"github.com/julienschmidt/httprouter"
)

type genericResponseMessage struct {
	Message string `json:"message"`
}

func (s *server) handleCreateDeck(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	collection := s.dbClient.Database("cardsdb").Collection("decks")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	responseBody, responseCode, err := api.HandleCreateDeck(r, ps, collection, ctx)
	w.WriteHeader(responseCode)
	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		json.NewEncoder(w).Encode(genericResponseMessage{err.Error()})
		return
	}
	json.NewEncoder(w).Encode(responseBody)
}
