package products

import (
	"kwadw0/gocommerce/internal/json"
	"log"
	"net/http"
)

type handler struct {
	service Service
}

func NewHandler(s Service) *handler {
	return &handler{service: s}
}

func (h *handler) GetAllProducts(w http.ResponseWriter, r *http.Request) {
	err := h.service.GetAllProducts(r.Context())
	if err != nil {
		log.Println("Error getting products", err)
		json.WriteJson(w, http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
		return
	}
	products := []string{"Product 1", "Product 2", "Product 3"}

	json.WriteJson(w, http.StatusOK, products)
}