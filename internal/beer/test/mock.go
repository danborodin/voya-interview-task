package test

import (
	backendbeer "interview-go/backend/client"
	"interview-go/internal/beer"
)

type mockService struct {
	GetAllBeersFunc       func() ([]backendbeer.BeerResponse, error)
	GetFilteredBeersFunc  func(filters string) ([]backendbeer.BeerResponse, error)
	GetDefaultFiltersFunc func() beer.BeerFilter

	LastFiltersString string
}

func (m *mockService) GetAllBeers() ([]backendbeer.BeerResponse, error) {
	if m.GetAllBeersFunc != nil {
		return m.GetAllBeersFunc()
	}
	return nil, nil
}

func (m *mockService) GetFilteredBeers(filters string) ([]backendbeer.BeerResponse, error) {
	m.LastFiltersString = filters
	if m.GetFilteredBeersFunc != nil {
		return m.GetFilteredBeersFunc(filters)
	}
	return nil, nil
}

func (m *mockService) GetDefaultFilters() beer.BeerFilter {
	if m.GetDefaultFiltersFunc != nil {
		return m.GetDefaultFiltersFunc()
	}
	return beer.BeerFilter{}
}
