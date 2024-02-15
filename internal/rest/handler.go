package rest

import "github.com/go-chi/chi/v5"

type Hanlder interface {
	Register(router *chi.Mux)
}
