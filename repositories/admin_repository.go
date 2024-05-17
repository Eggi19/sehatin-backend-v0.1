package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/tsanaativa/sehatin-backend-v0.1/constants"
	"github.com/tsanaativa/sehatin-backend-v0.1/custom_errors"
	"github.com/tsanaativa/sehatin-backend-v0.1/entities"
	"github.com/jackc/pgx/v5/pgconn"
)

type AdminRepoOpt struct {
	Db *sql.DB
}

type AdminRepository interface {
	CreateOne(ctx context.Context, admin entities.Admin) error
	DeleteOne(ctx context.Context, id int64) error
	FindOneById(ctx context.Context, id int64) (*entities.Admin, error)
	FindAll(ctx context.Context, params entities.AdminParams) ([]entities.Admin, int, error)
	FindOneByEmail(ctx context.Context, email string) (*entities.Admin, error)
}

type AdminRepositoryPostgres struct {
	db *sql.DB
}

func NewAdminRepositoryPostgers(adminOpt *AdminRepoOpt) AdminRepository {
	return &AdminRepositoryPostgres{
		db: adminOpt.Db,
	}
}

func (r *AdminRepositoryPostgres) FindOneByEmail(ctx context.Context, email string) (*entities.Admin, error) {
	a := entities.Admin{}

	var err error

	tx := extractTx(ctx)
	if tx != nil {
		err = tx.QueryRowContext(ctx, qFindAdminByEmail, email).Scan(&a.Id, &a.Name, &a.Email, &a.Password)
	} else {
		err = r.db.QueryRowContext(ctx, qFindAdminByEmail, email).Scan(&a.Id, &a.Name, &a.Email, &a.Password)
	}

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, custom_errors.NotFound(err)
		}
		return nil, err
	}

	return &a, nil
}

func (r *AdminRepositoryPostgres) CreateOne(ctx context.Context, admin entities.Admin) error {
	values := []interface{}{}
	values = append(values, admin.Name)
	values = append(values, admin.Email)
	values = append(values, admin.Password)

	var err error

	tx := extractTx(ctx)
	if tx != nil {
		_, err = tx.ExecContext(ctx, qCreateOneAdmin, values...)
	} else {
		_, err = r.db.ExecContext(ctx, qCreateOneAdmin, values...)
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

func (r *AdminRepositoryPostgres) DeleteOne(ctx context.Context, id int64) error {
	var err error
	var stmt *sql.Stmt

	tx := extractTx(ctx)
	if tx != nil {
		stmt, err = tx.PrepareContext(ctx, qDeleteOneAdmin)
	} else {
		stmt, err = r.db.PrepareContext(ctx, qDeleteOneAdmin)
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

func (r *AdminRepositoryPostgres) FindOneById(ctx context.Context, id int64) (*entities.Admin, error) {
	admin := entities.Admin{}

	var err error

	tx := extractTx(ctx)
	if tx != nil {
		err = tx.QueryRowContext(ctx, qFindOneAdminById, id).Scan(
			&admin.Id, &admin.Name, &admin.Email,
		)
	} else {
		err = r.db.QueryRowContext(ctx, qFindOneAdminById, id).Scan(
			&admin.Id, &admin.Name, &admin.Email,
		)
	}

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, custom_errors.NotFound(err)
		}
		return nil, err
	}

	return &admin, nil
}

func (r *AdminRepositoryPostgres) FindAll(ctx context.Context, params entities.AdminParams) ([]entities.Admin, int, error) {
	admins := []entities.Admin{}

	var totalRows int

	var sb strings.Builder
	sb.WriteString(qCountTotalRows)
	sb.WriteString(qAdminColl)
	sb.WriteString(qAdminCommand)

	var sbTotalRows strings.Builder
	sbTotalRows.WriteString(qCountTotalRows)
	sbTotalRows.WriteString(qAdminCommand)

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
		params.Sort = `ASC `
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
		admin := entities.Admin{}
		err := rows.Scan(&totalRows, &admin.Id, &admin.Name, &admin.Email)
		if err != nil {
			return nil, 0, err
		}
		admins = append(admins, admin)
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

	return admins, totalRows, nil
}
