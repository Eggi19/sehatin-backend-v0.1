package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/tsanaativa/sehatin-backend-v0.1/constants"
	"github.com/tsanaativa/sehatin-backend-v0.1/custom_errors"
	"github.com/tsanaativa/sehatin-backend-v0.1/entities"
	"github.com/jackc/pgx/v5/pgconn"
)

type DoctorRepoOpt struct {
	Db *sql.DB
}

type DoctorRepository interface {
	FindOneByEmail(ctx context.Context, email string) (*entities.Doctor, error)
	FindOneById(ctx context.Context, doctorId int64) (*entities.Doctor, error)
	UpdateOne(ctx context.Context, doctor entities.Doctor) error
	Delete(ctx context.Context, doctorId int64) error
	FindAll(ctx context.Context, params entities.DoctorParams, isPublic bool) ([]entities.Doctor, int, error)
	DoctorVerificationToken(ctx context.Context, doctorId int64, token string, exp time.Time) error
	CreateOneDoctor(ctx context.Context, u entities.Doctor) (*entities.Doctor, error)
	VerifyDoctor(ctx context.Context, email string) error
	UpdateIsOnline(ctx context.Context, doctor entities.Doctor) error
	UpdatePassword(ctx context.Context, doctorId int64, newPassword string) error
	FindPasswordById(ctx context.Context, doctorId int64) (*entities.Doctor, error)
}

type DoctorRepositoryPostgres struct {
	db *sql.DB
}

func NewDoctorRepositoryPostgres(drOpt *DoctorRepoOpt) DoctorRepository {
	return &DoctorRepositoryPostgres{
		db: drOpt.Db,
	}
}

func (r *DoctorRepositoryPostgres) FindOneByEmail(ctx context.Context, email string) (*entities.Doctor, error) {
	d := entities.Doctor{}

	var err error

	tx := extractTx(ctx)
	if tx != nil {
		err = tx.QueryRowContext(ctx, qFindDoctorByEmail, email).Scan(&d.Id, &d.Name, &d.Email, &d.Password, &d.IsVerified, &d.IsGoogle, &d.IsOnline, &d.ProfilePicture)
	} else {
		err = r.db.QueryRowContext(ctx, qFindDoctorByEmail, email).Scan(&d.Id, &d.Name, &d.Email, &d.Password, &d.IsVerified, &d.IsGoogle, &d.IsOnline, &d.ProfilePicture)
	}

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, custom_errors.NotFound(err)
		}
		return nil, err
	}

	return &d, nil
}

func (r *DoctorRepositoryPostgres) FindOneById(ctx context.Context, doctorId int64) (*entities.Doctor, error) {
	d := entities.Doctor{Specialist: &entities.DoctorSpecialist{}}

	var err error

	tx := extractTx(ctx)
	if tx != nil {
		err = tx.QueryRowContext(ctx, qFindDoctorById, doctorId).Scan(&d.Id, &d.Name, &d.Email, &d.Certificate, &d.IsOnline, &d.IsVerified, &d.IsGoogle, &d.Fee, &d.WorkStartYear,
			&d.Specialist.Id, &d.Specialist.Name, &d.ProfilePicture)
	} else {
		err = r.db.QueryRowContext(ctx, qFindDoctorById, doctorId).Scan(&d.Id, &d.Name, &d.Email, &d.Certificate, &d.IsOnline, &d.IsVerified, &d.IsGoogle, &d.Fee, &d.WorkStartYear,
			&d.Specialist.Id, &d.Specialist.Name, &d.ProfilePicture)
	}

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, custom_errors.NotFound(err)
		}
		return nil, err
	}

	return &d, nil
}

func (r *DoctorRepositoryPostgres) UpdateOne(ctx context.Context, doctor entities.Doctor) error {
	values := []interface{}{}
	values = append(values, doctor.Id)
	values = append(values, doctor.Name)
	values = append(values, doctor.Fee)
	values = append(values, doctor.WorkStartYear)
	values = append(values, doctor.Specialist.Id)

	var err error
	var stmt *sql.Stmt

	numberOfArgs := 6

	var sb strings.Builder
	sb.WriteString(qUpdateDoctorColl)
	if doctor.ProfilePicture.Valid {
		sb.WriteString(`profile_picture = `)
		sb.WriteString(fmt.Sprintf(`$%d ,`, numberOfArgs))
		values = append(values, doctor.ProfilePicture)
		numberOfArgs++
	}
	if doctor.Certificate.Valid {
		sb.WriteString(`certificate = `)
		sb.WriteString(fmt.Sprintf(`$%d ,`, numberOfArgs))
		values = append(values, doctor.Certificate)
		numberOfArgs++
	}
	sb.WriteString(qUpdateDoctorCommand)

	tx := extractTx(ctx)
	if tx != nil {
		stmt, err = tx.PrepareContext(ctx, sb.String())
	} else {
		stmt, err = r.db.PrepareContext(ctx, sb.String())
	}

	if err != nil {
		return err
	}

	res, err := stmt.ExecContext(ctx, values...)
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

func (r *DoctorRepositoryPostgres) Delete(ctx context.Context, doctorId int64) error {
	var err error
	var stmt *sql.Stmt

	tx := extractTx(ctx)
	if tx != nil {
		stmt, err = tx.PrepareContext(ctx, qDeleteDoctor)
	} else {
		stmt, err = r.db.PrepareContext(ctx, qDeleteDoctor)
	}

	if err != nil {
		return err
	}

	res, err := stmt.ExecContext(ctx, doctorId)
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

func (r *DoctorRepositoryPostgres) FindAll(ctx context.Context, params entities.DoctorParams, isPublic bool) ([]entities.Doctor, int, error) {
	doctors := []entities.Doctor{}
	var totalRows int

	var sb strings.Builder
	sb.WriteString(qCountTotalRows)
	sb.WriteString(qDoctorColl)
	sb.WriteString(qDoctorCommands)

	var sbTotalRows strings.Builder
	sbTotalRows.WriteString(qCountTotalRows)
	sbTotalRows.WriteString(qDoctorCommands)

	values := []interface{}{}
	valuesCountTotal := []interface{}{}

	numberOfArgs := 1

	if params.SpecialistId != 0 {
		sb.WriteString(`AND ds.id = `)
		sb.WriteString(fmt.Sprintf(`$%d `, numberOfArgs))
		values = append(values, params.SpecialistId)

		sbTotalRows.WriteString(`AND ds.id = `)
		sbTotalRows.WriteString(fmt.Sprintf(`$%d `, numberOfArgs))
		valuesCountTotal = append(valuesCountTotal, params.SpecialistId)

		numberOfArgs++
	}

	if isPublic {
		sb.WriteString(`AND d.is_verified = true `)
	}

	if params.Keyword != "" {
		sb.WriteString(`AND d.name ILIKE '%`)
		sb.WriteString(params.Keyword)
		sb.WriteString(`%' `)

		sbTotalRows.WriteString(`AND d.name ILIKE '%`)
		sbTotalRows.WriteString(params.Keyword)
		sbTotalRows.WriteString(`%' `)
	}

	if params.IsOnline != "" {
		sb.WriteString(fmt.Sprintf(`AND d.is_online = %s `, params.IsOnline))

		sbTotalRows.WriteString(fmt.Sprintf(`AND d.is_online = %s `, params.IsOnline))
	}

	var sortBy string
	switch params.SortBy {
	case "experience":
		sortBy = `d.work_start_year `
	case "fee":
		sortBy = `d.fee `
	default:
		sortBy = `d.is_online `
	}
	sb.WriteString(fmt.Sprintf(`ORDER BY %s `, sortBy))
	sbTotalRows.WriteString(fmt.Sprintf(`ORDER BY %s `, sortBy))

	if params.Sort == "" {
		params.Sort = `DESC `
	}
	sb.WriteString(fmt.Sprintf(`%s `, params.Sort))
	sbTotalRows.WriteString(fmt.Sprintf(`%s `, params.Sort))

	sb.WriteString(`LIMIT `)
	sb.WriteString(fmt.Sprintf(`$%d `, numberOfArgs))
	values = append(values, params.Limit)
	numberOfArgs++

	sb.WriteString(`OFFSET `)
	sb.WriteString(fmt.Sprintf(`$%d`, numberOfArgs))
	values = append(values, params.Limit*(params.Page-1))
	numberOfArgs++

	rows, err := r.db.QueryContext(ctx, sb.String(), values...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	for rows.Next() {
		d := entities.Doctor{Specialist: &entities.DoctorSpecialist{}}
		err := rows.Scan(&totalRows, &d.Id, &d.Name, &d.Email, &d.Certificate, &d.IsOnline, &d.IsVerified, &d.IsGoogle, &d.Fee, &d.WorkStartYear,
			&d.Specialist.Id, &d.Specialist.Name, &d.ProfilePicture)
		if err != nil {
			if err == sql.ErrNoRows {
				return nil, 0, nil
			}
			return nil, 0, err
		}
		doctors = append(doctors, d)
	}

	if totalRows == 0 {
		err := r.db.QueryRowContext(ctx, sbTotalRows.String(), valuesCountTotal...).Scan(&totalRows)
		if err != nil {
			if err == sql.ErrNoRows {
				return nil, 0, nil
			}
			return nil, 0, err
		}
	}

	return doctors, totalRows, nil
}

func (r *DoctorRepositoryPostgres) DoctorVerificationToken(ctx context.Context, doctorId int64, token string, exp time.Time) error {
	var err error

	tx := extractTx(ctx)
	if tx != nil {
		_, err = tx.ExecContext(ctx, qCreateDoctorVerificationToken, token, exp, doctorId)
	} else {
		_, err = r.db.ExecContext(ctx, qCreateDoctorVerificationToken, token, exp, doctorId)
	}

	if err != nil {
		return err
	}

	return nil
}

func (r *DoctorRepositoryPostgres) CreateOneDoctor(ctx context.Context, d entities.Doctor) (*entities.Doctor, error) {
	newD := entities.Doctor{}

	values := []interface{}{}
	values = append(values, d.Name)
	values = append(values, d.Email)
	values = append(values, d.Password)
	values = append(values, d.Fee)
	values = append(values, d.Certificate)
	values = append(values, d.WorkStartYear)
	values = append(values, d.Specialist.Id)

	var err error

	tx := extractTx(ctx)
	if tx != nil {
		err = tx.QueryRowContext(ctx, qCreateOneDoctor, values...).Scan(&newD.Id)
	} else {
		err = r.db.QueryRowContext(ctx, qCreateOneDoctor, values...).Scan(&newD.Id)
	}

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == constants.ViolatesUniqueConstraintPgErrCode {
			return nil, custom_errors.BadRequest(err, constants.UserEmailNotUniqueErrMsg)
		}
		if errors.As(err, &pgErr) && pgErr.Code == constants.VioletesForeignKeyConstraintPgErrCode {
			return nil, custom_errors.BadRequest(err, constants.SpecialistNotFoundErrMsg)
		}
		return nil, err
	}

	return &newD, nil
}

func (r *DoctorRepositoryPostgres) VerifyDoctor(ctx context.Context, email string) error {
	var err error

	tx := extractTx(ctx)
	if tx != nil {
		_, err = tx.ExecContext(ctx, qVerifyDoctor, email)
	} else {
		_, err = r.db.ExecContext(ctx, qVerifyDoctor, email)
	}

	if err != nil {
		return err
	}

	return nil
}

func (r *DoctorRepositoryPostgres) UpdateIsOnline(ctx context.Context, doctor entities.Doctor) error {
	var err error
	var stmt *sql.Stmt

	tx := extractTx(ctx)
	if tx != nil {
		stmt, err = tx.PrepareContext(ctx, qUpdateIsOnlineDoctor)
	} else {
		stmt, err = r.db.PrepareContext(ctx, qUpdateIsOnlineDoctor)
	}

	if err != nil {
		return err
	}

	res, err := stmt.ExecContext(ctx, doctor.Id, doctor.IsOnline)
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

func (r *DoctorRepositoryPostgres) UpdatePassword(ctx context.Context, doctorId int64, newPassword string) error {
	var err error
	var stmt *sql.Stmt

	tx := extractTx(ctx)
	if tx != nil {
		stmt, err = tx.PrepareContext(ctx, qUpadateDoctorPassword)
	} else {
		stmt, err = r.db.PrepareContext(ctx, qUpadateDoctorPassword)
	}

	if err != nil {
		return err
	}

	res, err := stmt.ExecContext(ctx, doctorId, newPassword)
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

func (r *DoctorRepositoryPostgres) FindPasswordById(ctx context.Context, doctorId int64) (*entities.Doctor, error) {
	d := entities.Doctor{}

	var err error

	tx := extractTx(ctx)
	if tx != nil {
		err = tx.QueryRowContext(ctx, qFindDoctorPasswordById, doctorId).Scan(
			&d.Id, &d.Password)
	} else {
		err = r.db.QueryRowContext(ctx, qFindDoctorPasswordById, doctorId).Scan(
			&d.Id, &d.Password)
	}

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, custom_errors.NotFound(err)
		}
		return nil, err
	}

	return &d, nil
}
