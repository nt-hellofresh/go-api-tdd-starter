package internal

type ProductID string

type Product struct {
	ID    ProductID
	Name  string  `json:"name"`
	Price float32 `json:"price"`
}
