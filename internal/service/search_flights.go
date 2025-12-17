package service

import (
	"context"
	"log"
	"strconv"
	"strings"
	"sync"
	"time"

	"bookcabin/internal/common"
	"bookcabin/internal/domain"
	"bookcabin/internal/infra"
)

type SearchFlightsUseCase struct {
	Providers []FlightProvider
	Cache     *infra.Cache
}

func (uc *SearchFlightsUseCase) Execute(
	ctx context.Context,
	req domain.SearchRequest,
) ([]domain.Flight, error) {

	cacheKey := searchCacheKey(req)

	if v, ok := uc.Cache.Get(cacheKey); ok {
		if cached, ok := v.([]domain.Flight); ok {
			return uc.filterAndSort(cached, req)
		}
	}

	var wg sync.WaitGroup
	result := make(chan []domain.Flight, len(uc.Providers))

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	for _, p := range uc.Providers {
		wg.Add(1)
		go func(p FlightProvider) {
			defer wg.Done()

			flights, err := p.Search(req)
			if err != nil {
				log.Printf("[WARN] provider %s failed: %v", p.Name(), err)
				return
			}
			result <- flights
		}(p)
	}

	go func() {
		wg.Wait()
		close(result)
	}()

	var allFlights []domain.Flight
	for flights := range result {
		for _, f := range flights {
			if !isValidFlight(f) {
				continue
			}
			allFlights = append(allFlights, f)
		}
	}

	uc.Cache.Set(cacheKey, allFlights, 3*time.Minute)

	return uc.filterAndSort(allFlights, req)
}

func (uc *SearchFlightsUseCase) filterAndSort(
	flights []domain.Flight,
	req domain.SearchRequest,
) ([]domain.Flight, error) {

	filter, err := buildFlightFilterFromRequest(req)
	if err != nil {
		return nil, err
	}

	filtered := filterFlights(flights, req, filter)

	if req.SortBy != "" {
		common.SortFlights(filtered, req.SortBy)
	}

	return filtered, nil
}

func buildFlightFilterFromRequest(req domain.SearchRequest) (domain.FlightFilter, error) {
	filter := domain.FlightFilter{
		MinPrice:    req.MinPrice,
		MaxPrice:    req.MaxPrice,
		MaxStops:    req.MaxStops,
		Airlines:    common.NormalizeAirlines(req.Airlines),
		MaxDuration: req.MaxDuration,
	}

	// parse base date
	baseDate, err := time.Parse("2006-01-02", req.DepartureDate)
	if err != nil {
		return filter, err
	}

	// departure window
	if req.EarliestDep != "" {
		t, err := time.Parse("15:04", req.EarliestDep)
		if err != nil {
			return filter, err
		}
		filter.EarliestDeparture = time.Date(
			baseDate.Year(), baseDate.Month(), baseDate.Day(),
			t.Hour(), t.Minute(), 0, 0, time.Local,
		)
	}

	if req.LatestDep != "" {
		t, err := time.Parse("15:04", req.LatestDep)
		if err != nil {
			return filter, err
		}
		filter.LatestDeparture = time.Date(
			baseDate.Year(), baseDate.Month(), baseDate.Day(),
			t.Hour(), t.Minute(), 59, 0, time.Local,
		)
	}

	// arrival window (optional)
	if req.EarliestArr != "" {
		t, err := time.Parse("15:04", req.EarliestArr)
		if err != nil {
			return filter, err
		}
		filter.EarliestArrival = time.Date(
			baseDate.Year(), baseDate.Month(), baseDate.Day(),
			t.Hour(), t.Minute(), 0, 0, time.Local,
		)
	}

	if req.LatestArr != "" {
		t, err := time.Parse("15:04", req.LatestArr)
		if err != nil {
			return filter, err
		}
		filter.LatestArrival = time.Date(
			baseDate.Year(), baseDate.Month(), baseDate.Day(),
			t.Hour(), t.Minute(), 59, 0, time.Local,
		)
	}

	return filter, nil
}

func filterFlights(
	flights []domain.Flight,
	req domain.SearchRequest,
	filter domain.FlightFilter,
) []domain.Flight {

	var res []domain.Flight

	reqDate, err := time.Parse("2006-01-02", req.DepartureDate)
	if err != nil {
		return res
	}

	for _, f := range flights {
		// origin / destination
		if !strings.EqualFold(f.Origin, req.Origin) ||
			!strings.EqualFold(f.Destination, req.Destination) {
			continue
		}

		// same calendar date (use flight timezone)
		fy, fm, fd := f.DepartureTime.Date()
		ry, rm, rd := reqDate.Date()
		if fy != ry || fm != rm || fd != rd {
			continue
		}

		// price
		if filter.MinPrice > 0 && f.PriceIDR < filter.MinPrice {
			continue
		}
		if filter.MaxPrice > 0 && f.PriceIDR > filter.MaxPrice {
			continue
		}

		// stops
		if filter.MaxStops >= 0 && f.Stops > filter.MaxStops {
			continue
		}

		// duration
		if filter.MaxDuration > 0 && f.DurationMin > filter.MaxDuration {
			continue
		}

		// airline (use IATA code)
		if len(filter.Airlines) > 0 && !contains(filter.Airlines, f.AirlineCode) {
			continue
		}

		// departure window
		if !filter.EarliestDeparture.IsZero() &&
			f.DepartureTime.Before(filter.EarliestDeparture) {
			continue
		}

		if !filter.LatestDeparture.IsZero() &&
			f.DepartureTime.After(filter.LatestDeparture) {
			continue
		}

		// arrival window
		if !filter.EarliestArrival.IsZero() &&
			f.ArrivalTime.Before(filter.EarliestArrival) {
			continue
		}

		if !filter.LatestArrival.IsZero() &&
			f.ArrivalTime.After(filter.LatestArrival) {
			continue
		}

		res = append(res, f)
	}

	return res
}

func contains(arr []string, s string) bool {
	for _, a := range arr {
		if a == s {
			return true
		}
	}
	return false
}

func isValidFlight(f domain.Flight) bool {
	return !f.DepartureTime.IsZero() && !f.ArrivalTime.IsZero() && f.ArrivalTime.After(f.DepartureTime)
}

func searchCacheKey(req domain.SearchRequest) string {
	parts := []string{
		req.Origin,
		req.Destination,
		req.DepartureDate,
		req.CabinClass,
		strconv.Itoa(req.Passengers),
	}

	return strings.Join(parts, "|")
}
