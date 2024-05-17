package dtos

import (
	"mime/multipart"

	"github.com/tsanaativa/sehatin-backend-v0.1/entities"
)

type UserResponse struct {
	Id             int64                 `json:"id"`
	Name           string                `json:"name"`
	Email          string                `json:"email"`
	BirthDate      *string               `json:"birth_date"`
	Gender         GenderResponse        `json:"gender"`
	IsVerified     bool                  `json:"is_verified"`
	IsGoogle       bool                  `json:"is_google"`
	IsOnline       bool                  `json:"is_online"`
	ProfilePicture *string               `json:"profile_picture"`
	UserAddresses  []UserAddressResponse `json:"addresses"`
}

type UserResponses struct {
	Pagination PaginationResponse `json:"pagination_info"`
	Users      []UserResponse     `json:"users"`
}

type GenderResponse struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
}

type GenderUpdateRequest struct {
	Id int64 `json:"id" binding:"required"`
}

type UserUpdateRequest struct {
	Name           string                `form:"name" binding:"required"`
	BirthDate      string                `form:"birth_date" binding:"required,datetime=2006-01-02"`
	GenderId       int64                 `form:"gender_id" binding:"required"`
	ProfilePicture *multipart.FileHeader `form:"profile_picture"`
}

func ConvertToGenderResponse(gender *entities.Gender) *GenderResponse {
	return &GenderResponse{
		Id:   gender.Id,
		Name: gender.Name,
	}
}

func ConvertToUserResponse(user entities.User) UserResponse {
	userResponse := UserResponse{
		Id:             user.Id,
		Name:           user.Name,
		Email:          user.Email,
		Gender:         *ConvertToGenderResponse(user.Gender),
		IsVerified:     user.IsVerified,
		IsGoogle:       user.IsGoogle,
		IsOnline:       user.IsOnline,
		ProfilePicture: nil,
		UserAddresses:  ConvertToUserAddressResponses(user.Address),
	}
	if user.ProfilePicture.Valid {
		userResponse.ProfilePicture = &user.ProfilePicture.String
	}
	if user.BirthDate.Valid {
		userResponse.BirthDate = &user.BirthDate.String
	}

	return userResponse
}

func ConvertToUserResponses(users []entities.User, pagination entities.PaginationInfo) *UserResponses {
	userResponses := []UserResponse{}

	for _, user := range users {
		userResponses = append(userResponses, ConvertToUserResponse(user))
	}

	return &UserResponses{
		Pagination: *ConvertToPaginationResponse(pagination),
		Users:      userResponses,
	}
}
