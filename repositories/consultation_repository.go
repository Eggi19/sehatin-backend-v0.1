package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/tsanaativa/sehatin-backend-v0.1/custom_errors"
	"github.com/tsanaativa/sehatin-backend-v0.1/entities"
)

type ConsultationRepoOpts struct {
	Db *sql.DB
}

type ConsultationRepository interface {
	FindById(ctx context.Context, id int64) (*entities.Consultation, error)
	FindAllByUserId(ctx context.Context, userId int64, params entities.ConsultationParams) ([]entities.Consultation, int, error)
	FindAllByDoctorId(ctx context.Context, userId int64, params entities.ConsultationParams) ([]entities.Consultation, int, error)
	CreateOne(ctx context.Context, consultation entities.Consultation) (*entities.Consultation, error)
	UpdateEndedAt(ctx context.Context, consultationId int64) error
	CreateCertificate(ctx context.Context, consultationId int64, certificateUrl string) error
	CreatePrescriptionItems(ctx context.Context, prescriptionData entities.PrescriptionData) error
	UpdatePrescription(ctx context.Context, consultationId int64, prescriptionUrl string) error
	FindAllPrescribedProductsById(ctx context.Context, consultationId int64) ([]int64, []int, error)
}

type ConsultationRepositoryPostgres struct {
	db *sql.DB
}

func NewConsultationRepositoryPostgres(cOpts *ConsultationRepoOpts) ConsultationRepository {
	return &ConsultationRepositoryPostgres{
		db: cOpts.Db,
	}
}

func (r *ConsultationRepositoryPostgres) FindById(ctx context.Context, id int64) (*entities.Consultation, error) {
	c := entities.Consultation{
		Doctor:        entities.Doctor{Specialist: &entities.DoctorSpecialist{}},
		User:          entities.User{Gender: &entities.Gender{}},
		PatientGender: entities.Gender{},
	}

	var err error

	tx := extractTx(ctx)
	if tx != nil {
		err = tx.QueryRowContext(ctx, qFindConsultationById, id).Scan(&c.Id, &c.Doctor.Id, &c.Doctor.Name, &c.Doctor.ProfilePicture, &c.Doctor.IsOnline, &c.Doctor.Specialist.Id, &c.Doctor.Specialist.Name, &c.User.Id, &c.User.Name, &c.User.ProfilePicture, &c.PatientGender.Id, &c.PatientGender.Name, &c.PatientName, &c.PatientBirthDate, &c.CertificateUrl, &c.PrescriptionUrl, &c.EndedAt, &c.CreatedAt)
	} else {
		err = r.db.QueryRowContext(ctx, qFindConsultationById, id).Scan(&c.Id, &c.Doctor.Id, &c.Doctor.Name, &c.Doctor.ProfilePicture, &c.Doctor.IsOnline, &c.Doctor.Specialist.Id, &c.Doctor.Specialist.Name, &c.User.Id, &c.User.Name, &c.User.ProfilePicture, &c.PatientGender.Id, &c.PatientGender.Name, &c.PatientName, &c.PatientBirthDate, &c.CertificateUrl, &c.PrescriptionUrl, &c.EndedAt, &c.CreatedAt)
	}

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, custom_errors.NotFound(err)
		}
		return nil, err
	}

	return &c, nil
}

func (r *ConsultationRepositoryPostgres) FindAllByUserId(ctx context.Context, userId int64, params entities.ConsultationParams) ([]entities.Consultation, int, error) {
	consultations := []entities.Consultation{}
	var total_rows int

	var qb strings.Builder
	qb.WriteString(qCountTotalRows)
	qb.WriteString(qConsultationColl)
	qb.WriteString(qConsultationCommands)
	qb.WriteString(qConsultationByUserIdCommands)

	var qbCountTotalRows strings.Builder
	qbCountTotalRows.WriteString(qCountTotalRows)
	qbCountTotalRows.WriteString(qConsultationCommands)
	qbCountTotalRows.WriteString(qConsultationByUserIdCommands)

	values := []interface{}{}
	valuesCountTotal := []interface{}{}

	valCount := 2
	values = append(values, userId)
	valuesCountTotal = append(valuesCountTotal, userId)

	if params.Status != "" {
		var condition string
		switch params.Status {
		case "ongoing":
			condition = "ended_at IS NULL"
		case "completed":
			condition = "ended_at IS NOT NULL"
		}

		if condition != "" {
			qb.WriteString(fmt.Sprintf(`AND %s `, condition))
			qbCountTotalRows.WriteString(fmt.Sprintf(`AND %s`, condition))
		}
	}

	qb.WriteString(` ORDER BY c.created_at DESC `)

	qb.WriteString(`LIMIT `)
	qb.WriteString(fmt.Sprintf(`$%d `, valCount))
	values = append(values, params.Limit)
	valCount++

	qb.WriteString(`OFFSET `)
	qb.WriteString(fmt.Sprintf(`$%d `, valCount))
	values = append(values, params.Limit*(params.Page-1))

	rows, err := r.db.QueryContext(ctx, qb.String(), values...)

	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	for rows.Next() {
		c := entities.Consultation{
			Doctor:        entities.Doctor{Specialist: &entities.DoctorSpecialist{}},
			User:          entities.User{Gender: &entities.Gender{}},
			PatientGender: entities.Gender{},
		}

		err := rows.Scan(&total_rows, &c.Id, &c.Doctor.Id, &c.Doctor.Name, &c.Doctor.ProfilePicture, &c.Doctor.IsOnline, &c.Doctor.Specialist.Id, &c.Doctor.Specialist.Name, &c.User.Id, &c.User.Name, &c.User.ProfilePicture, &c.PatientGender.Id, &c.PatientGender.Name, &c.PatientName, &c.PatientBirthDate, &c.CertificateUrl, &c.PrescriptionUrl, &c.EndedAt, &c.CreatedAt)
		if err != nil {
			return nil, 0, err
		}

		consultations = append(consultations, c)
	}

	err = rows.Err()
	if err != nil && err != sql.ErrNoRows {
		return nil, 0, err
	}

	if total_rows == 0 {
		err := r.db.QueryRowContext(ctx, qbCountTotalRows.String(), valuesCountTotal...).Scan(&total_rows)
		if err != nil && err != sql.ErrNoRows {
			return nil, 0, err
		}
	}

	return consultations, total_rows, nil
}

func (r *ConsultationRepositoryPostgres) FindAllByDoctorId(ctx context.Context, userId int64, params entities.ConsultationParams) ([]entities.Consultation, int, error) {
	consultations := []entities.Consultation{}
	var total_rows int

	var qb strings.Builder
	qb.WriteString(qCountTotalRows)
	qb.WriteString(qConsultationColl)
	qb.WriteString(qConsultationCommands)
	qb.WriteString(qConsultationByDoctorIdCommands)

	var qbCountTotalRows strings.Builder
	qbCountTotalRows.WriteString(qCountTotalRows)
	qbCountTotalRows.WriteString(qConsultationCommands)
	qbCountTotalRows.WriteString(qConsultationByDoctorIdCommands)

	values := []interface{}{}
	valuesCountTotal := []interface{}{}

	valCount := 2
	values = append(values, userId)
	valuesCountTotal = append(valuesCountTotal, userId)

	if params.Status != "" {
		var condition string
		switch params.Status {
		case "ongoing":
			condition = "ended_at IS NULL"
		case "completed":
			condition = "ended_at IS NOT NULL"
		}

		if condition != "" {
			qb.WriteString(fmt.Sprintf(`AND %s `, condition))
			qbCountTotalRows.WriteString(fmt.Sprintf(`AND %s`, condition))
		}
	}

	qb.WriteString(` ORDER BY c.created_at DESC `)

	qb.WriteString(`LIMIT `)
	qb.WriteString(fmt.Sprintf(`$%d `, valCount))
	values = append(values, params.Limit)
	valCount++

	qb.WriteString(`OFFSET `)
	qb.WriteString(fmt.Sprintf(`$%d `, valCount))
	values = append(values, params.Limit*(params.Page-1))

	rows, err := r.db.QueryContext(ctx, qb.String(), values...)

	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	for rows.Next() {
		c := entities.Consultation{
			Doctor:        entities.Doctor{Specialist: &entities.DoctorSpecialist{}},
			User:          entities.User{Gender: &entities.Gender{}},
			PatientGender: entities.Gender{},
		}

		err := rows.Scan(&total_rows, &c.Id, &c.Doctor.Id, &c.Doctor.Name, &c.Doctor.ProfilePicture, &c.Doctor.IsOnline, &c.Doctor.Specialist.Id, &c.Doctor.Specialist.Name, &c.User.Id, &c.User.Name, &c.User.ProfilePicture, &c.PatientGender.Id, &c.PatientGender.Name, &c.PatientName, &c.PatientBirthDate, &c.CertificateUrl, &c.PrescriptionUrl, &c.EndedAt, &c.CreatedAt)
		if err != nil {
			return nil, 0, err
		}

		consultations = append(consultations, c)
	}

	err = rows.Err()
	if err != nil && err != sql.ErrNoRows {
		return nil, 0, err
	}

	if total_rows == 0 {
		err := r.db.QueryRowContext(ctx, qbCountTotalRows.String(), valuesCountTotal...).Scan(&total_rows)
		if err != nil && err != sql.ErrNoRows {
			return nil, 0, err
		}
	}

	return consultations, total_rows, nil
}

func (r *ConsultationRepositoryPostgres) CreateOne(ctx context.Context, consultation entities.Consultation) (*entities.Consultation, error) {
	var err error

	tx := extractTx(ctx)
	if tx != nil {
		err = tx.QueryRowContext(ctx, qCreateOneConsultation, consultation.Doctor.Id, consultation.User.Id, consultation.PatientGender.Id, consultation.PatientName, consultation.PatientBirthDate).Scan(&consultation.Id)
	} else {
		err = r.db.QueryRowContext(ctx, qCreateOneConsultation, consultation.Doctor.Id, consultation.User.Id, consultation.PatientGender.Id, consultation.PatientName, consultation.PatientBirthDate).Scan(&consultation.Id)
	}

	if err != nil {
		return nil, err
	}

	return &consultation, nil
}

func (r *ConsultationRepositoryPostgres) UpdateEndedAt(ctx context.Context, consultationId int64) error {
	var err error
	var stmt *sql.Stmt

	tx := extractTx(ctx)
	if tx != nil {
		stmt, err = tx.PrepareContext(ctx, qUpdateEndedAtConsultation)
	} else {
		stmt, err = r.db.PrepareContext(ctx, qUpdateEndedAtConsultation)
	}

	if err != nil {
		return err
	}

	res, err := stmt.ExecContext(ctx, consultationId)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return custom_errors.NotFound(sql.ErrNoRows)
	}

	return nil
}

func (r *ConsultationRepositoryPostgres) CreateCertificate(ctx context.Context, consultationId int64, certificateUrl string) error {
	var err error
	var stmt *sql.Stmt

	tx := extractTx(ctx)
	if tx != nil {
		stmt, err = tx.PrepareContext(ctx, qUpdateCertificateConsultation)
	} else {
		stmt, err = r.db.PrepareContext(ctx, qUpdateCertificateConsultation)
	}

	if err != nil {
		return err
	}

	res, err := stmt.ExecContext(ctx, consultationId, certificateUrl)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return custom_errors.NotFound(sql.ErrNoRows)
	}

	return nil
}

func (r *ConsultationRepositoryPostgres) CreatePrescriptionItems(ctx context.Context, prescriptionData entities.PrescriptionData) error {
	var err error

	var qb strings.Builder
	values := []interface{}{}
	qb.WriteString(qCreatePrescriptionItems)
	valCount := 0

	for i := 0; i < len(prescriptionData.Products); i++ {
		values = append(values, prescriptionData.ConsultationId)
		values = append(values, prescriptionData.Products[i].Id)
		values = append(values, prescriptionData.Quantities[i])
		qb.WriteString(fmt.Sprintf("($%d, $%d, $%d)", valCount+1, valCount+2, valCount+3))
		if i != len(prescriptionData.Products)-1 {
			qb.WriteString(", ")
		}
		valCount += 3
	}

	tx := extractTx(ctx)
	if tx != nil {
		err = tx.QueryRowContext(ctx, qb.String(), values...).Scan()
	} else {
		err = r.db.QueryRowContext(ctx, qb.String(), values...).Scan()
	}

	if err != nil && err != sql.ErrNoRows {
		return err
	}

	return nil
}

func (r *ConsultationRepositoryPostgres) UpdatePrescription(ctx context.Context, consultationId int64, prescriptionUrl string) error {
	var err error
	var stmt *sql.Stmt

	tx := extractTx(ctx)
	if tx != nil {
		stmt, err = tx.PrepareContext(ctx, qUpdatePrescriptionConsultation)
	} else {
		stmt, err = r.db.PrepareContext(ctx, qUpdatePrescriptionConsultation)
	}

	if err != nil {
		return err
	}

	res, err := stmt.ExecContext(ctx, consultationId, prescriptionUrl)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return custom_errors.NotFound(sql.ErrNoRows)
	}

	return nil
}

func (r *ConsultationRepositoryPostgres) FindAllPrescribedProductsById(ctx context.Context, consultationId int64) ([]int64, []int, error) {
	itemIds := []int64{}
	itemQuantities := []int{}

	var err error
	var rows *sql.Rows

	tx := extractTx(ctx)
	if tx != nil {
		rows, err = tx.QueryContext(ctx, qPrescriptionProducts, consultationId)
	} else {
		rows, err = r.db.QueryContext(ctx, qPrescriptionProducts, consultationId)
	}
	defer rows.Close()

	for rows.Next() {
		var productId int64
		var quantity int
		rows.Scan(&productId, &quantity)
		itemIds = append(itemIds, productId)
		itemQuantities = append(itemQuantities, quantity)
	}

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil, custom_errors.NotFound(err)
		}
		return nil, nil, err
	}

	return itemIds, itemQuantities, nil
}
