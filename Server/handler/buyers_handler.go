package handler

import (
	"net/http"

	"../model"
	"../repository"
	jsoniter "github.com/json-iterator/go"
)

type BuyersHandler struct {
	buyerRepository repository.BuyerRepositoryDGraph
}

type BuyerDetailedInformation struct {
	Buyer               model.Buyer     `endpoint:"buyer,omitempty"`
	BuyerWithSameIP     []model.Buyer   `endpoint:"buyersWithSameIp,omitempty"`
	RecommendedProducts []model.Product `endpoint:"productsRecommended,omitempty"`
}

func (b *BuyersHandler) GetBuyersBasicInformation(w http.ResponseWriter, r *http.Request) {
	respondwithJSON(w, http.StatusOK, b.buyerRepository.FetchBasicInformation())
}

func (b *BuyersHandler) GetBuyerDetailedInformation(w http.ResponseWriter, r *http.Request) {
	buyer := b.buyerRepository.FindByIdWithTransactions("49688cb8")
	buyersWithSameIP := b.buyerRepository.FindBuyersWithSameIP("49688cb8")
	recommendedProducts := b.buyerRepository.FindRecommendedProducts("49688cb8")
	buyerDetailed := BuyerDetailedInformation{
		Buyer:               *buyer,
		BuyerWithSameIP:     buyersWithSameIP,
		RecommendedProducts: recommendedProducts,
	}
	respondwithJSON(w, http.StatusOK, buyerDetailed)
}

func respondwithJSON(w http.ResponseWriter, code int, payload interface{}) {
	var endpointCaseJSON = jsoniter.Config{TagKey: "predicate"}.Froze()
	response, _ := endpointCaseJSON.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
