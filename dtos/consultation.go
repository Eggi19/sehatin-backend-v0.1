package dtos

import (
	"time"

	"github.com/tsanaativa/sehatin-backend-v0.1/entities"
)

type CertificateRequest struct {
	StartDate  string `json:"start_date" binding:"required,datetime=2006-01-02"`
	EndDate    string `json:"end_date" binding:"required,datetime=2006-01-02"`
	Diagnosis  string `json:"diagnosis" binding:"required"`
	PatientAge int    `json:"patient_age" binding:"required"`
}

type ConsultationRequest struct {
	DoctorId         int64  `json:"doctor_id" binding:"required"`
	PatientGenderId  int64  `json:"patient_gender_id" binding:"required"`
	PatientName      string `json:"patient_name" binding:"required"`
	PatientBirthDate string `json:"patient_birth_date" binding:"required,datetime=2006-01-02"`
}

type ConsultationResponse struct {
	Id               int64          `json:"id"`
	Doctor           DoctorResponse `json:"doctor"`
	User             UserResponse   `json:"user"`
	PatientGender    GenderResponse `json:"patient_gender"`
	PatientName      string         `json:"patient_name"`
	PatientBirthDate string         `json:"patient_birth_date"`
	CertificateUrl   *string        `json:"certificate_url"`
	PrescriptionUrl  *string        `json:"prescription_url"`
	EndedAt          *time.Time     `json:"ended_at"`
	CreatedAt        time.Time      `json:"created_at"`
	Chats            []ChatResponse `json:"chats,omitempty"`
}

type ConsultationResponses struct {
	Pagination    PaginationResponse     `json:"pagination_info"`
	Consultations []ConsultationResponse `json:"consultations"`
}

func ConvertToConsultationResponse(consultation entities.Consultation) ConsultationResponse {
	consultationResponse := ConsultationResponse{
		Id:               consultation.Id,
		Doctor:           ConvertToDoctorResponse(consultation.Doctor),
		User:             ConvertToUserResponse(consultation.User),
		PatientGender:    *ConvertToGenderResponse(&consultation.PatientGender),
		PatientName:      consultation.PatientName,
		PatientBirthDate: consultation.PatientBirthDate,
		CertificateUrl:   nil,
		PrescriptionUrl:  nil,
		EndedAt:          nil,
		CreatedAt:        consultation.CreatedAt,
		Chats:            ConvertToChatResponses(consultation.Chats),
	}

	if consultation.CertificateUrl.Valid {
		consultationResponse.CertificateUrl = &consultation.CertificateUrl.String
	}

	if consultation.PrescriptionUrl.Valid {
		consultationResponse.PrescriptionUrl = &consultation.PrescriptionUrl.String
	}

	if consultation.EndedAt.Valid {
		consultationResponse.EndedAt = &consultation.EndedAt.Time
	}

	return consultationResponse
}

func ConvertToConsultationResponses(consultations []entities.Consultation, pagination entities.PaginationInfo) ConsultationResponses {
	consultationsResponses := []ConsultationResponse{}

	for _, c := range consultations {
		consultationsResponses = append(consultationsResponses, ConvertToConsultationResponse(c))
	}

	return ConsultationResponses{
		Pagination:    *ConvertToPaginationResponse(pagination),
		Consultations: consultationsResponses,
	}
}

type CertificateUrlResponse struct {
	CertificateUrl string `json:"certificate_url"`
}

type PrescriptionRequest struct {
	Products   []int64 `json:"products"`
	Quantities []int   `json:"quantities"`
	PatientAge int     `json:"patient_age" binding:"required"`
}

type PrescriptionUrlResponse struct {
	PrescriptionUrl string `json:"prescription_url"`
}
