package domain

import "time"

type SortOption string

const (
	SortPriceAsc     SortOption = "price_asc"
	SortPriceDesc    SortOption = "price_desc"
	SortDurationAsc  SortOption = "duration_asc"
	SortDurationDesc SortOption = "duration_desc"
	SortDepartureAsc SortOption = "departure_asc"
	SortArrivalAsc   SortOption = "arrival_asc"
	SortBestValue    SortOption = "best_value"
)

type Flight struct {
	FlightCode     string
	Airline        string
	AirlineCode    string
	Origin         string
	Destination    string
	DepartureTime  time.Time
	ArrivalTime    time.Time
	DurationMin    int
	Stops          int
	PriceIDR       int64
	AvailableSeats int
	Aircraft       string
	Baggage        string
	Amenities      []string
}

type FlightFilter struct {
	MinPrice          int64
	MaxPrice          int64
	MaxStops          int
	EarliestDeparture time.Time
	LatestDeparture   time.Time
	EarliestArrival   time.Time
	LatestArrival     time.Time
	Airlines          []string
	MaxDuration       int
}
