package common

import (
	"bookcabin/internal/domain"
	"errors"
	"sort"
	"strconv"
	"strings"
	"time"
)

// Parse time with explicit timezone (e.g. Asia/Jakarta, Asia/Makassar)
func ParseTimeWithTZ(datetime, tz string) (time.Time, error) {
	loc, err := time.LoadLocation(tz)
	if err != nil {
		return time.Time{}, err
	}
	return time.ParseInLocation("2006-01-02T15:04:05", datetime, loc)
}

func ParseFlexibleTime(value string) (time.Time, error) {
	layouts := []string{
		time.RFC3339,                // 2006-01-02T15:04:05Z07:00
		"2006-01-02T15:04:05-07:00", // fallback with colon
		"2006-01-02T15:04:05-0700",  // fallback without colon
		"2006-01-02T15:04:05",       // naive datetime
	}

	for _, l := range layouts {
		if t, err := time.Parse(l, value); err == nil {
			return t, nil
		}
	}

	return time.Time{}, errors.New("unsupported time format: " + value)
}

// Normalize price to IDR (supports string or numeric)
func ParsePriceToIDR(value interface{}, currency string) (int64, error) {
	switch v := value.(type) {
	case int:
		return int64(v), nil
	case float64:
		return int64(v), nil
	case string:
		clean := strings.ReplaceAll(v, ",", "")
		clean = strings.TrimSpace(clean)
		value, err := strconv.Atoi(clean)
		return int64(value), err
	default:
		return 0, errors.New("unsupported price format")
	}
}

func SortFlights(flights []domain.Flight, sortBy string) {
	switch strings.ToLower(sortBy) {
	case "price_asc":
		sort.Slice(flights, func(i, j int) bool { return flights[i].PriceIDR < flights[j].PriceIDR })
	case "price_desc":
		sort.Slice(flights, func(i, j int) bool { return flights[i].PriceIDR > flights[j].PriceIDR })
	case "duration_asc":
		sort.Slice(flights, func(i, j int) bool { return flights[i].DurationMin < flights[j].DurationMin })
	case "duration_desc":
		sort.Slice(flights, func(i, j int) bool { return flights[i].DurationMin > flights[j].DurationMin })
	case "departure_asc":
		sort.Slice(flights, func(i, j int) bool { return flights[i].DepartureTime.Before(flights[j].DepartureTime) })
	case "arrival_asc":
		sort.Slice(flights, func(i, j int) bool { return flights[i].ArrivalTime.Before(flights[j].ArrivalTime) })
	case "best_value":
		sort.Slice(flights, func(i, j int) bool { return bestValueScore(flights[i]) < bestValueScore(flights[j]) })
	default:
		// base value default
		sort.Slice(flights, func(i, j int) bool { return bestValueScore(flights[i]) < bestValueScore(flights[j]) })
	}
}

func bestValueScore(f domain.Flight) int64 {
	score := int64(0)
	score += f.PriceIDR / 10000
	score += int64(f.DurationMin)
	score += int64(f.Stops * 100)
	return score
}

var airlineAlias = map[string]string{
	"GARUDA":           "GA",
	"GARUDA INDONESIA": "GA",
	"LION":             "JT",
	"LION AIR":         "JT",
	"AIRASIA":          "QZ",
	"QZ":               "QZ",
	"GA":               "GA",
	"JT":               "JT",
	"ID":               "ID",
	"BATIK AIR":        "ID",
}

func NormalizeAirlines(airlines []string) []string {
	set := map[string]struct{}{}

	for _, a := range airlines {
		key := strings.ToUpper(strings.TrimSpace(a))
		if code, ok := airlineAlias[key]; ok {
			set[code] = struct{}{}
		}
	}

	out := make([]string, 0, len(set))
	for k := range set {
		out = append(out, k)
	}
	return out
}
