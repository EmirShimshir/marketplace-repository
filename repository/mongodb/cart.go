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

type MongoCartRepo struct{
	db *mongo.Collection
}

func NewCartRepo(db *mongo.Database) *MongoCartRepo {
	collection := db.Collection(CartProductCollection)
	indexModel := mongo.IndexModel{
		Keys:    bson.D{{"cart_id", 1}, {"product_id", 1}},
		Options: options.Index().SetUnique(true),
	}

	_, err := collection.Indexes().CreateOne(context.Background(), indexModel)
	if err != nil {
		log.Fatalf("unable to create cart product collection index, %v", err)
	}

	return &MongoCartRepo{
		db: db.Collection(CartCollection),
	}
}

func (c *MongoCartRepo) GetCartByID(ctx context.Context, cartID domain.ID) (domain.Cart, error) {
	var mgCart entity.MgCart
	if err := c.db.FindOne(ctx, bson.M{"_id": cartID}).Decode(&mgCart); err != nil {
		if err == mongo.ErrNoDocuments {
			return domain.Cart{}, errors.Wrap(domain.ErrNotExist, err.Error())
		}
		return domain.Cart{}, errors.Wrap(domain.ErrPersistenceFailed, err.Error())
	}

	cursor, err := c.db.Database().Collection(CartProductCollection).Find(ctx, bson.M{"cart_id":cartID})
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return domain.Cart{}, errors.Wrap(domain.ErrNotExist, err.Error())
		}
		return domain.Cart{}, errors.Wrap(domain.ErrPersistenceFailed, err.Error())
	}

	var mgCartItems []entity.MgCartItem
	err = cursor.All(ctx, &mgCartItems)
	if err != nil {
		return domain.Cart{}, err
	}

	cartItems := make([]domain.CartItem, len(mgCartItems))
	for i, cartItem := range mgCartItems {
		cartItems[i] = cartItem.ToDomain()
	}

	cart := mgCart.ToDomain()
	cart.Items = cartItems

	return cart, nil
}

func (c *MongoCartRepo) UpdateCart(ctx context.Context, cart domain.Cart) (domain.Cart, error) {
	var mgCart = entity.NewMgCart(cart)
	_, err := c.db.ReplaceOne(ctx, bson.M{"_id": mgCart.ID}, mgCart)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return domain.Cart{}, errors.Wrap(domain.ErrNotExist, err.Error())
		}
		return domain.Cart{}, errors.Wrap(domain.ErrPersistenceFailed, err.Error())
	}

	return c.GetCartByID(ctx, cart.ID)
}

func (c *MongoCartRepo) ClearCart(ctx context.Context, cartID domain.ID) error {
	_, err := c.db.Database().Collection(CartProductCollection).DeleteMany(ctx, bson.M{"cart_id": cartID})
	if err != nil {
		return errors.Wrap(domain.ErrDeleteFailed, err.Error())
	}
	return nil
}

func (c *MongoCartRepo) GetCartItemByID(ctx context.Context, cartItemID domain.ID) (domain.CartItem, error) {
	var mgCartItem entity.MgCartItem
	if err := c.db.Database().Collection(CartProductCollection).FindOne(ctx, bson.M{"_id": cartItemID}).Decode(&mgCartItem); err != nil {
		if err == mongo.ErrNoDocuments {
			return domain.CartItem{}, errors.Wrap(domain.ErrNotExist, err.Error())
		}
		return domain.CartItem{}, errors.Wrap(domain.ErrPersistenceFailed, err.Error())
	}

	return mgCartItem.ToDomain(), nil
}

func (c *MongoCartRepo) CreateCartItem(ctx context.Context, cartItem domain.CartItem) (domain.CartItem, error) {
	var mgCartItem = entity.NewMgCartItem(cartItem)
	_, err := c.db.Database().Collection(CartProductCollection).InsertOne(ctx, mgCartItem)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return domain.CartItem{}, errors.Wrap(domain.ErrDuplicate, err.Error())
		}
		return domain.CartItem{}, errors.Wrap(domain.ErrPersistenceFailed, err.Error())
	}

	return c.GetCartItemByID(ctx, cartItem.ID)
}

func (c *MongoCartRepo) UpdateCartItem(ctx context.Context, cartItem domain.CartItem) (domain.CartItem, error) {
	var mgCartItem = entity.NewMgCartItem(cartItem)
	_, err := c.db.Database().Collection(CartProductCollection).ReplaceOne(ctx, bson.M{"_id": mgCartItem.ID}, mgCartItem)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return domain.CartItem{}, errors.Wrap(domain.ErrNotExist, err.Error())
		}
		return domain.CartItem{}, errors.Wrap(domain.ErrPersistenceFailed, err.Error())
	}

	return c.GetCartItemByID(ctx, cartItem.ID)
}

func (c *MongoCartRepo) DeleteCartItem(ctx context.Context, cartItemID domain.ID) error {
	_, err := c.db.Database().Collection(CartProductCollection).DeleteOne(ctx, bson.M{"_id": cartItemID})
	if err != nil {
		return errors.Wrap(domain.ErrDeleteFailed, err.Error())
	}
	return nil
}