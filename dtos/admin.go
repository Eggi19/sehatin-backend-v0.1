package dtos

import "github.com/tsanaativa/sehatin-backend-v0.1/entities"

type AdminResponse struct {
	Id    int64  `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type AdminRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `form:"password" binding:"excludes= ,containsany=abcdefghijklmnopqrstuvwxyz,containsany=ABCDEFGHIJKLMNOPQRSTUVWXYZ,containsany=1234567890,containsany=!#$%&'()*+0x2C-./:\"\\;<=>?@[]^_{0x7C}~,min=8,max=128"`
}

type AdminResponses struct {
	Pagination PaginationResponse `json:"pagination_info"`
	Admins     []AdminResponse    `json:"admin"`
}

func ConvertToAdminResponse(admin *entities.Admin) *AdminResponse {
	return &AdminResponse{
		Id:    admin.Id,
		Name:  admin.Name,
		Email: admin.Email,
	}
}

func ConvertToAdminResponses(admins []entities.Admin, pagination entities.PaginationInfo) *AdminResponses {
	adminResponses := []AdminResponse{}

	for _, admin := range admins {
		adminResponses = append(adminResponses, *ConvertToAdminResponse(&admin))
	}

	return &AdminResponses{
		Pagination: *ConvertToPaginationResponse(pagination),
		Admins:     adminResponses,
	}
}
