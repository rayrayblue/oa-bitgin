package usecase

import (
	"errors"
	"fmt"
	"math"
	"oa-bitgin/pkg/domain"
	"time"
)

type cashierUsecase struct {
	TotalAmount  int64 // total amount of money that cashier has
	userRepo     domain.UserRepository
	activityRepo domain.ActivityRepository
	productRepo  domain.ProductRepository
}

func NewCashierUsecase(userRepo domain.UserRepository, activityRepo domain.ActivityRepository, productRepo domain.ProductRepository) domain.CashierUsecase {
	return &cashierUsecase{
		userRepo:     userRepo,
		activityRepo: activityRepo,
		productRepo:  productRepo,
	}
}

func (c *cashierUsecase) GetTotalAmount() int64 {
	return c.TotalAmount
}

func (c *cashierUsecase) NewUser(name string, memberLevel int) (int, error) {
	user := domain.User{
		Name: name,
		Member: domain.Member{
			Level:                   memberLevel,
			BuyTokenDefaultDiscount: c.userRepo.GetDefaultBuyTokenDiscount(memberLevel),
		},
	}
	user.Account.Token.Store(0)
	user.Account.Point.Store(0)
	return c.userRepo.NewUser(user)
}

func (c *cashierUsecase) BuyToken(userID int, token int64) (int, error) {
	user, err := c.userRepo.GetUser(userID)
	if err != nil {
		return -1, err
	}

	user.BuyToken(int(token))
	rtn := token * int64(user.Member.BuyTokenDefaultDiscount) / 100
	c.TotalAmount += rtn
	fmt.Println("[MSG] Need to charge: ", rtn)
	return int(rtn), nil
}

func (c *cashierUsecase) AddPoint(userID int, point int64) error {
	user, err := c.userRepo.GetUser(userID)
	if err != nil {
		return err
	}

	user.AddPoint(int(point))
	return nil
}

func (c *cashierUsecase) NewBuyTokenActivity(memberLevel int, startTime time.Time, endTime time.Time, discount int) (int, error) {
	a := domain.BuyTokenActivity{
		MemberLevel:      memberLevel,
		BuyTokenDiscount: discount,
	}
	_ = a.SetPeriod(startTime, endTime)
	aID, _ := c.activityRepo.AddBuyTokenActivity(a)
	return aID, nil
}

func (c *cashierUsecase) BuyTokenWithActivity(userID int, token int64) (int, error) {
	user, err := c.userRepo.GetUser(userID)
	if err != nil {
		return -1, err
	}

	activities, err := c.activityRepo.ListBuyTokenActivity()
	if err != nil {
		return -1, err
	}

	// find the best price for user
	find := false
	bestPrice := int64(math.MaxInt64)
	for _, a := range activities {
		if a.IsInPeriod(time.Now()) && a.MemberLevel == user.Member.Level {
			find = true
			tmp := token * int64(a.BuyTokenDiscount) / 100
			if tmp < bestPrice {
				bestPrice = tmp
			}
		}
	}
	if find == true {
		user.BuyToken(int(token))
		fmt.Println("[MSG] Need to charge: ", bestPrice)
		c.TotalAmount += bestPrice
		return int(bestPrice), nil
	}

	user.BuyToken(int(token))
	bestPrice = token * int64(user.Member.BuyTokenDefaultDiscount) / 100
	c.TotalAmount += bestPrice
	fmt.Println("[MSG] Need to charge (without activity because no activity matched): ", token*int64(user.Member.BuyTokenDefaultDiscount)/100)
	return int(bestPrice), nil
}

func (c *cashierUsecase) GetUserToken(userID int) (int, error) {
	user, err := c.userRepo.GetUser(userID)
	if err != nil {
		fmt.Println(fmt.Sprintf("[MSG] User %d not found", userID))
		return -1, err
	}
	fmt.Println(fmt.Sprintf("[MSG] User %d has %d token", user.ID, user.GetToken()))
	return user.GetToken(), err
}

func (c *cashierUsecase) GetUserPoint(userID int) (int, error) {
	user, err := c.userRepo.GetUser(userID)
	if err != nil {
		fmt.Println(fmt.Sprintf("[MSG] User %d not found", userID))
		return -1, err
	}
	fmt.Println(fmt.Sprintf("[MSG] User %d has %d point", user.ID, user.GetPoint()))
	return user.GetPoint(), err
}

func (c *cashierUsecase) BuyProduct(userID int, productID int) (int, error) {
	user, err := c.userRepo.GetUser(userID)
	if err != nil {
		fmt.Println(fmt.Sprintf("[MSG] User %d not found", userID))
		return -1, err
	}

	product, err := c.productRepo.GetProduct(productID)
	if err != nil {
		fmt.Println(fmt.Sprintf("[MSG] Product %d not found", productID))
		return -1, err
	}

	if user.GetToken() < product.Price {
		fmt.Println(fmt.Sprintf("[MSG] User %d has not enough token to buy product %d", userID, productID))
		return -1, errors.New("not enough token")
	}

	user.UseToken(product.Price)
	fmt.Println(fmt.Sprintf("[MSG] User %d has bought product %d use price %d", userID, productID, product.Price))
	return product.Price, nil
}

func (c *cashierUsecase) NewProduct(name string, price int) (int, error) {
	p := domain.Product{
		Name:  name,
		Price: price,
	}

	id, err := c.productRepo.AddProduct(p)
	if err != nil {
		return -1, err
	}
	return id, nil
}

func (c *cashierUsecase) NewBuyProductActivity(startTime time.Time, endTime time.Time, discount int) (int, error) {
	a := domain.BuyProductActivity{
		PointDiscount: discount,
	}
	_ = a.SetPeriod(startTime, endTime)
	aID, _ := c.activityRepo.AddBuyProductActivity(a)
	return aID, nil
}

func (c *cashierUsecase) BuyProductWithActivity(userID int, productID int, activityID int) (int, error) {
	user, err := c.userRepo.GetUser(userID)
	if err != nil {
		fmt.Println(fmt.Sprintf("[MSG] User %d not found", userID))
		return -1, err
	}

	product, err := c.productRepo.GetProduct(productID)
	if err != nil {
		fmt.Println(fmt.Sprintf("[MSG] Product %d not found", productID))
		return -1, err
	}

	activity, err := c.activityRepo.GetBuyProductActivity(activityID)
	if err != nil {
		fmt.Println(fmt.Sprintf("[MSG] Activity %d not found", activityID))
		return -1, err
	}

	needPoint := product.Price * (100 - activity.GetPointDiscount()) / 100
	if user.GetPoint() < needPoint {
		fmt.Println(fmt.Sprintf("[MSG] User %d has not enough point to buy product %d", userID, productID))
		return -1, errors.New("not enough point")
	}

	// 平台後來新增了另一個收費模式，如果有VIP身份扣100點以上折抵，另外享再九折優惠。
	var needToken int
	if user.Member.Level > 0 && needPoint > 100 {
		needToken = (product.Price - needPoint) * 90 / 100
	} else {
		needToken = product.Price - needPoint
	}

	if user.GetToken() < needToken {
		fmt.Println(fmt.Sprintf("[MSG] User %d has not enough token to buy product %d", userID, productID))
		return -1, errors.New("not enough token")
	}

	user.UsePoint(needPoint)
	user.UseToken(needToken)
	fmt.Println(fmt.Sprintf("[MSG] User %d has bought product %d use price %d, point %d", userID, productID, needToken, needPoint))
	return needToken, nil
}
