package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/AbhilashJN/cards/api"
	"github.com/AbhilashJN/cards/database"
	"github.com/julienschmidt/httprouter"
)

type genericResponseMessage struct {
	Message string `json:"message"`
}

func (s *server) handleCreateDeck(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	collection := s.dbClient.Collection("decks")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	dc := &database.DeckCRUDOperator{Collection: collection}
	responseBody, responseCode, err := api.HandleCreateDeck(r, ps, dc, ctx)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(responseCode)
	e := json.NewEncoder(w)
	if err != nil {
		e.Encode(genericResponseMessage{err.Error()})
	} else {
		e.Encode(responseBody)
	}
	log.Println(r.Method, r.URL.Path, responseCode)

}

func (s *server) handleGetDeck(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	collection := s.dbClient.Collection("decks")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	dc := &database.DeckCRUDOperator{Collection: collection}
	responseBody, responseCode, err := api.HandleGetDeck(r, ps, dc, ctx)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(responseCode)
	e := json.NewEncoder(w)
	if err != nil {
		e.Encode(genericResponseMessage{err.Error()})
	} else {
		e.Encode(responseBody)
	}
	log.Println(r.Method, r.URL.Path, responseCode)
}

func (s *server) handleDrawCards(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	collection := s.dbClient.Collection("decks")
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	dc := &database.DeckCRUDOperator{Collection: collection}
	responseBody, responseCode, err := api.HandleDrawCards(r, ps, dc, ctx)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(responseCode)
	e := json.NewEncoder(w)
	if err != nil {
		e.Encode(genericResponseMessage{err.Error()})
	} else {
		e.Encode(responseBody)
	}
	log.Println(r.Method, r.URL.Path, responseCode)
}
