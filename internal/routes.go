package internal

import (
	"database/sql"
	"fmt"
	"net/http"
)

type HandlerFunc func(w http.ResponseWriter, r *http.Request) error

func ToHandler(fn HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := fn(w, r); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "there was an error handling the request: %s", err)
		}
	}
}

func ConfigureRoutes(router *http.ServeMux, db *sql.DB) {
	repo := NewProductRepository(db)
	registerRoutes(router, repo)
}

func registerRoutes(router *http.ServeMux, repo AbstractProductRepository) {
	handler := NewProductHandler(repo)

	router.HandleFunc("GET /products", ToHandler(handler.ListProducts))
	router.HandleFunc("GET /products/{product_id}", ToHandler(handler.GetProduct))
	router.HandleFunc("PUT /products/{product_id}", ToHandler(handler.SaveProduct))
}
