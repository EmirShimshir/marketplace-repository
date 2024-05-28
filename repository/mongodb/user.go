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

type MongoUserRepo struct{
	db *mongo.Collection
}

func NewUserRepo(db *mongo.Database) *MongoUserRepo {
	collection := db.Collection(UserCollection)
	indexModel := mongo.IndexModel{
		Keys:    bson.D{{"email", 1}},
		Options: options.Index().SetUnique(true),
	}

	_, err := collection.Indexes().CreateOne(context.Background(), indexModel)
	if err != nil {
		log.Fatalf("unable to create user collection index, %v", err)
	}

	return &MongoUserRepo{
		db: collection,
	}
}

func (u *MongoUserRepo) Get(ctx context.Context, limit, offset int64) ([]domain.User, error) {
	cursor, err := u.db.Find(ctx, bson.M{}, options.Find().SetSkip(offset).SetLimit(limit))
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.Wrap(domain.ErrNotExist, err.Error())
		}
		return nil, errors.Wrap(domain.ErrPersistenceFailed, err.Error())
	}

	var mgUsersArray []entity.MgUser
	err = cursor.All(ctx, &mgUsersArray)
	if err != nil {
		return nil, err
	}

	users := make([]domain.User, len(mgUsersArray))
	for i, user := range mgUsersArray {
		users[i] = user.ToDomain()
	}

	return users, nil
}



func (u *MongoUserRepo) GetByID(ctx context.Context, userID domain.ID) (domain.User, error) {
	result := u.db.FindOne(ctx, bson.M{"_id": userID})

	var mgUser entity.MgUser
	if err := result.Decode(&mgUser); err != nil {
		if err == mongo.ErrNoDocuments {
			return domain.User{}, errors.Wrap(domain.ErrNotExist, err.Error())
		}
		return domain.User{}, errors.Wrap(domain.ErrPersistenceFailed, err.Error())
	}
	return mgUser.ToDomain(), nil
}

func (u *MongoUserRepo) GetByEmail(ctx context.Context, email string) (domain.User, error) {
	result := u.db.FindOne(ctx, bson.M{"email":email})

	var mgUser entity.MgUser
	if err := result.Decode(&mgUser); err != nil {
		if err == mongo.ErrNoDocuments {
			return domain.User{}, errors.Wrap(domain.ErrNotExist, err.Error())
		}
		return domain.User{}, errors.Wrap(domain.ErrPersistenceFailed, err.Error())
	}
	return mgUser.ToDomain(), nil
}

func (u *MongoUserRepo) Create(ctx context.Context, user domain.User) (domain.User, error) {
	session, err := u.db.Database().Client().StartSession()
	if err != nil {
		return domain.User{}, nil
	}
	defer session.EndSession(ctx)

	err = mongo.WithSession(ctx, session, func(sessionContext mongo.SessionContext) error {
		var mgCart = entity.NewMgCart(domain.Cart{ID: user.CartID, Price: 0})
		_, err := u.db.Database().Collection(CartCollection).InsertOne(sessionContext, mgCart)
		if err != nil {
			if mongo.IsDuplicateKeyError(err) {
				return errors.Wrap(domain.ErrDuplicate, err.Error())
			}
			return errors.Wrap(domain.ErrPersistenceFailed, err.Error())
		}

		var mgUser = entity.NewMgUser(user)
		_, err = u.db.InsertOne(sessionContext, mgUser)
		if err != nil {
			if mongo.IsDuplicateKeyError(err) {
				return errors.Wrap(domain.ErrDuplicate, err.Error())
			}
			return errors.Wrap(domain.ErrPersistenceFailed, err.Error())
		}

		return nil
	})
	if err != nil {
		return domain.User{}, err
	}

	return u.GetByID(ctx, user.ID)
}

func (u *MongoUserRepo) Update(ctx context.Context, user domain.User) (domain.User, error) {
	var mgUser = entity.NewMgUser(user)
	_, err := u.db.ReplaceOne(ctx, bson.M{"_id": mgUser.ID}, mgUser)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return domain.User{}, errors.Wrap(domain.ErrNotExist, err.Error())
		}
		return domain.User{}, errors.Wrap(domain.ErrPersistenceFailed, err.Error())
	}

	return u.GetByID(ctx, user.ID)
}

func (u *MongoUserRepo) Delete(ctx context.Context, userID domain.ID) error {
	_, err := u.db.DeleteOne(ctx, bson.M{"_id": userID})
	if err != nil {
		return errors.Wrap(domain.ErrDeleteFailed, err.Error())
	}
	return nil
}
