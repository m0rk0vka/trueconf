package service

import (
	"encoding/json"
	"io/fs"
	"os"
	"refactoring/internal/app/entity"
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
