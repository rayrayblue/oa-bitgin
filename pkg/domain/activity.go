package domain

import "time"

type Activity interface {
	SetID(id int)
	GetID() int
	GetPeriod() (time.Time, time.Time)
	SetPeriod(startTime, endTime time.Time) error
	IsInPeriod(time time.Time) bool
}

type activity struct {
	ID        int
	StartDate time.Time
	EndDate   time.Time
}

type BuyTokenActivity struct {
	activity
	MemberLevel      int
	BuyTokenDiscount int
}

func (a *activity) IsInPeriod(time time.Time) bool {
	return a.StartDate.Before(time) && a.EndDate.After(time)
}

func (a *activity) SetID(id int) {
	a.ID = id
}

func (a *activity) GetID() int {
	return a.ID
}

func (a *activity) SetPeriod(startTime, endTime time.Time) error {
	a.StartDate = startTime
	a.EndDate = endTime
	return nil
}

func (a *activity) GetPeriod() (time.Time, time.Time) {
	return a.StartDate, a.EndDate
}

func (t *BuyTokenActivity) GetMemberLevel() int {
	return t.MemberLevel
}

func (t *BuyTokenActivity) GetBuyTokenDiscount() int {
	return t.BuyTokenDiscount
}

type BuyProductActivity struct {
	activity
	PointDiscount int
}

func (t *BuyProductActivity) GetPointDiscount() int {
	return t.PointDiscount
}

type ActivityRepository interface {
	AddBuyTokenActivity(activity BuyTokenActivity) (int, error)
	ListBuyTokenActivity() ([]BuyTokenActivity, error)
	AddBuyProductActivity(activity BuyProductActivity) (int, error)
	GetBuyProductActivity(id int) (BuyProductActivity, error)
}
