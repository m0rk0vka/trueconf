package service

import (
	"encoding/json"
	"fmt"
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

func (s *Service) GetUserStore() (entity.UserStore, error) {
	b, err := os.ReadFile(entity.STORE_FILE)
	if err != nil {
		return entity.UserStore{}, fmt.Errorf("While getting store: %v", err)
	}
	us := entity.UserStore{}
	if err := json.Unmarshal(b, &us); err != nil {
		return entity.UserStore{}, fmt.Errorf("While getting store: %v", err)
	}
	return us, nil
}

func (s *Service) GetUser(id string) (entity.User, error) {
	us, err := s.GetUserStore()
	if err != nil {
		return entity.User{}, fmt.Errorf("While getting user: %v", err)
	}

	if !isUserExist(id, us) {
		return entity.User{}, errors.USER_NOT_FOUND
	}
	return us.List[id], nil
}

func (s *Service) save(us entity.UserStore) error {
	b, err := json.Marshal(&us)
	if err != nil {
		return fmt.Errorf("While saving store: %v", err)
	}
	if err := os.WriteFile(entity.STORE_FILE, b, fs.ModePerm); err != nil {
		return fmt.Errorf("While saving store: %v", err)
	}

	return nil
}

func (s *Service) CreateUser(r entity.CreateUserRequest) (string, error) {
	us, err := s.GetUserStore()
	if err != nil {
		return "", fmt.Errorf("While creating user: %v", err)
	}
	us.Increment++
	u := entity.User{
		CreatedAt:   time.Now(),
		DisplayName: r.DisplayName,
		Email:       r.Email,
	}

	id := strconv.Itoa(us.Increment)
	us.List[id] = u

	if err := s.save(us); err != nil {
		return "", fmt.Errorf("While creating user: %v", err)
	}

	return id, nil
}

func isUserExist(id string, us entity.UserStore) bool {
	_, ok := us.List[id]
	log.Printf("If id %v exist: %v", id, ok)
	return ok
}

func (s *Service) UpdateUser(id string, r entity.UpdateUserRequest) error {
	us, err := s.GetUserStore()
	if err != nil {
		return fmt.Errorf("While updating user: %v", err)
	}
	if !isUserExist(id, us) {
		return errors.USER_NOT_FOUND
	}

	u := us.List[id]
	if r.DisplayName != "" {
		u.DisplayName = r.DisplayName
	}
	us.List[id] = u

	if err := s.save(us); err != nil {
		return fmt.Errorf("While updating user: %v", err)
	}

	return nil
}

func (s *Service) DeleteUser(id string) error {
	us, err := s.GetUserStore()
	if err != nil {
		return fmt.Errorf("While deleting user: %v", err)
	}
	if !isUserExist(id, us) {
		return errors.USER_NOT_FOUND
	}

	delete(us.List, id)

	if err := s.save(us); err != nil {
		return fmt.Errorf("While deleting user: %v", err)
	}

	return nil
}
