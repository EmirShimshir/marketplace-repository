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

type MongoShopRepo struct{
	db *mongo.Collection
}

func NewShopRepo(db *mongo.Database) *MongoShopRepo {
	collection := db.Collection(ShopProductCollection)
	indexModel := mongo.IndexModel{
		Keys:    bson.D{{"shop_id", 1}, {"product_id", 1}},
		Options: options.Index().SetUnique(true),
	}

	_, err := collection.Indexes().CreateOne(context.Background(), indexModel)
	if err != nil {
		log.Fatalf("unable to create cart product collection index, %v", err)
	}

	return &MongoShopRepo{
		db: db.Collection(ShopCollection),
	}
}

func (s *MongoShopRepo) GetShops(ctx context.Context, limit, offset int64) ([]domain.Shop, error) {
	cursor, err := s.db.Find(ctx, bson.M{}, options.Find().SetSkip(offset).SetLimit(limit))
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.Wrap(domain.ErrNotExist, err.Error())
		}
		return nil, errors.Wrap(domain.ErrPersistenceFailed, err.Error())
	}

	var mgShopsArray []entity.MgShop
	err = cursor.All(ctx, &mgShopsArray)
	if err != nil {
		return nil, err
	}

	shops := make([]domain.Shop, len(mgShopsArray))
	for i, shop := range mgShopsArray {
		shops[i] = shop.ToDomain()
		shopItems, err := s.getShopItemsByShopID(ctx, shops[i].ID)
		if err != nil {
			return nil, err
		}
		shops[i].Items = shopItems
	}

	return shops, nil
}

func (s *MongoShopRepo) GetShopByID(ctx context.Context, shopID domain.ID) (domain.Shop, error) {
	result := s.db.FindOne(ctx, bson.M{"_id": shopID})

	var mgShop entity.MgShop
	if err := result.Decode(&mgShop); err != nil {
		if err == mongo.ErrNoDocuments {
			return domain.Shop{}, errors.Wrap(domain.ErrNotExist, err.Error())
		}
		return domain.Shop{}, errors.Wrap(domain.ErrPersistenceFailed, err.Error())
	}

	shop := mgShop.ToDomain()
	shopItems, err := s.getShopItemsByShopID(ctx, shop.ID)
	if err != nil {
		return domain.Shop{}, err
	}
	shop.Items = shopItems

	return shop, nil
}

func (s *MongoShopRepo) GetShopBySellerID(ctx context.Context, sellerID domain.ID) ([]domain.Shop, error) {
	cursor, err := s.db.Find(ctx, bson.M{"seller_id":sellerID})
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.Wrap(domain.ErrNotExist, err.Error())
		}
		return nil, errors.Wrap(domain.ErrPersistenceFailed, err.Error())
	}

	var mgShopsArray []entity.MgShop
	err = cursor.All(ctx, &mgShopsArray)
	if err != nil {
		return nil, err
	}

	shops := make([]domain.Shop, len(mgShopsArray))
	for i, shop := range mgShopsArray {
		shops[i] = shop.ToDomain()
		shopItems, err := s.getShopItemsByShopID(ctx, shops[i].ID)
		if err != nil {
			return nil, err
		}
		shops[i].Items = shopItems
	}

	return shops, nil
}

func (s *MongoShopRepo) CreateShop(ctx context.Context, shop domain.Shop) (domain.Shop, error) {
	var mgShop = entity.NewMgShop(shop)
	_, err := s.db.InsertOne(ctx, mgShop)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return domain.Shop{}, errors.Wrap(domain.ErrDuplicate, err.Error())
		}
		return domain.Shop{}, errors.Wrap(domain.ErrPersistenceFailed, err.Error())
	}

	return s.GetShopByID(ctx, shop.ID)
}

func (s *MongoShopRepo) UpdateShop(ctx context.Context, shop domain.Shop) (domain.Shop, error) {
	var mgShop = entity.NewMgShop(shop)
	_, err := s.db.ReplaceOne(ctx, bson.M{"_id": mgShop.ID}, mgShop)
	if err != nil {
		if err == mongo.ErrNoDocuments {
		return domain.Shop{}, errors.Wrap(domain.ErrNotExist, err.Error())
	}
		return domain.Shop{}, errors.Wrap(domain.ErrPersistenceFailed, err.Error())
	}


	return s.GetShopByID(ctx, shop.ID)
}

func (s *MongoShopRepo) DeleteShop(ctx context.Context, shopID domain.ID) error {
	_, err := s.db.DeleteOne(ctx, bson.M{"_id": shopID})
	if err != nil {
		return errors.Wrap(domain.ErrDeleteFailed, err.Error())
	}
	return nil
}

func (s *MongoShopRepo) GetShopItems(ctx context.Context, limit, offset int64) ([]domain.ShopItem, error) {
	cursor, err := s.db.Database().Collection(ShopProductCollection).Find(ctx, bson.M{}, options.Find().SetSkip(offset).SetLimit(limit))
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.Wrap(domain.ErrNotExist, err.Error())
		}
		return nil, errors.Wrap(domain.ErrPersistenceFailed, err.Error())
	}

	var mgShopItems []entity.MgShopItem
	err = cursor.All(ctx, &mgShopItems)
	if err != nil {
		return nil, err
	}

	shopItems := make([]domain.ShopItem, len(mgShopItems))
	for i, shopItem := range mgShopItems {
		shopItems[i] = shopItem.ToDomain()
	}
	return shopItems, nil
}

func (s *MongoShopRepo) GetShopItemByID(ctx context.Context, shopItemID domain.ID) (domain.ShopItem, error) {
	result := s.db.Database().Collection(ShopProductCollection).FindOne(ctx, bson.M{"_id": shopItemID})

	var mgShopItem entity.MgShopItem
	if err := result.Decode(&mgShopItem); err != nil {
		if err == mongo.ErrNoDocuments {
			return domain.ShopItem{}, errors.Wrap(domain.ErrNotExist, err.Error())
		}
		return domain.ShopItem{}, errors.Wrap(domain.ErrPersistenceFailed, err.Error())
	}

	return mgShopItem.ToDomain(), nil
}

func (s *MongoShopRepo) GetShopItemByProductID(ctx context.Context, productID domain.ID) (domain.ShopItem, error) {
	result := s.db.Database().Collection(ShopProductCollection).FindOne(ctx, bson.M{"product_id": productID})

	var mgShopItem entity.MgShopItem
	if err := result.Decode(&mgShopItem); err != nil {
		if err == mongo.ErrNoDocuments {
			return domain.ShopItem{}, errors.Wrap(domain.ErrNotExist, err.Error())
		}
		return domain.ShopItem{}, errors.Wrap(domain.ErrPersistenceFailed, err.Error())
	}

	return mgShopItem.ToDomain(), nil
}

func (s *MongoShopRepo) CreateShopItem(ctx context.Context, shopItem domain.ShopItem, product domain.Product) (domain.ShopItem, error) {
	session, err := s.db.Database().Client().StartSession()
	if err != nil {
		return domain.ShopItem{}, nil
	}
	defer session.EndSession(ctx)

	err = mongo.WithSession(ctx, session, func(sessionContext mongo.SessionContext) error {
		var mgProduct = entity.NewMgProduct(product)
		_, err := s.db.Database().Collection(ProductCollection).InsertOne(sessionContext, mgProduct)
		if err != nil {
			if mongo.IsDuplicateKeyError(err) {
				return errors.Wrap(domain.ErrDuplicate, err.Error())
			}
			return errors.Wrap(domain.ErrPersistenceFailed, err.Error())
		}

		var mgShopItem = entity.NewMgShopItem(shopItem)
		_, err = s.db.Database().Collection(ShopProductCollection).InsertOne(sessionContext, mgShopItem)
		if err != nil {
			if mongo.IsDuplicateKeyError(err) {
				return errors.Wrap(domain.ErrDuplicate, err.Error())
			}
			return errors.Wrap(domain.ErrPersistenceFailed, err.Error())
		}

		return nil
	})
	if err != nil {
		return domain.ShopItem{}, err
	}

	return s.GetShopItemByID(ctx, shopItem.ID)
}

func (s *MongoShopRepo) UpdateShopItem(ctx context.Context, shopItem domain.ShopItem) (domain.ShopItem, error) {
	var mgShopItem = entity.NewMgShopItem(shopItem)
	_, err := s.db.Database().Collection(ShopProductCollection).ReplaceOne(ctx, bson.M{"_id": mgShopItem.ID}, mgShopItem)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return domain.ShopItem{}, errors.Wrap(domain.ErrNotExist, err.Error())
		}
		return domain.ShopItem{}, errors.Wrap(domain.ErrPersistenceFailed, err.Error())
	}

	return s.GetShopItemByID(ctx, shopItem.ID)
}

func (s *MongoShopRepo) DeleteShopItem(ctx context.Context, shopItemID domain.ID) error {
	_, err := s.db.Database().Collection(ShopProductCollection).DeleteOne(ctx, bson.M{"_id": shopItemID})
	if err != nil {
		return errors.Wrap(domain.ErrDeleteFailed, err.Error())
	}
	return nil
}

func (s *MongoShopRepo) getShopItemsByShopID(ctx context.Context, shopID domain.ID) ([]domain.ShopItem, error) {
	cursor, err := s.db.Database().Collection(ShopProductCollection).Find(ctx, bson.M{"shop_id": shopID})
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.Wrap(domain.ErrNotExist, err.Error())
		}
		return nil, errors.Wrap(domain.ErrPersistenceFailed, err.Error())
	}

	var mgShopItems []entity.MgShopItem
	err = cursor.All(ctx, &mgShopItems)
	if err != nil {
		return nil, err
	}

	shopItems := make([]domain.ShopItem, len(mgShopItems))
	for i, shopItem := range mgShopItems {
		shopItems[i] = shopItem.ToDomain()
	}
	return shopItems, nil
}