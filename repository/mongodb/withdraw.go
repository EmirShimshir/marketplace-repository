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

type MongoWithdrawRepo struct {
	db *mongo.Collection
}

func NewWithdrawRepo(db *mongo.Database) *MongoWithdrawRepo {
	return &MongoWithdrawRepo{
		db: db.Collection(WithdrawCollection),
	}
}

func (w *MongoWithdrawRepo) Get(ctx context.Context, limit, offset int64) ([]domain.Withdraw, error) {
	cursor, err := w.db.Find(ctx, bson.M{}, options.Find().SetSkip(offset).SetLimit(limit))
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.Wrap(domain.ErrNotExist, err.Error())
		}
		return nil, errors.Wrap(domain.ErrPersistenceFailed, err.Error())
	}

	var mgWithdrawArray []entity.MgWithdraw
	err = cursor.All(ctx, &mgWithdrawArray)
	if err != nil {
		return nil, err
	}

	withdraws := make([]domain.Withdraw, len(mgWithdrawArray))
	for i, withdraw := range mgWithdrawArray {
		withdraws[i] = withdraw.ToDomain()
	}

	return withdraws, nil
}

func (w *MongoWithdrawRepo) GetByID(ctx context.Context, withdrawID domain.ID) (domain.Withdraw, error) {
	result := w.db.FindOne(ctx, bson.M{"_id": withdrawID})

	var mgWithdraw entity.MgWithdraw
	if err := result.Decode(&mgWithdraw); err != nil {
		if err == mongo.ErrNoDocuments {
			return domain.Withdraw{}, errors.Wrap(domain.ErrNotExist, err.Error())
		}
		return domain.Withdraw{}, errors.Wrap(domain.ErrPersistenceFailed, err.Error())
	}
	return mgWithdraw.ToDomain(), nil
}

func (w *MongoWithdrawRepo) GetByShopID(ctx context.Context, shopID domain.ID) ([]domain.Withdraw, error) {
	cursor, err := w.db.Find(ctx, bson.M{"shop_id": shopID})
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.Wrap(domain.ErrNotExist, err.Error())
		}
		return nil, errors.Wrap(domain.ErrPersistenceFailed, err.Error())
	}

	var mgWithdrawArray []entity.MgWithdraw
	err = cursor.All(ctx, &mgWithdrawArray)
	if err != nil {
		return nil, err
	}

	withdraws := make([]domain.Withdraw, len(mgWithdrawArray))
	for i, withdraw := range mgWithdrawArray {
		withdraws[i] = withdraw.ToDomain()
	}

	return withdraws, nil
}

func (w *MongoWithdrawRepo) Create(ctx context.Context, withdraw domain.Withdraw) (domain.Withdraw, error) {
	var mgWithdraw = entity.NewMgWithdraw(withdraw)
	_, err := w.db.InsertOne(ctx, mgWithdraw)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return domain.Withdraw{}, errors.Wrap(domain.ErrDuplicate, err.Error())
		}
		return domain.Withdraw{}, errors.Wrap(domain.ErrPersistenceFailed, err.Error())
	}

	return w.GetByID(ctx, withdraw.ID)
}

func (w *MongoWithdrawRepo) Update(ctx context.Context, withdraw domain.Withdraw) (domain.Withdraw, error) {
	var mgWithdraw = entity.NewMgWithdraw(withdraw)
	_, err := w.db.ReplaceOne(ctx, bson.M{"_id": mgWithdraw.ID}, mgWithdraw)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return domain.Withdraw{}, errors.Wrap(domain.ErrNotExist, err.Error())
		}
		return domain.Withdraw{}, errors.Wrap(domain.ErrPersistenceFailed, err.Error())
	}

	return w.GetByID(ctx, withdraw.ID)
}

func (w *MongoWithdrawRepo) Delete(ctx context.Context, withdrawID domain.ID) error {
	_, err := w.db.DeleteOne(ctx, bson.M{"_id": withdrawID})
	if err != nil {
		return errors.Wrap(domain.ErrDeleteFailed, err.Error())
	}
	return nil
}
