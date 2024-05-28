package mongodb

import (
	"context"
	"github.com/EmirShimshir/marketplace-core/domain"
	"github.com/EmirShimshir/marketplace-repository/repository/mongodb"
	"github.com/EmirShimshir/marketplace-repository/repository/mongodb/entity"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/mongo"
	"testing"
)

var withdraws = []domain.Withdraw{
	domain.Withdraw{
		ID:      domain.ID("30e18bc1-4354-4937-9a3b-03cf0b702ad1"),
		ShopID:  domain.ID("30e18bc1-4354-4937-9a3b-03cf0b7027b1"),
		Comment: "comment",
		Sum:     9999,
		Status:  domain.WithdrawStatusDone,
	},
}

var createdWithdraw = domain.Withdraw{
	ID:      domain.ID("30e18bc1-4354-4937-9a3b-03cf0b702ad2"),
	ShopID:  domain.ID("30e18bc1-4354-4937-9a3b-03cf0b7027b1"),
	Comment: "comment new",
	Sum:     999,
	Status:  domain.WithdrawStatusStart,
}

var updatedWithdraw = domain.Withdraw{
	ID:      domain.ID("30e18bc1-4354-4937-9a3b-03cf0b702ad1"),
	ShopID:  domain.ID("30e18bc1-4354-4937-9a3b-03cf0b7027b1"),
	Comment: "comment 2",
	Sum:     99992,
	Status:  domain.WithdrawStatusDone,
}

func InitWithdrawsMongoDB(ctx context.Context, db *mongo.Database) error {
	for _, withdraw := range withdraws {
		var mgWithdraw = entity.NewMgWithdraw(withdraw)
		_, err := db.Collection(mongodb.WithdrawCollection).InsertOne(ctx, mgWithdraw)
		if err != nil {
			return err
		}
	}

	return nil
}

func TestWithdrawRepository(t *testing.T) {
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

	err = InitWithdrawsMongoDB(ctx, db)
	if err != nil {
		t.Fatal(err)
	}

	t.Run("test Get", func(t *testing.T) {
		repo := mongodb.NewWithdrawRepo(db)
		found, err := repo.Get(ctx, 2, 0)
		if err != nil {
			t.Errorf("failed get: %v", err)
		}
		require.Equal(t, len(withdraws), len(found))
		for i := range withdraws {
			require.Equal(t, withdraws[i], found[i])
		}
	})

	t.Run("test GetByID", func(t *testing.T) {
		repo := mongodb.NewWithdrawRepo(db)
		found, err := repo.GetByID(ctx, withdraws[0].ID)
		if err != nil {
			t.Errorf("failed get with id: %v", err)
		}
		require.Equal(t, found, withdraws[0])
	})

	t.Run("test GetByShopID", func(t *testing.T) {
		repo := mongodb.NewWithdrawRepo(db)
		found, err := repo.GetByShopID(ctx, withdraws[0].ShopID)
		if err != nil {
			t.Errorf("failed get with shop id: %v", err)
		}
		require.Equal(t, len(withdraws), len(found))
		for i := range withdraws {
			require.Equal(t, withdraws[i], found[i])
		}
	})

	t.Run("test create", func(t *testing.T) {
		repo := mongodb.NewWithdrawRepo(db)
		withdraw, err := repo.Create(ctx, createdWithdraw)
		if err != nil {
			t.Errorf("failed to create: %v", err)
		}
		require.Equal(t, withdraw, createdWithdraw)
	})

	t.Run("test update", func(t *testing.T) {
		repo := mongodb.NewWithdrawRepo(db)
		withdraw, err := repo.Update(ctx, updatedWithdraw)
		if err != nil {
			t.Errorf("failed to update: %v", err)
		}
		require.Equal(t, withdraw, updatedWithdraw)
	})

	t.Run("test delete user", func(t *testing.T) {
		repo := mongodb.NewWithdrawRepo(db)
		err = repo.Delete(ctx, withdraws[0].ID)
		if err != nil {
			t.Errorf("failed to delete withdraw: %v", err)
		}
	})
}
