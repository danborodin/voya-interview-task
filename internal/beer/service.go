package beer

import (
	"errors"
	"fmt"
	backendbeer "interview-go/backend/client"
	"interview-go/config"
	"interview-go/internal/cache"
	"log"

	"golang.org/x/time/rate"
)

type Service interface {
	GetAllBeers() ([]backendbeer.BeerResponse, error)
	GetFilteredBeers(filters string) ([]backendbeer.BeerResponse, error)
	GetDefaultFilters() BeerFilter
}

type service struct {
	cache       cache.Cache
	client      backendbeer.Client
	rateLimiter *rate.Limiter // for api rate limit simulation
}

type BeerFilter struct {
	IncludeIpa   bool
	Year         int
	HasFood      string
	AbvSortOrder string
}

var (
	ErrRateLimitExceeded = errors.New("api rate limit exceeded")
)

func NewService(client backendbeer.Client, cfg *config.Configuration) Service {
	return &service{
		cache:       cache.NewInMemory(cfg.Cache.TTL, cfg.Cache.ClearTicker),
		client:      client,
		rateLimiter: rate.NewLimiter(rate.Every(cfg.ApiRateLimit.Rate), cfg.ApiRateLimit.Burst),
	}
}

func (s *service) GetAllBeers() ([]backendbeer.BeerResponse, error) {
	return s.client.ListBeers()
}

func (s *service) GetFilteredBeers(filters string) ([]backendbeer.BeerResponse, error) {
	if filters != "" {
		cachedValue, err := s.cache.Get(filters)
		if err != nil {
			if err != cache.ErrCacheMiss && err != cache.ErrTTLExpired {
				return nil, err
			}
		}
		if cachedValue != nil {
			beers, ok := cachedValue.([]backendbeer.BeerResponse)
			if !ok {
				return nil, errors.New("malformed data type in cache")
			}
			return beers, nil
		}
	}

	// simulating api rate limit
	if !s.rateLimiter.Allow() {
		return nil, errors.New("rate limit exceeded")
	}

	beers, err := s.client.ListBeers()
	if err != nil {
		return nil, err
	}

	err = s.cache.Set(filters, beers)
	if err != nil {
		// do not return here, just log it
		log.Println(err)
	}

	return beers, nil
}

func (s *service) GetDefaultFilters() BeerFilter {
	filter := BeerFilter{
		IncludeIpa:   true,
		Year:         2015,
		HasFood:      "wolf",
		AbvSortOrder: "asc",
	}

	return filter
}

func (bf *BeerFilter) String() string {
	return fmt.Sprintf("%t%d%s%s", bf.IncludeIpa, bf.Year, bf.HasFood, bf.AbvSortOrder)
}
