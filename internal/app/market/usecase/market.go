package usecase

import (
	"errors"
	"fmt"
	"hackfest-uc/internal/app/market/repository"
	"hackfest-uc/internal/domain/dto"
	"hackfest-uc/internal/domain/entity"
	"hackfest-uc/internal/infra/supabase"
	"path/filepath"

	"github.com/google/uuid"
)

type MarketUsecaseItf interface {
	CreateProduct(param dto.CreateProduct) (dto.ProductResponse, error)
	GetAllProducts(page, size int) ([]dto.ProductResponse, error)
	GetProductByID(productID string) (*dto.ProductResponse, error)
}

type MarketUsecase struct {
	marketRepo repository.MarketMySQLItf
	sb         supabase.SupabaseItf
}

func NewMarketUsecase(marketRepo repository.MarketMySQLItf, sb supabase.SupabaseItf) MarketUsecaseItf {
	return &MarketUsecase{
		marketRepo: marketRepo,
		sb:         sb,
	}
}

func (m *MarketUsecase) CreateProduct(param dto.CreateProduct) (dto.ProductResponse, error) {
	productId := uuid.New().String()

	// Validate file
	ext := filepath.Ext(param.PhotoIMG.Filename)
	if ext != ".jpg" && ext != ".jpeg" && ext != ".png" {
		return dto.ProductResponse{}, errors.New("invalid file type, only jpg, jpeg, png are allowed")
	}
	param.PhotoIMG.Filename = fmt.Sprintf("%v%v", productId, ext)

	// Upload photo
	newPhotoLink, err := m.sb.Upload(param.PhotoIMG)
	if err != nil {
		return dto.ProductResponse{}, err
	}

	// Determine weight filter
	var weightFilter entity.Berat
	switch {
	case param.ProductWeight < 1:
		weightFilter = entity.BB1
	case param.ProductWeight >= 1 && param.ProductWeight < 5:
		weightFilter = entity.BB2
	case param.ProductWeight >= 5 && param.ProductWeight <= 10:
		weightFilter = entity.BB3
	case param.ProductWeight > 10:
		weightFilter = entity.BB4
	}

	// Create product

	marketPost := entity.Market{
		ProductId:           productId,
		StoreId:             param.StoreId,
		ProductName:         param.ProductName,
		ProductPrice:        param.ProductPrice,
		ProductWeight:       param.ProductWeight,
		ProductWeightFilter: weightFilter,
		ProductType:         entity.Jenis(param.ProductType),
		ProductUsage:        entity.Kegunaan(param.ProductUsage),
		Composition:         entity.Composition(param.Composition),
		Description:         param.Description,
		PhotoUrl:            newPhotoLink,
	}

	createdProduct, err := m.marketRepo.CreateProduct(marketPost)
	if err != nil {
		return dto.ProductResponse{}, err
	}

	return dto.ProductResponse{
		ProductId:           createdProduct.ProductId,
		StoreId:             createdProduct.Store.StoreId,
		StoreName:           createdProduct.Store.StoreName,
		ProductName:         createdProduct.ProductName,
		ProductPrice:        createdProduct.ProductPrice,
		ProductWeight:       createdProduct.ProductWeight,
		ProductType:         string(createdProduct.ProductType),
		ProductWeightFilter: string(createdProduct.ProductWeightFilter),
		ProductUsage:        string(createdProduct.ProductUsage),
		Composition:         string(createdProduct.Composition),
		Description:         createdProduct.Description,
		PhotoUrl:            createdProduct.PhotoUrl,
	}, nil
}

func (m *MarketUsecase) GetAllProducts(page, size int) ([]dto.ProductResponse, error) {
	products, err := m.marketRepo.GetAllProducts(page, size)
	if err != nil {
		return nil, fmt.Errorf("failed to get products: %w", err)
	}

	var responses []dto.ProductResponse
	for _, product := range products {
		responses = append(responses, m.convertToProductResponse(&product))
	}

	return responses, nil
}

func (m *MarketUsecase) GetProductByID(productID string) (*dto.ProductResponse, error) {
	product, err := m.marketRepo.GetProductByID(productID)
	if err != nil {
		return nil, fmt.Errorf("product not found: %w", err)
	}

	response := m.convertToProductResponse(product)
	return &response, nil
}

func (m *MarketUsecase) convertToProductResponse(product *entity.Market) dto.ProductResponse {
	return dto.ProductResponse{
		ProductId:           product.ProductId,
		StoreId:             product.Store.StoreId,
		StoreName:           product.Store.StoreName,
		ProductName:         product.ProductName,
		ProductPrice:        product.ProductPrice,
		ProductWeight:       product.ProductWeight,
		ProductType:         string(product.ProductType),
		ProductWeightFilter: string(product.ProductWeightFilter),
		ProductUsage:        string(product.ProductUsage),
		Composition:         string(product.Composition),
		Description:         product.Description,
		PhotoUrl:            product.PhotoUrl,
	}
}
