package dtos

import (
	"time"

	"github.com/tsanaativa/sehatin-backend-v0.1/entities"
	"github.com/shopspring/decimal"
)

type CreateCartItemRequest struct {
	Quantity          int   `json:"quantity" binding:"required"`
	PharmacyProductId int64 `json:"pharmacy_product_id" binding:"required"`
}

type UpdateCartItemRequest struct {
	Quantity int `json:"quantity" binding:"required"`
}

type CartItemResponse struct {
	PharmacyId      int64           `json:"pharmacy_id"`
	PharmacyName    string          `json:"pharmacy_name"`
	Id              int64           `json:"id"`
	ProductName     string          `json:"product_name"`
	ProductPicture  string          `json:"product_picture"`
	SellingUnit     string          `json:"selling_unit"`
	Price           decimal.Decimal `json:"price"`
	Quantity        int             `json:"quantity"`
	UpdatedAt       time.Time       `json:"updated_at"`
	ShippingMethods ShippingMethod  `json:"shipping_methods"`
	SlugId          string          `json:"slug_id"`
	TotalStock      int             `json:"total_stock"`
	IsAvailable     bool            `json:"is_available"`
	Weight          int             `json:"weight"`
}

func ConvertToCartItemResponse(req entities.CartItem) *CartItemResponse {
	return &CartItemResponse{
		PharmacyId:     req.PharmacyId,
		PharmacyName:   req.PharmacyName,
		Id:             req.Id,
		ProductName:    req.ProductName,
		ProductPicture: req.ProductPicture,
		SellingUnit:    req.SellingUnit,
		Price:          req.Price,
		Quantity:       req.Quantity,
		UpdatedAt:      req.UpdatedAt,
		SlugId:         req.SlugId,
		TotalStock:     req.TotalStock,
		IsAvailable:    req.IsAvailable,
		Weight:         req.Weight,
	}
}

func ConvertToCartItemResponses(req []entities.CartItem) []CartItemResponse {
	result := []CartItemResponse{}

	for i := 0; i < len(req); i++ {
		item := ConvertToCartItemResponse(req[i])

		result = append(result, *item)
	}
	return result
}
