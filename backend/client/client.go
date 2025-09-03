package client

import (
	"math"
	"math/rand"
	"time"

	"github.com/brianvoe/gofakeit/v7"
)

type Client interface {
	ListBeers() ([]BeerResponse, error)
}

type FakeBeerClient struct {
	count int
}

func NewFakeBeerClient(count int) *FakeBeerClient {
	if count <= 0 {
		count = 50
	}

	gofakeit.Seed(time.Now().UnixNano())
	return &FakeBeerClient{count: count}
}

func (c *FakeBeerClient) ListBeers() ([]BeerResponse, error) {
	out := make([]BeerResponse, 0, c.count)
	for i := 0; i < c.count; i++ {
		out = append(out, c.fakeBeer(i+1))
	}
	return out, nil
}

func (c *FakeBeerClient) fakeBeer(id int) BeerResponse {
	year := rand.Intn(34-4+1) + (time.Now().Year() - 34)
	month := time.Month(rand.Intn(12) + 1)
	firstBrewed := time.Date(year, month, 1, 0, 0, 0, 0, time.UTC).Format("2006-01")

	abv := math.Round((1.0+rand.Float64()*14.0)*10) / 10

	foods := make([]string, rand.Intn(5)+1)
	for i := range foods {
		foods[i] = gofakeit.MinecraftAnimal()
	}

	return BeerResponse{
		ID:          id,
		Name:        gofakeit.BeerName(),
		Tagline:     gofakeit.BeerStyle(),
		FirstBrewed: firstBrewed,
		Description: gofakeit.Sentence(12),
		ABV:         abv,
		Ingredients: Ingredients{
			Malt: []Malt{
				{Name: "Extra Pale", Amount: Amount{Value: 5, Unit: "kilograms"}},
			},
			Hops: []Hops{
				{Name: "Cascade", Amount: Amount{Value: 25, Unit: "grams"}, Add: "start", Attribute: "bitter"},
				{Name: "Citra", Amount: Amount{Value: 25, Unit: "grams"}, Add: "end", Attribute: "aroma"},
			},
			Yeast: "Wyeast 1056 - American Ale™",
		},
		FoodPairing:   foods,
		BrewersTips:   gofakeit.Sentence(8),
		ContributedBy: gofakeit.Name(),
	}
}
