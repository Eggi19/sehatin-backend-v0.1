package dtos

import "github.com/tsanaativa/sehatin-backend-v0.1/entities"

type SpecialistResponse struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
}

type SpecialistResponses struct {
	SpecialistResponse []SpecialistResponse `json:"specialists"`
}

func ConvertToSpecialistResponse(specialist *entities.DoctorSpecialist) *SpecialistResponse {
	return &SpecialistResponse{
		Id:   specialist.Id.Int64,
		Name: specialist.Name.String,
	}
}

func ConvertToSpecialistResponses(specialists []entities.DoctorSpecialist) SpecialistResponses {
	specialistsResponses := []SpecialistResponse{}

	for _, sp := range specialists {
		specialistsResponses = append(specialistsResponses, *ConvertToSpecialistResponse(&sp))
	}

	return SpecialistResponses{
		SpecialistResponse: specialistsResponses,
	}
}
