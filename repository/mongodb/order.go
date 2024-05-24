package mongodb

import (
	"context"
	"github.com/EmirShimshir/marketplace-core/domain"
	"github.com/EmirShimshir/marketplace-repository/repository/mongodb/entity"
	"github.com/pkg/errors"
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
	cursor, err := o.db.Find(ctx, bson.M{"customer_id": customerID})
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.Wrap(domain.ErrNotExist, err.Error())
		}
		return nil, errors.Wrap(domain.ErrPersistenceFailed, err.Error())
	}
	var mgOrderCustomersArray []entity.MgOrderCustomer
	err = cursor.All(ctx, &mgOrderCustomersArray)
	if err != nil {
		return nil, err
	}

	orderCustomers := make([]domain.OrderCustomer, len(mgOrderCustomersArray))
	for i := range orderCustomers {
		orderCustomers[i] = mgOrderCustomersArray[i].ToDomain()
		var err error
		orderCustomers[i], err = o.GetOrderCustomerByID(ctx, orderCustomers[i].ID)
		if err != nil {
			return nil, err
		}
	}

	return orderCustomers, nil
}

func (o *MongoOrderRepo) GetOrderCustomerByID(ctx context.Context, orderCustomerID domain.ID) (domain.OrderCustomer, error) {
	result := o.db.FindOne(ctx, bson.M{"_id": orderCustomerID})
	var mgOrderCustomer entity.MgOrderCustomer
	if err := result.Decode(&mgOrderCustomer); err != nil {
		if err == mongo.ErrNoDocuments {
			return domain.OrderCustomer{}, errors.Wrap(domain.ErrNotExist, err.Error())
		}
		return domain.OrderCustomer{}, errors.Wrap(domain.ErrPersistenceFailed, err.Error())
	}
	orderCustomer := mgOrderCustomer.ToDomain()

	cursor, err := o.db.Database().Collection(OrderShopCollection).Find(ctx, bson.M{"order_customer_id": orderCustomerID})
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return domain.OrderCustomer{}, errors.Wrap(domain.ErrNotExist, err.Error())
		}
		return domain.OrderCustomer{}, errors.Wrap(domain.ErrPersistenceFailed, err.Error())
	}
	var mgOrderShopsArray []entity.MgOrderShop
	err = cursor.All(ctx, &mgOrderShopsArray)
	if err != nil {
		return domain.OrderCustomer{}, err
	}

	orderShops := make([]domain.OrderShop, len(mgOrderShopsArray))
	for i := range orderShops {
		orderShops[i] = mgOrderShopsArray[i].ToDomain()
		var err error
		orderShops[i].OrderShopItems, err = o.getOrderShopItemsByOrderShopID(ctx, orderShops[i].ID)
		if err != nil {
			return domain.OrderCustomer{}, err
		}
	}
	orderCustomer.OrderShops = orderShops

	return orderCustomer, nil
}

func (o *MongoOrderRepo) getOrderShopItemsByOrderShopID(ctx context.Context, orderShopID domain.ID) ([]domain.OrderShopItem, error) {
	cursor, err := o.db.Database().Collection(OrderShopProductCollection).Find(ctx, bson.M{"order_shop_id": orderShopID})
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.Wrap(domain.ErrNotExist, err.Error())
		}
		return nil, errors.Wrap(domain.ErrPersistenceFailed, err.Error())
	}

	var mgOrderShopItems []entity.MgOrderShopItem
	err = cursor.All(ctx, &mgOrderShopItems)
	if err != nil {
		return nil, err
	}

	OrderShopItems := make([]domain.OrderShopItem, len(mgOrderShopItems))
	for i, OrderShopItem := range mgOrderShopItems {
		OrderShopItems[i] = OrderShopItem.ToDomain()
	}
	return OrderShopItems, nil
}

func (o *MongoOrderRepo) GetOrderShopByID(ctx context.Context, orderShopID domain.ID) (domain.OrderShop, error) {
	result := o.db.Database().Collection(OrderShopCollection).FindOne(ctx, bson.M{"_id": orderShopID})
	var mgOrderShop entity.MgOrderShop
	if err := result.Decode(&mgOrderShop); err != nil {
		if err == mongo.ErrNoDocuments {
			return domain.OrderShop{}, errors.Wrap(domain.ErrNotExist, err.Error())
		}
		return domain.OrderShop{}, errors.Wrap(domain.ErrPersistenceFailed, err.Error())
	}

	orderShopItems, err := o.getOrderShopItemsByOrderShopID(ctx, orderShopID)
	if err != nil {
		return domain.OrderShop{}, err
	}

	orderShop := mgOrderShop.ToDomain()
	orderShop.OrderShopItems = orderShopItems

	return orderShop, nil
}

func (o *MongoOrderRepo) GetNoNotifiedOrderShops(ctx context.Context) ([]domain.OrderShop, error) {
	cursor, err := o.db.Database().Collection(OrderShopCollection).Find(ctx, bson.M{"notified": false})
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.Wrap(domain.ErrNotExist, err.Error())
		}
		return nil, errors.Wrap(domain.ErrPersistenceFailed, err.Error())
	}

	var mgOrderShopsArray []entity.MgOrderShop
	err = cursor.All(ctx, &mgOrderShopsArray)
	if err != nil {
		return nil, err
	}

	orderShops := make([]domain.OrderShop, len(mgOrderShopsArray))
	for i := range orderShops {
		orderShops[i] = mgOrderShopsArray[i].ToDomain()
		var err error
		orderShops[i].OrderShopItems, err = o.getOrderShopItemsByOrderShopID(ctx, orderShops[i].ID)
		if err != nil {
			return nil, err
		}
	}

	return orderShops, nil
}

func (o *MongoOrderRepo) getMgEntities(orderCustomer domain.OrderCustomer) (entity.MgOrderCustomer, []entity.MgOrderShop, []entity.MgOrderShopItem) {
	mgOrderCustomer := entity.NewMgOrderCustomer(orderCustomer)
	mgOrderShops := make([]entity.MgOrderShop, 0)
	mgOrderShopItems := make([]entity.MgOrderShopItem, 0)
	for _, orderShop := range orderCustomer.OrderShops {
		mgOrderShops = append(mgOrderShops, entity.NewMgOrderShop(orderShop))
		for _, item := range orderShop.OrderShopItems {
			mgOrderShopItems = append(mgOrderShopItems, entity.NewMgOrderShopItem(item))
		}
	}
	return mgOrderCustomer, mgOrderShops, mgOrderShopItems
}

func (o *MongoOrderRepo) CreateOrderCustomer(ctx context.Context, orderCustomer domain.OrderCustomer) (domain.OrderCustomer, error) {
	session, err := o.db.Database().Client().StartSession()
	if err != nil {
		return domain.OrderCustomer{}, nil
	}
	defer session.EndSession(ctx)

	err = mongo.WithSession(ctx, session, func(sessionContext mongo.SessionContext) error {
		mgOrderCustomer, mgOrderShops, mgOrderShopItems := o.getMgEntities(orderCustomer)
		err = o.txInsertOrderCustomer(ctx, mgOrderCustomer)
		if err != nil {
			return err
		}
		for _, mgOrderShop := range mgOrderShops {
			err = o.txInsertOrderShop(ctx, mgOrderShop)
			if err != nil {
				return err
			}
		}
		for _, pgOrderShopItem := range mgOrderShopItems {
			err = o.txUpdateShopItem(ctx, pgOrderShopItem)
			if err != nil {
				return err
			}
			err = o.txInsertOrderShopItem(ctx, pgOrderShopItem)
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return domain.OrderCustomer{}, err
	}

	return o.GetOrderCustomerByID(ctx, orderCustomer.ID)
}

func (o *MongoOrderRepo) GetOrderShopByShopID(ctx context.Context, shopID domain.ID) ([]domain.OrderShop, error) {
	cursor, err := o.db.Database().Collection(OrderShopCollection).Find(ctx, bson.M{"shop_id": shopID})
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.Wrap(domain.ErrNotExist, err.Error())
		}
		return nil, errors.Wrap(domain.ErrPersistenceFailed, err.Error())
	}

	var mgOrderShopsArray []entity.MgOrderShop
	err = cursor.All(ctx, &mgOrderShopsArray)
	if err != nil {
		return nil, err
	}

	orderShops := make([]domain.OrderShop, len(mgOrderShopsArray))
	for i := range orderShops {
		orderShops[i] = mgOrderShopsArray[i].ToDomain()
		var err error
		orderShops[i].OrderShopItems, err = o.getOrderShopItemsByOrderShopID(ctx, orderShops[i].ID)
		if err != nil {
			return nil, err
		}
	}

	return orderShops, nil
}

func (o *MongoOrderRepo) UpdateOrderShop(ctx context.Context, orderShop domain.OrderShop) (domain.OrderShop, error) {
	var mgOrderShop = entity.NewMgOrderShop(orderShop)
	_, err := o.db.Database().Collection(OrderShopCollection).ReplaceOne(ctx, bson.M{"_id": mgOrderShop.ID}, mgOrderShop)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return domain.OrderShop{}, errors.Wrap(domain.ErrNotExist, err.Error())
		}
		return domain.OrderShop{}, errors.Wrap(domain.ErrPersistenceFailed, err.Error())
	}

	return o.GetOrderShopByID(ctx, orderShop.ID)
}

func (o *MongoOrderRepo) UpdatePaymentStatus(ctx context.Context, orderCustomerID domain.ID) error {
	updateQuery := bson.M{}
	updateQuery["payed"] = true

	_, err := o.db.UpdateOne(ctx, bson.M{"_id": orderCustomerID}, bson.M{"$set": updateQuery})
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return errors.Wrap(domain.ErrNotExist, err.Error())
		}
		return errors.Wrap(domain.ErrPersistenceFailed, err.Error())
	}

	return nil
}

func (o *MongoOrderRepo) txInsertOrderCustomer(ctx context.Context, customer entity.MgOrderCustomer) error {
	_, err := o.db.InsertOne(ctx, customer)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return errors.Wrap(domain.ErrDuplicate, err.Error())
		}
		return errors.Wrap(domain.ErrPersistenceFailed, err.Error())
	}
	return nil
}

func (o *MongoOrderRepo) txInsertOrderShop(ctx context.Context, shop entity.MgOrderShop) error {
	_, err := o.db.Database().Collection(OrderShopCollection).InsertOne(ctx, shop)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return errors.Wrap(domain.ErrDuplicate, err.Error())
		}
		return errors.Wrap(domain.ErrPersistenceFailed, err.Error())
	}
	return nil
}

func (o *MongoOrderRepo) txUpdateShopItem(ctx context.Context, item entity.MgOrderShopItem) error {
	_, err := o.db.Database().Collection(OrderShopProductCollection).ReplaceOne(ctx, bson.M{"_id": item.ID}, item)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return errors.Wrap(domain.ErrNotExist, err.Error())
		}
		return errors.Wrap(domain.ErrPersistenceFailed, err.Error())
	}
	return nil
}

func (o *MongoOrderRepo) txInsertOrderShopItem(ctx context.Context, item entity.MgOrderShopItem) error {
	_, err := o.db.Database().Collection(OrderShopProductCollection).InsertOne(ctx, item)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return errors.Wrap(domain.ErrDuplicate, err.Error())
		}
		return errors.Wrap(domain.ErrPersistenceFailed, err.Error())
	}
	return nil
}
