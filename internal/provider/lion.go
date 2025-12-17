package provider

import (
	"bookcabin/internal/common"
	"bookcabin/internal/domain"
	"encoding/json"
	"errors"
	"net/http"
)

type LionResponse struct {
	Success bool     `json:"success"`
	Data    LionData `json:"data"`
}

type LionData struct {
	AvailableFlights []LionFlight `json:"available_flights"`
}

type LionFlight struct {
	ID         string        `json:"id"`
	Carrier    LionCarrier   `json:"carrier"`
	Route      LionRoute     `json:"route"`
	Schedule   LionSchedule  `json:"schedule"`
	FlightTime int           `json:"flight_time"`
	IsDirect   bool          `json:"is_direct"`
	StopCount  int           `json:"stop_count,omitempty"`
	Layovers   []LionLayover `json:"layovers,omitempty"`
	Pricing    LionPricing   `json:"pricing"`
	SeatsLeft  int           `json:"seats_left"`
	PlaneType  string        `json:"plane_type"`
	Services   LionServices  `json:"services"`
}

type LionCarrier struct {
	Name string `json:"name"`
	IATA string `json:"iata"`
}

type LionRoute struct {
	From LionAirport `json:"from"`
	To   LionAirport `json:"to"`
}

type LionAirport struct {
	Code string `json:"code"`
	Name string `json:"name"`
	City string `json:"city"`
}

type LionSchedule struct {
	Departure         string `json:"departure"`
	DepartureTimezone string `json:"departure_timezone"`
	Arrival           string `json:"arrival"`
	ArrivalTimezone   string `json:"arrival_timezone"`
}

type LionPricing struct {
	Total    int    `json:"total"`
	Currency string `json:"currency"`
	FareType string `json:"fare_type"`
}

type LionServices struct {
	WifiAvailable bool        `json:"wifi_available"`
	MealsIncluded bool        `json:"meals_included"`
	Baggage       LionBaggage `json:"baggage_allowance"`
}

type LionBaggage struct {
	Cabin string `json:"cabin"`
	Hold  string `json:"hold"`
}

type LionLayover struct {
	Airport         string `json:"airport"`
	DurationMinutes int    `json:"duration_minutes"`
}

type LionAirProvider struct {
	BaseURL string
	Client  *http.Client
}

func (l *LionAirProvider) Name() string {
	return "Lion Air"
}

func (l *LionAirProvider) Search(req domain.SearchRequest) ([]domain.Flight, error) {
	var lionRaw LionResponse
	resp, err := l.Client.Get(l.BaseURL + "/lion/search")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(&lionRaw); err != nil {
		return nil, err
	}

	if !lionRaw.Success {
		return nil, errors.New("lion air api returned failure")
	}

	flights := []domain.Flight{}
	for _, r := range lionRaw.Data.AvailableFlights {

		dep, err := common.ParseTimeWithTZ(
			r.Schedule.Departure,
			r.Schedule.DepartureTimezone,
		)
		if err != nil {
			continue
		}

		arr, err := common.ParseTimeWithTZ(
			r.Schedule.Arrival,
			r.Schedule.ArrivalTimezone,
		)
		if err != nil {
			continue
		}

		priceIDR, err := common.ParsePriceToIDR(
			r.Pricing.Total,
			r.Pricing.Currency,
		)
		if err != nil {
			continue
		}

		stops := 0
		if !r.IsDirect {
			stops = r.StopCount
		}

		flights = append(flights, domain.Flight{
			FlightCode:     r.ID,
			Airline:        r.Carrier.Name,
			AirlineCode:    r.Carrier.IATA,
			Origin:         r.Route.From.Code,
			Destination:    r.Route.To.Code,
			DepartureTime:  dep,
			ArrivalTime:    arr,
			DurationMin:    r.FlightTime,
			Stops:          stops,
			PriceIDR:       priceIDR,
			AvailableSeats: r.SeatsLeft,
			Aircraft:       r.PlaneType,
			Baggage: r.Services.Baggage.Cabin + " cabin, " +
				r.Services.Baggage.Hold + " checked",
		})
	}

	return flights, nil
}
