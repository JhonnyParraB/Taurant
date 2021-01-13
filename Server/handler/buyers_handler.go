package handler

import (
	"net/http"

	"../model"
	"../repository"
	"github.com/go-chi/chi"
)

//BuyersHandler is used at router
type BuyersHandler struct {
	buyerRepository repository.BuyerRepositoryDGraph
}

type buyerDetailedInformation struct {
	Buyer               model.Buyer     `endpoint:"buyer,omitempty"`
	BuyerWithSameIP     []model.Buyer   `endpoint:"buyers_with_same_ip,omitempty"`
	RecommendedProducts []model.Product `endpoint:"recommended_products,omitempty"`
}

//GetBuyersBasicInformation is used at router
func (b *BuyersHandler) GetBuyersBasicInformation(w http.ResponseWriter, r *http.Request) {
	buyers, err := b.buyerRepository.FetchBasicInformation()
	if err != nil {
		handleInternalServerError(err, w)
		return
	}
	respondwithJSON(w, http.StatusOK, buyers)
}

//GetBuyerDetailedInformation is used at router
func (b *BuyersHandler) GetBuyerDetailedInformation(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	buyer, err := b.buyerRepository.FindByIdWithTransactions(id)
	if err != nil {
		handleInternalServerError(err, w)
		return
	}
	if buyer != nil {
		buyersWithSameIP, err := b.buyerRepository.FindBuyersWithSameIP(id)
		if err != nil {
			handleInternalServerError(err, w)
			return
		}
		recommendedProducts, err := b.buyerRepository.FindRecommendedProducts(id)
		if err != nil {
			handleInternalServerError(err, w)
			return
		}
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
