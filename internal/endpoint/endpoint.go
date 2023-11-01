package endpoint

import (
	"encoding/json"
	"io/fs"
	"io/ioutil"
	"net/http"
	"refactoring/internal/entity"
	"refactoring/internal/errors"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

type Endpoint struct {
}

func New() *Endpoint {
	return &Endpoint{}
}

const store = `./data/users.json`

func SearchUsers(w http.ResponseWriter, r *http.Request) {
	f, _ := ioutil.ReadFile(store)
	s := entity.UserStore{}
	_ = json.Unmarshal(f, &s)

	render.JSON(w, r, s.List)
}

type CreateUserRequest struct {
	DisplayName string `json:"display_name"`
	Email       string `json:"email"`
}

func (c *CreateUserRequest) Bind(r *http.Request) error { return nil }

func CreateUser(w http.ResponseWriter, r *http.Request) {
	f, _ := ioutil.ReadFile(store)
	s := entity.UserStore{}
	_ = json.Unmarshal(f, &s)

	request := CreateUserRequest{}

	if err := render.Bind(r, &request); err != nil {
		_ = render.Render(w, r, errors.BadRequest(err.Error()))
		return
	}

	s.Increment++
	u := entity.User{
		CreatedAt:   time.Now(),
		DisplayName: request.DisplayName,
		Email:       request.DisplayName,
	}

	id := strconv.Itoa(s.Increment)
	s.List[id] = u

	b, _ := json.Marshal(&s)
	_ = ioutil.WriteFile(store, b, fs.ModePerm)

	render.Status(r, http.StatusCreated)
	render.JSON(w, r, map[string]interface{}{
		"user_id": id,
	})
}

func GetUser(w http.ResponseWriter, r *http.Request) {
	f, _ := ioutil.ReadFile(store)
	s := entity.UserStore{}
	_ = json.Unmarshal(f, &s)

	id := chi.URLParam(r, "id")

	render.JSON(w, r, s.List[id])
}

type UpdateUserRequest struct {
	DisplayName string `json:"display_name"`
}

func (c *UpdateUserRequest) Bind(r *http.Request) error { return nil }

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	f, _ := ioutil.ReadFile(store)
	s := entity.UserStore{}
	_ = json.Unmarshal(f, &s)

	request := UpdateUserRequest{}

	if err := render.Bind(r, &request); err != nil {
		_ = render.Render(w, r, errors.BadRequest(err.Error()))
		return
	}

	id := chi.URLParam(r, "id")

	if _, ok := s.List[id]; !ok {
		_ = render.Render(w, r, errors.NotFound(errors.USER_NOT_FOUND))
		return
	}

	u := s.List[id]
	u.DisplayName = request.DisplayName
	s.List[id] = u

	b, _ := json.Marshal(&s)
	_ = ioutil.WriteFile(store, b, fs.ModePerm)

	render.Status(r, http.StatusNoContent)
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	f, _ := ioutil.ReadFile(store)
	s := entity.UserStore{}
	_ = json.Unmarshal(f, &s)

	id := chi.URLParam(r, "id")

	if _, ok := s.List[id]; !ok {
		_ = render.Render(w, r, errors.NotFound(errors.USER_NOT_FOUND))
		return
	}

	delete(s.List, id)

	b, _ := json.Marshal(&s)
	_ = ioutil.WriteFile(store, b, fs.ModePerm)

	render.Status(r, http.StatusNoContent)
}
