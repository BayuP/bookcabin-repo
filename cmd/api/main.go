package main

import (
	"bookcabin/internal/handler"
	"bookcabin/internal/infra"
	"bookcabin/internal/mock"
	"bookcabin/internal/provider"
	"bookcabin/internal/service"
	"net/http"

	_ "bookcabin/docs"

	httpSwagger "github.com/swaggo/http-swagger"
)

// @title BOOKCABIN Flight Search API
// @version 1.0
// @description Flight search aggregation service
// @host localhost:8080
// @BasePath /
func main() {

	mock.MockAirAsiaServer()
	mock.MockBatikServer()
	mock.MockLionServer()
	mock.MockGarudaServer()
	cache := infra.NewCache()

	uc := &service.SearchFlightsUseCase{
		Providers: []service.FlightProvider{
			&provider.AirAsiaProvider{BaseURL: "http://127.0.0.1:8081", Client: provider.NewHTTPClient()},
			&provider.BatikProvider{BaseURL: "http://127.0.0.1:8082", Client: provider.NewHTTPClient()},
			&provider.GarudaProvider{BaseURL: "http://127.0.0.1:8083", Client: provider.NewHTTPClient()},
			&provider.LionAirProvider{BaseURL: "http://127.0.0.1:8084", Client: provider.NewHTTPClient()},
		},
		Cache: cache,
	}

	h := handler.NewFlightHandler(uc)

	http.HandleFunc("/search", h.Search)
	http.Handle("/swagger/", httpSwagger.WrapHandler)
	http.ListenAndServe(":8080", nil)
}
