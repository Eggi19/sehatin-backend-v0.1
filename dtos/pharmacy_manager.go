package dtos

import (
	"mime/multipart"

	"github.com/tsanaativa/sehatin-backend-v0.1/entities"
)

type PharmacyManagerResponse struct {
	Id          int64  `json:"id,omitempty"`
	Name        string `json:"name,omitempty"`
	Email       string `json:"email,omitempty"`
	PhoneNumber string `json:"phone_number,omitempty"`
	Logo        string `json:"logo,omitempty"`
}

type PharmacyMangerUpdateRequest struct {
	Name        string                `form:"name" binding:"required"`
	PhoneNumber string                `form:"phone_number" binding:"required"`
	Logo        *multipart.FileHeader `form:"logo"`
}

type PharmacyManagerResponses struct {
	Pagination       PaginationResponse        `json:"pagination_info"`
	PharmacyManagers []PharmacyManagerResponse `json:"pharmacy_managers"`
}

func ConvertToPharmacyManagerResponse(pharmacyManager *entities.PharmacyManager) *PharmacyManagerResponse {
	return &PharmacyManagerResponse{
		Id:          pharmacyManager.Id,
		Name:        pharmacyManager.Name,
		Email:       pharmacyManager.Email,
		PhoneNumber: pharmacyManager.PhoneNumber,
		Logo:        pharmacyManager.Logo,
	}
}

func ConvertToPharmacyManagerResponses(pharmacyManager []entities.PharmacyManager, pagination entities.PaginationInfo) *PharmacyManagerResponses {
	pharmacyManagersResponses := []PharmacyManagerResponse{}

	for _, pm := range pharmacyManager {
		pharmacyManagersResponses = append(pharmacyManagersResponses, *ConvertToPharmacyManagerResponse(&pm))
	}

	return &PharmacyManagerResponses{
		Pagination:       *ConvertToPaginationResponse(pagination),
		PharmacyManagers: pharmacyManagersResponses,
	}
}
