package domain

type Product struct {
	ID    int
	Name  string
	Price int
}

type ProductRepository interface {
	AddProduct(product Product) (int, error)
	GetProduct(id int) (Product, error)
}
