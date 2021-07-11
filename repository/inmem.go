package repository

import (
	"errors"
	"github.com/cecilozaur/tnemesum/domain"
	"sync"
)

var (
	emptyCity       = domain.City{}
	errCityNotFound = errors.New("not found")
)

type InMemRepo struct {
	items sync.Map
}

func NewInMem() *InMemRepo {
	return &InMemRepo{
		items: sync.Map{},
	}
}

func (m *InMemRepo) GetItems() []domain.City {
	items := make([]domain.City, 0)
	m.items.Range(func(key, value interface{}) bool {
		items = append(items, value.(domain.City))
		return true
	})

	return items
}

func (m *InMemRepo) Get(key uint64) (domain.City, error) {
	item, ok := m.items.Load(key)
	if !ok {
		return emptyCity, errCityNotFound
	}

	return item.(domain.City), nil
}

func (m *InMemRepo) GetForecast(key uint64) (domain.Forecast, error) {
	item, err := m.Get(key)
	if err != nil {
		return nil, err
	}

	return item.Forecast, nil
}

func (m *InMemRepo) Store(key uint64, item domain.City) {
	m.items.Store(key, item)
}

func (m *InMemRepo) UpdateForecast(key uint64, forecast domain.Forecast) bool {
	item, err := m.Get(key)
	if err != nil {
		return false
	}

	item.Forecast = forecast

	m.items.Store(key, item)

	return true
}
