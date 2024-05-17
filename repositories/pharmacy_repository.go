package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"math"
	"strings"

	"github.com/tsanaativa/sehatin-backend-v0.1/constants"
	"github.com/tsanaativa/sehatin-backend-v0.1/custom_errors"
	"github.com/tsanaativa/sehatin-backend-v0.1/dtos"
	"github.com/tsanaativa/sehatin-backend-v0.1/entities"
	"github.com/jackc/pgx/v5/pgconn"
)

type PharmacyRepoOpts struct {
	Db *sql.DB
}

type PharmacyRepository interface {
	CreateOne(ctx context.Context, pharmacy entities.Pharmacy) (*entities.Pharmacy, error)
	FindOneById(ctx context.Context, pharmacyId int64) (*entities.Pharmacy, error)
	UpdateOne(ctx context.Context, pharmacy entities.Pharmacy) error
	DeleteById(ctx context.Context, pharmacyId int64) error
	FindAllByPharmacyManagerId(ctx context.Context, pharmacyManagerId int64, params entities.PharmacyParams) ([]entities.Pharmacy, int, error)
	FindNearestPharmacy(ctx context.Context, req entities.NearestPharmacyParams) (*entities.Pharmacy, error)
	FindPharmacyByProduct(ctx context.Context, req entities.PharmacyByProductParams) ([]dtos.PharmacyResponse, error)
	FindNearestPharmacyProductByProductId(ctx context.Context, req entities.PharmacyByProductParams) (int64, error)
	FindNearestPharmacyMostBought(ctx context.Context, req entities.NearestPharmacyParams) (*entities.Pharmacy, error)
}

type PharmacyRepositoryPostgres struct {
	db *sql.DB
}

func NewPharmacyRepositoryPostgres(pharmacyOpt *PharmacyRepoOpts) PharmacyRepository {
	return &PharmacyRepositoryPostgres{
		db: pharmacyOpt.Db,
	}
}

func (r *PharmacyRepositoryPostgres) CreateOne(ctx context.Context, pharmacy entities.Pharmacy) (*entities.Pharmacy, error) {
	newPharmacy := entities.Pharmacy{}

	values := []interface{}{}
	values = append(values, pharmacy.Name)
	values = append(values, pharmacy.PharmacyManager.Id)
	values = append(values, pharmacy.OperationalHour)
	values = append(values, pharmacy.OperationalDay)
	values = append(values, pharmacy.PharmacistName)
	values = append(values, pharmacy.PharmacistLicenseNumber)
	values = append(values, pharmacy.PharmacistPhoneNumber)

	var err error

	tx := extractTx(ctx)
	if tx != nil {
		err = tx.QueryRowContext(ctx, qCreateOnePharmacy, values...).Scan(&newPharmacy.Id)
	} else {
		err = r.db.QueryRowContext(ctx, qCreateOnePharmacy, values...).Scan(&newPharmacy.Id)
	}

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == constants.ViolatesUniqueConstraintPgErrCode {
			return nil, custom_errors.BadRequest(err, constants.PharmacyNotUniqueErrMg)
		}
		return nil, err
	}

	return &newPharmacy, nil
}

func (r *PharmacyRepositoryPostgres) FindOneById(ctx context.Context, pharmacyId int64) (*entities.Pharmacy, error) {
	p := entities.Pharmacy{}

	var err error

	tx := extractTx(ctx)
	if tx != nil {
		err = tx.QueryRowContext(ctx, qFindPharmacyById, pharmacyId).Scan(
			&p.Id, &p.Name, &p.OperationalHour, &p.OperationalDay, &p.PharmacistName, &p.PharmacistLicenseNumber, &p.PharmacistPhoneNumber,
			&p.PharmacyManager.Id, &p.PharmacyManager.Name, &p.PharmacyManager.Email, &p.PharmacyManager.PhoneNumber, &p.PharmacyManager.Logo,
			&p.PharmacyAddress.Id, &p.PharmacyAddress.PharmacyId, &p.PharmacyAddress.City, &p.PharmacyAddress.Province, &p.PharmacyAddress.Address, &p.PharmacyAddress.District,
			&p.PharmacyAddress.SubDistrict, &p.PharmacyAddress.PostalCode, &p.PharmacyAddress.Coordinate,
		)
	} else {
		err = r.db.QueryRowContext(ctx, qFindPharmacyById, pharmacyId).Scan(&p.Id, &p.Name, &p.OperationalHour, &p.OperationalDay, &p.PharmacistName, &p.PharmacistLicenseNumber, &p.PharmacistPhoneNumber,
			&p.PharmacyManager.Id, &p.PharmacyManager.Name, &p.PharmacyManager.Email, &p.PharmacyManager.PhoneNumber, &p.PharmacyManager.Logo,
			&p.PharmacyAddress.Id, &p.PharmacyAddress.PharmacyId, &p.PharmacyAddress.City, &p.PharmacyAddress.Province, &p.PharmacyAddress.Address, &p.PharmacyAddress.District,
			&p.PharmacyAddress.SubDistrict, &p.PharmacyAddress.PostalCode, &p.PharmacyAddress.Coordinate,
		)
	}

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, custom_errors.NotFound(err)
		}
		return nil, err
	}

	return &p, err
}

func (r *PharmacyRepositoryPostgres) UpdateOne(ctx context.Context, pharmacy entities.Pharmacy) error {
	values := []interface{}{}
	values = append(values, pharmacy.Id)
	values = append(values, pharmacy.Name)
	values = append(values, pharmacy.OperationalHour)
	values = append(values, pharmacy.OperationalDay)
	values = append(values, pharmacy.PharmacistName)
	values = append(values, pharmacy.PharmacistLicenseNumber)
	values = append(values, pharmacy.PharmacistPhoneNumber)

	var err error
	var stmt *sql.Stmt

	tx := extractTx(ctx)
	if tx != nil {
		stmt, err = tx.PrepareContext(ctx, qUpdatePharmacy)
	} else {
		stmt, err = r.db.PrepareContext(ctx, qUpdatePharmacy)
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

func (r *PharmacyRepositoryPostgres) DeleteById(ctx context.Context, pharmacyId int64) error {
	var err error
	var stmt *sql.Stmt

	tx := extractTx(ctx)
	if tx != nil {
		stmt, err = tx.PrepareContext(ctx, qDeletePharmacyById)
	} else {
		stmt, err = r.db.PrepareContext(ctx, qDeletePharmacyById)
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

func (r *PharmacyRepositoryPostgres) FindAllByPharmacyManagerId(ctx context.Context, pharmacyManagerId int64, params entities.PharmacyParams) ([]entities.Pharmacy, int, error) {
	pharmacies := []entities.Pharmacy{}
	var totalRows int

	var sb strings.Builder
	sb.WriteString(qCountTotalRows)
	sb.WriteString(qPharmacyColl)
	sb.WriteString(qPharmacyCommands)

	var sbTotalRows strings.Builder
	sbTotalRows.WriteString(qCountTotalRows)
	sbTotalRows.WriteString(qPharmacyCommands)

	values := []interface{}{}
	valuesTotalArgument := []interface{}{}

	numberOfArgs := 2
	values = append(values, pharmacyManagerId)
	valuesTotalArgument = append(valuesTotalArgument, pharmacyManagerId)

	if params.Keyword != "" {
		sb.WriteString(` AND p.name ILIKE '%`)
		sb.WriteString(params.Keyword)
		sb.WriteString(`%' `)

		sbTotalRows.WriteString(` AND p.name ILIKE '%`)
		sbTotalRows.WriteString(params.Keyword)
		sbTotalRows.WriteString(`%' `)
	}

	var sortBy string
	switch params.SortBy {
	case "name":
		sortBy = `p.name `
	default:
		sortBy = `p.id `
	}
	sb.WriteString(fmt.Sprintf(`ORDER BY %s `, sortBy))
	sbTotalRows.WriteString(fmt.Sprintf(`ORDER BY %s `, sortBy))

	if params.Sort == "" {
		params.Sort = `ASC`
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
		p := entities.Pharmacy{}
		err := rows.Scan(&totalRows,
			&p.Id, &p.Name, &p.OperationalHour, &p.OperationalDay, &p.PharmacistName, &p.PharmacistLicenseNumber, &p.PharmacistPhoneNumber,
			&p.PharmacyManager.Id, &p.PharmacyManager.Name, &p.PharmacyManager.Email, &p.PharmacyManager.PhoneNumber, &p.PharmacyManager.Logo,
			&p.PharmacyAddress.Id, &p.PharmacyAddress.PharmacyId, &p.PharmacyAddress.City, &p.PharmacyAddress.Province, &p.PharmacyAddress.Address, &p.PharmacyAddress.District,
			&p.PharmacyAddress.SubDistrict, &p.PharmacyAddress.PostalCode, &p.PharmacyAddress.Coordinate,
		)
		if err != nil {
			return nil, 0, err
		}
		pharmacies = append(pharmacies, p)
	}

	err = rows.Err()
	if err != nil && err != sql.ErrNoRows {
		return nil, 0, err
	}

	if totalRows == 0 {
		err := r.db.QueryRowContext(ctx, sbTotalRows.String(), valuesTotalArgument...).Scan(&totalRows)
		if err != nil && err != sql.ErrNoRows {
			return nil, 0, err
		}
	}

	return pharmacies, totalRows, nil
}

func (r *PharmacyRepositoryPostgres) FindNearestPharmacy(ctx context.Context, req entities.NearestPharmacyParams) (*entities.Pharmacy, error) {
	pharmacy := entities.Pharmacy{}

	var err error

	tx := extractTx(ctx)
	if tx != nil {
		err = tx.QueryRowContext(ctx, qFindNearestPharmacy, req.Longitude, req.Latitude, req.Radius).Scan(&pharmacy.Id)
	} else {
		err = r.db.QueryRowContext(ctx, qFindNearestPharmacy, req.Longitude, req.Latitude, req.Radius).Scan(&pharmacy.Id)
	}

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &pharmacy, nil
}

func (r *PharmacyRepositoryPostgres) FindPharmacyByProduct(ctx context.Context, req entities.PharmacyByProductParams) ([]dtos.PharmacyResponse, error) {
	pharmacies := []dtos.PharmacyResponse{}

	var err error
	var rows *sql.Rows

	tx := extractTx(ctx)
	if tx != nil {
		rows, err = tx.QueryContext(ctx, qFindPharmacyByProduct, req.Longitude, req.Latitude, req.Radius, req.ProductId)
	} else {
		rows, err = r.db.QueryContext(ctx, qFindPharmacyByProduct, req.Longitude, req.Latitude, req.Radius, req.ProductId)
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		pharmacy := dtos.PharmacyResponse{}
		err := rows.Scan(&pharmacy.Id, &pharmacy.Name, &pharmacy.Distance)
		if err != nil {
			return nil, err
		}
		ratio := math.Pow(10, float64(2))
		pharmacy.Distance = math.Round(pharmacy.Distance*ratio) / ratio

		pharmacies = append(pharmacies, pharmacy)
	}

	return pharmacies, nil
}

func (r *PharmacyRepositoryPostgres) FindNearestPharmacyProductByProductId(ctx context.Context, req entities.PharmacyByProductParams) (int64, error) {
	var pharmacyProductId int64

	var err error

	tx := extractTx(ctx)
	if tx != nil {
		err = tx.QueryRowContext(ctx, qFindNearestPharmacyProductByProduct, req.Longitude, req.Latitude, req.Radius, req.ProductId).Scan(&pharmacyProductId)
	} else {
		err = r.db.QueryRowContext(ctx, qFindNearestPharmacyProductByProduct, req.Longitude, req.Latitude, req.Radius, req.ProductId).Scan(&pharmacyProductId)
	}

	if err != nil {
		if err == sql.ErrNoRows {
			return 0, nil
		}
		return 0, err
	}

	return pharmacyProductId, nil
}

func (r *PharmacyRepositoryPostgres) FindNearestPharmacyMostBought(ctx context.Context, req entities.NearestPharmacyParams) (*entities.Pharmacy, error) {
	pharmacy := entities.Pharmacy{}

	var err error

	tx := extractTx(ctx)
	if tx != nil {
		err = tx.QueryRowContext(ctx, qFindNearestPharmacyMostBought, req.Longitude, req.Latitude, req.Radius).Scan(&pharmacy.Id)
	} else {
		err = r.db.QueryRowContext(ctx, qFindNearestPharmacyMostBought, req.Longitude, req.Latitude, req.Radius).Scan(&pharmacy.Id)
	}

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &pharmacy, nil
}
