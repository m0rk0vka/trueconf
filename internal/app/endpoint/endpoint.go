package endpoint

import (
	"log"
	"net/http"
	"refactoring/internal/app/entity"
	"refactoring/internal/app/errors"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	validator "github.com/go-playground/validator/v10"
)

type Service interface {
	GetUserStore() (entity.UserStore, error)
	CreateUser(entity.CreateUserRequest) (string, error)
	UpdateUser(string, entity.UpdateUserRequest) error
	GetUser(string) (entity.User, error)
	DeleteUser(string) error
}

type Endpoint struct {
	s Service
	l *log.Logger
}

func New(s Service, l *log.Logger) *Endpoint {
	return &Endpoint{
		s: s,
		l: l,
	}
}

func (e *Endpoint) SearchUsers(w http.ResponseWriter, r *http.Request) {
	e.l.Println("Search users handler func")
	us, err := e.s.GetUserStore()
	if err != nil {
		_ = render.Render(w, r, errors.InternalServerError(err.Error()))
		return
	}

	e.l.Println("Successifully get all users")

	render.JSON(w, r, us.List)
}

func (e *Endpoint) CreateUser(w http.ResponseWriter, r *http.Request) {
	e.l.Println("Create user handler func")
	req := entity.CreateUserRequest{}
	if err := render.Bind(r, &req); err != nil {
		_ = render.Render(w, r, errors.BadRequest(err.Error()))
		return
	}

	e.l.Printf("Request: %v\n", req)
	v := validator.New()
	if err := v.Struct(req); err != nil {
		_ = render.Render(w, r, errors.BadRequest(err.Error()))
		return
	}

	id, err := e.s.CreateUser(req)
	if err != nil {
		_ = render.Render(w, r, errors.InternalServerError(err.Error()))
		return
	}

	e.l.Printf("Successifully create user with id: %v\n", id)

	render.Status(r, http.StatusCreated)
	render.JSON(w, r, map[string]interface{}{
		"user_id": id,
	})
}

func (e *Endpoint) GetUser(w http.ResponseWriter, r *http.Request) {
	e.l.Println("Get user by id handler func")

	id := chi.URLParam(r, "id")
	e.l.Printf("Id: %v\n", id)

	u, err := e.s.GetUser(id)
	if err != nil {
		switch err {
		case errors.USER_NOT_FOUND:
			_ = render.Render(w, r, errors.NotFound(err.Error()))
		default:
			_ = render.Render(w, r, errors.InternalServerError(err.Error()))
		}
		return
	}

	e.l.Printf("Successifully get user: %+v\n", u)

	render.JSON(w, r, u)
}

func (e *Endpoint) UpdateUser(w http.ResponseWriter, r *http.Request) {
	e.l.Println("Update user handler func")
	req := entity.UpdateUserRequest{}

	if err := render.Bind(r, &req); err != nil {
		_ = render.Render(w, r, errors.BadRequest(err.Error()))
		return
	}
	e.l.Printf("Request: %v\n", req)

	v := validator.New()
	if err := v.Struct(req); err != nil {
		_ = render.Render(w, r, errors.BadRequest(err.Error()))
		return
	}

	id := chi.URLParam(r, "id")
	e.l.Printf("Id: %v\n", id)

	if err := e.s.UpdateUser(id, req); err != nil {
		switch err {
		case errors.USER_NOT_FOUND:
			_ = render.Render(w, r, errors.NotFound(err.Error()))
		default:
			_ = render.Render(w, r, errors.InternalServerError(err.Error()))
		}
		return

	}

	e.l.Println("Successifully update user")

	render.NoContent(w, r)
}

func (e *Endpoint) DeleteUser(w http.ResponseWriter, r *http.Request) {
	e.l.Println("Delete user handler func")

	id := chi.URLParam(r, "id")
	e.l.Printf("Id: %v\n", id)

	if err := e.s.DeleteUser(id); err != nil {
		switch err {
		case errors.USER_NOT_FOUND:
			_ = render.Render(w, r, errors.NotFound(err.Error()))
		default:
			_ = render.Render(w, r, errors.InternalServerError(err.Error()))
		}
		return
	}

	e.l.Println("Successifully delete user")

	render.NoContent(w, r)
}
