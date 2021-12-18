package main

import (
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func main() {
	s := &server{
		router: httprouter.New(),
	}
	s.initRouter()
	log.Fatal(http.ListenAndServe(":8080", s))
}
