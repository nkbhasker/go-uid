package main

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func InitRouter(userController UserController) http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Post("/user", userController.CreateUser)
	r.Get("/user/{id}", userController.GetUser)

	return r
}
