package handlers

import (
	"fmt"
	"strconv"
	"time"

	"github.com/Anvarsha-k/Product-Inventory-System/internal/db"
	"github.com/Anvarsha-k/Product-Inventory-System/internal/models"
	"github.com/Anvarsha-k/Product-Inventory-System/internal/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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

// ListProducts with Pagiantion logic included

func ListProducts(c *fiber.Ctx) error {
	pageStr := c.Query("page", "1")
	limitStr := c.Query("limit", "10")
	page, _ := strconv.Atoi(pageStr)
	limit, _ := strconv.Atoi(limitStr)
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}
	offset := (page - 1) * limit

	var prods []models.Product
	if err := db.DB.Preload("Variants.Options").Preload("SubVariants").
		Limit(limit).Offset(offset).Find(&prods).Error; err != nil {
		return utils.Error(c, fiber.StatusInternalServerError, "DB_ERROR", err.Error(), nil)
	}
	var total int64
	db.DB.Model(&models.Product{}).Count(&total)

	return utils.Success(c, fiber.StatusOK, fiber.Map{
		"products": prods,
		"page":     page,
		"limit":    limit,
		"total":    total,
	})
}

// StockAdjustRequest

type StockAdjustRequest struct {
	SubVariantID    string `json:"sub_variant_id"`
	Quantity        string `json:"quantity"`
	TransactionType string `json:"transaction_type"`
}

// AdjustStock

func AdjustStock(c *fiber.Ctx) error {
	var req StockAdjustRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.Error(c, fiber.StatusBadRequest, "VALIDATION_ERROR", err.Error(), nil)
	}
	if req.SubVariantID == "" || req.Quantity == "" || req.TransactionType == "" {
		return utils.Error(c, fiber.StatusBadRequest, "VALIDATION_ERROR", "sub_variant_id, quantity, transaction_type required", nil)
	}

	qty, err := decimal.NewFromString(req.Quantity)
	if err != nil || qty.IsNegative() || qty.IsZero() {
		return utils.Error(c, fiber.StatusBadRequest, "INVALID_QUANTITY", "quantity must be positive decimal", nil)
	}

	err = db.DB.Transaction(func(tx *gorm.DB) error {
		var sub models.SubVariant
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("id = ?", req.SubVariantID).First(&sub).Error; err != nil {
			return err
		}
		var product models.Product
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("id = ?", sub.ProductID).First(&product).Error; err != nil {
			return err
		}

		if req.TransactionType == "IN" {
			sub.Stock = sub.Stock.Add(qty)
			product.TotalStock = product.TotalStock.Add(qty)
		} else if req.TransactionType == "OUT" {
			if sub.Stock.Sub(qty).IsNegative() {
				return fmt.Errorf("insufficient subvariant stock")
			}
			sub.Stock = sub.Stock.Sub(qty)
			product.TotalStock = product.TotalStock.Sub(qty)
			if product.TotalStock.IsNegative() {
				return fmt.Errorf("product total would become negative")
			}
		} else {
			return fmt.Errorf("invalid transaction type")
		}

		if err := tx.Save(&sub).Error; err != nil {
			return err
		}
		if err := tx.Model(&models.Product{}).Where("id = ?", product.ID).Update("total_stock", product.TotalStock).Error; err != nil {
			return err
		}
		st := models.StockTransaction{
			ProductID:       product.ID,
			SubVariantID:    sub.ID,
			Quantity:        qty,
			TransactionType: req.TransactionType,
			TransactionDate: time.Now(),
		}
		if err := tx.Create(&st).Error; err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return utils.Error(c, fiber.StatusNotFound, "NOT_FOUND", err.Error(), nil)
		}
		if err.Error() == "insufficient subvariant stock" || err.Error() == "product total would become negative" {
			return utils.Error(c, fiber.StatusBadRequest, "INSUFFICIENT_STOCK", err.Error(), nil)
		}
		return utils.Error(c, fiber.StatusInternalServerError, "DB_ERROR", err.Error(), nil)
	}

	return utils.Success(c, fiber.StatusOK, fiber.Map{"message": "stock adjusted"})
}

// StockReport - GET /api/stock/report?from=...&to=...&product_id=...

func StockReport(c *fiber.Ctx) error {
	fromStr := c.Query("from")
	toStr := c.Query("to")
	productID := c.Query("product_id")

	var from time.Time
	var to time.Time
	var err error
	if fromStr != "" {
		from, err = time.Parse(time.RFC3339, fromStr)
		if err != nil {
			return utils.Error(c, fiber.StatusBadRequest, "INVALID_DATE", "invalid from date", nil)
		}
	}
	if toStr != "" {
		to, err = time.Parse(time.RFC3339, toStr)
		if err != nil {
			return utils.Error(c, fiber.StatusBadRequest, "INVALID_DATE", "invalid to date", nil)
		}
	}

	q := db.DB.Model(&models.StockTransaction{})
	if !from.IsZero() {
		q = q.Where("transaction_date >= ?", from)
	}
	if !to.IsZero() {
		q = q.Where("transaction_date <= ?", to)
	}
	if productID != "" {
		q = q.Where("product_id = ?", productID)
	}

	var txs []models.StockTransaction
	if err := q.Order("transaction_date desc").Find(&txs).Error; err != nil {
		return utils.Error(c, fiber.StatusInternalServerError, "DB_ERROR", err.Error(), nil)
	}

	type Summary struct {
		ProductID string          `json:"product_id"`
		StockIn   decimal.Decimal `json:"stock_in"`
		StockOut  decimal.Decimal `json:"stock_out"`
		Net       decimal.Decimal `json:"net"`
	}
	summaries := map[string]*Summary{}
	for _, t := range txs {
		pid := t.ProductID.String()
		if _, ok := summaries[pid]; !ok {
			summaries[pid] = &Summary{ProductID: pid, StockIn: decimal.Zero, StockOut: decimal.Zero, Net: decimal.Zero}
		}
		if t.TransactionType == "IN" {
			summaries[pid].StockIn = summaries[pid].StockIn.Add(t.Quantity)
			summaries[pid].Net = summaries[pid].Net.Add(t.Quantity)
		} else {
			summaries[pid].StockOut = summaries[pid].StockOut.Add(t.Quantity)
			summaries[pid].Net = summaries[pid].Net.Sub(t.Quantity)
		}
	}

	out := make([]Summary, 0, len(summaries))
	for _, v := range summaries {
		out = append(out, *v)
	}

	return utils.Success(c, fiber.StatusOK, fiber.Map{
		"transactions": txs,
		"summary":      out,
	})
}