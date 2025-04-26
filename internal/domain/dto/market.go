package dto

import (
	"mime/multipart"
)

type Jenis string
type Kegunaan string
type Composition string
type Berat string

const (
	KomposPadat Jenis = "kompos padat"
	KomposCair  Jenis = "kompos cair"
	Vermikompos Jenis = "vermi kompos"
	Bokashi     Jenis = "bokashi"
)

const (
	TanamanHerbal     Kegunaan = "tanaman herbal"
	SayurBuah         Kegunaan = "sayur dan buah"
	TanamanHias       Kegunaan = "tanaman hias"
	TanamanPerkebunan Kegunaan = "tanaman perkebunan"
)

const (
	Organik       Composition = "organik"
	NonGMO        Composition = "non-gmo"
	Probiotik     Composition = "probiotik"
	PestisidaFree Composition = "pestisida free"
)

const (
	BB1 Berat = "<1kg"
	BB2 Berat = "1kg - 5kg"
	BB3 Berat = "5kg - 10kg"
	BB4 Berat = ">10kg"
)

type CreateProduct struct {
	ProductId     string                `form:"product_id" binding:"required"`
	StoreId       string                `form:"store_id" binding:"required"`
	ProductName   string                `form:"product_name" binding:"required"`
	ProductPrice  uint64                `form:"product_price" binding:"required"`
	ProductWeight uint64                `form:"product_weight" binding:"required"`
	ProductType   Jenis                 `form:"product_type" binding:"required"`
	ProductUsage  Kegunaan              `form:"product_usage" binding:"required"`
	Composition   Composition           `form:"composition" binding:"required"`
	Description   string                `form:"description"`
	PhotoIMG      *multipart.FileHeader `form:"photo_img" binding:"required"`
	PhotoUrl      string                `form:"photo_url"`
}

type GetAllProduct struct {
	ProductId    string `json:"product_id"`
	StoreName    string `json:"store_name"`
	ProductName  string `json:"product_name"`
	ProductPrice uint64 `json:"product_price"`
	PhotoUrl     string `json:"photo_url"`
}

type ProductResponse struct {
	ProductId           string `json:"product_id"`
	StoreId             string `json:"store_id"`
	StoreName           string `json:"store_name"`
	ProductName         string `json:"product_name"`
	ProductPrice        uint64 `json:"product_price"`
	ProductWeight       uint64 `json:"product_weight"`
	ProductType         string `json:"product_type"`
	ProductWeightFilter string `json:"product_weight_filter"`
	ProductUsage        string `json:"product_usage"`
	Composition         string `json:"composition"`
	Description         string `json:"description"`
	PhotoUrl            string `json:"photo_url"`
}

type PatchProductRequest struct {
	ProductName   *string               `form:"product_name"`
	ProductPrice  *uint64               `form:"product_price"`
	ProductWeight *uint64               `form:"product_weight"`
	ProductType   *Jenis                `form:"product_type"`
	ProductUsage  *Kegunaan             `form:"product_usage"`
	Composition   *Composition          `form:"composition"`
	Description   *string               `form:"description"`
	PhotoIMG      *multipart.FileHeader `form:"photo_img"`
}
