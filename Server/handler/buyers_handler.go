package handler

import (
	"net/http"

	"../model"
	"../repository"
	"github.com/go-chi/chi"
)

type BuyersHandler struct {
	buyerRepository repository.BuyerRepositoryDGraph
}

type buyerDetailedInformation struct {
	Buyer               model.Buyer     `endpoint:"buyer,omitempty"`
	BuyerWithSameIP     []model.Buyer   `endpoint:"buyers_with_same_ip,omitempty"`
	RecommendedProducts []model.Product `endpoint:"recommended_products,omitempty"`
}

func (b *BuyersHandler) GetBuyersBasicInformation(w http.ResponseWriter, r *http.Request) {
	buyers, err := b.buyerRepository.FetchBasicInformation()
	handleInternalServerError(err, w)
	respondwithJSON(w, http.StatusOK, buyers)
}

func (b *BuyersHandler) GetBuyerDetailedInformation(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	buyer, err := b.buyerRepository.FindByIdWithTransactions(id)
	handleInternalServerError(err, w)
	if buyer != nil {
		buyersWithSameIP, err := b.buyerRepository.FindBuyersWithSameIP(id)
		handleInternalServerError(err, w)
		recommendedProducts, err := b.buyerRepository.FindRecommendedProducts(id)
		handleInternalServerError(err, w)
		buyerDetailed := buyerDetailedInformation{
			Buyer:               *buyer,
			BuyerWithSameIP:     buyersWithSameIP,
			RecommendedProducts: recommendedProducts,
		}
		respondwithJSON(w, http.StatusOK, buyerDetailed)
	} else {
		respondwithJSON(w, http.StatusNotFound, errorMessage{
			Code:    buyerNotFoundCode,
			Details: errorsDetails[buyerNotFoundCode],
		})
	}
}
