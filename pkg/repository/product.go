package repository

import (
	"errors"
	"oa-bitgin/pkg/domain"
	"sync/atomic"
)

type productRepository struct {
	IDCounter atomic.Value
	Product   map[int]domain.Product
}

func (p *productRepository) init() {
	p.IDCounter.Store(0)
	p.Product = make(map[int]domain.Product)
}

func NewProductRepository() domain.ProductRepository {
	store := &productRepository{}
	store.init()
	return store
}

func (p *productRepository) AddProduct(product domain.Product) (int, error) {
	id := p.IDCounter.Load().(int)
	id++
	p.IDCounter.Store(id)
	product.ID = id
	p.Product[id] = product
	return id, nil
}

func (p *productRepository) GetProduct(id int) (domain.Product, error) {
	if product, ok := p.Product[id]; ok {
		return product, nil
	} else {
		return domain.Product{}, errors.New("product not found")
	}
}
