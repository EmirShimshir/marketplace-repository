package entity

import (
	"github.com/EmirShimshir/marketplace-core/domain"
)

const (
	MgProductElectronic = "Electronic"
	MgProductFashion    = "Fashion"
	MgProductHome       = "Home"
	MgProductHealth     = "Health"
	MgProductSport      = "Sport"
	MgProductBooks      = "Books"
)

type MgProduct struct {
	ID          string `bson:"_id"`
	Name        string    `bson:"name"`
	Description string    `bson:"description"`
	Price       int64     `bson:"price"`
	Category    string    `bson:"category"`
	PhotoUrl    string    `bson:"photo_url"`
}

func (u *MgProduct) ToDomain() domain.Product {
	var productCategory domain.ProductCategory
	switch u.Category {
	case MgProductElectronic:
		productCategory = domain.ElectronicCategory
	case MgProductFashion:
		productCategory = domain.FashionCategory
	case MgProductHome:
		productCategory = domain.HomeCategory
	case MgProductHealth:
		productCategory = domain.HealthCategory
	case MgProductSport:
		productCategory = domain.SportCategory
	case MgProductBooks:
		productCategory = domain.BooksCategory
	}
	return domain.Product{
		ID:          domain.ID(u.ID),
		Name:        u.Name,
		Description: u.Description,
		Price:       u.Price,
		Category:    productCategory,
		PhotoUrl:    u.PhotoUrl,
	}
}

func NewMgProduct(product domain.Product) MgProduct {
	var productCategory string
	switch product.Category {
	case domain.ElectronicCategory:
		productCategory = MgProductElectronic
	case domain.FashionCategory:
		productCategory = MgProductFashion
	case domain.HomeCategory:
		productCategory = MgProductHome
	case domain.HealthCategory:
		productCategory = MgProductHealth
	case domain.SportCategory:
		productCategory = MgProductSport
	case domain.BooksCategory:
		productCategory = MgProductBooks
	}
	return MgProduct{
		ID:          product.ID.String(),
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
		Category:    productCategory,
		PhotoUrl:    product.PhotoUrl,
	}
}
