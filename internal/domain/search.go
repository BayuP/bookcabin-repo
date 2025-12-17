package domain

type SearchRequest struct {
	Origin        string `json:"origin"`
	Destination   string `json:"destination"`
	DepartureDate string `json:"departure_date"`
	Passengers    int    `json:"passenger"`
	CabinClass    string `json:"cabin_class"`

	// filter
	MinPrice    int64    `json:"min_price,omitempty"`
	MaxPrice    int64    `json:"max_price,omitempty"`
	MaxStops    int      `json:"max_stops,omitempty"`
	Airlines    []string `json:"airlines,omitempty"`
	MaxDuration int      `json:"max_duration,omitempty"`
	EarliestDep string   `json:"earliest_departure,omitempty"`
	LatestDep   string   `json:"latest_departure,omitempty"`
	EarliestArr string   `json:"earliest_arrival,omitempty"`
	LatestArr   string   `json:"latest_arrival,omitempty"`

	// sort
	SortBy string `json:"sort_by,omitempty"`
}

type SearchResult struct {
	Flights            []Flight
	CacheHit           bool
	ProvidersQueried   int
	ProvidersSucceeded int
	ProvidersFailed    int
}
