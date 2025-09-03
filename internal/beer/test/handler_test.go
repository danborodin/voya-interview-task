package test

import (
	"encoding/json"
	"errors"
	backendbeer "interview-go/backend/client"
	"interview-go/internal/beer"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
)

func setupEcho() *echo.Echo {
	e := echo.New()
	return e
}

func TestFilteredBeers_DefaultFilteringAndSorting(t *testing.T) {
	e := setupEcho()

	defaultFilters := beer.BeerFilter{
		IncludeIpa:   true,
		Year:         2015,
		HasFood:      "wolf",
		AbvSortOrder: "asc",
	}

	svc := &mockService{
		GetDefaultFiltersFunc: func() beer.BeerFilter { return defaultFilters },
		GetFilteredBeersFunc: func(filters string) ([]backendbeer.BeerResponse, error) {
			return []backendbeer.BeerResponse{
				{ID: 1, Name: "Ruby IPA", FirstBrewed: "2016-01", ABV: 6.0, FoodPairing: []string{"wolf", "steak"}},
				{ID: 2, Name: "Lager", FirstBrewed: "2017-05", ABV: 4.5, FoodPairing: []string{"wolf"}},
				{ID: 3, Name: "Pale Ale", FirstBrewed: "2018-03", ABV: 5.2, FoodPairing: []string{"pizza"}},
				{ID: 4, Name: "Imperial IPA", FirstBrewed: "2014-12", ABV: 8.5, FoodPairing: []string{"wolf"}},
			}, nil
		},
	}

	h := beer.NewHandler(svc)
	req := httptest.NewRequest(http.MethodGet, "/getFiltered", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	require.NoError(t, h.FilteredBeers(c))
	require.Equal(t, http.StatusOK, rec.Code)

	var got []backendbeer.BeerResponse
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &got))

	require.Len(t, got, 1)
	require.Equal(t, 1, got[0].ID)
}

func TestFilteredBeers_QueryOverrides(t *testing.T) {
	e := setupEcho()

	defaultFilters := beer.BeerFilter{IncludeIpa: true, Year: 2015, HasFood: "wolf", AbvSortOrder: "asc"}

	svc := &mockService{
		GetDefaultFiltersFunc: func() beer.BeerFilter { return defaultFilters },
		GetFilteredBeersFunc: func(filters string) ([]backendbeer.BeerResponse, error) {
			return []backendbeer.BeerResponse{
				{ID: 10, Name: "Dark Lager", FirstBrewed: "2020-01", ABV: 7.2, FoodPairing: []string{"fish"}},
				{ID: 11, Name: "Summer Ale", FirstBrewed: "2021-06", ABV: 4.0, FoodPairing: []string{"fish"}},
				{ID: 12, Name: "Tropical IPA", FirstBrewed: "2022-03", ABV: 5.5, FoodPairing: []string{"fish"}},
			}, nil
		},
	}

	h := beer.NewHandler(svc)
	req := httptest.NewRequest(http.MethodGet, "/getFiltered?includeIpa=false&year=2020&hasFood=fish&abvSortOrder=desc", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	require.NoError(t, h.FilteredBeers(c))
	require.Equal(t, http.StatusOK, rec.Code)

	var got []backendbeer.BeerResponse
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &got))

	require.Len(t, got, 2)
	require.Equal(t, 12, got[0].ID) // higher ABV first due to desc
	require.Equal(t, 11, got[1].ID)
}

func TestFilteredBeers_EmptyServiceResponse(t *testing.T) {
	e := setupEcho()
	defaultFilters := beer.BeerFilter{IncludeIpa: true, Year: 2015, HasFood: "wolf", AbvSortOrder: "asc"}
	svc := &mockService{
		GetDefaultFiltersFunc: func() beer.BeerFilter { return defaultFilters },
		GetFilteredBeersFunc: func(filters string) ([]backendbeer.BeerResponse, error) {
			return []backendbeer.BeerResponse{}, nil
		},
	}
	h := beer.NewHandler(svc)
	req := httptest.NewRequest(http.MethodGet, "/getFiltered", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := h.FilteredBeers(c)
	require.NoError(t, err)
	require.Equal(t, http.StatusNoContent, rec.Code)
}

func TestFilteredBeers_ServiceError(t *testing.T) {
	e := setupEcho()
	svc := &mockService{
		GetDefaultFiltersFunc: func() beer.BeerFilter { return beer.BeerFilter{} },
		GetFilteredBeersFunc: func(filters string) ([]backendbeer.BeerResponse, error) {
			return nil, errors.New("some error")
		},
	}
	h := beer.NewHandler(svc)
	req := httptest.NewRequest(http.MethodGet, "/getFiltered", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := h.FilteredBeers(c)
	require.Error(t, err)
	eHTTPErr, ok := err.(*echo.HTTPError)
	require.True(t, ok)
	require.Equal(t, http.StatusInternalServerError, eHTTPErr.Code)
}

func TestFilteredBeers_InvalidQueryParams(t *testing.T) {
	e := setupEcho()
	svc := &mockService{
		GetDefaultFiltersFunc: func() beer.BeerFilter { return beer.BeerFilter{} },
	}
	h := beer.NewHandler(svc)

	cases := []string{
		"includeIpa=notabool",
		"year=notanint",
	}

	for _, q := range cases {
		req := httptest.NewRequest(http.MethodGet, "/getFiltered?"+q, nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := h.FilteredBeers(c)
		require.Error(t, err, q)
		var httpErr *echo.HTTPError
		require.True(t, errors.As(err, &httpErr), q)
		require.Equal(t, http.StatusBadRequest, httpErr.Code, q)
	}
}
