package server

import (
	handler "interview-go/internal/beer"

	"github.com/labstack/echo/v4"
)

func BeerRoutes(g *echo.Group, h handler.HTTPHandler) {
	g.GET("/getAll", h.ListAllBeers)
	g.GET("/getFiltered", h.FilteredBeers)
}
