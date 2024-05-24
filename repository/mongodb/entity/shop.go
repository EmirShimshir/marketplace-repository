package entity

import (
	"github.com/EmirShimshir/marketplace-core/domain"
)

type MgShop struct {
	ID          string `bson:"_id"`
	SellerID    string `bson:"seller_id"`
	Name        string    `bson:"name"`
	Description string    `bson:"description"`
	Requisites  string    `bson:"requisites"`
	Email       string    `bson:"email"`
}

func (s *MgShop) ToDomain() domain.Shop {
	return domain.Shop{
		ID:          domain.ID(s.ID),
		SellerID:    domain.ID(s.SellerID),
		Name:        s.Name,
		Description: s.Description,
		Requisites:  s.Requisites,
		Email:       s.Email,
	}
}

func NewMgShop(shop domain.Shop) MgShop {
	return MgShop{
		ID:          shop.ID.String(),
		SellerID:    shop.SellerID.String(),
		Name:        shop.Name,
		Description: shop.Description,
		Requisites:  shop.Requisites,
		Email:       shop.Email,
	}
}

type MgShopItem struct {
	ID        string `bson:"_id"`
	ShopID    string `bson:"shop_id"`
	ProductID string `bson:"product_id"`
	Quantity  int64     `bson:"quantity"`
}

func (si *MgShopItem) ToDomain() domain.ShopItem {
	return domain.ShopItem{
		ID:        domain.ID(si.ID),
		ShopID:    domain.ID(si.ShopID),
		ProductID: domain.ID(si.ProductID),
		Quantity:  si.Quantity,
	}
}

func NewMgShopItem(shopItem domain.ShopItem) MgShopItem {
	return MgShopItem{
		ID:        shopItem.ID.String(),
		ShopID:    shopItem.ShopID.String(),
		ProductID: shopItem.ProductID.String(),
		Quantity:  shopItem.Quantity,
	}
}
