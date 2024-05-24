package entity

import (
	"github.com/EmirShimshir/marketplace-core/domain"
	"time"
)

type MgOrderCustomer struct {
	ID         string    `bson:"_id"`
	CustomerID string    `bson:"customer_id"`
	Address    string    `bson:"address"`
	CreatedAt  time.Time `bson:"created_at"`
	TotalPrice int64     `bson:"total_price"`
	Payed      bool      `bson:"payed"`
}

func (oc *MgOrderCustomer) ToDomain() domain.OrderCustomer {
	return domain.OrderCustomer{
		ID:         domain.ID(oc.ID),
		CustomerID: domain.ID(oc.CustomerID),
		Address:    oc.Address,
		CreatedAt:  oc.CreatedAt,
		TotalPrice: oc.TotalPrice,
		Payed:      oc.Payed,
	}
}

func NewMgOrderCustomer(orderCustomer domain.OrderCustomer) MgOrderCustomer {
	return MgOrderCustomer{
		ID:         orderCustomer.ID.String(),
		CustomerID: orderCustomer.CustomerID.String(),
		Address:    orderCustomer.Address,
		CreatedAt:  orderCustomer.CreatedAt,
		TotalPrice: orderCustomer.TotalPrice,
		Payed:      orderCustomer.Payed,
	}
}

const (
	MgOrderShopStart = "Start"
	MgOrderShopReady = "Ready"
	MgOrderShopDone  = "Done"
)

type MgOrderShop struct {
	ID              string `bson:"_id"`
	ShopID          string `bson:"shop_id"`
	OrderCustomerID string `bson:"order_customer_id"`
	Status          string `bson:"status"`
	Notified        bool   `bson:"notified"`
}

func (os *MgOrderShop) ToDomain() domain.OrderShop {
	var orderShopStatus domain.OrderShopStatus
	switch os.Status {
	case MgOrderShopStart:
		orderShopStatus = domain.OrderShopStatusStart
	case MgOrderShopReady:
		orderShopStatus = domain.OrderShopStatusReady
	case MgOrderShopDone:
		orderShopStatus = domain.OrderShopStatusDone
	}

	return domain.OrderShop{
		ID:              domain.ID(os.ID),
		ShopID:          domain.ID(os.ShopID),
		OrderCustomerID: domain.ID(os.OrderCustomerID),
		Status:          orderShopStatus,
		Notified:        os.Notified,
	}
}

func NewMgOrderShop(orderShop domain.OrderShop) MgOrderShop {
	var orderShopStatus string
	switch orderShop.Status {
	case domain.OrderShopStatusStart:
		orderShopStatus = MgOrderShopStart
	case domain.OrderShopStatusReady:
		orderShopStatus = MgOrderShopReady
	case domain.OrderShopStatusDone:
		orderShopStatus = MgOrderShopDone
	}

	return MgOrderShop{
		ID:              orderShop.ID.String(),
		ShopID:          orderShop.ShopID.String(),
		OrderCustomerID: orderShop.OrderCustomerID.String(),
		Status:          orderShopStatus,
		Notified:        orderShop.Notified,
	}
}

type MgOrderShopItem struct {
	ID          string `bson:"_id"`
	OrderShopID string `bson:"order_shop_id"`
	ProductID   string `bson:"product_id"`
	Quantity    int64  `bson:"quantity"`
}

func (osi *MgOrderShopItem) ToDomain() domain.OrderShopItem {
	return domain.OrderShopItem{
		ID:          domain.ID(osi.ID),
		OrderShopID: domain.ID(osi.OrderShopID),
		ProductID:   domain.ID(osi.ProductID),
		Quantity:    osi.Quantity,
	}
}

func NewMgOrderShopItem(orderShopItem domain.OrderShopItem) MgOrderShopItem {
	return MgOrderShopItem{
		ID:          orderShopItem.ID.String(),
		OrderShopID: orderShopItem.OrderShopID.String(),
		ProductID:   orderShopItem.ProductID.String(),
		Quantity:    orderShopItem.Quantity,
	}
}
