package mongodb

import (
	"context"
	"github.com/EmirShimshir/marketplace-core/domain"
	"github.com/EmirShimshir/marketplace-repository/repository/mongodb/entity"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoProductRepo struct{
	db *mongo.Collection
}

func NewProductRepo(db *mongo.Database) *MongoProductRepo {
	return &MongoProductRepo{
		db: db.Collection(ProductCollection),
	}
}

func (p *MongoProductRepo) Get(ctx context.Context, limit, offset int64) ([]domain.Product, error) {
	cursor, err := p.db.Find(ctx, bson.M{}, options.Find().SetSkip(offset).SetLimit(limit))
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.Wrap(domain.ErrNotExist, err.Error())
		}
		return nil, errors.Wrap(domain.ErrPersistenceFailed, err.Error())
	}

	var mgProductsArray []entity.MgProduct
	err = cursor.All(ctx, &mgProductsArray)
	if err != nil {
		return nil, err
	}

	products := make([]domain.Product, len(mgProductsArray))
	for i, product := range mgProductsArray {
		products[i] = product.ToDomain()
	}

	return products, nil
}

func (p *MongoProductRepo) GetByID(ctx context.Context, productID domain.ID) (domain.Product, error) {
	result := p.db.FindOne(ctx, bson.M{"_id": productID})

	var mgProduct entity.MgProduct
	if err := result.Decode(&mgProduct); err != nil {
		if err == mongo.ErrNoDocuments {
			return domain.Product{}, errors.Wrap(domain.ErrNotExist, err.Error())
		}
		return domain.Product{}, errors.Wrap(domain.ErrPersistenceFailed, err.Error())
	}
	return mgProduct.ToDomain(), nil
}

func (p *MongoProductRepo) Create(ctx context.Context, product domain.Product) (domain.Product, error) {
	var mgProduct = entity.NewMgProduct(product)
	_, err := p.db.InsertOne(ctx, mgProduct)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return domain.Product{}, errors.Wrap(domain.ErrDuplicate, err.Error())
		}
		return domain.Product{}, errors.Wrap(domain.ErrPersistenceFailed, err.Error())
	}

	return p.GetByID(ctx, product.ID)
}

func (p *MongoProductRepo) Update(ctx context.Context, product domain.Product) (domain.Product, error) {
	var mgProduct = entity.NewMgProduct(product)
	_, err := p.db.ReplaceOne(ctx, bson.M{"_id": mgProduct.ID}, mgProduct)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return domain.Product{}, errors.Wrap(domain.ErrNotExist, err.Error())
		}
		return domain.Product{}, errors.Wrap(domain.ErrPersistenceFailed, err.Error())
	}

	return p.GetByID(ctx, product.ID)
}

func (p *MongoProductRepo) Delete(ctx context.Context, productID domain.ID) error {
	_, err := p.db.DeleteOne(ctx, bson.M{"_id": productID})
	if err != nil {
		return errors.Wrap(domain.ErrDeleteFailed, err.Error())
	}
	return nil
}
