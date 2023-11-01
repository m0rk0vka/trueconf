package app

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"refactoring/internal/app/endpoint"
	"refactoring/internal/app/service"
	"time"

	"refactoring/internal/app/config"
	"refactoring/pkg/logger"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

var flagConfig = flag.String("config", "./config/local.yml", "path to the config file")

type App struct {
	e *endpoint.Endpoint
	s *service.Service
	l *log.Logger
	r *chi.Mux
}

func New() (*App, error) {
	a := &App{}

	l, err := logger.New()
	if err != nil {
		return nil, fmt.Errorf("New app: %v", err)
	}

	a.l = l

	a.s = service.New(a.l)

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
	a.l.Println("Running app")
	flag.Parse()
	a.l.Println("flags:", flag.Args())

	cfg, err := config.Load(*flagConfig)
	if err != nil {
		return fmt.Errorf("Failed to load config")
	}
	a.l.Printf("Config: %+v\n", cfg)

	addr := fmt.Sprintf(":%v", cfg.ServerPort)
	a.l.Printf("Address: %+v\n", addr)

	if err := http.ListenAndServe(addr, a.r); err != nil {
		return fmt.Errorf("Failed to start http server: %v", err)
	}

	return nil
}
