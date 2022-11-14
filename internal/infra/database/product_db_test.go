package database

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/caiocp/go-api/internal/entities"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestCreateNewProduct(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{})
	if err != nil {
		t.Error(err)
	}

	db.AutoMigrate(&entities.Product{})

	product, err := entities.NewProduct("Product 1", 10.0)
	assert.NoError(t, err)

	productDB := NewProduct(db)

	err = productDB.Create(product)
	assert.NoError(t, err)
	assert.NotEmpty(t, product.ID)
}

func TestFindAllProducts(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{})
	if err != nil {
		t.Error(err)
	}

	db.AutoMigrate(&entities.Product{})

	for i := 0; i < 13; i++ {
		product, err := entities.NewProduct(fmt.Sprintf("Product %d", i), rand.Float64()*100.0)
		assert.NoError(t, err)

		db.Create(product)
	}

	productDB := NewProduct(db)
	products, err := productDB.FindAll(1, 5, "asc")
	assert.NoError(t, err)
	assert.Len(t, products, 5)
	assert.Equal(t, "Product 0", products[0].Name)
	assert.Equal(t, "Product 4", products[4].Name)

	products, err = productDB.FindAll(2, 5, "asc")
	assert.NoError(t, err)
	assert.Len(t, products, 5)
	assert.Equal(t, "Product 5", products[0].Name)
	assert.Equal(t, "Product 9", products[4].Name)

	products, err = productDB.FindAll(3, 5, "asc")
	assert.NoError(t, err)
	assert.Len(t, products, 3)
	assert.Equal(t, "Product 10", products[0].Name)
	assert.Equal(t, "Product 12", products[2].Name)
}

func TestFindProductByID(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{})
	if err != nil {
		t.Error(err)
	}

	db.AutoMigrate(&entities.Product{})

	product, err := entities.NewProduct("Product 1", 10.0)
	assert.NoError(t, err)

	db.Create(product)

	productDB := NewProduct(db)
	product, err = productDB.FindByID(product.ID.String())
	assert.NoError(t, err)
	assert.Equal(t, "Product 1", product.Name)
	assert.Equal(t, 10.0, product.Price)
}

func TestUpdateProduct(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{})
	if err != nil {
		t.Error(err)
	}

	db.AutoMigrate(&entities.Product{})

	product, err := entities.NewProduct("Product 1", 10.0)
	assert.NoError(t, err)

	db.Create(product)

	productDB := NewProduct(db)
	product, err = productDB.FindByID(product.ID.String())
	assert.NoError(t, err)

	product.Name = "Product 2"
	product.Price = 20.0

	err = productDB.Update(product)
	assert.NoError(t, err)

	product, err = productDB.FindByID(product.ID.String())
	assert.NoError(t, err)
	assert.Equal(t, "Product 2", product.Name)
	assert.Equal(t, 20.0, product.Price)
}

func TestDeleteProduct(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{})
	if err != nil {
		t.Error(err)
	}

	db.AutoMigrate(&entities.Product{})

	product, err := entities.NewProduct("Product 1", 10.0)
	assert.NoError(t, err)

	db.Create(product)

	productDB := NewProduct(db)
	err = productDB.Delete(product.ID.String())
	assert.NoError(t, err)

	_, err = productDB.FindByID(product.ID.String())
	assert.Error(t, err)
}
