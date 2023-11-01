package main

import (
	"net/http"
	"time"

	"refactoring/internal/app/endpoint"
	"refactoring/internal/app/service"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	//r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(time.Now().String()))
	})

	s := service.New()
	e := endpoint.New(s)

	r.Route("/api", func(r chi.Router) {
		r.Route("/v1", func(r chi.Router) {
			r.Route("/users", func(r chi.Router) {
				r.Get("/", e.SearchUsers)
				r.Post("/", e.CreateUser)

				r.Route("/{id}", func(r chi.Router) {
					r.Get("/", e.GetUser)
					r.Patch("/", e.UpdateUser)
					r.Delete("/", e.DeleteUser)
				})
			})
		})
	})

	http.ListenAndServe(":3333", r)
}
