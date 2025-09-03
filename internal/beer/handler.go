package beer

import (
	backendbeer "interview-go/backend/client"
	"net/http"
	"slices"
	"sort"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
)

type HTTPHandler interface {
	ListAllBeers(c echo.Context) error
	//TODO
	FilteredBeers(c echo.Context) error
}

type beerHandler struct {
	service Service
}

func (h *beerHandler) FilteredBeers(c echo.Context) error {
	var err error

	filters := h.service.GetDefaultFilters()

	includeIpa := c.QueryParam("includeIpa")
	if includeIpa != "" {
		filters.IncludeIpa, err = strconv.ParseBool(includeIpa)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}
	}

	year := c.QueryParam("year")
	if year != "" {
		filters.Year, err = strconv.Atoi(year)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}
	}

	food := c.QueryParam("hasFood")
	if food != "" {
		filters.HasFood = food
	}

	abvSortOrder := c.QueryParam("abvSortOrder")
	if abvSortOrder != "" {
		filters.AbvSortOrder = abvSortOrder
	}

	var filteredResp []backendbeer.BeerResponse

	resp, err := h.service.GetFilteredBeers(filters.String())
	if err != nil {
		if err == ErrRateLimitExceeded {
			return echo.NewHTTPError(http.StatusTooManyRequests, err)
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	if len(resp) == 0 {
		return c.NoContent(http.StatusNoContent)
	}

	for _, v := range resp {
		if filters.IncludeIpa && !strings.Contains(strings.ToLower(v.Name), "ipa") {
			continue
		}

		if !(extractYear(v.FirstBrewed) > filters.Year) {
			continue
		}

		if !slices.Contains(v.FoodPairing, strings.ToLower(filters.HasFood)) {
			continue
		}

		filteredResp = append(filteredResp, v)
	}

	if strings.ToLower(filters.AbvSortOrder) == "desc" {
		sort.Slice(filteredResp, func(i, j int) bool {
			return filteredResp[i].ABV > filteredResp[j].ABV
		})
	} else if strings.ToLower(filters.AbvSortOrder) == "asc" {
		sort.Slice(filteredResp, func(i, j int) bool {
			return filteredResp[i].ABV < filteredResp[j].ABV
		})
	}

	return c.JSON(http.StatusOK, filteredResp)

}

func NewHandler(service Service) HTTPHandler {
	return &beerHandler{
		service: service,
	}
}

func (h *beerHandler) ListAllBeers(c echo.Context) error {
	resp, err := h.service.GetAllBeers()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	if len(resp) == 0 {
		return c.NoContent(http.StatusNoContent)
	}

	return c.JSON(http.StatusOK, resp)
}

func extractYear(firstBrewed string) int {
	// Acceptă"YYYY-MM"
	s := strings.TrimSpace(firstBrewed)
	if len(s) < 4 {
		return 0
	}
	// ia doar primele 4 caractere (dacă sunt cifre)
	yearStr := s[:4]
	if y, err := strconv.Atoi(yearStr); err == nil {
		return y
	}
	return 0
}
