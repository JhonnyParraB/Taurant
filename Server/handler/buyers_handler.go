package handler

import (
	"net/http"

	"../model"
	"../repository"
	"github.com/go-chi/chi"
	jsoniter "github.com/json-iterator/go"
)

type BuyersHandler struct {
	buyerRepository repository.BuyerRepositoryDGraph
}

type BuyerDetailedInformation struct {
	Buyer               model.Buyer     `endpoint:"buyer,omitempty"`
	BuyerWithSameIP     []model.Buyer   `endpoint:"buyers_with_same_ip,omitempty"`
	RecommendedProducts []model.Product `endpoint:"recommended_products,omitempty"`
}

func (b *BuyersHandler) GetBuyersBasicInformation(w http.ResponseWriter, r *http.Request) {
	respondwithJSON(w, http.StatusOK, b.buyerRepository.FetchBasicInformation())
}

func (b *BuyersHandler) GetBuyerDetailedInformation(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	buyer := b.buyerRepository.FindByIdWithTransactions(id)
	if buyer != nil {
		buyersWithSameIP := b.buyerRepository.FindBuyersWithSameIP(id)
		recommendedProducts := b.buyerRepository.FindRecommendedProducts(id)
		buyerDetailed := BuyerDetailedInformation{
			Buyer:               *buyer,
			BuyerWithSameIP:     buyersWithSameIP,
			RecommendedProducts: recommendedProducts,
		}
		respondwithJSON(w, http.StatusOK, buyerDetailed)
	}
}

func respondwithJSON(w http.ResponseWriter, code int, payload interface{}) {
	var endpointCaseJSON = jsoniter.Config{TagKey: "endpoint"}.Froze()
	response, _ := endpointCaseJSON.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
