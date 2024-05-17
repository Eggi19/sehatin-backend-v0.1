package dtos

import "github.com/tsanaativa/sehatin-backend-v0.1/entities"

type PharmacyCreateRequest struct {
	Name                    string  `json:"name" binding:"required"`
	OperationalHour         string  `json:"operational_hour" binding:"required"`
	OperationalDay          string  `json:"operational" binding:"required"`
	PharmacistName          string  `json:"pharmacist_name" binding:"required"`
	PharmacistLicenseNumber string  `json:"pharmacist_license_number" binding:"required"`
	PharmacistPhoneNumber   string  `json:"pharmacist_phone_number" binding:"required"`
	City                    string  `json:"city" binding:"required"`
	Province                string  `json:"province" binding:"required"`
	Address                 string  `json:"address" binding:"required"`
	District                string  `json:"district" binding:"required"`
	SubDistrict             string  `json:"sub_district" binding:"required"`
	PostalCode              string  `json:"postal_code" binding:"required"`
	Longitude               string  `json:"longitude" binding:"required"`
	Latitude                string  `json:"latitude" binding:"required"`
	OfficialShippingId      []int64 `json:"official_shipping_id" binding:"required"`
	NonOfficialShippingId   []int64 `json:"non_official_shipping_id"`
}

type PharmacyUpdateRequest struct {
	Name                    string  `json:"name" binding:"required"`
	OperationalHour         string  `json:"operational_hour" binding:"required"`
	OperationalDay          string  `json:"operational" binding:"required"`
	PharmacistName          string  `json:"pharmacist_name" binding:"required"`
	PharmacistLicenseNumber string  `json:"pharmacist_license_number" binding:"required"`
	PharmacistPhoneNumber   string  `json:"pharmacist_phone_number" binding:"required"`
	City                    string  `json:"city" binding:"required"`
	Province                string  `json:"province" binding:"required"`
	Address                 string  `json:"address" binding:"required"`
	District                string  `json:"district" binding:"required"`
	SubDistrict             string  `json:"sub_district" binding:"required"`
	PostalCode              string  `json:"postal_code" binding:"required"`
	Longitude               string  `json:"longitude" binding:"required"`
	Latitude                string  `json:"latitude" binding:"required"`
	OfficialShippingId      []int64 `json:"official_shipping_id" binding:"required"`
	NonOfficialShippingId   []int64 `json:"non_official_shipping_id" binding:"required"`
}

type PharmacyResponse struct {
	Id                      int64                   `json:"id,omitempty"`
	PharmacyManager         PharmacyManagerResponse `json:"pharmacy_manager,omitempty"`
	Name                    string                  `json:"name,omitempty"`
	OperationalHour         string                  `json:"operational_hour,omitempty"`
	OperationalDay          string                  `json:"operational_day,omitempty"`
	PharmacistName          string                  `json:"pharmacist_name,omitempty"`
	PharmacistLicenseNumber string                  `json:"pharmacist_license_number,omitempty"`
	PharmacistPhoneNumber   string                  `json:"pharmacist_phone_number,omitempty"`
	Distance                float64                 `json:"distance,omitempty"`
	ShippingMethods         ShippingMethod          `json:"shipping_methods,omitempty"`
	PharmacyAddress         PharmacyAddressResponse `json:"pharmacy_address,omitempty"`
}

type PharmacyAddressResponse struct {
	Id          int64  `json:"id,omitempty"`
	PharmacyId  int64  `json:"pharmacy_id,omitempty"`
	City        string `json:"city,omitempty"`
	Province    string `json:"province,omitempty"`
	Address     string `json:"address,omitempty"`
	District    string `json:"district,omitempty"`
	SubDistrict string `json:"sub_district,omitempty"`
	PostalCode  string `json:"postal_code,omitempty"`
	Coordinate  string `json:"coordinate,omitempty"`
}

type PharmacyResponses struct {
	Pagination PaginationResponse `json:"pagination_info"`
	Pharmacies []PharmacyResponse `json:"pharmacies"`
}

func ConvertToPharmacyResponse(pharmacy *entities.Pharmacy) *PharmacyResponse {
	return &PharmacyResponse{
		Id:                      pharmacy.Id,
		PharmacyManager:         *ConvertToPharmacyManagerResponse(&pharmacy.PharmacyManager),
		Name:                    pharmacy.Name,
		OperationalHour:         pharmacy.OperationalHour,
		OperationalDay:          pharmacy.OperationalDay,
		PharmacistName:          pharmacy.PharmacistName,
		PharmacistLicenseNumber: pharmacy.PharmacistLicenseNumber,
		PharmacistPhoneNumber:   pharmacy.PharmacistPhoneNumber,
		PharmacyAddress: PharmacyAddressResponse{
			Id:          pharmacy.PharmacyAddress.Id,
			PharmacyId:  pharmacy.PharmacyAddress.PharmacyId,
			City:        pharmacy.PharmacyAddress.City,
			Province:    pharmacy.PharmacyAddress.Province,
			Address:     pharmacy.PharmacyAddress.Address,
			District:    pharmacy.PharmacyAddress.District,
			SubDistrict: pharmacy.PharmacyAddress.SubDistrict,
			PostalCode:  pharmacy.PharmacyAddress.PostalCode,
			Coordinate:  pharmacy.PharmacyAddress.Coordinate,
		},
		ShippingMethods: ConvertToShippingMethodDto(pharmacy.OfficialShippingMethod, pharmacy.NonOfficialShippingMethod),
	}
}

func ConvertToPharmacyResponses(pharmacies []entities.Pharmacy, pagination entities.PaginationInfo) *PharmacyResponses {
	pharmacyResponses := []PharmacyResponse{}

	for _, pharmacy := range pharmacies {
		pharmacyResponses = append(pharmacyResponses, *ConvertToPharmacyResponse(&pharmacy))
	}

	return &PharmacyResponses{
		Pagination: *ConvertToPaginationResponse(pagination),
		Pharmacies: pharmacyResponses,
	}
}
func ConvertToPharmacyAddressResponse(req entities.PharmacyAddress) *PharmacyAddressResponse {
	return &PharmacyAddressResponse{
		Id:          req.Id,
		PharmacyId:  req.PharmacyId,
		City:        req.City,
		Province:    req.Province,
		Address:     req.Address,
		District:    req.District,
		SubDistrict: req.SubDistrict,
		PostalCode:  req.PostalCode,
		Coordinate:  req.Coordinate,
	}
}
