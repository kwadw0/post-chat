package products

import (
	"kwadw0/gocommerce/internal/json"
	"log"
	"log/slog"
	"net/http"
)

type handler struct {
	service Service
}

func NewHandler(s Service) *handler {
	return &handler{service: s}
}

func (h *handler) GetAllProducts(w http.ResponseWriter, r *http.Request) {
	products, err := h.service.GetAllProducts(r.Context())
	if err != nil {
		log.Println("Error getting products", err)
		json.WriteJson(w, http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
		return
	}

	json.WriteJson(w, http.StatusOK, products)
}


// AddProduct godoc
// @Summary      Create a new product
// @Description  Takes a JSON payload and stores a new toy in our database
// @Tags         products
// @Accept       json
// @Produce      json
// @Param        product  body      CreateProduct  true  "Product to create"
// @Success      201      {object}  ProductResponse
// @Failure      400      {object}  map[string]string
// @Router       /products [post]
func (h *handler) AddProduct(w http.ResponseWriter, r *http.Request) {
	var payload CreateProduct
	if err := json.ReadJson(w, r, &payload); err != nil {
		slog.Error("Error reading json", "error", err)
		json.WriteJson(w, http.StatusBadRequest, map[string]string{"error": "Invalid json"})
		return
	}

	product, err := h.service.AddProduct(r.Context(), payload)
	if err != nil {
		log.Println("Error adding product", err)
		json.WriteJson(w, http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
		return
	}

	json.WriteJson(w, http.StatusOK, product)
}