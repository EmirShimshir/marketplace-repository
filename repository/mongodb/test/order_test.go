package mongodb

import (
	"context"
	"github.com/EmirShimshir/marketplace-core/domain"
	"github.com/EmirShimshir/marketplace-repository/repository/mongodb"
	"github.com/EmirShimshir/marketplace-repository/repository/mongodb/entity"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/mongo"
	"testing"
	"time"
)

var orderShopItems = []domain.OrderShopItem{
	domain.OrderShopItem{
		ID:          domain.ID("30e18bc1-4354-4937-9a3b-03cf0b70eee1"),
		OrderShopID: domain.ID("30e18bc1-4354-4937-9a3b-03cf0b702ee1"),
		ProductID:   domain.ID("30e18bc1-4354-4937-9a3b-03cf0b7027a1"),
		Quantity:    1,
	},
}

var orderShops = []domain.OrderShop{
	domain.OrderShop{
		ID:              domain.ID("30e18bc1-4354-4937-9a3b-03cf0b702ee1"),
		ShopID:          domain.ID("30e18bc1-4354-4937-9a3b-03cf0b7027b1"),
		OrderCustomerID: domain.ID("30e18bc1-4354-4937-9a3b-03cf0b702ae1"),
		Status:          domain.OrderShopStatusStart,
		Notified:        false,
		OrderShopItems:  orderShopItems,
	},
}

var orderCustomers = []domain.OrderCustomer{
	domain.OrderCustomer{
		ID:         domain.ID("30e18bc1-4354-4937-9a3b-03cf0b702ae1"),
		CustomerID: domain.ID("30e18bc1-4354-4937-9a3b-03cf0b7027cc"),
		Address:    "Pushkina 1-2-3",
		CreatedAt:  time.Date(2022, 10, 10, 11, 30, 30, 0, time.UTC),
		OrderShops: orderShops,
	},
}

var createdOrderShopItems = []domain.OrderShopItem{
	domain.OrderShopItem{
		ID:          domain.ID("30e18bc1-4354-4937-9a3b-03cf0beeeeee"),
		OrderShopID: domain.ID("30e18bc1-4354-4937-9a3b-03cf0b7eeeee"),
		ProductID:   domain.ID("30e18bc1-4354-4937-9a3b-03cf0b7027a1"),
		Quantity:    3,
	},
}

var createdorderShops = []domain.OrderShop{
	domain.OrderShop{
		ID:              domain.ID("30e18bc1-4354-4937-9a3b-03cf0b7eeeee"),
		ShopID:          domain.ID("30e18bc1-4354-4937-9a3b-03cf0b7027b1"),
		OrderCustomerID: domain.ID("30e18bc1-4354-4937-9a3b-03cf0b70eeee"),
		Status:          domain.OrderShopStatusStart,
		Notified:        true,
		OrderShopItems:  createdOrderShopItems,
	},
}

var createdOrderCustomers = []domain.OrderCustomer{
	domain.OrderCustomer{
		ID:         domain.ID("30e18bc1-4354-4937-9a3b-03cf0b70eeee"),
		CustomerID: domain.ID("30e18bc1-4354-4937-9a3b-03cf0b7027cc"),
		Address:    "Pushkina 1-2-4",
		CreatedAt:  time.Date(2024, 10, 10, 11, 30, 30, 0, time.UTC),
		OrderShops: createdorderShops,
	},
}

func InitOrderCustomersMongoDB(ctx context.Context, db *mongo.Database) error {
	for _, orderCustomer := range orderCustomers {
		var mgOrderCustomer = entity.NewMgOrderCustomer(orderCustomer)
		_, err := db.Collection(mongodb.OrderCustomerCollection).InsertOne(ctx, mgOrderCustomer)
		if err != nil {
			return err
		}
	}

	return nil
}

func InitOrderShopsMongoDB(ctx context.Context, db *mongo.Database) error {
	for _, orderShop := range orderShops {
		var mgOrderShop = entity.NewMgOrderShop(orderShop)
		_, err := db.Collection(mongodb.OrderShopCollection).InsertOne(ctx, mgOrderShop)
		if err != nil {
			return err
		}
	}

	return nil
}

func InitOrderShopItemsMongoDB(ctx context.Context, db *mongo.Database) error {
	for _, orderShopItem := range orderShopItems {
		var mgOrderShopItem = entity.NewMgOrderShopItem(orderShopItem)
		_, err := db.Collection(mongodb.OrderShopProductCollection).InsertOne(ctx, mgOrderShopItem)
		if err != nil {
			return err
		}
	}

	return nil
}

func TestOrderRepository(t *testing.T) {
	ctx := context.Background()
	container, err := newMongoContainer(ctx)
	if err != nil {
		t.Fatal(err)
	}

	// Clean up the container after the test is complete
	t.Cleanup(func() {
		if err := container.Terminate(ctx); err != nil {
			t.Fatalf("failed to terminate container: %s", err)
		}
	})

	url, err := container.ConnectionString(ctx)
	if err != nil {
		t.Fatal(err)
	}

	db, err := newMongoDB(ctx, url)
	if err != nil {
		t.Fatal(err)
	}

	err = InitOrderCustomersMongoDB(ctx, db)
	if err != nil {
		t.Fatal(err)
	}
	err = InitOrderShopsMongoDB(ctx, db)
	if err != nil {
		t.Fatal(err)
	}
	err = InitOrderShopItemsMongoDB(ctx, db)
	if err != nil {
		t.Fatal(err)
	}

	t.Run("test GetOrderCustomerByID", func(t *testing.T) {
		repo := mongodb.NewOrderRepo(db)
		found, err := repo.GetOrderCustomerByID(ctx, orderCustomers[0].ID)
		if err != nil {
			t.Errorf("failed to GetOrderCustomerByID: %v", err)
		}

		require.Equal(t, orderCustomers[0], found)
	})
	t.Run("test CreateOrderCustomer", func(t *testing.T) {
		repo := mongodb.NewOrderRepo(db)
		found, err := repo.CreateOrderCustomer(ctx, createdOrderCustomers[0])
		if err != nil {
			t.Errorf("failed to CreateOrderCustomer: %v", err)
		}

		require.Equal(t, createdOrderCustomers[0], found)
	})
}
