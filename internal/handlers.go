package internal

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"net/http"
)

type ProductHandler struct {
	repo AbstractProductRepository
}

func NewProductHandler(repo AbstractProductRepository) *ProductHandler {
	return &ProductHandler{
		repo: repo,
	}
}

func (h *ProductHandler) ListProducts(w http.ResponseWriter, r *http.Request) error {
	products, err := h.repo.ListProducts(r.Context())

	if err != nil {
		return err
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(products); err != nil {
		return err
	}

	return nil
}

func (h *ProductHandler) GetProduct(w http.ResponseWriter, r *http.Request) error {
	productID := r.PathValue("product_id")

	product, err := h.repo.GetByID(r.Context(), ProductID(productID))

	if err != nil {
		return err
	}

	if err := json.NewEncoder(w).Encode(product); err != nil {
		return err
	}

	w.WriteHeader(http.StatusOK)
	return nil
}

func (h *ProductHandler) SaveProduct(w http.ResponseWriter, r *http.Request) error {
	var product Product

	if err := unmarshal(r, &product); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, err = fmt.Fprint(w, err.Error())
		return err
	}

	if err := h.repo.Save(r.Context(), &product); err != nil {
		return err
	}

	w.WriteHeader(http.StatusOK)
	_, err := fmt.Fprintf(w, "product saved successfully")
	return err
}

func unmarshal(r *http.Request, product *Product) error {
	productID := r.PathValue("product_id")

	if err := json.NewDecoder(r.Body).Decode(product); err != nil {
		return errors.Wrap(err, "could not decode product from request body")
	}

	product.ID = ProductID(productID)

	return nil
}
