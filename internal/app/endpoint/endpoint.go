package endpoint

import (
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
}

func New(s Service) *Endpoint {
	return &Endpoint{
		s: s,
	}
}

func (e *Endpoint) SearchUsers(w http.ResponseWriter, r *http.Request) {
	us, err := e.s.GetUserStore()
	if err != nil {
		_ = render.Render(w, r, errors.InternalServerError(err.Error()))
		return
	}
	render.JSON(w, r, us.List)
}

func (e *Endpoint) CreateUser(w http.ResponseWriter, r *http.Request) {
	req := entity.CreateUserRequest{}
	if err := render.Bind(r, &req); err != nil {
		_ = render.Render(w, r, errors.BadRequest(err.Error()))
		return
	}

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

	render.Status(r, http.StatusCreated)
	render.JSON(w, r, map[string]interface{}{
		"user_id": id,
	})
}

func (e *Endpoint) GetUser(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

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

	render.JSON(w, r, u)
}

func (e *Endpoint) UpdateUser(w http.ResponseWriter, r *http.Request) {
	req := entity.UpdateUserRequest{}

	if err := render.Bind(r, &req); err != nil {
		_ = render.Render(w, r, errors.BadRequest(err.Error()))
		return
	}

	id := chi.URLParam(r, "id")

	if err := e.s.UpdateUser(id, req); err != nil {
		switch err {
		case errors.USER_NOT_FOUND:
			_ = render.Render(w, r, errors.NotFound(err.Error()))
		default:
			_ = render.Render(w, r, errors.InternalServerError(err.Error()))
		}
		return

	}

	render.NoContent(w, r)
}

func (e *Endpoint) DeleteUser(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	if err := e.s.DeleteUser(id); err != nil {
		switch err {
		case errors.USER_NOT_FOUND:
			_ = render.Render(w, r, errors.NotFound(err.Error()))
		default:
			_ = render.Render(w, r, errors.InternalServerError(err.Error()))
		}
		return
	}

	render.NoContent(w, r)
}
