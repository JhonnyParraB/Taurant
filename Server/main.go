package main

import (
	"log"
	"net/http"
)

func main() {
	router := getRouter()
	log.Println("Server listen at : 4200")
	err := http.ListenAndServe(":4200", router)
	if err != nil {
		log.Fatal(err)
	}
}
