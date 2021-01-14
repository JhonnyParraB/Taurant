package router

import (
	"Taurant/handler"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func GetRouter() *chi.Mux {
	router := chi.NewRouter()
	router.Use(middleware.Recoverer)
	router.Use(middleware.Logger)

	var dayDataLoaderHandler = handler.NewLoadDayDataHandler()
	var buyersHandler handler.BuyersHandler

	router.Post("/load-day-data/{date}", dayDataLoaderHandler.LoadDayData)

	router.Get("/buyers", buyersHandler.GetBuyersBasicInformation)
	router.Get("/buyers/{id}", buyersHandler.GetBuyerDetailedInformation)

	return router
}
