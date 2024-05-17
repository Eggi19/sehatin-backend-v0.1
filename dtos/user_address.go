package dtos

import "github.com/tsanaativa/sehatin-backend-v0.1/entities"

type UserAddressResponse struct {
	Id          int64  `json:"id,omitempty"`
	UserId      int64  `json:"user_id,omitempty"`
	CityId      int    `json:"city_id,omitempty"`
	City        string `json:"city,omitempty"`
	Province    string `json:"province,omitempty"`
	Address     string `json:"address,omitempty"`
	District    string `json:"district,omitempty"`
	SubDistrict string `json:"sub_district,omitempty"`
	PostalCode  string `json:"postal_code,omitempty"`
	Coordinate  string `json:"coordinate,omitempty"`
	IsMain      bool   `json:"is_main,omitempty"`
}

type UserAddressCreateRequest struct {
	CityId      int    `json:"city_id" binding:"required"`
	City        string `json:"city" binding:"required"`
	Province    string `json:"province" binding:"required"`
	Address     string `json:"address" binding:"required"`
	District    string `json:"district" binding:"required"`
	SubDistrict string `json:"sub_district" binding:"required"`
	PostalCode  string `json:"postal_code" binding:"required"`
	Longitude   string `json:"longitude" binding:"required"`
	Latitude    string `json:"latitude" binding:"required"`
	IsMain      *bool  `json:"is_main" binding:"required"`
}

type UserAddressUpdateRequest struct {
	CityId      int    `json:"city_id" binding:"required"`
	City        string `json:"city" binding:"required"`
	Province    string `json:"province" binding:"required"`
	Address     string `json:"address" binding:"required"`
	District    string `json:"district" binding:"required"`
	SubDistrict string `json:"sub_district" binding:"required"`
	PostalCode  string `json:"postal_code" binding:"required"`
	Latitude    string `json:"latitude" binding:"required"`
	Longitude   string `json:"longitude" binding:"required"`
	IsMain      *bool  `json:"is_main" binding:"required"`
}

func ConvertToUserAddressResponse(address *entities.UserAddress) *UserAddressResponse {
	return &UserAddressResponse{
		Id:          address.Id,
		UserId:      address.UserId,
		City:        address.City,
		Province:    address.Province,
		Address:     address.Address,
		District:    address.District,
		SubDistrict: address.SubDistrict,
		PostalCode:  address.PostalCode,
		Coordinate:  address.Coordinate,
		IsMain:      address.IsMain,
		CityId:      address.CityId,
	}
}

func ConvertToUserAddressResponses(addresses []entities.UserAddress) []UserAddressResponse {
	userAddressResponses := []UserAddressResponse{}

	for _, address := range addresses {
		userAddressResponses = append(userAddressResponses, *ConvertToUserAddressResponse(&address))
	}

	return userAddressResponses
}
