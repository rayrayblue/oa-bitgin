package domain

import "time"

type CashierUsecase interface {
	NewUser(name string, memberLevel int) (int, error)
	NewBuyTokenActivity(memberLevel int, startTime time.Time, endTime time.Time, discount int) (int, error)
	NewBuyProductActivity(startTime time.Time, endTime time.Time, discount int) (int, error)

	BuyToken(userID int, token int64) (int, error)
	BuyTokenWithActivity(userID int, token int64) (int, error)

	AddPoint(userID int, token int64) error

	GetUserToken(userID int) (int, error)
	GetUserPoint(userID int) (int, error)

	NewProduct(name string, price int) (int, error)
	BuyProduct(userID int, productID int) (int, error)
	BuyProductWithActivity(userID int, productID int, activityID int) (int, error)

	GetTotalAmount() int64
}
