package dtos

import "github.com/tsanaativa/sehatin-backend-v0.1/entities"

type MutationSatusResponse struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
}

func ConvertToMutationStatusResponse(mutation *entities.MutationSatus) *MutationSatusResponse {
	return &MutationSatusResponse{Id: mutation.Id, Name: mutation.Name}
}

func ConvertToMutationStatusResponses(mutations []entities.MutationSatus) []MutationSatusResponse {
	mutationResponses := []MutationSatusResponse{}

	for _, mutation := range mutations {
		mutationResponses = append(mutationResponses, *ConvertToMutationStatusResponse(&mutation))
	}

	return mutationResponses
}
