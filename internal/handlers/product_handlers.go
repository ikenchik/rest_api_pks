package handlers

import (
	"encoding/json"
	"net/http"
	"rest_api_pks/internal/models"
	"rest_api_pks/internal/repository"
	"strconv"
	"strings"
)

type ProductHandler struct {
	repo *repository.ProductRepository
}

func NewProductHandler(repo *repository.ProductRepository) *ProductHandler {
	return &ProductHandler{repo: repo}
}

func (h *ProductHandler) GetProductsHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	// Получение параметров фильтрации
	searchQuery := query.Get("search")
	minPrice := query.Get("min_price")
	maxPrice := query.Get("max_price")
	sortBy := query.Get("sort_by")
	sortOrder := query.Get("sort_order")

	// Преобразование параметров в нужные типы
	var minPriceFloat, maxPriceFloat float64
	var err error
	if minPrice != "" {
		minPriceFloat, err = strconv.ParseFloat(minPrice, 64)
		if err != nil {
			http.Error(w, "Invalid min_price", http.StatusBadRequest)
			return
		}
	}
	if maxPrice != "" {
		maxPriceFloat, err = strconv.ParseFloat(maxPrice, 64)
		if err != nil {
			http.Error(w, "Invalid max_price", http.StatusBadRequest)
			return
		}
	}

	// Получение продуктов с учетом фильтрации и сортировки
	products, err := h.repo.GetFilteredProducts(r.Context(), searchQuery, minPriceFloat, maxPriceFloat, sortBy, sortOrder)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(products)
}

func (h *ProductHandler) GetProductByIDHandler(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/products/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid Product ID", http.StatusBadRequest)
		return
	}

	product, err := h.repo.GetProductByID(r.Context(), id)
	if err != nil {
		http.Error(w, "Product not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(product)
}

func (h *ProductHandler) CreateProductHandler(w http.ResponseWriter, r *http.Request) {
	var newProduct models.Product
	err := json.NewDecoder(r.Body).Decode(&newProduct)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.repo.CreateProduct(r.Context(), &newProduct)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(newProduct)
}

func (h *ProductHandler) UpdateProductHandler(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/products/update/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid Product ID", http.StatusBadRequest)
		return
	}

	var updatedProduct models.Product
	err = json.NewDecoder(r.Body).Decode(&updatedProduct)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	updatedProduct.ID = id
	err = h.repo.UpdateProduct(r.Context(), &updatedProduct)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedProduct)
}

func (h *ProductHandler) DeleteProductHandler(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/products/delete/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid Product ID", http.StatusBadRequest)
		return
	}

	err = h.repo.DeleteProduct(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *ProductHandler) UpdateProductQuantityHandler(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/products/quantity/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid Product ID", http.StatusBadRequest)
		return
	}

	var updatedQuantity struct {
		Quantity int `json:"quantity"`
	}
	err = json.NewDecoder(r.Body).Decode(&updatedQuantity)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.repo.UpdateProductQuantity(r.Context(), id, updatedQuantity.Quantity)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedQuantity)
}

func (h *ProductHandler) ToggleFavoriteHandler(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/products/favorite/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid Product ID", http.StatusBadRequest)
		return
	}

	err = h.repo.ToggleFavorite(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *ProductHandler) ToggleCartHandler(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/products/cart/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid Product ID", http.StatusBadRequest)
		return
	}

	err = h.repo.ToggleCart(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
