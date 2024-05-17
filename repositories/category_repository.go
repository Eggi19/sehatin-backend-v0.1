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

type CategoryRepoOpts struct {
	Db *sql.DB
}

type CategoryRepository interface {
	FindAll(ctx context.Context, params entities.CategoryParams) ([]entities.Category, int, error)
	FindById(ctx context.Context, categoryId int64) (*entities.Category, error)
	CreateOne(ctx context.Context, category entities.Category) error
	UpdateOne(ctx context.Context, category entities.Category) error
	DeleteOne(ctx context.Context, categoryId int64) error
	GetProductCategory(ctx context.Context, productId int64) ([]entities.Category, error)
}

type CategoryRepositoryPostgres struct {
	db *sql.DB
}

func NewCategoryRepositoryPostgres(cOpts *CategoryRepoOpts) CategoryRepository {
	return &CategoryRepositoryPostgres{
		db: cOpts.Db,
	}
}

func (r *CategoryRepositoryPostgres) FindAll(ctx context.Context, params entities.CategoryParams) ([]entities.Category, int, error) {
	categories := []entities.Category{}

	var totalRows int

	var sb strings.Builder
	sb.WriteString(qCountTotalRows)
	sb.WriteString(qCategoryColl)
	sb.WriteString(qCategoryCommands)

	var sbTotalRows strings.Builder
	sbTotalRows.WriteString(qCountTotalRows)
	sbTotalRows.WriteString(qCategoryCommands)

	values := []interface{}{}
	valuesCountTotal := []interface{}{}

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

	if params.Limit != 0 {
		sb.WriteString(`LIMIT `)
		sb.WriteString(fmt.Sprintf(`$%d `, numberOfArgs))
		values = append(values, params.Limit)
		numberOfArgs++
	}

	if params.Page != 0 {
		sb.WriteString(`OFFSET `)
		sb.WriteString(fmt.Sprintf(`$%d`, numberOfArgs))
		values = append(values, params.Limit*(params.Page-1))
		numberOfArgs++
	}

	rows, err := r.db.QueryContext(ctx, sb.String(), values...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	for rows.Next() {
		category := entities.Category{}

		err := rows.Scan(&totalRows, &category.Id, &category.Name)
		if err != nil {
			if err == sql.ErrNoRows {
				return nil, 0, custom_errors.NotFound(err)
			}
			return nil, 0, err
		}
		categories = append(categories, category)
	}

	if totalRows == 0 {
		err := r.db.QueryRowContext(ctx, sbTotalRows.String(), valuesCountTotal...).Scan(&totalRows)
		if err != nil {
			if err == sql.ErrNoRows {
				return nil, 0, custom_errors.NotFound(err)
			}
			return nil, 0, err
		}
	}

	return categories, totalRows, nil
}

func (r *CategoryRepositoryPostgres) FindById(ctx context.Context, categoryId int64) (*entities.Category, error) {
	c := entities.Category{}

	var err error

	tx := extractTx(ctx)
	if tx != nil {
		err = tx.QueryRowContext(ctx, qFindOneCategoryById, categoryId).Scan(&c.Id, &c.Name)
	} else {
		err = r.db.QueryRowContext(ctx, qFindOneCategoryById, categoryId).Scan(&c.Id, &c.Name)
	}

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, custom_errors.NotFound(err)
		}
		return nil, err
	}

	return &c, nil
}

func (r *CategoryRepositoryPostgres) CreateOne(ctx context.Context, category entities.Category) error {
	var err error

	tx := extractTx(ctx)
	if tx != nil {
		err = tx.QueryRowContext(ctx, qCreateOneCategory, category.Name).Scan(&category.Id)
	} else {
		err = r.db.QueryRowContext(ctx, qCreateOneCategory, category.Name).Scan(&category.Id)
	}

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == constants.ViolatesUniqueConstraintPgErrCode {
			return custom_errors.BadRequest(err, constants.CategoryNameNotUniqueErrMsg)
		}
		return err
	}

	return nil
}

func (r *CategoryRepositoryPostgres) UpdateOne(ctx context.Context, category entities.Category) error {
	var err error
	var stmt *sql.Stmt

	tx := extractTx(ctx)
	if tx != nil {
		stmt, err = tx.PrepareContext(ctx, qUpdateOneCategory)
	} else {
		stmt, err = r.db.PrepareContext(ctx, qUpdateOneCategory)
	}

	if err != nil {
		return err
	}

	res, err := stmt.ExecContext(ctx, category.Id, category.Name)
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

func (r *CategoryRepositoryPostgres) DeleteOne(ctx context.Context, categoryId int64) error {
	var err error
	var stmt *sql.Stmt

	tx := extractTx(ctx)
	if tx != nil {
		stmt, err = tx.PrepareContext(ctx, qDeleteOneCategory)
	} else {
		stmt, err = r.db.PrepareContext(ctx, qDeleteOneCategory)
	}

	if err != nil {
		return err
	}

	res, err := stmt.ExecContext(ctx, categoryId)
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

func (r *CategoryRepositoryPostgres) GetProductCategory(ctx context.Context, productId int64) ([]entities.Category, error) {
	categories := []entities.Category{}

	var err error
	var rows *sql.Rows

	tx := extractTx(ctx)
	if tx != nil {
		rows, err = tx.QueryContext(ctx, qFindProductCategories, productId)
	} else {
		rows, err = r.db.QueryContext(ctx, qFindProductCategories, productId)
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		category := entities.Category{}
		err := rows.Scan(&category.Id, &category.Name)
		if err != nil {
			return nil, err
		}

		categories = append(categories, category)
	}

	return categories, nil
}
