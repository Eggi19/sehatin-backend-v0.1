package dtos

import (
	"github.com/tsanaativa/sehatin-backend-v0.1/entities"
	"github.com/shopspring/decimal"
)

type MostBoughtUserResponse struct {
	ProductId         int64           `json:"product_id"`
	PharmacyProductId int64           `json:"pharmacy_product_id"`
	ProductPicture    string          `json:"product_picture"`
	Name              string          `json:"name"`
	Price             decimal.Decimal `json:"price"`
	SellingUnit       string          `json:"selling_unit"`
	SlugId            string          `json:"slug_id"`
	Total             int             `json:"total,omitempty"`
}

type MostBoughtUserResponses struct {
	Paginations             PaginationResponse       `json:"pagination_info"`
	MostBoughtUserResponses []MostBoughtUserResponse `json:"most_boughts"`
}

func ConvertMostBoughtUserResponse(mb *entities.MostBoughtUser) *MostBoughtUserResponse {
	return &MostBoughtUserResponse{
		ProductId:         mb.PharmacyProduct.Product.Id,
		PharmacyProductId: mb.PharmacyProduct.Id,
		ProductPicture:    mb.PharmacyProduct.Product.ProductPicture,
		Name:              mb.PharmacyProduct.Product.Name,
		Price:             mb.PharmacyProduct.Price,
		SellingUnit:       mb.PharmacyProduct.Product.SellingUnit,
		SlugId:            mb.PharmacyProduct.Product.SlugId,
	}
}

func ConvertMostBoughtUserResponses(mbs []entities.MostBoughtUser, pagination entities.PaginationInfo) *MostBoughtUserResponses {
	mbResponses := []MostBoughtUserResponse{}

	for _, mb := range mbs {
		mbResponses = append(mbResponses, *ConvertMostBoughtUserResponse(&mb))
	}

	return &MostBoughtUserResponses{
		Paginations:             *ConvertToPaginationResponse(pagination),
		MostBoughtUserResponses: mbResponses,
	}
}
