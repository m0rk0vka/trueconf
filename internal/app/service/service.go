package service

import (
	"encoding/json"
	"os"
	"refactoring/internal/app/entity"
)

type Service struct{}

func New() *Service {
	return &Service{}
}

func (s *Service) GetUsers() entity.UserStore {
	f, _ := os.ReadFile(entity.STORE_FILE)
	us := entity.UserStore{}
	_ = json.Unmarshal(f, &us)
	return us
}
