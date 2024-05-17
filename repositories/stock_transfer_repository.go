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

type StockTransferRepoOpts struct {
	Db *sql.DB
}

type StockTransferRepository interface {
	CreateOne(ctx context.Context, st entities.StockTransfer) (*entities.StockTransfer, error)
	FindAll(ctx context.Context, params entities.StockTransferParams) ([]entities.StockTransfer, int, error)
	UpdateMutationStatus(ctx context.Context, id, mutationStatusId int64) error
	FindOneById(ctx context.Context, id int64) (*entities.StockTransfer, error)
}

type StockTransferRepositoryPostgres struct {
	db *sql.DB
}

func NewStockTransferRepositoryPostgres(strOpts *StockTransferRepoOpts) StockTransferRepository {
	return &StockTransferRepositoryPostgres{
		db: strOpts.Db,
	}
}

func (r *StockTransferRepositoryPostgres) CreateOne(ctx context.Context, st entities.StockTransfer) (*entities.StockTransfer, error) {
	newSt := entities.StockTransfer{
		PharmacySender: entities.Pharmacy{
			PharmacyManager:           entities.PharmacyManager{},
			PharmacyAddress:           entities.PharmacyAddress{},
			OfficialShippingMethod:    []entities.OfficialShippingMethod{},
			NonOfficialShippingMethod: []entities.NonOfficialShippingMethod{},
			CreatedAt:                 time.Time{},
			UpdatedAt:                 time.Time{},
			DeletedAt:                 sql.NullTime{},
		},
		PharmacyReceiver: entities.Pharmacy{
			Id:                        0,
			PharmacyManager:           entities.PharmacyManager{},
			PharmacyAddress:           entities.PharmacyAddress{},
			OfficialShippingMethod:    []entities.OfficialShippingMethod{},
			NonOfficialShippingMethod: []entities.NonOfficialShippingMethod{},
		},
		MutationStatus: entities.MutationSatus{},
		Product: entities.Product{
			ProductForm:           entities.ProductForm{},
			ProductClassification: entities.ProductClassification{},
			Manufacture:           entities.Manufacture{},
			Categories:            []entities.Category{},
		},
		Quantity: 0,
	}

	values := []interface{}{}
	values = append(values, st.PharmacySender.Id)
	values = append(values, st.PharmacyReceiver.Id)
	values = append(values, st.MutationStatus.Id)
	values = append(values, st.Product.Id)
	values = append(values, st.Quantity)

	var err error

	tx := extractTx(ctx)
	if tx != nil {
		err = tx.QueryRowContext(ctx, qCreateOneStockTransfer, values...).Scan(
			&newSt.Id, &newSt.PharmacySender.Id, &newSt.PharmacyReceiver.Id, &newSt.MutationStatus.Id, &newSt.Product.Id, &newSt.Quantity,
		)
	} else {
		err = r.db.QueryRowContext(ctx, qCreateOneStockTransfer, values...).Scan(
			&newSt.Id, &newSt.PharmacySender.Id, &newSt.PharmacyReceiver.Id, &newSt.MutationStatus.Id, &newSt.Product.Id, &newSt.Quantity,
		)
	}

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == constants.ViolatesUniqueConstraintPgErrCode {
			return nil, custom_errors.BadRequest(err, constants.PharmacyNotUniqueErrMg)
		}
		return nil, err
	}

	return &newSt, nil
}

func (r *StockTransferRepositoryPostgres) FindAll(ctx context.Context, params entities.StockTransferParams) ([]entities.StockTransfer, int, error) {
	stockTransfers := []entities.StockTransfer{}
	var totalRows int

	var sb strings.Builder
	sb.WriteString(qCountTotalRows)
	sb.WriteString(qStockTrasnferColl)
	sb.WriteString(qStockTrasnferCommands)

	var sbTotalRows strings.Builder
	sbTotalRows.WriteString(qCountTotalRows)
	sbTotalRows.WriteString(qPharmacyCommands)

	values := []interface{}{}
	valuesTotalArgument := []interface{}{}

	numberOfArgs := 1

	var sortBy string
	switch params.SortBy {
	case "updated":
		sortBy = `st.updated_at `
	default:
		sortBy = `st.id `
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
		st := entities.StockTransfer{}
		err := rows.Scan(&totalRows,
			&st.Id, &st.PharmacySender.Id, &st.PharmacySender.Name, &st.PharmacyReceiver.Id, &st.PharmacyReceiver.Name,
			&st.MutationStatus.Id, &st.MutationStatus.Name, &st.Product.Id, &st.Product.Name, &st.Quantity, &st.UpdatedAt,
		)
		if err != nil {
			return nil, 0, err
		}
		stockTransfers = append(stockTransfers, st)
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

	return stockTransfers, totalRows, nil
}

func (r *StockTransferRepositoryPostgres) UpdateMutationStatus(ctx context.Context, id, mutationStatusId int64) error {
	var err error
	var stmt *sql.Stmt

	tx := extractTx(ctx)
	if tx != nil {
		stmt, err = tx.PrepareContext(ctx, qUpdateOneMutationStatusId)
	} else {
		stmt, err = r.db.PrepareContext(ctx, qUpdateOneMutationStatusId)
	}

	if err != nil {
		return err
	}

	res, err := stmt.ExecContext(ctx, id, mutationStatusId)
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

func (r *StockTransferRepositoryPostgres) FindOneById(ctx context.Context, id int64) (*entities.StockTransfer, error) {
	st := entities.StockTransfer{
		PharmacySender:   entities.Pharmacy{PharmacyManager: entities.PharmacyManager{}},
		PharmacyReceiver: entities.Pharmacy{PharmacyManager: entities.PharmacyManager{}},
		MutationStatus:   entities.MutationSatus{},
		Product:          entities.Product{},
	}

	var err error

	tx := extractTx(ctx)
	if tx != nil {
		err = tx.QueryRowContext(ctx, qFindOneStockTransfer, id).Scan(
			&st.Id, &st.PharmacySender.Id, &st.PharmacyReceiver.Id, &st.MutationStatus.Id, &st.Product.Id, &st.Quantity,
		)
	} else {
		err = r.db.QueryRowContext(ctx, qFindOneStockTransfer, id).Scan(
			&st.Id, &st.PharmacySender.Id, &st.PharmacyReceiver.Id, &st.MutationStatus.Id, &st.Product.Id, &st.Quantity,
		)
	}

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, custom_errors.NotFound(err)
		}
		return nil, err
	}

	return &st, err
}
