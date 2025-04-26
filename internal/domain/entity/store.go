package entity

type Store struct {
	StoreId   string `json:"store_id" gorm:"primaryKey;type:varchar(36)"`
	StoreName string `json:"store_name" gorm:"type:varchar(100);not null"`
}
