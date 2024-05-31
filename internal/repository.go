package internal

import (
	"context"
	"database/sql"
	"github.com/pkg/errors"
)

type AbstractProductRepository interface {
	Save(ctx context.Context, product *Product) error
	GetByID(ctx context.Context, ID ProductID) (*Product, error)
	ListProducts(ctx context.Context) ([]*Product, error)
}

type ProductRepository struct {
	db *sql.DB
}

func NewProductRepository(db *sql.DB) *ProductRepository {
	return &ProductRepository{
		db: db,
	}
}

func (r *ProductRepository) Save(ctx context.Context, product *Product) error {
	_, err := r.db.ExecContext(
		ctx,
		"INSERT INTO product (id, name, price) VALUES ($1, $2, $3)",
		product.ID,
		product.Name,
		product.Price,
	)
	return err
}

func (r *ProductRepository) GetByID(ctx context.Context, ID ProductID) (*Product, error) {
	var product Product
	err := r.db.QueryRowContext(
		ctx,
		"SELECT id, name, price FROM product WHERE id = $1",
		ID,
	).Scan(&product.ID, &product.Name, &product.Price)
	return &product, err
}

func (r *ProductRepository) ListProducts(ctx context.Context) ([]*Product, error) {
	products := make([]*Product, 0)

	rows, err := r.db.QueryContext(ctx, "SELECT id, name, price FROM product")

	if err != nil {
		return nil, errors.Wrap(err, "could not query products")
	}

	for rows.Next() {
		var product Product
		err := rows.Scan(&product.ID, &product.Name, &product.Price)
		if err != nil {
			return nil, errors.Wrap(err, "could not scan product")
		}
		products = append(products, &product)
	}

	return products, nil
}
