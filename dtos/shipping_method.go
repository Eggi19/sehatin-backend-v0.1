package dtos

import "github.com/tsanaativa/sehatin-backend-v0.1/entities"

type OfficialShippingMethod struct {
	Id   int64  `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

type NonOfficialShippingMethod struct {
	Id          int64  `json:"id,omitempty"`
	Name        string `json:"name,omitempty"`
	Courier     string `json:"courier,omitempty"`
	Service     string `json:"service,omitempty"`
	Description string `json:"description,omitempty"`
}

type ShippingMethod struct {
	Official    []OfficialShippingMethod    `json:"official"`
	NonOfficial []NonOfficialShippingMethod `json:"non_official"`
}

type OfficialShippingFeeRequest struct {
	UserAddressId            int64 `json:"user_address_id" binding:"required"`
	PharmacyId               int64 `json:"pharmacy_id" binding:"required"`
	OfficialShippingMethodId int64 `json:"official_shipping_method_id" binding:"required"`
}

type NonOfficialShippingFeeRequest struct {
	UserAddressId               int64   `json:"user_address_id" binding:"required"`
	PharmacyId                  int64   `json:"pharmacy_id" binding:"required"`
	TotalWeight                 float64 `json:"total_weight" binding:"required"`
	NonOfficialShippingMethodId int64   `json:"non_official_shipping_method_id" binding:"required"`
}

type ShippingCostResponse struct {
	Cost float64 `json:"cost"`
}

func ConvertToShippingMethodDto(officials []entities.OfficialShippingMethod, nonOfficials []entities.NonOfficialShippingMethod) ShippingMethod {
	officialMap := []OfficialShippingMethod{}
	nonOfficialMap := []NonOfficialShippingMethod{}

	for _, shipping := range officials {
		var official OfficialShippingMethod
		official.Id = shipping.Id
		official.Name = shipping.Name
		officialMap = append(officialMap, official)
	}

	for _, shipping := range nonOfficials {
		var nonOfficial NonOfficialShippingMethod
		nonOfficial.Id = shipping.Id
		nonOfficial.Name = shipping.Name
		nonOfficial.Courier = shipping.Courier
		nonOfficial.Service = shipping.Service
		nonOfficial.Description = shipping.Description
		nonOfficialMap = append(nonOfficialMap, nonOfficial)
	}

	return ShippingMethod{
		Official:    officialMap,
		NonOfficial: nonOfficialMap,
	}
}
