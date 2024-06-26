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

var products = []domain.Product{
	domain.Product{
		ID:          domain.ID("30e18bc1-4354-4937-9a3b-03cf0b7027a1"),
		Name:        "iphone 15",
		Description: "apple IOS",
		Price:       129990,
		Category:    domain.ElectronicCategory,
		PhotoUrl:    "photo/1.png",
	},
	domain.Product{
		ID:          domain.ID("30e18bc1-4354-4937-9a3b-03cf0b7027a2"),
		Name:        "harry potter",
		Description: "Rouling",
		Price:       2990,
		Category:    domain.BooksCategory,
		PhotoUrl:    "photo/2.png",
	},
}

var createdProduct = domain.Product{
	ID:          domain.ID("30e18bc1-4354-4937-9a3b-03cf0b7027a3"),
	Name:        "new",
	Description: "new",
	Price:       129990,
	Category:    domain.ElectronicCategory,
	PhotoUrl:    "photo/new.png",
}

var updatedProduct = domain.Product{
	ID:          domain.ID("30e18bc1-4354-4937-9a3b-03cf0b7027a1"),
	Name:        "iphone 15",
	Description: "apple IOS 17",
	Price:       129990,
	Category:    domain.ElectronicCategory,
	PhotoUrl:    "photo/1.png",
}

func InitProductsMongoDB(ctx context.Context, db *mongo.Database) error {
	for _, product := range products {
		var mgProduct = entity.NewMgProduct(product)
		_, err := db.Collection(mongodb.ProductCollection).InsertOne(ctx, mgProduct)
		if err != nil {
			return err
		}
	}

	return nil
}

func TestProductRepository(t *testing.T) {
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

	err = InitProductsMongoDB(ctx, db)
	if err != nil {
		t.Fatal(err)
	}

	t.Run("test get products", func(t *testing.T) {
		repo := mongodb.NewProductRepo(db)
		found, err := repo.Get(ctx, 2, 0)
		if err != nil {
			t.Errorf("failed to get products: %v", err)
		}
		require.Equal(t, len(products), len(found))
		for i := range users {
			require.Equal(t, products[i], found[i])
		}
	})

	t.Run("test get product by id", func(t *testing.T) {
		repo := mongodb.NewProductRepo(db)
		product, err := repo.GetByID(ctx, products[0].ID)
		if err != nil {
			t.Errorf("failed to get product with id: %v", err)
		}
		require.Equal(t, product, products[0])
	})

	t.Run("test create product", func(t *testing.T) {
		repo := mongodb.NewProductRepo(db)
		product, err := repo.Create(ctx, createdProduct)
		if err != nil {
			t.Errorf("failed to create product: %v", err)
		}
		require.Equal(t, product, createdProduct)
	})

	t.Run("test update product", func(t *testing.T) {
		repo := mongodb.NewProductRepo(db)
		product, err := repo.Update(ctx, updatedProduct)
		if err != nil {
			t.Errorf("failed to create product: %v", err)
		}
		require.Equal(t, product, updatedProduct)
	})

	t.Run("test delete product", func(t *testing.T) {
		repo := mongodb.NewProductRepo(db)
		err = repo.Delete(ctx, products[0].ID)
		if err != nil {
			t.Errorf("failed to delete product: %v", err)
		}
	})
}
