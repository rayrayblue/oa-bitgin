package repository

import (
	"errors"
	"oa-bitgin/pkg/domain"
	"sync/atomic"
)

// Use to store all activity, use id as key
type activityStore struct {
	BuyTokenActivitiesIDCounter   atomic.Value
	BuyProductActivitiesIDCounter atomic.Value
	BuyTokenActivities            map[int]domain.BuyTokenActivity
	BuyProductActivities          map[int]domain.BuyProductActivity
}

func (s *activityStore) init() {
	s.BuyTokenActivitiesIDCounter.Store(0)
	s.BuyProductActivitiesIDCounter.Store(0)
	s.BuyTokenActivities = make(map[int]domain.BuyTokenActivity)
	s.BuyProductActivities = make(map[int]domain.BuyProductActivity)
}

type activityRepository struct {
	store *activityStore
}

func NewActivityRepository() *activityRepository {
	store := &activityStore{}
	store.init()
	return &activityRepository{store: store}
}

func (a *activityRepository) ListBuyTokenActivity() ([]domain.BuyTokenActivity, error) {
	var rtn []domain.BuyTokenActivity
	for _, v := range a.store.BuyTokenActivities {
		rtn = append(rtn, v)
	}
	if len(rtn) == 0 {
		return make([]domain.BuyTokenActivity, 0), nil
	}
	return rtn, nil
}

func (a *activityRepository) AddBuyTokenActivity(activity domain.BuyTokenActivity) (int, error) {
	id := a.store.BuyProductActivitiesIDCounter.Load().(int)
	id++
	a.store.BuyProductActivitiesIDCounter.Store(id)
	activity.SetID(id)
	a.store.BuyTokenActivities[id] = activity
	return id, nil
}

func (a *activityRepository) AddBuyProductActivity(activity domain.BuyProductActivity) (int, error) {
	id := a.store.BuyProductActivitiesIDCounter.Load().(int)
	id++
	a.store.BuyProductActivitiesIDCounter.Store(id)
	activity.SetID(id)
	a.store.BuyProductActivities[id] = activity
	return id, nil
}

func (a *activityRepository) GetBuyProductActivity(id int) (domain.BuyProductActivity, error) {
	if activity, ok := a.store.BuyProductActivities[id]; ok {
		return activity, nil
	} else {
		return domain.BuyProductActivity{}, errors.New("activity not found")
	}
}
