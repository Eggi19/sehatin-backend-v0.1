package dtos

import (
	"mime/multipart"

	"github.com/tsanaativa/sehatin-backend-v0.1/entities"
)

type DoctorUpdateRequest struct {
	Name           string                `form:"name" binding:"required"`
	Fee            int                   `form:"fee" binding:"required"`
	Certificate    *multipart.FileHeader `form:"certificate"`
	WorkStartYear  int                   `form:"work_start_year" binding:"required"`
	SpecialistId   int64                 `form:"specialist_id" binding:"required"`
	ProfilePicture *multipart.FileHeader `form:"profile_picture"`
}

type DoctorResponse struct {
	Id             int64               `json:"id"`
	Name           string              `json:"name"`
	Email          string              `json:"email"`
	Certificate    *string             `json:"certificate"`
	IsOnline       bool                `json:"is_online"`
	IsVerified     bool                `json:"is_verified"`
	IsGoogle       bool                `json:"is_google"`
	Fee            *int                `json:"fee"`
	WorkStartYear  *int                `json:"work_start_year"`
	Specialist     *SpecialistResponse `json:"specialist"`
	ProfilePicture *string             `json:"profile_picture"`
}

type DoctorResponses struct {
	Pagination PaginationResponse `json:"pagination_info"`
	Doctors    []DoctorResponse   `json:"doctors"`
}

func ConvertToDoctorResponse(doctor entities.Doctor) DoctorResponse {
	doctorResponse := DoctorResponse{
		Id:             doctor.Id,
		Name:           doctor.Name,
		Email:          doctor.Email,
		Certificate:    nil,
		IsOnline:       doctor.IsOnline,
		IsVerified:     doctor.IsVerified,
		IsGoogle:       doctor.IsGoogle,
		Fee:            nil,
		WorkStartYear:  nil,
		Specialist:     nil,
		ProfilePicture: nil,
	}

	if doctor.Specialist.Id.Valid {
		doctorResponse.Specialist = ConvertToSpecialistResponse(doctor.Specialist)
	}
	if doctor.Certificate.Valid {
		doctorResponse.Certificate = &doctor.Certificate.String
	}
	if doctor.Fee.Valid {
		feeInt := int(doctor.Fee.Int64)
		doctorResponse.Fee = &feeInt
	}
	if doctor.ProfilePicture.Valid {
		doctorResponse.ProfilePicture = &doctor.ProfilePicture.String
	}
	if doctor.WorkStartYear.Valid {
		startYearInt := int(doctor.WorkStartYear.Int64)
		doctorResponse.WorkStartYear = &startYearInt
	}

	return doctorResponse
}

func ConvertToDoctorResponses(doctors []entities.Doctor, pagination entities.PaginationInfo) *DoctorResponses {
	doctorResponses := []DoctorResponse{}

	for _, doctor := range doctors {
		doctorResponses = append(doctorResponses, ConvertToDoctorResponse(doctor))
	}

	return &DoctorResponses{
		Pagination: *ConvertToPaginationResponse(pagination),
		Doctors:    doctorResponses,
	}
}
