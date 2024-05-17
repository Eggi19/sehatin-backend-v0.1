package dtos

import "github.com/tsanaativa/sehatin-backend-v0.1/entities"

type ProductFieldResponse struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
}

func ConvertToProductFieldResponse(form *entities.ProductForm, classification *entities.ProductClassification, manufacture *entities.Manufacture) *ProductFieldResponse {
	if form != nil {
		return &ProductFieldResponse{Id: form.Id, Name: form.Name}
	}
	if classification != nil {
		return &ProductFieldResponse{Id: classification.Id, Name: classification.Name}
	}
	if manufacture != nil {
		return &ProductFieldResponse{Id: manufacture.Id, Name: manufacture.Name}
	}
	return nil
}

func ConvertToProductFieldResponses(forms []entities.ProductForm, classifications []entities.ProductClassification, manufactures []entities.Manufacture) []ProductFieldResponse {
	productFieldResponses := []ProductFieldResponse{}

	if forms != nil {
		for _, form := range forms {
			productFieldResponses = append(productFieldResponses, *ConvertToProductFieldResponse(&form, nil, nil))
		}
		return productFieldResponses
	}

	if classifications != nil {
		for _, classification := range classifications {
			productFieldResponses = append(productFieldResponses, *ConvertToProductFieldResponse(nil, &classification, nil))
		}
		return productFieldResponses
	}

	if manufactures != nil {
		for _, manufacture := range manufactures {
			productFieldResponses = append(productFieldResponses, *ConvertToProductFieldResponse(nil, nil, &manufacture))
		}
		return productFieldResponses
	}
	return nil
}
