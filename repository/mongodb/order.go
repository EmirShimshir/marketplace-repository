package mongodb

import (
	"context"
	"github.com/EmirShimshir/marketplace-core/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

type MongoOrderRepo struct {
	db *mongo.Collection
}

func NewOrderRepo(db *mongo.Database) *MongoOrderRepo {
	collection := db.Collection(OrderShopCollection)
	indexModel := mongo.IndexModel{
		Keys:    bson.D{{"shop_id", 1}, {"order_customer_id", 1}},
		Options: options.Index().SetUnique(true),
	}

	_, err := collection.Indexes().CreateOne(context.Background(), indexModel)
	if err != nil {
		log.Fatalf("unable to create OrderShopCollection index, %v", err)
	}

	collection = db.Collection(OrderShopProductCollection)
	indexModel = mongo.IndexModel{
		Keys:    bson.D{{"order_shop_id", 1}, {"product_id", 1}},
		Options: options.Index().SetUnique(true),
	}

	_, err = collection.Indexes().CreateOne(context.Background(), indexModel)
	if err != nil {
		log.Fatalf("unable to create OrderShopProductCollection index, %v", err)
	}

	return &MongoOrderRepo{
		db: db.Collection(OrderCustomerCollection),
	}
}

func (o *MongoOrderRepo) GetOrderCustomerByCustomerID(ctx context.Context, customerID domain.ID) ([]domain.OrderCustomer, error) {
	//TODO implement me
	panic("implement me")
}

func (o *MongoOrderRepo) GetOrderCustomerByID(ctx context.Context, orderCustomerID domain.ID) (domain.OrderCustomer, error) {
	//TODO implement me
	panic("implement me")
}

func (o *MongoOrderRepo) GetOrderShopByID(ctx context.Context, orderShopID domain.ID) (domain.OrderShop, error) {
	//TODO implement me
	panic("implement me")
}

func (o *MongoOrderRepo) GetNoNotifiedOrderShops(ctx context.Context) ([]domain.OrderShop, error) {
	//TODO implement me
	panic("implement me")
}

func (o *MongoOrderRepo) CreateOrderCustomer(ctx context.Context, orderCustomer domain.OrderCustomer) (domain.OrderCustomer, error) {
	//TODO implement me
	panic("implement me")
}

func (o *MongoOrderRepo) GetOrderShopByShopID(ctx context.Context, shopID domain.ID) ([]domain.OrderShop, error) {
	//TODO implement me
	panic("implement me")
}

func (o *MongoOrderRepo) UpdateOrderShop(ctx context.Context, orderShop domain.OrderShop) (domain.OrderShop, error) {
	//TODO implement me
	panic("implement me")
}

func (o *MongoOrderRepo) UpdatePaymentStatus(ctx context.Context, orderCustomerID domain.ID) error {
	//TODO implement me
	panic("implement me")
}
