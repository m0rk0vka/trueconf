package endpoint

import (
	"net/http"
	"refactoring/internal/app/entity"
	"refactoring/internal/app/errors"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

type Service interface {
	GetUserStore() entity.UserStore
	Save(entity.UserStore)
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
	us := e.s.GetUserStore()
	render.JSON(w, r, us.List)
}

type CreateUserRequest struct {
	DisplayName string `json:"display_name"`
	Email       string `json:"email"`
}

func (c *CreateUserRequest) Bind(r *http.Request) error { return nil }

func (e *Endpoint) CreateUser(w http.ResponseWriter, r *http.Request) {
	us := e.s.GetUserStore()

	request := CreateUserRequest{}

	if err := render.Bind(r, &request); err != nil {
		_ = render.Render(w, r, errors.BadRequest(err.Error()))
		return
	}

	us.Increment++
	u := entity.User{
		CreatedAt:   time.Now(),
		DisplayName: request.DisplayName,
		Email:       request.DisplayName,
	}

	id := strconv.Itoa(us.Increment)
	us.List[id] = u

	e.s.Save(us)

	render.Status(r, http.StatusCreated)
	render.JSON(w, r, map[string]interface{}{
		"user_id": id,
	})
}

func (e *Endpoint) GetUser(w http.ResponseWriter, r *http.Request) {
	us := e.s.GetUserStore()
	id := chi.URLParam(r, "id")

	render.JSON(w, r, us.List[id])
}

type UpdateUserRequest struct {
	DisplayName string `json:"display_name"`
}

func (c *UpdateUserRequest) Bind(r *http.Request) error { return nil }

func (e *Endpoint) UpdateUser(w http.ResponseWriter, r *http.Request) {
	us := e.s.GetUserStore()

	req := entity.UpdateUserRequest{}

	if err := render.Bind(r, &req); err != nil {
		_ = render.Render(w, r, errors.BadRequest(err.Error()))
		return
	}

	id := chi.URLParam(r, "id")

	if _, ok := us.List[id]; !ok {
		_ = render.Render(w, r, errors.NotFound(errors.USER_NOT_FOUND))
		return
	}

	u := us.List[id]
	u.DisplayName = req.DisplayName
	us.List[id] = u

	e.s.Save(us)

	render.Status(r, http.StatusNoContent)
}

func (e *Endpoint) DeleteUser(w http.ResponseWriter, r *http.Request) {
	us := e.s.GetUserStore()

	id := chi.URLParam(r, "id")

	if _, ok := us.List[id]; !ok {
		_ = render.Render(w, r, errors.NotFound(errors.USER_NOT_FOUND))
		return
	}

	delete(us.List, id)

	e.s.Save(us)

	render.Status(r, http.StatusNoContent)
}
