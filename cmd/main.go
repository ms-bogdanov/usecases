package main

import (
	"github.com/go-chi/chi"
	"github.com/ptflp/godecoder"
	"net/http"
	"usecases/internal/controller"
	"usecases/internal/service"
)

func main() {
	r := chi.NewRouter()

	c := controller.NewController(controller.NewResponder(godecoder.NewDecoder(), nil), service.NewService())

	r.Post("/api/address/search", c.SearchHandler)
	r.Post("/api/address/geocode", c.GeocodeHandler)

	http.ListenAndServe(":8080", r)
}
