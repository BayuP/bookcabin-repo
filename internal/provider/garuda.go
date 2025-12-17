package provider

import (
	"bookcabin/internal/common"
	"bookcabin/internal/domain"
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

type GarudaResponse struct {
	Status  string         `json:"status"`
	Flights []GarudaFlight `json:"flights"`
}

type GarudaFlight struct {
	FlightID       string          `json:"flight_id"`
	Airline        string          `json:"airline"`
	AirlineCode    string          `json:"airline_code"`
	Departure      GarudaLocation  `json:"departure"`
	Arrival        GarudaLocation  `json:"arrival"`
	DurationMin    int             `json:"duration_minutes"`
	Stops          int             `json:"stops"`
	Aircraft       string          `json:"aircraft"`
	Price          GarudaPrice     `json:"price"`
	Segments       []GarudaSegment `json:"segments,omitempty"`
	AvailableSeats int             `json:"available_seats"`
	FareClass      string          `json:"fare_class"`
	Baggage        GarudaBaggage   `json:"baggage"`
	Amenities      []string        `json:"amenities,omitempty"`
}

type GarudaLocation struct {
	Airport  string `json:"airport"`
	City     string `json:"city,omitempty"`
	Time     string `json:"time"`
	Terminal string `json:"terminal,omitempty"`
}

type GarudaPrice struct {
	Amount   int    `json:"amount"`
	Currency string `json:"currency"`
}

type GarudaBaggage struct {
	CarryOn int `json:"carry_on"`
	Checked int `json:"checked"`
}

type GarudaSegment struct {
	FlightNumber    string       `json:"flight_number"`
	Departure       SegmentPoint `json:"departure"`
	Arrival         SegmentPoint `json:"arrival"`
	DurationMinutes int          `json:"duration_minutes"`
	LayoverMinutes  int          `json:"layover_minutes,omitempty"`
}

type SegmentPoint struct {
	Airport string `json:"airport"`
	Time    string `json:"time"`
}

type GarudaProvider struct {
	BaseURL string
	Client  *http.Client
}

func (g *GarudaProvider) Name() string { return "Garuda Indonesia" }

func (g *GarudaProvider) Search(req domain.SearchRequest) ([]domain.Flight, error) {
	var garudaRaw GarudaResponse
	resp, err := g.Client.Get(g.BaseURL + "/garuda/search")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(bodyBytes, &garudaRaw); err != nil {
		return nil, err
	}

	if garudaRaw.Status != StatusSuccess {
		return nil, errors.New("garuda api returned failure")
	}

	flights := make([]domain.Flight, 0)
	for _, r := range garudaRaw.Flights {
		dep, _ := common.ParseFlexibleTime(r.Departure.Time)
		arr, _ := common.ParseFlexibleTime(r.Arrival.Time)
		price, _ := common.ParsePriceToIDR(r.Price.Amount, r.Price.Currency)

		flights = append(flights, domain.Flight{
			FlightCode:     r.FlightID,
			Airline:        r.Airline,
			AirlineCode:    r.AirlineCode,
			Origin:         r.Departure.Airport,
			Destination:    r.Arrival.Airport,
			DepartureTime:  dep,
			ArrivalTime:    arr,
			DurationMin:    r.DurationMin,
			Stops:          r.Stops,
			PriceIDR:       price,
			AvailableSeats: r.AvailableSeats,
			Aircraft:       r.Aircraft,
			Amenities:      r.Amenities,
		})
	}

	return flights, nil
}
