package app

import (
	"fmt"
	"log"
	"net/http"
	"refactoring/internal/app/endpoint"
	"refactoring/internal/app/service"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type App struct {
	e *endpoint.Endpoint
	s *service.Service
	r *chi.Mux
}

func New() (*App, error) {
	a := &App{}

	a.s = service.New()

	a.e = endpoint.New(a.s)

	a.r = chi.NewRouter()

	a.r.Use(middleware.RequestID)
	a.r.Use(middleware.RealIP)
	a.r.Use(middleware.Logger)
	//a.r.Use(middleware.Recoverer)
	a.r.Use(middleware.Timeout(60 * time.Second))

	a.r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(time.Now().String()))
	})

	a.r.Route("/api", func(r chi.Router) {
		r.Route("/v1", func(r chi.Router) {
			r.Route("/users", func(r chi.Router) {
				r.Get("/", a.e.SearchUsers)
				r.Post("/", a.e.CreateUser)

				r.Route("/{id}", func(r chi.Router) {
					r.Get("/", a.e.GetUser)
					r.Patch("/", a.e.UpdateUser)
					r.Delete("/", a.e.DeleteUser)
				})
			})
		})
	})

	return a, nil
}

func (a *App) Run() error {
	log.Println("Server running")

	if err := http.ListenAndServe(":3333", a.r); err != nil {
		return fmt.Errorf("Failed to start http server: %v", err)
	}

	return nil
}
