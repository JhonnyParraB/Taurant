package main

import (
	"log"
	"net/http"

	"./handler"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func main() {
	router := getRouter()
	log.Println("Server listen at : 4200")
	err := http.ListenAndServe(":4200", router)
	if err != nil {
		log.Fatal(err)
	}
}

func getRouter() *chi.Mux {
	router := chi.NewRouter()
	router.Use(middleware.Recoverer)
	router.Use(middleware.Logger)

	var dayDataLoaderHandler handler.LoadDayDataHandler
	var buyersHandler handler.BuyersHandler
	dayDataLoaderHandler.Init()
	router.Post("/load-day-data/{date}", dayDataLoaderHandler.LoadDayData)

	//buyers
	router.Get("/buyers", buyersHandler.GetBuyersBasicInformation)
	router.Get("/buyers/{id}", buyersHandler.GetBuyerDetailedInformation)
	return router
}
