package handlers

import (
	"fmt"
	"time"

	"github.com/Anvarsha-k/Product-Inventory-System/internal/db"
	"github.com/Anvarsha-k/Product-Inventory-System/internal/models"
	"github.com/Anvarsha-k/Product-Inventory-System/internal/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"

)

// CreateProductRequest structure
type CreateProductRequest struct {
	ProductID    int64  `json:"product_id"`
	ProductCode  string `json:"product_code"`
	ProductName  string `json:"product_name"`
	ProductImage string `json:"product_image"`

	Variants []struct {
		Name    string   `json:"name"`
		Options []string `json:"options"`
	} `json:"variants"`

	SubVariants []struct {
		OptionIDs []string `json:"option_ids"`
		SKU       string   `json:"sku"`
		Stock     string   `json:"stock"`
	} `json:"sub_variants"`
}

// CreateProduct - POST /api/products

func CreateProduct(c *fiber.Ctx) error {
	var req CreateProductRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.Error(c, fiber.StatusBadRequest, "VALIDATION_ERROR", err.Error(), nil)
	}

	if req.ProductID == 0 || req.ProductCode == "" || req.ProductName == "" {
		return utils.Error(c, fiber.StatusBadRequest, "VALIDATION_ERROR", "product_id, product_code, product_name required", nil)
	}

	if err := db.DB.Transaction(func(tx *gorm.DB) error {
		p := models.Product{
			ProductID:    req.ProductID,
			ProductCode:  req.ProductCode,
			ProductName:  req.ProductName,
			ProductImage: req.ProductImage,
			CreatedDate:  time.Now(),
			Active:       true,
		}
		if err := tx.Create(&p).Error; err != nil {
			return err
		}

		for _, v := range req.Variants {
			variant := models.Variant{ProductID: p.ID, Name: v.Name}
			if err := tx.Create(&variant).Error; err != nil {
				return err
			}
			for _, opt := range v.Options {
				vo := models.VariantOption{VariantID: variant.ID, Value: opt}
				if err := tx.Create(&vo).Error; err != nil {
					return err
				}
			}
		}

		totalStock := decimal.Zero
		for _, sv := range req.SubVariants {
			stockDec := decimal.Zero
			if sv.Stock != "" {
				d, err := decimal.NewFromString(sv.Stock)
				if err != nil {
					return fmt.Errorf("invalid stock for sku %s", sv.SKU)
				}
				stockDec = d
			}
			sub := models.SubVariant{
				ProductID: p.ID,
				OptionIDs: sv.OptionIDs,
				SKU:       sv.SKU,
				Stock:     stockDec,
			}
			if err := tx.Create(&sub).Error; err != nil {
				return err
			}
			if !stockDec.IsZero() {
				st := models.StockTransaction{
					ProductID:       p.ID,
					SubVariantID:    sub.ID,
					Quantity:        stockDec,
					TransactionType: "IN",
					TransactionDate: time.Now(),
				}
				if err := tx.Create(&st).Error; err != nil {
					return err
				}
				totalStock = totalStock.Add(stockDec)
			}
		}

		if err := tx.Model(&models.Product{}).Where("id = ?", p.ID).Update("total_stock", totalStock).Error; err != nil {
			return err
		}
		return nil
	}); err != nil {
		return utils.Error(c, fiber.StatusInternalServerError, "DB_ERROR", err.Error(), nil)
	}

	return utils.Success(c, fiber.StatusCreated, fiber.Map{"message": "created"})
}

