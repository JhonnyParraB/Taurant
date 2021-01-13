package main

import (
	"log"
	"net/http"

	Router "./router"
)

func main() {
	router := Router.GetRouter()
	log.Println("Server listen at : 4200")
	err := http.ListenAndServe(":4200", router)
	if err != nil {
		log.Fatal(err)
	}
}
