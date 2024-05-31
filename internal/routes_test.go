package internal

import (
	"context"
	"encoding/json"
	"github.com/oklog/ulid/v2"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type mock struct {
	generatedIDs  []ProductID
	simulateError bool
}

func (m *mock) SimulateError() {
	m.simulateError = true
}

func (m *mock) Save(_ context.Context, _ *Product) error {
	if m.simulateError {
		m.simulateError = false
		return errors.New("could not save product")
	}
	return nil
}

func (m *mock) GetByID(_ context.Context, ID ProductID) (*Product, error) {
	return mockProduct(ID), nil
}

func (m *mock) ListProducts(_ context.Context) ([]*Product, error) {
	products := make([]*Product, 0)

	for _, id := range m.generatedIDs {
		products = append(products, mockProduct(id))
	}

	return products, nil
}

func mockRepo() *mock {
	return &mock{
		generatedIDs: []ProductID{
			ProductID(ulid.Make().String()),
			ProductID(ulid.Make().String()),
			ProductID(ulid.Make().String()),
		},
	}
}

func mockProduct(ID ProductID) *Product {
	return &Product{
		ID:    ID,
		Name:  "mock product",
		Price: 14.95,
	}
}

func TestProductHandler(t *testing.T) {
	handler := http.NewServeMux()
	repository := mockRepo()
	registerRoutes(handler, repository)

	t.Run("ListProducts", func(t *testing.T) {
		request := httptest.NewRequest(http.MethodGet, "/products", nil)
		resp := httptest.NewRecorder()

		handler.ServeHTTP(resp, request)

		assert.Equal(t, http.StatusOK, resp.Code)

		expectedProducts, err := repository.ListProducts(context.TODO())
		assert.NoError(t, err)

		products := make([]*Product, 0)
		assert.NoError(t, json.NewDecoder(resp.Body).Decode(&products))
		assert.Equal(t, expectedProducts, products)
	})

	t.Run("GetProduct", func(t *testing.T) {
		productID := ulid.Make().String()

		request := httptest.NewRequest(http.MethodGet, "/products/"+productID, nil)
		resp := httptest.NewRecorder()

		handler.ServeHTTP(resp, request)

		assert.Equal(t, http.StatusOK, resp.Code)

		var product Product
		assert.NoError(t, json.NewDecoder(resp.Body).Decode(&product))
		assert.Equal(t, mockProduct(ProductID(productID)), &product)
	})

	t.Run("SaveNewProduct", func(t *testing.T) {
		t.Run("happy path", func(t *testing.T) {
			productID := ulid.Make().String()

			body := strings.NewReader(`{"name": "test", "price": 14.95}`)
			request := httptest.NewRequest(http.MethodPut, "/products/"+productID, body)
			resp := httptest.NewRecorder()

			handler.ServeHTTP(resp, request)

			assert.Equal(t, http.StatusOK, resp.Code)
			assert.Equal(t, "product saved successfully", resp.Body.String())
		})

		t.Run("invalid payload", func(t *testing.T) {
			productID := ulid.Make().String()

			body := strings.NewReader("")
			request := httptest.NewRequest(http.MethodPut, "/products/"+productID, body)
			resp := httptest.NewRecorder()

			handler.ServeHTTP(resp, request)

			assert.Equal(t, http.StatusBadRequest, resp.Code)
			assert.Equal(t, "could not decode product from request body: EOF", resp.Body.String())
		})

		t.Run("repository error", func(t *testing.T) {
			productID := ulid.Make().String()
			repository.SimulateError()

			body := strings.NewReader(`{"name": "test", "price": 14.95}`)
			request := httptest.NewRequest(http.MethodPut, "/products/"+productID, body)
			resp := httptest.NewRecorder()

			handler.ServeHTTP(resp, request)

			assert.Equal(t, http.StatusInternalServerError, resp.Code)
			assert.Equal(t, "there was an error handling the request: could not save product", resp.Body.String())
		})
	})

}
