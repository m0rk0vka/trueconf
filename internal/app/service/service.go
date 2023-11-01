package service

import (
	"encoding/json"
	"io/fs"
	"os"
	"refactoring/internal/app/entity"
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

func (s *Service) Save(us entity.UserStore) {
	b, _ := json.Marshal(&us)
	_ = os.WriteFile(entity.STORE_FILE, b, fs.ModePerm)
}
