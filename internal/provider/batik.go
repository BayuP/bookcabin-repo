package provider

import (
	"bookcabin/internal/common"
	"bookcabin/internal/domain"
	"encoding/json"
	"errors"
	"net/http"
)

// ===== RAW BATIK RESPONSE =====
type BatikAirResponse struct {
	Code    int              `json:"code"`
	Message string           `json:"message"`
	Results []BatikAirFlight `json:"results"`
}

type BatikAirFlight struct {
	FlightNumber      string               `json:"flightNumber"`
	AirlineName       string               `json:"airlineName"`
	AirlineIATA       string               `json:"airlineIATA"`
	Origin            string               `json:"origin"`
	Destination       string               `json:"destination"`
	DepartureDateTime string               `json:"departureDateTime"`
	ArrivalDateTime   string               `json:"arrivalDateTime"`
	TravelTime        string               `json:"travelTime"`
	NumberOfStops     int                  `json:"numberOfStops"`
	Connections       []BatikAirConnection `json:"connections,omitempty"`
	Fare              BatikAirFare         `json:"fare"`
	SeatsAvailable    int                  `json:"seatsAvailable"`
	AircraftModel     string               `json:"aircraftModel"`
	BaggageInfo       string               `json:"baggageInfo"`
	OnboardServices   []string             `json:"onboardServices"`
}

type BatikAirConnection struct {
	StopAirport  string `json:"stopAirport"`
	StopDuration string `json:"stopDuration"`
}

type BatikAirFare struct {
	BasePrice  int    `json:"basePrice"`
	Taxes      int    `json:"taxes"`
	TotalPrice int    `json:"totalPrice"`
	Currency   string `json:"currencyCode"`
	Class      string `json:"class"`
}

type BatikProvider struct {
	BaseURL string
	Client  *http.Client
}

func (b *BatikProvider) Name() string { return "Batik Air" }

func (b *BatikProvider) Search(req domain.SearchRequest) ([]domain.Flight, error) {
	resp, err := b.Client.Get(b.BaseURL + "/batik/search")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var batikRaw BatikAirResponse
	if err := json.NewDecoder(resp.Body).Decode(&batikRaw); err != nil {
		return nil, err
	}

	if batikRaw.Code != CodeSuccess {
		return nil, errors.New("batik air api returned failure")
	}

	flights := make([]domain.Flight, 0)
	for _, r := range batikRaw.Results {
		dep, _ := common.ParseFlexibleTime(r.DepartureDateTime)
		arr, _ := common.ParseFlexibleTime(r.ArrivalDateTime)

		price, _ := common.ParsePriceToIDR(r.Fare.TotalPrice, r.Fare.Currency)

		flights = append(flights, domain.Flight{
			FlightCode:     r.FlightNumber,
			Airline:        r.AirlineName,
			AirlineCode:    r.AirlineIATA,
			Origin:         r.Origin,
			Destination:    r.Destination,
			DepartureTime:  dep,
			ArrivalTime:    arr,
			DurationMin:    int(arr.Sub(dep).Minutes()),
			Stops:          r.NumberOfStops,
			PriceIDR:       price,
			AvailableSeats: r.SeatsAvailable,
			Aircraft:       r.AircraftModel,
			Baggage:        r.BaggageInfo,
			Amenities:      r.OnboardServices,
		})
	}

	return flights, nil
}
