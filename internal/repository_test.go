package internal

import (
	"context"
	"github.com/oklog/ulid/v2"
	"github.com/stretchr/testify/assert"
	"playground/internal/database"
	"playground/internal/testdata"
	"testing"
)

func TestIntegrationDatabase(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode")
	}

	db, err := database.ConnectDB()

	if err != nil {
		t.Error(err)
		return
	}

	defer func() {
		if err := db.Close(); err != nil {
			t.Fatalf("could not close the database: %v", err)
		}
	}()

	t.Run("can ping the database", func(t *testing.T) {
		if err := db.Ping(); err != nil {
			t.Error(err)
			return
		}

		f := testdata.NewTestFixture(t, db)

		repository := NewProductRepository(db)

		f.Run("get product from database", func(t *testing.T) {
			productID := ProductID("123")
			ctx := context.TODO()

			product, err := repository.GetByID(ctx, productID)

			assert.NoError(t, err)
			assert.Equal(t, &Product{
				ID:    productID,
				Name:  "Noise cancelling headphones",
				Price: 295.00,
			}, product)
		})

		f.Run("save product to database", func(t *testing.T) {
			productID := ProductID(ulid.Make().String())
			ctx := context.TODO()

			assert.NoError(t, repository.Save(ctx, &Product{
				ID:    productID,
				Name:  "random product",
				Price: 5.45,
			}))

			product, err := repository.GetByID(ctx, productID)

			assert.NoError(t, err)
			assert.Equal(t, &Product{
				ID:    productID,
				Name:  "random product",
				Price: 5.45,
			}, product)
		})

		f.Run("list all products from database", func(t *testing.T) {
			ctx := context.TODO()

			expectedProductList := []*Product{
				{
					ID:    ProductID("123"),
					Name:  "Noise cancelling headphones",
					Price: 295.00,
				},
				{
					ID:    ProductID("456"),
					Name:  "Mechanical keyboard",
					Price: 150.00,
				},
				{
					ID:    ProductID("789"),
					Name:  "Wireless mouse",
					Price: 79.95,
				},
			}

			products, err := repository.ListProducts(ctx)

			assert.NoError(t, err)
			assert.Equal(t, expectedProductList, products)
		})
	})

}
