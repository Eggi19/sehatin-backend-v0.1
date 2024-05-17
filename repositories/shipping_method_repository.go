package repositories

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/tsanaativa/sehatin-backend-v0.1/constants"
	"github.com/tsanaativa/sehatin-backend-v0.1/custom_errors"
	"github.com/tsanaativa/sehatin-backend-v0.1/dtos"
	"github.com/tsanaativa/sehatin-backend-v0.1/entities"
	"github.com/tsanaativa/sehatin-backend-v0.1/utils"
	"github.com/jackc/pgx/v5/pgconn"
)

type ShippingMethodRepoOpt struct {
	Db *sql.DB
}

type ShippingMethodRepository interface {
	GetOfficialShippingMethod(ctx context.Context, pharmacyId int64) ([]entities.OfficialShippingMethod, error)
	GetNonOfficialShippingMethod(ctx context.Context, pharmacyId int64) ([]entities.NonOfficialShippingMethod, error)
	CreatePharmacyOfficialShippingMethod(ctx context.Context, shippingMethod entities.PharmacyShippingMethod) (*entities.PharmacyShippingMethod, error)
	CreatePharmacyNonOfficialShippingMethod(ctx context.Context, shippingMethod entities.PharmacyShippingMethod) (*entities.PharmacyShippingMethod, error)
	UpdatePharmacyOfficialShippingMethod(ctx context.Context, shippingMethod entities.PharmacyShippingMethod) error
	UpdatePharmacyNonOfficialShippingMethod(ctx context.Context, shippingMethod entities.PharmacyShippingMethod) error
	DeleteShippingMethodByPharmacyId(ctx context.Context, pharmacyId int64) error
	DeleteAllOfficialShippingMethod(ctx context.Context, pharmacyId int64) error
	GetOfficialShippingFee(ctx context.Context, shippingId int64) (*entities.OfficialShippingMethod, error)
	GetNonOfficialFee(originCityId int, destinationCityId int, weight float64, courier string) (*dtos.RajaOngkirResponse, error)
	GetNonOfficialShippingService(ctx context.Context, shippingId int64) (*entities.NonOfficialShippingMethod, error)
}

type ShippingMethodRepositoryPostgres struct {
	db *sql.DB
}

func NewShippingMethodRepositoryPostgres(smOpt *ShippingMethodRepoOpt) *ShippingMethodRepositoryPostgres {
	return &ShippingMethodRepositoryPostgres{
		db: smOpt.Db,
	}
}

func (r *ShippingMethodRepositoryPostgres) GetOfficialShippingMethod(ctx context.Context, pharmacyId int64) ([]entities.OfficialShippingMethod, error) {
	shippingMethods := []entities.OfficialShippingMethod{}

	var err error
	var rows *sql.Rows

	tx := extractTx(ctx)
	if tx != nil {
		rows, err = tx.QueryContext(ctx, qFindPharmacyOfficialShippingMethod, pharmacyId)
	} else {
		rows, err = r.db.QueryContext(ctx, qFindPharmacyOfficialShippingMethod, pharmacyId)
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		shippingMethod := entities.OfficialShippingMethod{}
		err := rows.Scan(&shippingMethod.Id, &shippingMethod.Name, &shippingMethod.Fee)
		if err != nil {
			return nil, err
		}

		shippingMethods = append(shippingMethods, shippingMethod)
	}

	return shippingMethods, nil
}

func (r *ShippingMethodRepositoryPostgres) GetNonOfficialShippingMethod(ctx context.Context, pharmacyId int64) ([]entities.NonOfficialShippingMethod, error) {
	shippingMethods := []entities.NonOfficialShippingMethod{}

	var err error
	var rows *sql.Rows

	tx := extractTx(ctx)
	if tx != nil {
		rows, err = tx.QueryContext(ctx, qFindPharmacyNonOfficialShippingMethod, pharmacyId)
	} else {
		rows, err = r.db.QueryContext(ctx, qFindPharmacyNonOfficialShippingMethod, pharmacyId)
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		shippingMethod := entities.NonOfficialShippingMethod{}
		err := rows.Scan(&shippingMethod.Id, &shippingMethod.Name, &shippingMethod.Courier, &shippingMethod.Service, &shippingMethod.Description)
		if err != nil {
			return nil, err
		}

		shippingMethods = append(shippingMethods, shippingMethod)
	}

	return shippingMethods, nil
}

func (r *ShippingMethodRepositoryPostgres) CreatePharmacyOfficialShippingMethod(ctx context.Context, shippingMethod entities.PharmacyShippingMethod) (*entities.PharmacyShippingMethod, error) {
	newShippingMethod := entities.PharmacyShippingMethod{}

	values := []interface{}{}
	values = append(values, shippingMethod.OfficialShippingId)
	values = append(values, shippingMethod.PharmacyId)

	var err error

	tx := extractTx(ctx)
	if tx != nil {
		err = tx.QueryRowContext(ctx, qCreatePharmacyOfficialShippingMethod, values...).Scan(&newShippingMethod.Id)
	} else {
		err = r.db.QueryRowContext(ctx, qCreatePharmacyOfficialShippingMethod, values...).Scan(&newShippingMethod.Id)
	}

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == constants.ViolatesUniqueConstraintPgErrCode {
			return nil, custom_errors.BadRequest(err, constants.PharmacyNotUniqueErrMg)
		}
		return nil, err
	}

	return &newShippingMethod, nil
}

func (r *ShippingMethodRepositoryPostgres) CreatePharmacyNonOfficialShippingMethod(ctx context.Context, shippingMethod entities.PharmacyShippingMethod) (*entities.PharmacyShippingMethod, error) {
	newShippingMethod := entities.PharmacyShippingMethod{}

	values := []interface{}{}
	values = append(values, shippingMethod.NonOfficialShippingId)
	values = append(values, shippingMethod.PharmacyId)

	var err error

	tx := extractTx(ctx)
	if tx != nil {
		err = tx.QueryRowContext(ctx, qCreatePharmacyNonOfficialShippingMethod, values...).Scan(&newShippingMethod.Id)
	} else {
		err = r.db.QueryRowContext(ctx, qCreatePharmacyNonOfficialShippingMethod, values...).Scan(&newShippingMethod.Id)
	}

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == constants.ViolatesUniqueConstraintPgErrCode {
			return nil, custom_errors.BadRequest(err, constants.PharmacyNotUniqueErrMg)
		}
		return nil, err
	}

	return &newShippingMethod, nil
}

func (r *ShippingMethodRepositoryPostgres) UpdatePharmacyOfficialShippingMethod(ctx context.Context, shippingMethod entities.PharmacyShippingMethod) error {
	values := []interface{}{}
	values = append(values, shippingMethod.Id)
	values = append(values, shippingMethod.OfficialShippingId)
	values = append(values, shippingMethod.PharmacyId)

	var err error
	var stmt *sql.Stmt

	tx := extractTx(ctx)
	if tx != nil {
		stmt, err = tx.PrepareContext(ctx, qUpdatedOfficialPharmacyShippingMethod)
	} else {
		stmt, err = r.db.PrepareContext(ctx, qUpdatedOfficialPharmacyShippingMethod)
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

func (r *ShippingMethodRepositoryPostgres) UpdatePharmacyNonOfficialShippingMethod(ctx context.Context, shippingMethod entities.PharmacyShippingMethod) error {
	values := []interface{}{}
	values = append(values, shippingMethod.Id)
	values = append(values, shippingMethod.NonOfficialShippingId)
	values = append(values, shippingMethod.PharmacyId)

	var err error
	var stmt *sql.Stmt

	tx := extractTx(ctx)
	if tx != nil {
		stmt, err = tx.PrepareContext(ctx, qUpdatedNonOfficialPharmacyShippingMethod)
	} else {
		stmt, err = r.db.PrepareContext(ctx, qUpdatedNonOfficialPharmacyShippingMethod)
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

func (r *ShippingMethodRepositoryPostgres) DeleteAllOfficialShippingMethod(ctx context.Context, pharmacyId int64) error {
	var err error
	var stmt *sql.Stmt

	tx := extractTx(ctx)
	if tx != nil {
		stmt, err = tx.PrepareContext(ctx, qDeletedAllPharmacyShippingMethod)
	} else {
		stmt, err = r.db.PrepareContext(ctx, qDeletedAllPharmacyShippingMethod)
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

func (r *ShippingMethodRepositoryPostgres) DeleteShippingMethodByPharmacyId(ctx context.Context, pharmacyId int64) error {
	var err error
	var stmt *sql.Stmt

	tx := extractTx(ctx)
	if tx != nil {
		stmt, err = tx.PrepareContext(ctx, qDeletedShippingMethodByPharmacyId)
	} else {
		stmt, err = r.db.PrepareContext(ctx, qDeletedShippingMethodByPharmacyId)
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

func (r *ShippingMethodRepositoryPostgres) GetOfficialShippingFee(ctx context.Context, shippingId int64) (*entities.OfficialShippingMethod, error) {
	osm := entities.OfficialShippingMethod{}

	var err error

	tx := extractTx(ctx)
	if tx != nil {
		err = tx.QueryRowContext(ctx, qFindOfficialShippingFee, shippingId).Scan(&osm.Fee)
	} else {
		err = r.db.QueryRowContext(ctx, qFindOfficialShippingFee, shippingId).Scan(&osm.Fee)
	}

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, custom_errors.NotFound(err)
		}
		return nil, err
	}

	return &osm, nil
}

func (r *ShippingMethodRepositoryPostgres) GetNonOfficialFee(originCityId int, destinationCityId int, weight float64, courier string) (*dtos.RajaOngkirResponse, error){
	url := "https://api.rajaongkir.com/starter/cost"
	config, err := utils.ConfigInit()
	if err != nil {
		return nil, err
	}

	payload := strings.NewReader(fmt.Sprintf("origin=%d&destination=%d&weight=%.2f&courier=%s", originCityId, destinationCityId, weight, courier))

	request, err := http.NewRequest("POST", url, payload)
	if err != nil {
		return nil, err
	}

	request.Header.Add("key", config.RajaOngkirKey)
	request.Header.Add("content-type", "application/x-www-form-urlencoded")

	res, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var p dtos.RajaOngkirResponse

	err = json.Unmarshal(body, &p)
	if err != nil {
	  return nil, err
	}

	return &p, nil
}

func (r *ShippingMethodRepositoryPostgres) GetNonOfficialShippingService(ctx context.Context, shippingId int64) (*entities.NonOfficialShippingMethod, error) {
	nosm := entities.NonOfficialShippingMethod{}

	var err error

	tx := extractTx(ctx)
	if tx != nil {
		err = tx.QueryRowContext(ctx, qFindNonOfficialShippingMethod, shippingId).Scan(&nosm.Courier, &nosm.Service)
	} else {
		err = r.db.QueryRowContext(ctx, qFindNonOfficialShippingMethod, shippingId).Scan(&nosm.Courier, &nosm.Service)
	}

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, custom_errors.NotFound(err)
		}
		return nil, err
	}

	return &nosm, nil
}