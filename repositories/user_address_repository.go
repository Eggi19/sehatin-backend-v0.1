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

type UserAddressRepoOpts struct {
	Db *sql.DB
}

type UserAddressRepository interface {
	FindAllByUserId(ctx context.Context, userId int64) ([]entities.UserAddress, error)
	CreateOne(ctx context.Context, userAddress entities.UserAddress) (int64, error)
	UpdateOne(ctx context.Context, userAddress entities.UserAddress) error
	DeleteOne(ctx context.Context, addressId, userId int64) error
	FindById(ctx context.Context, addressId, userId int64) (*entities.UserAddress, error)
	UpdateIsMain(ctx context.Context, addressId int64) error
	UpdateIsMainFalse(ctx context.Context, userAddressId int64, userId int64) error
	FindMainByUserId(ctx context.Context, userId int64) (*entities.UserAddress, error)
	GetDistanceFromPharmacy(ctx context.Context, userAddressId int64, pharmacyId int64) (*float64, error)
}

type UserAddressRepositoryPostgres struct {
	db *sql.DB
}

func NewUserAddressRepositoryPostgres(aOpt *UserAddressRepoOpts) UserAddressRepository {
	return &UserAddressRepositoryPostgres{
		db: aOpt.Db,
	}
}

func (r *UserAddressRepositoryPostgres) FindAllByUserId(ctx context.Context, userId int64) ([]entities.UserAddress, error) {
	as := []entities.UserAddress{}

	var err error
	var rows *sql.Rows

	tx := extractTx(ctx)
	if tx != nil {
		rows, err = tx.QueryContext(ctx, qFindUserAddressByUserId, userId)
	} else {
		rows, err = r.db.QueryContext(ctx, qFindUserAddressByUserId, userId)
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		a := entities.UserAddress{}
		err = rows.Scan(
			&a.Id, &a.UserId, &a.CityId, &a.City, &a.Province, &a.Address, &a.District, &a.SubDistrict, &a.PostalCode, &a.Coordinate, &a.IsMain,
		)
		if err != nil {
			if err == sql.ErrNoRows {
				return nil, custom_errors.NotFound(err)
			}
			return nil, err
		}
		as = append(as, a)
	}

	return as, nil
}

func (r *UserAddressRepositoryPostgres) CreateOne(ctx context.Context, userAddress entities.UserAddress) (int64, error) {
	newUserAddress := entities.UserAddress{}

	values := []interface{}{}
	values = append(values, userAddress.UserId)
	values = append(values, userAddress.City)
	values = append(values, userAddress.Province)
	values = append(values, userAddress.Address)
	values = append(values, userAddress.District)
	values = append(values, userAddress.SubDistrict)
	values = append(values, userAddress.PostalCode)
	values = append(values, userAddress.Longitude)
	values = append(values, userAddress.Latitude)
	values = append(values, userAddress.IsMain)
	values = append(values, userAddress.CityId)

	var err error

	tx := extractTx(ctx)
	if tx != nil {
		err = tx.QueryRowContext(ctx, qCreateOneUserAddress, values...).Scan(&newUserAddress.Id)
	} else {
		err = r.db.QueryRowContext(ctx, qCreateOneUserAddress, values...).Scan(&newUserAddress.Id)
	}

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == constants.ViolatesUniqueConstraintPgErrCode {
			return 0, custom_errors.BadRequest(err, constants.UserEmailNotUniqueErrMsg)
		}
		return 0, err
	}

	return newUserAddress.Id, nil
}

func (r *UserAddressRepositoryPostgres) UpdateOne(ctx context.Context, userAddress entities.UserAddress) error {
	values := []interface{}{}
	values = append(values, userAddress.Id)
	values = append(values, userAddress.UserId)
	values = append(values, userAddress.City)
	values = append(values, userAddress.Province)
	values = append(values, userAddress.Address)
	values = append(values, userAddress.District)
	values = append(values, userAddress.SubDistrict)
	values = append(values, userAddress.PostalCode)
	values = append(values, userAddress.Longitude)
	values = append(values, userAddress.Latitude)
	values = append(values, userAddress.IsMain)
	values = append(values, userAddress.CityId)

	var err error
	var stmt *sql.Stmt

	tx := extractTx(ctx)
	if tx != nil {
		stmt, err = tx.PrepareContext(ctx, qUpdateUserAdrress)
	} else {
		stmt, err = r.db.PrepareContext(ctx, qUpdateUserAdrress)
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

func (r *UserAddressRepositoryPostgres) DeleteOne(ctx context.Context, addressId, userId int64) error {
	var err error
	var stmt *sql.Stmt

	tx := extractTx(ctx)
	if tx != nil {
		stmt, err = tx.PrepareContext(ctx, qDeleteUserAddress)
	} else {
		stmt, err = r.db.PrepareContext(ctx, qDeleteUserAddress)
	}

	if err != nil {
		return err
	}

	res, err := stmt.ExecContext(ctx, addressId, userId)
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

func (r *UserAddressRepositoryPostgres) FindById(ctx context.Context, userAddressId, userId int64) (*entities.UserAddress, error) {
	a := entities.UserAddress{}

	var err error

	tx := extractTx(ctx)
	if tx != nil {
		err = tx.QueryRowContext(ctx, qFindUserAddressById, userAddressId, userId).Scan(
			&a.Id, &a.UserId, &a.CityId, &a.City, &a.Province, &a.Address, &a.District, &a.SubDistrict, &a.PostalCode, &a.Coordinate, &a.IsMain,
		)
	} else {
		err = r.db.QueryRowContext(ctx, qFindUserAddressById, userAddressId, userId).Scan(
			&a.Id, &a.UserId, &a.CityId, &a.City, &a.Province, &a.Address, &a.District, &a.SubDistrict, &a.PostalCode, &a.Coordinate, &a.IsMain,
		)
	}

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, custom_errors.NotFound(err)
		}
		return nil, err
	}

	return &a, nil
}

func (r *UserAddressRepositoryPostgres) UpdateIsMain(ctx context.Context, userAddressId int64) error {
	var err error
	var stmt *sql.Stmt

	tx := extractTx(ctx)
	if tx != nil {
		stmt, err = tx.PrepareContext(ctx, qUpdateIsMain)
	} else {
		stmt, err = r.db.PrepareContext(ctx, qUpdateIsMain)
	}

	if err != nil {
		return err
	}

	res, err := stmt.ExecContext(ctx, userAddressId, userAddressId)
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

func (r *UserAddressRepositoryPostgres) UpdateIsMainFalse(ctx context.Context, userAddressId int64, userId int64) error {
	var err error
	var stmt *sql.Stmt

	tx := extractTx(ctx)
	if tx != nil {
		stmt, err = tx.PrepareContext(ctx, qUpdateIsMainFalse)
	} else {
		stmt, err = r.db.PrepareContext(ctx, qUpdateIsMainFalse)
	}

	if err != nil {
		return err
	}

	res, err := stmt.ExecContext(ctx, userAddressId, userId)
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

func (r *UserAddressRepositoryPostgres) GetDistanceFromPharmacy(ctx context.Context, userAddressId int64, pharmacyId int64) (*float64, error) {
	var distance float64

	var err error

	tx := extractTx(ctx)
	if tx != nil {
		err = tx.QueryRowContext(ctx, qFindUserPharmacyDistance, userAddressId, pharmacyId).Scan(&distance)
	} else {
		err = r.db.QueryRowContext(ctx, qFindUserPharmacyDistance, userAddressId, pharmacyId).Scan(&distance)
	}

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, custom_errors.NotFound(err)
		}
		return nil, err
	}

	return &distance, nil
}

func (r *UserAddressRepositoryPostgres) FindMainByUserId(ctx context.Context, userId int64) (*entities.UserAddress, error) {
	a := entities.UserAddress{}

	var err error

	tx := extractTx(ctx)
	if tx != nil {
		err = tx.QueryRowContext(ctx, qFindMainUserAddressByUserId, userId).Scan(
			&a.Id, &a.UserId, &a.CityId, &a.City, &a.Province, &a.Address, &a.District, &a.SubDistrict, &a.PostalCode, &a.Coordinate, &a.IsMain,
		)
	} else {
		err = r.db.QueryRowContext(ctx, qFindMainUserAddressByUserId, userId).Scan(
			&a.Id, &a.UserId, &a.CityId, &a.City, &a.Province, &a.Address, &a.District, &a.SubDistrict, &a.PostalCode, &a.Coordinate, &a.IsMain,
		)
	}

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, custom_errors.NotFound(err)
		}
		return nil, err
	}

	return &a, nil
}
