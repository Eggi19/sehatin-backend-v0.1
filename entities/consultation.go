package entities

import (
	"database/sql"
	"time"
)

type Consultation struct {
	Id                int64
	Doctor            Doctor
	User              User
	PatientGender     Gender
	PatientName       string
	PatientBirthDate  string
	CertificateUrl    sql.NullString
	PrescriptionUrl   sql.NullString
	PrescriptionItems []PrescriptionItem
	EndedAt           sql.NullTime
	CreatedAt         time.Time
	UpdatedAt         time.Time
	DeletedAt         sql.NullTime
	Chats             []Chat
}

type ConsultationParams struct {
	Limit  int
	Page   int
	Status string
}

type CertificateData struct {
	ConsultationId   int64
	Diagnosis        string
	StartDate        string
	EndDate          string
	PatientName      string
	PatientBirthDate string
	PatientGender    Gender
	PatientAge       int
	DoctorName       string
}

type PrescriptionData struct {
	ConsultationId   int64
	Products         []Product
	Quantities       []int
	PatientName      string
	PatientBirthDate string
	PatientGender    Gender
	PatientAge       int
	DoctorName       string
}
