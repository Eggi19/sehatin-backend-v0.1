package usecases

import (
	"context"
	"os"
	"strconv"
	"strings"

	"github.com/tsanaativa/sehatin-backend-v0.1/constants"
	"github.com/tsanaativa/sehatin-backend-v0.1/custom_errors"
	"github.com/tsanaativa/sehatin-backend-v0.1/entities"
	"github.com/tsanaativa/sehatin-backend-v0.1/repositories"
	"github.com/tsanaativa/sehatin-backend-v0.1/utils"
)

type ConsultationUsecaseOpts struct {
	ConsultationRepo repositories.ConsultationRepository
	DoctorRepo       repositories.DoctorRepository
	ChatRepo         repositories.ChatRepository
	ProductRepo      repositories.ProductRepository
	PharmacyRepo     repositories.PharmacyRepository
	UserAddressRepo  repositories.UserAddressRepository
	CartRepo         repositories.CartRepository
	UploadFile       utils.FileUploader
}

type ConsultationUsecase interface {
	GetAllConsultationByUser(ctx context.Context, userId int64, params entities.ConsultationParams) ([]entities.Consultation, *entities.PaginationInfo, error)
	GetAllConsultationByDoctor(ctx context.Context, doctorId int64, params entities.ConsultationParams) ([]entities.Consultation, *entities.PaginationInfo, error)
	GetConsultationById(ctx context.Context, consultationId int64) (*entities.Consultation, error)
	CreateConsultation(ctx context.Context, consultation entities.Consultation) (*entities.Consultation, error)
	EndConsultation(ctx context.Context, consultationId int64, userId int64) error
	CreateChat(ctx context.Context, chat entities.Chat, userId int64) error
	CreateCertificate(ctx context.Context, certificateData entities.CertificateData, doctorId int64) (string, error)
	CreatePrescription(ctx context.Context, prescriptionData entities.PrescriptionData, doctorId int64) (string, error)
	AddPrescriptionToCart(ctx context.Context, consultationId int64, userId int64) error
}

type ConsultationUsecaseImpl struct {
	ConsultationRepository repositories.ConsultationRepository
	DoctorRepository       repositories.DoctorRepository
	ChatRepository         repositories.ChatRepository
	ProductRepository      repositories.ProductRepository
	PharmacyRepository     repositories.PharmacyRepository
	UserAddressRepository  repositories.UserAddressRepository
	CartRepository         repositories.CartRepository
	UploadFile             utils.FileUploader
}

func NewConsultationUsecaseImpl(cuOpts *ConsultationUsecaseOpts) ConsultationUsecase {
	return &ConsultationUsecaseImpl{
		ConsultationRepository: cuOpts.ConsultationRepo,
		DoctorRepository:       cuOpts.DoctorRepo,
		ChatRepository:         cuOpts.ChatRepo,
		ProductRepository:      cuOpts.ProductRepo,
		PharmacyRepository:     cuOpts.PharmacyRepo,
		UserAddressRepository:  cuOpts.UserAddressRepo,
		CartRepository:         cuOpts.CartRepo,
		UploadFile:             cuOpts.UploadFile,
	}
}

func (u *ConsultationUsecaseImpl) GetAllConsultationByUser(ctx context.Context, userId int64, params entities.ConsultationParams) ([]entities.Consultation, *entities.PaginationInfo, error) {
	consultations, totalData, err := u.ConsultationRepository.FindAllByUserId(ctx, userId, params)
	if err != nil {
		return nil, nil, err
	}

	var totalPage int
	if params.Limit != 0 && params.Page != 0 {
		totalPage = totalData / params.Limit
		if totalData%params.Limit > 0 {
			totalPage++
		}
	}

	pagination := entities.PaginationInfo{
		Page:      params.Page,
		Limit:     params.Limit,
		TotalData: totalData,
		TotalPage: totalPage,
	}

	return consultations, &pagination, nil
}

func (u *ConsultationUsecaseImpl) GetAllConsultationByDoctor(ctx context.Context, doctorId int64, params entities.ConsultationParams) ([]entities.Consultation, *entities.PaginationInfo, error) {
	consultations, totalData, err := u.ConsultationRepository.FindAllByDoctorId(ctx, doctorId, params)
	if err != nil {
		return nil, nil, err
	}

	var totalPage int
	if params.Limit != 0 && params.Page != 0 {
		totalPage = totalData / params.Limit
		if totalData%params.Limit > 0 {
			totalPage++
		}
	}

	pagination := entities.PaginationInfo{
		Page:      params.Page,
		Limit:     params.Limit,
		TotalData: totalData,
		TotalPage: totalPage,
	}

	return consultations, &pagination, nil
}

func (u *ConsultationUsecaseImpl) GetConsultationById(ctx context.Context, consultationId int64) (*entities.Consultation, error) {
	consultation, err := u.ConsultationRepository.FindById(ctx, consultationId)
	if err != nil {
		return nil, err
	}

	chats, err := u.ChatRepository.FindAllConsultationChat(ctx, consultationId)
	if err != nil {
		return nil, err
	}

	consultation.Chats = chats

	return consultation, nil
}

func (u *ConsultationUsecaseImpl) CreateConsultation(ctx context.Context, consultation entities.Consultation) (*entities.Consultation, error) {
	doctor, err := u.DoctorRepository.FindOneById(ctx, consultation.Doctor.Id)
	if err != nil {
		return nil, err
	}

	if !doctor.IsVerified {
		return nil, custom_errors.DoctorIsNotVerified()
	}

	newC, err := u.ConsultationRepository.CreateOne(ctx, consultation)
	if err != nil {
		return nil, err
	}

	return newC, nil
}

func (u *ConsultationUsecaseImpl) EndConsultation(ctx context.Context, consultationId int64, userId int64) error {
	existingC, err := u.ConsultationRepository.FindById(ctx, consultationId)
	if err != nil {
		return err
	}

	if existingC.User.Id != userId {
		return custom_errors.Forbidden()
	}

	err = u.ConsultationRepository.UpdateEndedAt(ctx, consultationId)
	if err != nil {
		return err
	}

	return nil
}

func (u *ConsultationUsecaseImpl) CreateChat(ctx context.Context, chat entities.Chat, userId int64) error {
	consultation, err := u.ConsultationRepository.FindById(ctx, chat.ConsultationId)
	if err != nil {
		return err
	}

	if (chat.IsFromUser && (consultation.User.Id != userId)) || (!chat.IsFromUser && (consultation.Doctor.Id != userId)) {
		return custom_errors.Forbidden()
	}

	if chat.Type == "file" {
		file, _ := chat.File.Open()
		if chat.File.Size > 1000000 {
			return custom_errors.FileTooLarge()
		}

		fileExtension := strings.Split(chat.File.Filename, ".")[1]
		if fileExtension != "pdf" && fileExtension != "png" && fileExtension != "jpg" && fileExtension != "jpeg" {
			return custom_errors.FileNotPdf()
		}

		fileUrl, err := u.UploadFile.UploadFile(ctx, file)
		if err != nil {
			return custom_errors.UploadFile()
		}

		chat.Content = fileUrl
	}

	err = u.ChatRepository.CreateOne(ctx, chat)
	if err != nil {
		return err
	}

	return nil
}

func (u *ConsultationUsecaseImpl) CreatePrescription(ctx context.Context, prescriptionData entities.PrescriptionData, doctorId int64) (string, error) {
	consultation, err := u.ConsultationRepository.FindById(ctx, prescriptionData.ConsultationId)
	if err != nil {
		return "", err
	}

	if consultation.Doctor.Id != doctorId {
		return "", custom_errors.Forbidden()
	}

	err = u.ConsultationRepository.CreatePrescriptionItems(ctx, prescriptionData)
	if err != nil {
		return "", err
	}

	for i := 0; i < len(prescriptionData.Products); i++ {
		product, err := u.ProductRepository.FindOneById(ctx, prescriptionData.Products[i].Id)
		if err != nil {
			return "", err
		}
		prescriptionData.Products[i] = *product
	}

	prescriptionData.PatientName = consultation.PatientName
	prescriptionData.PatientGender = consultation.PatientGender
	prescriptionData.PatientBirthDate = consultation.PatientBirthDate
	prescriptionData.DoctorName = consultation.Doctor.Name

	pdfName, err := utils.GeneratePrescriptionPdf(prescriptionData)
	if err != nil {
		return "", err
	}

	fileUrl, err := u.UploadFile.UploadFile(ctx, pdfName)
	if err != nil {
		return "", custom_errors.UploadFile()
	}

	_ = os.Remove(pdfName)

	err = u.ConsultationRepository.UpdatePrescription(ctx, prescriptionData.ConsultationId, fileUrl)
	if err != nil {
		return "", err
	}

	return fileUrl, nil
}

func (u *ConsultationUsecaseImpl) CreateCertificate(ctx context.Context, certificateData entities.CertificateData, doctorId int64) (string, error) {
	consultation, err := u.ConsultationRepository.FindById(ctx, certificateData.ConsultationId)
	if err != nil {
		return "", err
	}

	if consultation.Doctor.Id != doctorId {
		return "", custom_errors.Forbidden()
	}

	certificateData.PatientName = consultation.PatientName
	certificateData.PatientGender = consultation.PatientGender
	certificateData.PatientBirthDate = consultation.PatientBirthDate
	certificateData.DoctorName = consultation.Doctor.Name

	pdfName, err := utils.GenerateCertificatePdf(certificateData)
	if err != nil {
		return "", err
	}

	fileUrl, err := u.UploadFile.UploadFile(ctx, pdfName)
	if err != nil {
		return "", custom_errors.UploadFile()
	}

	_ = os.Remove(pdfName)

	err = u.ConsultationRepository.CreateCertificate(ctx, certificateData.ConsultationId, fileUrl)
	if err != nil {
		return "", err
	}

	return fileUrl, nil
}

func (u *ConsultationUsecaseImpl) AddPrescriptionToCart(ctx context.Context, consultationId int64, userId int64) error {
	consultation, err := u.ConsultationRepository.FindById(ctx, consultationId)
	if err != nil {
		return err
	}

	if consultation.User.Id != userId {
		return custom_errors.Forbidden()
	}

	userMainAddress, err := u.UserAddressRepository.FindMainByUserId(ctx, userId)
	if err != nil {
		return err
	}
	longLat := userMainAddress.Coordinate[16 : len(userMainAddress.Coordinate)-1]
	longLatArr := strings.Split(longLat, " ")
	longitudeStr, latitudeStr := longLatArr[0], longLatArr[1]

	longitude, _ := strconv.ParseFloat(longitudeStr, 64)
	latitude, _ := strconv.ParseFloat(latitudeStr, 64)

	productIds, quantities, err := u.ConsultationRepository.FindAllPrescribedProductsById(ctx, consultationId)
	if err != nil {
		return err
	}

	pharmacyProductIds := []int64{}

	for i := 0; i < len(productIds); i++ {
		params := entities.PharmacyByProductParams{
			Longitude: longitude,
			Latitude:  latitude,
			Radius:    constants.DefaultRadius,
			ProductId: productIds[i],
		}

		pharmacyProductId, err := u.PharmacyRepository.FindNearestPharmacyProductByProductId(ctx, params)
		if err != nil {
			return err
		}

		pharmacyProductIds = append(pharmacyProductIds, pharmacyProductId)
	}

	for i := 0; i < len(pharmacyProductIds); i++ {
		cartEntity := entities.CartItem{
			UserId:            int64(userId),
			Quantity:          quantities[i],
			PharmacyProductId: pharmacyProductIds[i],
		}

		cart, err := u.CartRepository.FindCartItem(ctx, cartEntity)
		if err != nil && err.Error() != constants.ResponseMsgErrorNotFound {
			return err
		}

		if err != nil && err.Error() == constants.ResponseMsgErrorNotFound {
			err = u.CartRepository.CreateOneCartItem(ctx, cartEntity)
			if err != nil {
				return err
			}

			continue
		}

		cartEntity.Id = cart.Id

		err = u.CartRepository.IncreaseCartQuantity(ctx, cartEntity)
		if err != nil {
			return err
		}

	}

	return nil
}
