package entity

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
	BB1 Berat = "< 1kg"
	BB2 Berat = "1kg - 5kg"
	BB3 Berat = "5kg - 10kg"
	BB4 Berat = "> 10kg"
)

type Market struct {
	ProductId           string      `json:"product_id" gorm:"primaryKey;type:varchar(36)"`
	ProductName         string      `json:"product_name"`
	StoreId             string      `json:"store_id" gorm:"type:varchar(36);not null;index"`
	ProductPrice        uint64      `json:"product_price"`
	ProductWeight       uint64      `json:"product_weight"`
	ProductType         Jenis       `json:"product_type"`
	ProductWeightFilter Berat       `json:"product_weight_filter"`
	ProductUsage        Kegunaan    `json:"product_usage"`
	Composition         Composition `json:"composition"`
	Description         string      `json:"description"`
	PhotoUrl            string      `json:"photo_url"`

	Store Store `gorm:"foreignKey:StoreId;references:StoreId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"` // relasi ke Store-nya tetap
}
