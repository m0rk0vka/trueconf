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

type Service struct {
	l *log.Logger
}

func New(l *log.Logger) *Service {
	return &Service{
		l: l,
	}
}

func (s *Service) GetUserStore() (entity.UserStore, error) {
	s.l.Println("Get user store")
	b, err := os.ReadFile(entity.STORE_FILE)
	if err != nil {
		return entity.UserStore{}, fmt.Errorf("While getting store: %v", err)
	}
	us := entity.UserStore{}
	if err := json.Unmarshal(b, &us); err != nil {
		return entity.UserStore{}, fmt.Errorf("While getting store: %v", err)
	}
	s.l.Printf("User store: %+v\n", us)
	return us, nil
}

func (s *Service) GetUser(id string) (entity.User, error) {
	s.l.Printf("Get user by id: %v\n", id)
	us, err := s.GetUserStore()
	if err != nil {
		return entity.User{}, fmt.Errorf("While getting user: %v", err)
	}

	if !s.isUserExist(id, us) {
		return entity.User{}, errors.USER_NOT_FOUND
	}

	u := us.List[id]
	s.l.Printf("User with id=%v: %+v\n", id, u)

	return u, nil
}

func (s *Service) save(us entity.UserStore) error {
	s.l.Printf("Save user store:%+v\n", us)

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
	s.l.Printf("Create user with request: %+v", r)
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

	s.l.Printf("new user: %+v\n", u)

	id := strconv.Itoa(us.Increment)
	s.l.Printf("new id: %v\n", id)
	us.List[id] = u

	if err := s.save(us); err != nil {
		return "", fmt.Errorf("While creating user: %v", err)
	}

	return id, nil
}

func (s *Service) isUserExist(id string, us entity.UserStore) bool {
	s.l.Printf("Checking if exist user with id: %v\n", id)
	_, ok := us.List[id]
	s.l.Printf("If id %v exist: %v", id, ok)
	return ok
}

func (s *Service) UpdateUser(id string, r entity.UpdateUserRequest) error {
	s.l.Printf("Update user with id: %v, request:%v\n", id, r)
	us, err := s.GetUserStore()
	if err != nil {
		return fmt.Errorf("While updating user: %v", err)
	}
	if !s.isUserExist(id, us) {
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
	s.l.Printf("Deleting user with id: %v\n", id)
	us, err := s.GetUserStore()
	if err != nil {
		return fmt.Errorf("While deleting user: %v", err)
	}
	if !s.isUserExist(id, us) {
		return errors.USER_NOT_FOUND
	}

	delete(us.List, id)

	if err := s.save(us); err != nil {
		return fmt.Errorf("While deleting user: %v", err)
	}

	return nil
}
