package dtos

import (
	"github.com/tsanaativa/sehatin-backend-v0.1/entities"
	"github.com/shopspring/decimal"
)

type ProductResponse struct {
	ProductId         int64           `json:"product_id"`
	PharmacyProductId int64           `json:"pharmacy_product_id"`
	ProductPicture    string          `json:"product_picture"`
	Name              string          `json:"name"`
	Price             decimal.Decimal `json:"price"`
	SellingUnit       string          `json:"selling_unit"`
	SlugId            string          `json:"slug_id"`
	Total             int             `json:"total,omitempty"`
	Day               string          `json:"day,omitempty"`
	QuantitySold      int             `json:"quantity_sold,omitempty"`
}

type GetProductResponse struct {
	PaginationInfo PaginationResponse `json:"pagination_info"`
	Products       []ProductResponse  `json:"products"`
}

type ProductDetail struct {
	Id                    int64              `json:"id"`
	Name                  string             `json:"name"`
	GenericName           string             `json:"generic_name"`
	Content               string             `json:"content"`
	Description           string             `json:"description"`
	UnitInPack            string             `json:"unit_in_pack"`
	SellingUnit           string             `json:"selling_unit"`
	Weight                int                `json:"weight"`
	Height                int                `json:"height"`
	Length                int                `json:"length"`
	Width                 int                `json:"width"`
	ProductPicture        string             `json:"product_picture"`
	SlugId                string             `json:"slug_id"`
	ProductForm           string             `json:"product_form"`
	ProductClassification string             `json:"product_classification"`
	Manufacture           string             `json:"manufacture"`
	Category              []CategoryResponse `json:"categories"`
	Pharmacy              []PharmacyResponse `json:"pharmacies"`
	Price                 decimal.Decimal    `json:"price"`
	TotalStock            int                `json:"total_stock"`
}

type PharmacyProductRequest struct {
	Price       decimal.Decimal `json:"price" binding:"required"`
	TotalStock  *int            `json:"total_stock" binding:"required"`
	IsAvailable *bool           `json:"is_available" binding:"required"`
	ProductId   int64           `json:"product_id" binding:"required"`
	PharmacyId  int64           `json:"pharmacy_id" binding:"required"`
}

type ProductSummaryResponse struct {
	Id          int64              `json:"id"`
	Name        string             `json:"name"`
	SellingUnit string             `json:"selling_unit"`
	SlugId      string             `json:"slug_id"`
	Category    []CategoryResponse `json:"categories"`
}

type PharmacyProductResponse struct {
	Id          int64                   `json:"id"`
	Price       decimal.Decimal         `json:"price"`
	TotalStock  int                     `json:"total_stock"`
	IsAvailable bool                    `json:"is_available"`
	SlugId      string                  `json:"slug_id"`
	Product     ProductCategoryResponse `json:"product,omitempty"`
}

type PharmacyProductResponses struct {
	Pagination       PaginationResponse        `json:"pagination_info"`
	PharmacyProducts []PharmacyProductResponse `json:"pharmacy_products"`
}

type PharmacyProductItem struct {
	Id          int64                   `json:"id"`
	TotalStock  int                     `json:"total_stock"`
	IsAvailable bool                    `json:"is_available"`
	Price       decimal.Decimal         `json:"price"`
	Product     ProductCategoryResponse `json:"product"`
}

func ConvertToPharmacyProductResponse(pharmacyProduct entities.PharmacyProduct) *PharmacyProductResponse {
	return &PharmacyProductResponse{
		Id:          pharmacyProduct.Id,
		Price:       pharmacyProduct.Price,
		TotalStock:  pharmacyProduct.TotalStock,
		IsAvailable: pharmacyProduct.IsAvailable,
		Product: ProductCategoryResponse{
			Id:                    pharmacyProduct.Product.Id,
			Name:                  pharmacyProduct.Product.Name,
			GenericName:           pharmacyProduct.Product.GenericName,
			Content:               pharmacyProduct.Product.Content,
			Description:           pharmacyProduct.Product.Description,
			UnitInPack:            pharmacyProduct.Product.UnitInPack,
			SellingUnit:           pharmacyProduct.Product.SellingUnit,
			Weight:                pharmacyProduct.Product.Weight,
			Height:                pharmacyProduct.Product.Height,
			Length:                pharmacyProduct.Product.Length,
			Width:                 pharmacyProduct.Product.Width,
			ProductPicture:        pharmacyProduct.Product.ProductPicture,
			SlugId:                pharmacyProduct.Product.SlugId,
			ProductForm:           pharmacyProduct.Product.ProductForm.Name,
			ProductClassification: pharmacyProduct.Product.ProductClassification.Name,
			Manufacture:           pharmacyProduct.Product.Manufacture.Name,
			Categories:            ConvertToCategoryResponsesWithoutPagination(pharmacyProduct.Product.Categories),
		},
	}
}

func ConvertToPharmacyProductResponses(pharmacyProducts []entities.PharmacyProduct, pagination entities.PaginationInfo) *PharmacyProductResponses {
	pharmacyProductResponses := []PharmacyProductResponse{}

	for _, pharmacyProduct := range pharmacyProducts {
		pharmacyProductResponses = append(pharmacyProductResponses, *ConvertToPharmacyProductResponse(pharmacyProduct))
	}

	return &PharmacyProductResponses{PharmacyProducts: pharmacyProductResponses, Pagination: *ConvertToPaginationResponse(pagination)}
}

func ConvertToPharmacyProductItem(pp entities.PharmacyProduct) *PharmacyProductItem {
	return &PharmacyProductItem{
		Id:          pp.Id,
		TotalStock:  pp.TotalStock,
		IsAvailable: pp.IsAvailable,
		Price:       pp.Price,
		Product: ProductCategoryResponse{
			Id:                    pp.Product.Id,
			Name:                  pp.Product.Name,
			GenericName:           pp.Product.GenericName,
			Content:               pp.Product.Content,
			Description:           pp.Product.Description,
			UnitInPack:            pp.Product.UnitInPack,
			SellingUnit:           pp.Product.SellingUnit,
			Weight:                pp.Product.Weight,
			Height:                pp.Product.Height,
			Length:                pp.Product.Length,
			Width:                 pp.Product.Width,
			ProductPicture:        pp.Product.ProductPicture,
			SlugId:                pp.Product.SlugId,
			ProductForm:           pp.Product.ProductForm.Name,
			ProductClassification: pp.Product.ProductClassification.Name,
			Manufacture:           pp.Product.Manufacture.Name,
			Categories:            ConvertToCategoryResponsesWithoutPagination(pp.Product.Categories),
		},
	}
}
