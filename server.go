package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"go.mongodb.org/mongo-driver/mongo"
)

type server struct {
	router   *httprouter.Router
	dbClient *mongo.Database
}

func crashHandler(w http.ResponseWriter, r *http.Request, err interface{}) {
	log.Println(r.Method, r.URL.Path, err)
	w.WriteHeader(http.StatusInternalServerError)
	fmt.Fprintf(w, "Internal Server Error\n")
}

func Status(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, "UP!\n")
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println(r.Method, r.URL.Path)
	s.router.ServeHTTP(w, r)
}

func (s *server) initRouter() {
	s.router.PanicHandler = crashHandler
	s.router.GET("/status", Status)
	s.router.POST("/deck", s.handleCreateDeck)
	s.router.GET("/deck/:uuid", s.handleGetDeck)
	s.router.PATCH("/deck/:uuid", s.handleDrawCards)

}
