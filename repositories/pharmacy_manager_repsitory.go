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

type PharmacyManagerRepoOpt struct {
	Db *sql.DB
}

type PharmacyManagerRepository interface {
	FindOneByEmail(ctx context.Context, email string) (*entities.PharmacyManager, error)
	FindAll(ctx context.Context, params entities.PharmacyManagerParams) ([]entities.PharmacyManager, int, error)
	FindOneById(ctx context.Context, pharmacyManagerId int64) (*entities.PharmacyManager, error)
	UpdateOne(ctx context.Context, pharmacyManager entities.PharmacyManager) error
	DeleteById(ctx context.Context, pharmacyManagerId int64) error
	CreateOnePharmacyManager(ctx context.Context, p entities.PharmacyManager) error
}

type PharmacyManagerRepositoryPostgres struct {
	db *sql.DB
}

func NewPharmacyManagerRepositoryPostgres(pmOpt *PharmacyManagerRepoOpt) PharmacyManagerRepository {
	return &PharmacyManagerRepositoryPostgres{
		db: pmOpt.Db,
	}
}

func (r *PharmacyManagerRepositoryPostgres) FindOneByEmail(ctx context.Context, email string) (*entities.PharmacyManager, error) {
	pm := entities.PharmacyManager{}

	var err error

	tx := extractTx(ctx)
	if tx != nil {
		err = tx.QueryRowContext(ctx, qFindPharmacyManagerByEmail, email).Scan(&pm.Id, &pm.Name, &pm.Email, &pm.Password, &pm.Logo)
	} else {
		err = r.db.QueryRowContext(ctx, qFindPharmacyManagerByEmail, email).Scan(&pm.Id, &pm.Name, &pm.Email, &pm.Password, &pm.Logo)
	}

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, custom_errors.NotFound(err)
		}
		return nil, err
	}

	return &pm, nil
}

func (r *PharmacyManagerRepositoryPostgres) FindAll(ctx context.Context, params entities.PharmacyManagerParams) ([]entities.PharmacyManager, int, error) {
	pms := []entities.PharmacyManager{}
	var totalRows int

	var sb strings.Builder
	sb.WriteString(qCountTotalRows)
	sb.WriteString(qPharmacyManagerColl)
	sb.WriteString(qPharmacyManagerCommands)

	var sbTotalRows strings.Builder
	sbTotalRows.WriteString(qCountTotalRows)
	sbTotalRows.WriteString(qPharmacyManagerCommands)

	values := []interface{}{}
	valuesTotalArgument := []interface{}{}

	numberOfArgs := 1

	if params.Keyword != "" {
		sb.WriteString(`AND name ILIKE '%`)
		sb.WriteString(params.Keyword)
		sb.WriteString(`%' `)

		sbTotalRows.WriteString(`AND name ILIKE '%`)
		sbTotalRows.WriteString(params.Keyword)
		sbTotalRows.WriteString(`%' `)
	}

	var sortBy string
	switch params.SortBy {
	case "name":
		sortBy = `name `
	default:
		sortBy = `id `
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
		pm := entities.PharmacyManager{}
		err := rows.Scan(&totalRows, &pm.Id, &pm.Name, &pm.Email, &pm.PhoneNumber, &pm.Logo)
		if err != nil {
			return nil, 0, err
		}
		pms = append(pms, pm)
	}

	err = rows.Err()
	if err != nil && err != sql.ErrNoRows {
		return nil, 0, err
	}

	if totalRows == 0 {
		err := r.db.QueryRowContext(ctx, sbTotalRows.String(), valuesTotalArgument...).Scan(&totalRows)
		if err != nil {
			if err == sql.ErrNoRows {
				return nil, 0, custom_errors.NotFound(err)
			}
			return nil, 0, err
		}
	}

	return pms, totalRows, nil
}

func (r *PharmacyManagerRepositoryPostgres) FindOneById(ctx context.Context, pharmacyManagerId int64) (*entities.PharmacyManager, error) {
	pm := entities.PharmacyManager{
		Id:          pharmacyManagerId,
		Name:        "",
		Email:       "",
		Password:    "",
		PhoneNumber: "",
		Logo:        "",
		CreatedAt:   time.Time{},
		UpdatedAt:   time.Time{},
		DeletedAt:   sql.NullTime{},
	}

	var err error

	tx := extractTx(ctx)
	if tx != nil {
		err = tx.QueryRowContext(ctx, qFindPharmacyMangerById, pharmacyManagerId).Scan(
			&pm.Id, &pm.Name, &pm.Email, &pm.PhoneNumber, &pm.Logo,
		)
	} else {
		err = r.db.QueryRowContext(ctx, qFindPharmacyMangerById, pharmacyManagerId).Scan(
			&pm.Id, &pm.Name, &pm.Email, &pm.PhoneNumber, &pm.Logo,
		)
	}

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, custom_errors.NotFound(err)
		}
		return nil, err
	}

	return &pm, nil
}

func (r *PharmacyManagerRepositoryPostgres) UpdateOne(ctx context.Context, pharmacyManager entities.PharmacyManager) error {
	values := []interface{}{}
	values = append(values, pharmacyManager.Id)
	values = append(values, pharmacyManager.Name)
	values = append(values, pharmacyManager.PhoneNumber)

	var err error
	var stmt *sql.Stmt

	numberOfArgs := 4

	var sb strings.Builder
	sb.WriteString(qUpdatePharmacyManagerColl)
	if pharmacyManager.Logo != "" {
		sb.WriteString(`logo = `)
		sb.WriteString(fmt.Sprintf(`$%d ,`, numberOfArgs))
		values = append(values, pharmacyManager.Logo)
		numberOfArgs++
	}
	sb.WriteString(qUpdatePharmacyManagerCommand)

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

func (r *PharmacyManagerRepositoryPostgres) DeleteById(ctx context.Context, pharmacyManagerId int64) error {
	var err error
	var stmt *sql.Stmt

	tx := extractTx(ctx)
	if tx != nil {
		stmt, err = tx.PrepareContext(ctx, qDeletePharmacyManagerById)
	} else {
		stmt, err = r.db.PrepareContext(ctx, qDeletePharmacyManagerById)
	}

	if err != nil {
		return err
	}

	res, err := stmt.ExecContext(ctx, pharmacyManagerId)
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

func (r *PharmacyManagerRepositoryPostgres) CreateOnePharmacyManager(ctx context.Context, p entities.PharmacyManager) error {
	values := []interface{}{}
	values = append(values, p.Name)
	values = append(values, p.Email)
	values = append(values, p.Password)
	values = append(values, p.PhoneNumber)
	values = append(values, p.Logo)

	var err error

	tx := extractTx(ctx)
	if tx != nil {
		_, err = tx.ExecContext(ctx, qCreateOnePharmacyManager, values...)
	} else {
		_, err = r.db.ExecContext(ctx, qCreateOnePharmacyManager, values...)
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
