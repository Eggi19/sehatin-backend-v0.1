package repositories

import (
	"context"
	"database/sql"
	"errors"

	"github.com/tsanaativa/sehatin-backend-v0.1/constants"
	"github.com/tsanaativa/sehatin-backend-v0.1/custom_errors"
	"github.com/tsanaativa/sehatin-backend-v0.1/entities"
	"github.com/jackc/pgx/v5/pgconn"
)

type DoctorResetPasswordRepoOpts struct {
	Db *sql.DB
}

type DoctorResetPasswordRepository interface {
	CreateOne(ctx context.Context, userReset entities.UserResetPassword) (*entities.UserResetPassword, error)
	FindOneByToken(ctx context.Context, token string) (*entities.UserResetPassword, error)
	FindOneByUserId(ctx context.Context, userId int64) (*entities.UserResetPassword, error)
	DeleteOneById(ctx context.Context, id int64) error
}

type DoctorResetPasswordRepositoryPostgres struct {
	db *sql.DB
}

func NewDoctorResetPasswordRepositoryPostgres(drrOpts *DoctorResetPasswordRepoOpts) DoctorResetPasswordRepository {
	return &DoctorResetPasswordRepositoryPostgres{
		db: drrOpts.Db,
	}
}

func (r *DoctorResetPasswordRepositoryPostgres) CreateOne(ctx context.Context, userReset entities.UserResetPassword) (*entities.UserResetPassword, error) {
	newUserReset := entities.UserResetPassword{}

	values := []interface{}{}
	values = append(values, userReset.Token)
	values = append(values, userReset.UserId)
	values = append(values, userReset.ExpiredAt)

	var err error

	tx := extractTx(ctx)
	if tx != nil {
		err = tx.QueryRowContext(ctx, qCreateDoctorResetPassword, values...).Scan(&newUserReset.Id)
	} else {
		err = r.db.QueryRowContext(ctx, qCreateDoctorResetPassword, values...).Scan(&newUserReset.Id)
	}

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == constants.ViolatesUniqueConstraintPgErrCode {
			return nil, custom_errors.BadRequest(err, constants.PharmacyNotUniqueErrMg)
		}
		return nil, err
	}

	return &newUserReset, nil
}

func (r *DoctorResetPasswordRepositoryPostgres) FindOneByToken(ctx context.Context, token string) (*entities.UserResetPassword, error) {
	userReset := entities.UserResetPassword{}

	var err error

	tx := extractTx(ctx)
	if tx != nil {
		err = tx.QueryRowContext(ctx, qFindDoctorResetPasswordByToken, token).Scan(&userReset.Id, &userReset.Token, &userReset.UserId, &userReset.ExpiredAt)
	} else {
		err = r.db.QueryRowContext(ctx, qFindDoctorResetPasswordByToken, token).Scan(&userReset.Id, &userReset.Token, &userReset.UserId, &userReset.ExpiredAt)
	}

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, custom_errors.NotFound(err)
		}
		return nil, err
	}

	return &userReset, err
}

func (r *DoctorResetPasswordRepositoryPostgres) FindOneByUserId(ctx context.Context, userId int64) (*entities.UserResetPassword, error) {
	userReset := entities.UserResetPassword{}

	var err error

	tx := extractTx(ctx)
	if tx != nil {
		err = tx.QueryRowContext(ctx, qFindDoctorResetPasswordByDoctorId, userId).Scan(&userReset.Id, &userReset.Token, &userReset.UserId, &userReset.ExpiredAt)
	} else {
		err = r.db.QueryRowContext(ctx, qFindDoctorResetPasswordByDoctorId, userId).Scan(&userReset.Id, &userReset.Token, &userReset.UserId, &userReset.ExpiredAt)
	}

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, custom_errors.NotFound(err)
		}
		return nil, err
	}

	return &userReset, err
}

func (r *DoctorResetPasswordRepositoryPostgres) DeleteOneById(ctx context.Context, id int64) error {
	var err error
	var stmt *sql.Stmt

	tx := extractTx(ctx)
	if tx != nil {
		stmt, err = tx.PrepareContext(ctx, qDeleteDoctorResetPassword)
	} else {
		stmt, err = r.db.PrepareContext(ctx, qDeleteDoctorResetPassword)
	}

	if err != nil {
		return err
	}

	res, err := stmt.ExecContext(ctx, id)
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
