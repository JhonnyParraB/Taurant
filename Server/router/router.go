package router

import (
	"../handler"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func GetRouter() *chi.Mux {
	router := chi.NewRouter()
	router.Use(middleware.Recoverer)
	router.Use(middleware.Logger)

	var dayDataLoaderHandler handler.LoadDayDataHandler
	var buyersHandler handler.BuyersHandler
	dayDataLoaderHandler.Init()

	router.Post("/load-day-data/{date}", dayDataLoaderHandler.LoadDayData)

	router.Get("/buyers", buyersHandler.GetBuyersBasicInformation)
	router.Get("/buyers/{id}", buyersHandler.GetBuyerDetailedInformation)

	return router
}
