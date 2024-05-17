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

type PharmacyAddressRepoOpts struct {
	Db *sql.DB
}

type PharmacyAddressRepository interface {
	CreateOne(ctx context.Context, pharmacyAddress entities.PharmacyAddress) error
	UpdateOne(ctx context.Context, pharmacyAddress entities.PharmacyAddress) error
	DeleteOne(ctx context.Context, pharmacyId int64) error
	FindById(ctx context.Context, id int64) (*entities.PharmacyAddress, error)
}

type PharmacyAddressRepositoryPostgres struct {
	db *sql.DB
}

func NewPharmacyAddressRepositoryPostgres(paOpts *PharmacyAddressRepoOpts) PharmacyAddressRepository {
	return &PharmacyAddressRepositoryPostgres{
		db: paOpts.Db,
	}
}

func (r *PharmacyAddressRepositoryPostgres) CreateOne(ctx context.Context, pharmacyAddress entities.PharmacyAddress) error {
	newPharmacyAddress := entities.PharmacyAddress{}

	values := []interface{}{}
	values = append(values, pharmacyAddress.PharmacyId)
	values = append(values, pharmacyAddress.CityId)
	values = append(values, pharmacyAddress.City)
	values = append(values, pharmacyAddress.Province)
	values = append(values, pharmacyAddress.Address)
	values = append(values, pharmacyAddress.District)
	values = append(values, pharmacyAddress.SubDistrict)
	values = append(values, pharmacyAddress.PostalCode)
	values = append(values, pharmacyAddress.Longitude)
	values = append(values, pharmacyAddress.Latitude)

	var err error

	tx := extractTx(ctx)
	if tx != nil {
		err = tx.QueryRowContext(ctx, qCreatePharmacyAddess, values...).Scan(&newPharmacyAddress.Id)
	} else {
		err = r.db.QueryRowContext(ctx, qCreatePharmacyAddess, values...).Scan(&newPharmacyAddress.Id)
	}

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == constants.ViolatesUniqueConstraintPgErrCode {
			return custom_errors.BadRequest(err, constants.UserEmailNotUniqueErrMsg)
		}
		return err
	}

	return nil
}

func (r *PharmacyAddressRepositoryPostgres) UpdateOne(ctx context.Context, pharmacyAddress entities.PharmacyAddress) error {
	values := []interface{}{}
	values = append(values, pharmacyAddress.Id)
	values = append(values, pharmacyAddress.PharmacyId)
	values = append(values, pharmacyAddress.City)
	values = append(values, pharmacyAddress.Province)
	values = append(values, pharmacyAddress.Address)
	values = append(values, pharmacyAddress.District)
	values = append(values, pharmacyAddress.SubDistrict)
	values = append(values, pharmacyAddress.PostalCode)
	values = append(values, pharmacyAddress.Longitude)
	values = append(values, pharmacyAddress.Latitude)

	var err error
	var stmt *sql.Stmt

	tx := extractTx(ctx)
	if tx != nil {
		stmt, err = tx.PrepareContext(ctx, qUpdatePharmacyAdrress)
	} else {
		stmt, err = r.db.PrepareContext(ctx, qUpdatePharmacyAdrress)
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

func (r *PharmacyAddressRepositoryPostgres) DeleteOne(ctx context.Context, pharmacyId int64) error {
	var err error
	var stmt *sql.Stmt

	tx := extractTx(ctx)
	if tx != nil {
		stmt, err = tx.PrepareContext(ctx, qDeletePharmacyAddress)
	} else {
		stmt, err = r.db.PrepareContext(ctx, qDeletePharmacyAddress)
	}

	if err != nil {
		return err
	}

	res, err := stmt.ExecContext(ctx, pharmacyId)
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

func (r *PharmacyAddressRepositoryPostgres) FindById(ctx context.Context, id int64) (*entities.PharmacyAddress, error) {
	a := entities.PharmacyAddress{}

	var err error

	tx := extractTx(ctx)
	if tx != nil {
		err = tx.QueryRowContext(ctx, qFindPharmacyAddress, id).Scan(&a.CityId)
	} else {
		err = r.db.QueryRowContext(ctx, qFindPharmacyAddress, id).Scan(&a.CityId)
	}

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, custom_errors.NotFound(err)
		}
		return nil, err
	}

	return &a, nil
}
