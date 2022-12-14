package repository

import (
	"errors"
	"oa-bitgin/pkg/domain"
	"sync/atomic"
)

type userRepository struct {
	IDCounter atomic.Value
	Users     map[int]*domain.User
}

func (u *userRepository) init() {
	u.IDCounter.Store(0)
	u.Users = make(map[int]*domain.User)
}

func NewUserRepository() domain.UserRepository {
	store := &userRepository{}
	store.init()
	return store
}

func (u *userRepository) GetUser(id int) (*domain.User, error) {
	if user, ok := u.Users[id]; ok {
		return user, nil
	} else {
		return &domain.User{}, errors.New("user not found")
	}
}

func (u *userRepository) NewUser(user domain.User) (int, error) {
	id := u.IDCounter.Load().(int)
	id++
	u.IDCounter.Store(id)
	user.ID = id
	user.Account.Token.Store(0)
	user.Account.Point.Store(0)
	u.Users[id] = &user
	return id, nil
}

func (u *userRepository) GetDefaultBuyTokenDiscount(level int) int {
	switch level {
	case 1:
		return 95
	case 2:
		return 90
	case 3:
		return 85
	default:
		return 100
	}
}
