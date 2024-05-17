package dtos

import (
	"mime/multipart"

	"github.com/tsanaativa/sehatin-backend-v0.1/entities"
)

type ProductCreateRequest struct {
	Name                    string                `form:"name" binding:"required"`
	GenericName             string                `form:"generic_name" binding:"required"`
	Content                 string                `form:"content" binding:"required"`
	Description             string                `form:"description" binding:"required"`
	UnitInPack              string                `form:"unit_in_pack" binding:"required"`
	SellingUnit             string                `form:"selling_unit" binding:"required"`
	Weight                  int                   `form:"weight" binding:"required"`
	Height                  int                   `form:"height" binding:"required"`
	Length                  int                   `form:"length" binding:"required"`
	Width                   int                   `form:"width" binding:"required"`
	ProductPicture          *multipart.FileHeader `form:"product_picture"`
	SlugId                  string                `form:"slug_id" binding:"required"`
	ProductFormId           int64                 `form:"product_form_id" binding:"required"`
	ProductClassificationId int64                 `form:"product_classification_id" binding:"required"`
	ManufactureId           int64                 `form:"manufacture_id" binding:"required"`
	Categories              []int64               `form:"categories_id" binding:"required"`
}

type ProductCategoryResponse struct {
	Id                    int64              `json:"id"`
	Name                  string             `json:"name"`
	GenericName           string             `json:"generic_name,omitempty"`
	Content               string             `json:"content,"`
	Description           string             `json:"description"`
	UnitInPack            string             `json:"unit_in_pack"`
	SellingUnit           string             `json:"selling_unit"`
	Weight                int                `json:"weight"`
	Height                int                `json:"height"`
	Length                int                `json:"length"`
	Width                 int                `json:"width"`
	ProductPicture        string             `json:"product_picture"`
	SlugId                string             `json:"slug_id"`
	ProductForm           string             `json:"form"`
	ProductClassification string             `json:"classification"`
	Manufacture           string             `json:"manufacture"`
	Categories            []CategoryResponse `json:"categories"`
}

type CategoryIdRequest struct {
	Id int64 `form:"id" binding:"required"`
}

type ProductResponses struct {
	Pagination PaginationResponse        `json:"pagination_info"`
	Products   []ProductCategoryResponse `json:"products"`
}

func ConvertToProductResponse(product entities.Product) *ProductCategoryResponse {
	return &ProductCategoryResponse{
		Id:                    product.Id,
		Name:                  product.Name,
		GenericName:           product.GenericName,
		Content:               product.Content,
		Description:           product.Description,
		UnitInPack:            product.UnitInPack,
		SellingUnit:           product.SellingUnit,
		Weight:                product.Weight,
		Height:                product.Height,
		Length:                product.Length,
		Width:                 product.Width,
		ProductPicture:        product.ProductPicture,
		SlugId:                product.SlugId,
		ProductForm:           product.ProductForm.Name,
		ProductClassification: product.ProductClassification.Name,
		Manufacture:           product.Manufacture.Name,
		Categories:            ConvertToCategoryResponsesWithoutPagination(product.Categories),
	}
}

func ConvertToProductResponses(products []entities.Product, pagination entities.PaginationInfo) *ProductResponses {
	productResponses := []ProductCategoryResponse{}

	for _, product := range products {
		productResponses = append(productResponses, *ConvertToProductResponse(product))
	}

	return &ProductResponses{Products: productResponses, Pagination: *ConvertToPaginationResponse(pagination)}
}
