package domain

import (
	"sync/atomic"
)

type Account struct {
	Token atomic.Value // 平台幣
	Point atomic.Value // 平台點數
}

type Member struct {
	Level                   int // 0: Normal, 1: VIP1, 2: VIP2, 3: VIP3
	BuyTokenDefaultDiscount int // 1-100, VIP會員有平台幣優惠價格 (例如: VIP1: 95折，VIP2: 9折，VIP3: 85折，各個等級的折扣會依照活動做調整。)
}

type User struct {
	ID      int
	Name    string
	Account Account
	Member  Member
}

func (u *User) BuyToken(token int) int {
	u.Account.Token.Store(u.Account.Token.Load().(int) + token)
	rtn := u.Account.Token.Load().(int)
	return rtn
}

func (u *User) UseToken(token int) int {
	u.Account.Token.Store(u.Account.Token.Load().(int) - token)
	rtn := u.Account.Token.Load().(int)
	return rtn
}

func (u *User) AddPoint(point int) int {
	u.Account.Point.Store(u.Account.Point.Load().(int) + point)
	rtn := u.Account.Point.Load().(int)
	return rtn
}

func (u *User) UsePoint(point int) int {
	u.Account.Point.Store(u.Account.Point.Load().(int) - point)
	rtn := u.Account.Point.Load().(int)
	return rtn
}

func (u *User) GetPoint() int {
	return u.Account.Point.Load().(int)
}

func (u *User) GetToken() int {
	return u.Account.Token.Load().(int)
}

type UserRepository interface {
	GetUser(id int) (*User, error)
	NewUser(user User) (int, error)
	GetDefaultBuyTokenDiscount(level int) int
}
