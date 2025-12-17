package service

import "bookcabin/internal/domain"

type FlightProvider interface {
	Search(req domain.SearchRequest) ([]domain.Flight, error)
	Name() string
}
