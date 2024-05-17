package dtos

import (
	"mime/multipart"

	"github.com/tsanaativa/sehatin-backend-v0.1/constants"
	"github.com/tsanaativa/sehatin-backend-v0.1/entities"
)

type UserRegisterData struct {
	Name      string  `json:"name" binding:"required"`
	Email     string  `json:"email" binding:"required,email"`
	Password  *string `json:"password" binding:"required,excludes= ,containsany=abcdefghijklmnopqrstuvwxyz,containsany=ABCDEFGHIJKLMNOPQRSTUVWXYZ,containsany=1234567890,containsany=!#$%&'()*+0x2C-./:\"\\;<=>?@[]^_{0x7C}~,min=8,max=128"`
	BirthDate *string `json:"birth_date" binding:"required,datetime=2006-01-02"`
	GenderId  *int64  `json:"gender_id" binding:"required"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
	Role     string `json:"role" binding:"required"`
}

type UserLoginResponse struct {
	Id             int64                 `json:"id"`
	Role           string                `json:"role"`
	Name           string                `json:"name"`
	Email          string                `json:"email"`
	IsVerified     bool                  `json:"is_verified,omitempty"`
	IsGoogle       bool                  `json:"is_google"`
	IsOnline       bool                  `json:"is_online"`
	ProfilePicture *string               `json:"profile_picture"`
	Addresses      []UserAddressResponse `json:"addresses"`
}

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token,omitempty"`
}

type LoginResponse struct {
	Exp    string             `json:"exp"`
	Tokens TokenResponse      `json:"token"`
	User   *UserLoginResponse `json:"user,omitempty"`
}

type AvailableRole struct {
	User            *entities.User
	Doctor          *entities.Doctor
	PharmacyManager *entities.PharmacyManager
	Admin           *entities.Admin
}

type VerificationReq struct {
	Token    string `json:"token" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type ResendVerificationReq struct {
	Email string `json:"email" binding:"required,email"`
	Role  string `json:"role" binding:"required"`
}

type DoctorRegisterData struct {
	Name                string                `form:"name" binding:"required"`
	Email               string                `form:"email" binding:"required,email"`
	Password            string                `form:"password" binding:"required,excludes= ,containsany=abcdefghijklmnopqrstuvwxyz,containsany=ABCDEFGHIJKLMNOPQRSTUVWXYZ,containsany=1234567890,containsany=!#$%&'()*+0x2C-./:\"\\;<=>?@[]^_{0x7C}~,min=8,max=128"`
	Fee                 int                   `form:"fee" binding:"required"`
	Certificate         *multipart.FileHeader `form:"certificate"`
	WorkStartYear       int                   `form:"work_start_year" binding:"required"`
	DoctorSpecialistsId int64                 `form:"doctor_specialists_id" binding:"required"`
}

type PharmacyManagerData struct {
	Name        string                `form:"name" binding:"required"`
	Email       string                `form:"email" binding:"required,email"`
	Password    string                `form:"password" binding:"required,excludes= ,containsany=abcdefghijklmnopqrstuvwxyz,containsany=ABCDEFGHIJKLMNOPQRSTUVWXYZ,containsany=1234567890,containsany=!#$%&'()*+0x2C-./:\"\\;<=>?@[]^_{0x7C}~,min=8,max=128"`
	PhoneNumber string                `form:"phone_number" binding:"required"`
	Logo        *multipart.FileHeader `form:"logo"`
}

type GoogleAuthRequest struct {
	AuthCode *string `json:"auth_code" binding:"required"`
	Role     string  `json:"role" binding:"required"`
}

type GoogleAuthToken struct {
	IdToken     string `json:"id_token"`
	AccessToken string `json:"access_token"`
}

type GoogleAuthUserResponse struct {
	Name          string  `json:"name"`
	Email         string  `json:"email"`
	Picture       string  `json:"picture"`
	VerifiedEmail bool    `json:"verified_email"`
	Gender        *string `json:"gender,omitempty"`
	BirthDate     *string `json:"birth_date,omitempty"`
}

func ConvertToLoginResponse(role string, availableRole *AvailableRole) *LoginResponse {
	response := UserLoginResponse{}
	switch role {
	case constants.UserRole:
		response.Name = availableRole.User.Name
		response.Email = availableRole.User.Email
		response.IsVerified = availableRole.User.IsVerified
		response.IsGoogle = availableRole.User.IsGoogle
		response.IsOnline = availableRole.User.IsOnline
		response.ProfilePicture = nil
		if availableRole.User.ProfilePicture.Valid {
			response.ProfilePicture = &availableRole.User.ProfilePicture.String
		}
		response.Id = availableRole.User.Id
		response.Addresses = ConvertToUserAddressResponses(availableRole.User.Address)
	case constants.DoctorRole:
		response.Name = availableRole.Doctor.Name
		response.Email = availableRole.Doctor.Email
		response.IsVerified = availableRole.Doctor.IsVerified
		response.IsGoogle = availableRole.Doctor.IsGoogle
		response.IsOnline = availableRole.Doctor.IsOnline
		response.ProfilePicture = nil
		if availableRole.Doctor.ProfilePicture.Valid {
			response.ProfilePicture = &availableRole.Doctor.ProfilePicture.String
		}
		response.Id = availableRole.Doctor.Id
	case constants.PharmacyManagerRole:
		response.Name = availableRole.PharmacyManager.Name
		response.Email = availableRole.PharmacyManager.Email
		response.ProfilePicture = &availableRole.PharmacyManager.Logo
		response.Id = availableRole.PharmacyManager.Id
	case constants.AdminRole:
		response.Name = availableRole.Admin.Name
		response.Email = availableRole.Admin.Email
		response.Id = availableRole.Admin.Id
	}

	response.Role = role

	return &LoginResponse{
		User: &response,
	}
}
