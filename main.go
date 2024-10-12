package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/olahol/melody"
)

var m *TMap
var mel *melody.Melody

func init() {
	m = getNewMap(10, 10)
	mel = melody.New()
}

func main() {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/", HandleRoot)
	r.Get("/ws", HandleWS(mel))
	mel.HandleMessage(HandleMessage(m))

	fmt.Println("Listening on :8080")
	http.ListenAndServe(":8080", r)
}
