package repository

import (
	"hackfest-uc/internal/domain/entity"

	"gorm.io/gorm"
)

type MarketMySQLItf interface {
	CreateProduct(product entity.Market) (entity.Market, error)
	GetAllProducts() ([]entity.Market, error)
	GetProductByID(productID string) (*entity.Market, error)
	InitDummyStores() error
}

type MarketMySQL struct {
	db *gorm.DB
}

func NewMarketMySQL(db *gorm.DB) MarketMySQLItf {
	repo := &MarketMySQL{db}
	repo.InitDummyStores()
	return repo
}

func (m MarketMySQL) InitDummyStores() error {
	dummyStores := []entity.Store{
		{StoreId: "store1", StoreName: "Toko Pertanian Sejahtera"},
		{StoreId: "store2", StoreName: "Toko Organik Maju"},
		{StoreId: "store3", StoreName: "Toko Pupuk Alam"},
	}

	for _, store := range dummyStores {
		if err := m.db.FirstOrCreate(&store).Error; err != nil {
			return err
		}
	}
	return nil
}

func (m MarketMySQL) CreateProduct(product entity.Market) (entity.Market, error) {
	if err := m.db.Create(&product).Error; err != nil {
		return entity.Market{}, err
	}
	return product, nil
}

func (m MarketMySQL) GetAllProducts() ([]entity.Market, error) {
	var products []entity.Market
	if err := m.db.Preload("Store").Find(&products).Error; err != nil {
		return nil, err
	}
	return products, nil
}

func (m MarketMySQL) GetProductByID(productID string) (*entity.Market, error) {
	var product entity.Market
	if err := m.db.Preload("Store").Where("product_id = ?", productID).First(&product).Error; err != nil {
		return nil, err
	}
	return &product, nil
}
