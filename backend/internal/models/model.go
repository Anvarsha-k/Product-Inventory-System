package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type Product struct {
	ID           uuid.UUID  `gorm:"type:uuid;primaryKey" json:"id"`
	ProductID    int64      `gorm:"unique;not null" json:"product_id"`
	ProductCode  string     `gorm:"unique;not null" json:"product_code"`
	ProductName  string     `json:"product_name"`
	ProductImage string     `json:"product_image"`
	CreatedDate  time.Time  `gorm:"autoCreateTime" json:"created_date"`
	UpdatedDate  *time.Time `json:"updated_date"`
	CreatedUser  uuid.UUID  `gorm:"type:uuid" json:"created_user"`
	IsFavourite  bool       `json:"is_favourite"`
	Active       bool       `json:"active"`
	HSNCode      string     `json:"hsn_code"`
	TotalStock  decimal.Decimal `gorm:"type:numeric(20,8);default:0" json:"total_stock"`

	Variants    []Variant    `gorm:"foreignKey:ProductID" json:"variants,omitempty"`
	SubVariants []SubVariant `gorm:"foreignKey:ProductID" json:"sub_variants,omitempty"`
}

func (p *Product) BeforeCreate(tx *gorm.DB) (err error) {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	return
}

type Variant struct {
	ID        uuid.UUID       `gorm:"type:uuid;primaryKey" json:"id"`
	ProductID uuid.UUID       `gorm:"type:uuid" json:"product_id"`
	Name      string          `json:"name"`
	Options   []VariantOption `gorm:"foreignKey:VariantID" json:"options,omitempty"`
}

func (v *Variant) BeforeCreate(tx *gorm.DB) (err error) {
	if v.ID == uuid.Nil {
		v.ID = uuid.New()
	}
	return
}

type VariantOption struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	VariantID uuid.UUID `gorm:"type:uuid" json:"variant_id"`
	Value     string    `json:"value"`
}

func (vo *VariantOption) BeforeCreate(tx *gorm.DB) (err error) {
	if vo.ID == uuid.Nil {
		vo.ID = uuid.New()
	}
	return
}

type SubVariant struct {
	ID        uuid.UUID       `gorm:"type:uuid;primaryKey" json:"id"`
	ProductID uuid.UUID       `gorm:"type:uuid" json:"product_id"`
	OptionIDs pq.StringArray  `gorm:"type:text[]" json:"option_ids"`
	SKU       string          `gorm:"unique;not null" json:"sku"`
	Stock     decimal.Decimal `gorm:"type:numeric(20,8);default:0" json:"stock"`
}

func (sv *SubVariant) BeforeCreate(tx *gorm.DB) (err error) {
	if sv.ID == uuid.Nil {
		sv.ID = uuid.New()
	}
	return
}

type StockTransaction struct {
	ID              uuid.UUID       `gorm:"type:uuid;primaryKey" json:"id"`
	ProductID       uuid.UUID       `gorm:"type:uuid" json:"product_id"`
	SubVariantID    uuid.UUID       `gorm:"type:uuid" json:"sub_variant_id"`
	Quantity        decimal.Decimal `gorm:"type:numeric(20,8)" json:"quantity"`
	TransactionType string          `json:"transaction_type"`
	TransactionDate time.Time       `gorm:"autoCreateTime" json:"transaction_date"`
}

func (st *StockTransaction) BeforeCreate(tx *gorm.DB) (err error) {
	if st.ID == uuid.Nil {
		st.ID = uuid.New()
	}
	return
}
