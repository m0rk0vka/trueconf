package service

import (
	"encoding/json"
	"io/fs"
	"log"
	"os"
	"refactoring/internal/app/entity"
	"refactoring/internal/app/errors"
	"strconv"
	"time"
)

type Service struct{}

func New() *Service {
	return &Service{}
}

func (s *Service) GetUserStore() entity.UserStore {
	f, _ := os.ReadFile(entity.STORE_FILE)
	us := entity.UserStore{}
	_ = json.Unmarshal(f, &us)
	return us
}

func (s *Service) GetUser(id string) (entity.User, error) {
	us := s.GetUserStore()
	if !isUserExist(id, us) {
		return entity.User{}, errors.USER_NOT_FOUND
	}
	return us.List[id], nil
}

func (s *Service) Save(us entity.UserStore) {
	b, _ := json.Marshal(&us)
	_ = os.WriteFile(entity.STORE_FILE, b, fs.ModePerm)
}
func (s *Service) CreateUser(r entity.CreateUserRequest) string {
	us := s.GetUserStore()
	us.Increment++
	u := entity.User{
		CreatedAt:   time.Now(),
		DisplayName: r.DisplayName,
		Email:       r.DisplayName,
	}

	id := strconv.Itoa(us.Increment)
	us.List[id] = u

	s.Save(us)

	return id
}

func isUserExist(id string, us entity.UserStore) bool {
	_, ok := us.List[id]
	log.Printf("If id %v exist: %v", id, ok)
	return ok
}

func (s *Service) UpdateUser(id string, r entity.UpdateUserRequest) error {
	us := s.GetUserStore()
	if !isUserExist(id, us) {
		return errors.USER_NOT_FOUND
	}

	u := us.List[id]
	if r.DisplayName != "" {
		u.DisplayName = r.DisplayName
	}
	us.List[id] = u

	s.Save(us)

	return nil
}

func (s *Service) DeleteUser(id string) error {
	us := s.GetUserStore()
	if !isUserExist(id, us) {
		return errors.USER_NOT_FOUND
	}

	delete(us.List, id)

	s.Save(us)

	return nil
}
