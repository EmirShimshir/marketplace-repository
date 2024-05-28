package entity

import (
	"github.com/EmirShimshir/marketplace-core/domain"
)

type MgCart struct {
	ID    string `bson:"_id"`
	Price int64     `bson:"price"`
}

func (c *MgCart) ToDomain() domain.Cart {
	return domain.Cart{
		ID:    domain.ID(c.ID),
		Price: c.Price,
	}
}

func NewMgCart(cart domain.Cart) MgCart {
	return MgCart{
		ID: string(cart.ID),
		Price: cart.Price,
	}
}

type MgCartItem struct {
	ID        string `bson:"_id"`
	CartID    string `bson:"cart_id"`
	ProductID string `bson:"product_id"`
	Quantity  int64     `bson:"quantity"`
}

func (ci *MgCartItem) ToDomain() domain.CartItem {
	return domain.CartItem{
		ID: domain.ID(ci.ID),
		CartID:    domain.ID(ci.CartID),
		ProductID: domain.ID(ci.ProductID),
		Quantity:  ci.Quantity,
	}
}

func NewMgCartItem(cartItem domain.CartItem) MgCartItem {
	return MgCartItem{
		ID:        string(cartItem.ID),
		CartID:    string(cartItem.CartID),
		ProductID: string(cartItem.ProductID),
		Quantity:  cartItem.Quantity,
	}
}
