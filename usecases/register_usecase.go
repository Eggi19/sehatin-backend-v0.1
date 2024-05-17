package usecases

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/tsanaativa/sehatin-backend-v0.1/constants"
	"github.com/tsanaativa/sehatin-backend-v0.1/custom_errors"
	"github.com/tsanaativa/sehatin-backend-v0.1/dtos"
	"github.com/tsanaativa/sehatin-backend-v0.1/entities"
	"github.com/tsanaativa/sehatin-backend-v0.1/repositories"
	"github.com/tsanaativa/sehatin-backend-v0.1/utils"
)

type RegisterUsecaseOpts struct {
	HashAlgorithm       utils.Hasher
	EmailSender         utils.EmailSender
	Transactor          repositories.Transactor
	AuthTokenProvider   utils.AuthTokenProvider
	UploadFile          utils.FileUploader
	DoctorRepo          repositories.DoctorRepository
	UserRepo            repositories.UserRepository
	PharmacyManagerRepo repositories.PharmacyManagerRepository
}

type RegisterUsecase interface {
	RegisterUser(ctx context.Context, user entities.User) error
	RegisterUserWithTransaction(ctx context.Context, user entities.User) error
	RegisterDoctor(ctx context.Context, req dtos.DoctorRegisterData) error
	RegisterDoctorWithTransaction(ctx context.Context, doctor dtos.DoctorRegisterData) error
	RegisterPharmacyManager(ctx context.Context, req dtos.PharmacyManagerData) error
}

type RegisterUsecaseImpl struct {
	HashAlgorithm             utils.Hasher
	EmailSender               utils.EmailSender
	Transactor                repositories.Transactor
	AuthTokenProvider         utils.AuthTokenProvider
	UploadFile                utils.FileUploader
	DoctorRepository          repositories.DoctorRepository
	UserRepository            repositories.UserRepository
	PharmacyManagerRepository repositories.PharmacyManagerRepository
}

func NewRegisterUsecaseImpl(registerOpts *RegisterUsecaseOpts) *RegisterUsecaseImpl {
	return &RegisterUsecaseImpl{
		HashAlgorithm:             registerOpts.HashAlgorithm,
		EmailSender:               registerOpts.EmailSender,
		Transactor:                registerOpts.Transactor,
		AuthTokenProvider:         registerOpts.AuthTokenProvider,
		UploadFile:                registerOpts.UploadFile,
		DoctorRepository:          registerOpts.DoctorRepo,
		UserRepository:            registerOpts.UserRepo,
		PharmacyManagerRepository: registerOpts.PharmacyManagerRepo,
	}
}

func (u *RegisterUsecaseImpl) RegisterUser(ctx context.Context, user entities.User) error {
	config, err := utils.ConfigInit()
	if err != nil {
		return err
	}

	pwd := user.Password
	pwdHash, err := u.HashAlgorithm.HashPassword(pwd.String)
	if err != nil {
		return err
	}

	user.Password = *utils.ByteToNullString(pwdHash)

	newUser, err := u.UserRepository.CreateOneUser(ctx, user)
	if err != nil {
		return err
	}

	dataTokenMap := make(map[string]interface{})
	dataTokenMap[constants.UserEmail] = user.Email
	dataTokenMap[constants.Role] = constants.UserRole

	token, err := u.AuthTokenProvider.CreateAndSign(dataTokenMap)
	if err != nil {
		return err
	}

	message := fmt.Sprintf("%s%s", config.FrontendUrl, token.AccessToken)

	tokenExp := time.Now().Add(time.Hour * 1)

	err = u.UserRepository.UserVerificationToken(ctx, newUser.Id, token.AccessToken, tokenExp)
	if err != nil {
		return err
	}

	err = u.EmailSender.SendEmail(user.Email, message, constants.VerificationEmailSubject)
	if err != nil {
		return custom_errors.SendEmail()
	}

	return nil
}

func (u *RegisterUsecaseImpl) RegisterUserWithTransaction(ctx context.Context, user entities.User) error {
	_, err := u.Transactor.WithinTransaction(ctx, func(txCtx context.Context) (interface{}, error) {
		err := u.RegisterUser(txCtx, user)
		if err != nil {
			return nil, err
		}
		return nil, nil
	})
	if err != nil {
		return err
	}

	return nil
}

func (u *RegisterUsecaseImpl) RegisterDoctor(ctx context.Context, req dtos.DoctorRegisterData) error {
	doctor := entities.Doctor{
		Specialist: &entities.DoctorSpecialist{},
	}

	config, err := utils.ConfigInit()
	if err != nil {
		return err
	}

	pwd := req.Password
	pwdHash, err := u.HashAlgorithm.HashPassword(pwd)
	if err != nil {
		return err
	}

	file, _ := req.Certificate.Open()
	if req.Certificate.Size > 1000000 {
		return custom_errors.FileTooLarge()
	}
	if strings.Split(req.Certificate.Filename, ".")[1] != "pdf" {
		return custom_errors.FileNotPdf()
	}

	fileUrl, err := u.UploadFile.UploadFile(ctx, file)
	if err != nil {
		return custom_errors.UploadFile()
	}

	doctor.Name = req.Name
	doctor.Email = req.Email
	doctor.Password = *utils.ByteToNullString(pwdHash)
	doctor.Fee = *utils.Int64ToNullInt64(int64(req.Fee))
	doctor.Certificate = *utils.StringToNullString(fileUrl)
	doctor.WorkStartYear = *utils.Int64ToNullInt64(int64(req.WorkStartYear))
	doctor.Specialist.Id = *utils.Int64ToNullInt64(req.DoctorSpecialistsId)

	newDoctor, err := u.DoctorRepository.CreateOneDoctor(ctx, doctor)
	if err != nil {
		return err
	}

	dataTokenMap := make(map[string]interface{})
	dataTokenMap[constants.UserEmail] = doctor.Email
	dataTokenMap[constants.Role] = constants.DoctorRole

	token, err := u.AuthTokenProvider.CreateAndSign(dataTokenMap)
	if err != nil {
		return err
	}

	message := fmt.Sprintf("%s%s", config.FrontendUrl, token.AccessToken)

	tokenExp := time.Now().Add(time.Hour * 1)

	err = u.DoctorRepository.DoctorVerificationToken(ctx, newDoctor.Id, token.AccessToken, tokenExp)
	if err != nil {
		return err
	}

	err = u.EmailSender.SendEmail(doctor.Email, message, constants.VerificationEmailSubject)
	if err != nil {
		return custom_errors.SendEmail()
	}

	return nil
}

func (u *RegisterUsecaseImpl) RegisterDoctorWithTransaction(ctx context.Context, doctor dtos.DoctorRegisterData) error {
	_, err := u.Transactor.WithinTransaction(ctx, func(txCtx context.Context) (interface{}, error) {
		err := u.RegisterDoctor(txCtx, doctor)
		if err != nil {
			return nil, err
		}
		return nil, nil
	})
	if err != nil {
		return err
	}

	return nil
}

func (u *RegisterUsecaseImpl) RegisterPharmacyManager(ctx context.Context, req dtos.PharmacyManagerData) error {
	var pharmacyManager entities.PharmacyManager

	pwd := req.Password
	pwdHash, err := u.HashAlgorithm.HashPassword(pwd)
	if err != nil {
		return err
	}

	file, _ := req.Logo.Open()
	if req.Logo.Size > 500000 {
		return custom_errors.FileTooLarge()
	}
	if strings.Split(req.Logo.Filename, ".")[1] != "png" {
		return custom_errors.FileNotPng()
	}

	fileUrl, err := u.UploadFile.UploadFile(ctx, file)
	if err != nil {
		return custom_errors.UploadFile()
	}

	pharmacyManager.Name = req.Name
	pharmacyManager.Email = req.Email
	pharmacyManager.Password = string(pwdHash)
	pharmacyManager.PhoneNumber = req.PhoneNumber
	pharmacyManager.Logo = fileUrl

	err = u.PharmacyManagerRepository.CreateOnePharmacyManager(ctx, pharmacyManager)
	if err != nil {
		return err
	}

	message := fmt.Sprintf("password: %s", req.Password)

	err = u.EmailSender.SendEmail(pharmacyManager.Email, message, constants.CredentialEmailSubject)
	if err != nil {
		return custom_errors.SendEmail()
	}

	return nil
}
