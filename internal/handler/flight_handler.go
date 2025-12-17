package handler

import (
	"bookcabin/internal/domain"
	"bookcabin/internal/service"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type FlightHandler struct {
	FlightService *service.SearchFlightsUseCase
}

func NewFlightHandler(fs *service.SearchFlightsUseCase) FlightHandler {
	return FlightHandler{
		FlightService: fs,
	}
}

// SearchFlights godoc
// @Summary      Search flights
// @Description  Search flights with filters and sorting
// @Tags         Flights
// @Accept       json
// @Produce      json
//
// @Param origin query string true "Origin airport code (e.g. CGK)"
// @Param destination query string true "Destination airport code (e.g. DPS)"
// @Param departure_date query string true "Departure date (YYYY-MM-DD)"
// @Param passengers query int false "Number of passengers"
// @Param cabin_class query string false "Cabin class (economy, business)"
//
// @Param min_price query int false "Minimum price (IDR)"
// @Param max_price query int false "Maximum price (IDR)"
// @Param max_stops query int false "Maximum stops"
// @Param max_duration query int false "Maximum duration (minutes)"
// @Param airlines query string false "Airline codes (CSV or repeated), e.g. GA,ID"
//
// @Param earliest_departure query string false "Earliest departure time (HH:MM)"
// @Param latest_departure query string false "Latest departure time (HH:MM)"
// @Param earliest_arrival query string false "Earliest arrival time (HH:MM)"
// @Param latest_arrival query string false "Latest arrival time (HH:MM)"
//
// @Param sort_by query string false "Sort option" Enums(price_asc,price_desc,duration_asc,duration_desc,departure_asc,arrival_asc,best_value)
//
// @Success 200 {object} domain.FlightSearchResponse
// @Failure 400 {string} string "Bad Request"
// @Failure 405 {string} string "Method Not Allowed"
// @Failure 500 {string} string "Internal Server Error"
//
// @Router /search [get]
func (h *FlightHandler) Search(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	start := time.Now()
	q := r.URL.Query()

	req := domain.SearchRequest{
		Origin:        q.Get("origin"),
		Destination:   q.Get("destination"),
		DepartureDate: q.Get("departure_date"),
		CabinClass:    q.Get("cabin_class"),

		// optional filters (strings first)
		EarliestDep: q.Get("earliest_departure"),
		LatestDep:   q.Get("latest_departure"),
		EarliestArr: q.Get("earliest_arrival"),
		LatestArr:   q.Get("latest_arrival"),
		SortBy:      q.Get("sort_by"),
	}

	// passengers
	if v := q.Get("passengers"); v != "" {
		if p, err := strconv.Atoi(v); err == nil {
			req.Passengers = p
		}
	}

	// prices
	if v := q.Get("min_price"); v != "" {
		if p, err := strconv.ParseInt(v, 10, 64); err == nil {
			req.MinPrice = p
		}
	}

	if v := q.Get("max_price"); v != "" {
		if p, err := strconv.ParseInt(v, 10, 64); err == nil {
			req.MaxPrice = p
		}
	}

	// stops
	if v := q.Get("max_stops"); v != "" {
		if s, err := strconv.Atoi(v); err == nil {
			req.MaxStops = s
		}
	} else {
		req.MaxStops = -1 // default unset
	}

	// duration
	if v := q.Get("max_duration"); v != "" {
		if d, err := strconv.Atoi(v); err == nil {
			req.MaxDuration = d
		}
	}

	if airlines := q["airlines"]; len(airlines) > 0 {
		var parsed []string
		for _, a := range airlines {
			parsed = append(parsed, strings.Split(a, ",")...)
		}
		req.Airlines = parsed
	}

	if req.Origin == "" || req.Destination == "" || req.DepartureDate == "" {
		http.Error(w, "missing required query parameters", http.StatusBadRequest)
		return
	}

	result, err := h.FlightService.Execute(r.Context(), req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := domain.FlightSearchResponse{
		SearchCriteria: domain.SearchCriteria{
			Origin:        req.Origin,
			Destination:   req.Destination,
			DepartureDate: req.DepartureDate,
			Passengers:    req.Passengers,
			CabinClass:    req.CabinClass,
		},
		Metadata: domain.Metadata{
			TotalResults:       len(result.Flights),
			ProvidersQueried:   result.ProvidersQueried,
			ProvidersSucceeded: result.ProvidersSucceeded,
			ProvidersFailed:    result.ProvidersFailed,
			SearchTimeMS:       int(time.Since(start).Milliseconds()),
			CacheHit:           result.CacheHit,
		},
		Flights: result.Flights,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(resp)
}
