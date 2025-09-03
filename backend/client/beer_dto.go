package client

type BeerResponse struct {
	ID            int         `json:"id"`
	Name          string      `json:"name"`
	Tagline       string      `json:"tagline"`
	FirstBrewed   string      `json:"first_brewed"` // "YYYY-MM"
	Description   string      `json:"description"`
	ABV           float64     `json:"abv"`
	Ingredients   Ingredients `json:"ingredients"`
	FoodPairing   []string    `json:"food_pairing"`
	BrewersTips   string      `json:"brewers_tips"`
	ContributedBy string      `json:"contributed_by"`
}

type BeerRequest struct {
	BeerName     string `json:"beerName,omitempty"`
	Yeast        string `json:"yeast,omitempty"`
	BrewedBefore string `json:"brewedBefore,omitempty"` // format yyyy-mm
	BrewedAfter  string `json:"brewedAfter,omitempty"`  // format yyyy-mm
	Hops         string `json:"hops,omitempty"`
	Malt         string `json:"malt,omitempty"`
	Food         string `json:"food,omitempty"`
}

type Ingredients struct {
	Malt  []Malt `json:"malt"`
	Hops  []Hops `json:"hops"`
	Yeast string `json:"yeast"`
}

type Temp struct {
	Value float64 `json:"value"`
	Unit  string  `json:"unit"`
}

type Amount struct {
	Value float64 `json:"value"`
	Unit  string  `json:"unit"`
}

type MashTemp struct {
	Temp     Temp `json:"temp"`
	Duration *int `json:"duration,omitempty"` // optional (Integer in Java)
}

type Fermentation struct {
	Temp Temp `json:"temp"`
}

type Malt struct {
	Name   string `json:"name"`
	Amount Amount `json:"amount"`
}

type Hops struct {
	Name      string `json:"name"`
	Amount    Amount `json:"amount"`
	Add       string `json:"add"`
	Attribute string `json:"attribute"`
}
