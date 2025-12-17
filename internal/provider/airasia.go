package provider

import (
	"bookcabin/internal/common"
	"bookcabin/internal/domain"
	"encoding/json"
	"errors"
	"net/http"
)

type AirAsiaResponse struct {
	Status  string          `json:"status"`
	Flights []AirAsiaFlight `json:"flights"`
}

type AirAsiaFlight struct {
	FlightCode   string        `json:"flight_code"`
	Airline      string        `json:"airline"`
	FromAirport  string        `json:"from_airport"`
	ToAirport    string        `json:"to_airport"`
	DepartTime   string        `json:"depart_time"`
	ArriveTime   string        `json:"arrive_time"`
	DurationHrs  float64       `json:"duration_hours"`
	DirectFlight bool          `json:"direct_flight"`
	Stops        []AirAsiaStop `json:"stops,omitempty"`
	PriceIDR     int           `json:"price_idr"`
	Seats        int           `json:"seats"`
	CabinClass   string        `json:"cabin_class"`
	BaggageNote  string        `json:"baggage_note"`
}

type AirAsiaStop struct {
	Airport         string `json:"airport"`
	WaitTimeMinutes int    `json:"wait_time_minutes"`
}

type AirAsiaProvider struct {
	BaseURL string
	Client  *http.Client
}

func (a *AirAsiaProvider) Name() string { return "AirAsia" }

func (a *AirAsiaProvider) Search(req domain.SearchRequest) ([]domain.Flight, error) {
	resp, err := a.Client.Get(a.BaseURL + "/airasia/search")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var airAsiaRaw AirAsiaResponse
	if err := json.NewDecoder(resp.Body).Decode(&airAsiaRaw); err != nil {
		return nil, err
	}

	if airAsiaRaw.Status != StatusOK {
		return nil, errors.New("air asia api returned failure")
	}

	flights := []domain.Flight{}
	for _, r := range airAsiaRaw.Flights {
		dep, _ := common.ParseFlexibleTime(r.DepartTime)
		arr, _ := common.ParseFlexibleTime(r.ArriveTime)

		stops := 0
		if !r.DirectFlight {
			stops = 1
		}

		flights = append(flights, domain.Flight{
			FlightCode:     r.FlightCode,
			Airline:        r.Airline,
			AirlineCode:    r.FlightCode[:2],
			Origin:         r.FromAirport,
			Destination:    r.ToAirport,
			DepartureTime:  dep,
			ArrivalTime:    arr,
			DurationMin:    int(r.DurationHrs * 60),
			Stops:          stops,
			PriceIDR:       int64(r.PriceIDR),
			AvailableSeats: r.Seats,
			Baggage:        r.BaggageNote,
		})
	}

	return flights, nil
}
