package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"strings"

	"github.com/tsanaativa/sehatin-backend-v0.1/custom_errors"
	"github.com/tsanaativa/sehatin-backend-v0.1/entities"
)

type CartRepoOpts struct {
	Db *sql.DB
}

type CartRepository interface {
	CreateOneCartItem(ctx context.Context, req entities.CartItem) error
	IncreaseCartQuantity(ctx context.Context, req entities.CartItem) error
	DecreaseCartQuantity(ctx context.Context, req entities.CartItem) (*entities.CartItem, error)
	DeleteCartItem(ctx context.Context, cartId int64) error
	FindCartItem(ctx context.Context, req entities.CartItem) (*entities.CartItem, error)
	FindAllUserCartItem(ctx context.Context, userId int64) ([]entities.CartItem, error)
	FindPharmacyIdByCartId(ctx context.Context, id int64) (*entities.CartItem, error)
	CartBulkDelete(ctx context.Context, id []int64) error
}

type CartRepositoryPostgres struct {
	db *sql.DB
}

func NewCartRepositoryPostgres(cOpts *CartRepoOpts) CartRepository {
	return &CartRepositoryPostgres{
		db: cOpts.Db,
	}
}

func (r *CartRepositoryPostgres) CreateOneCartItem(ctx context.Context, req entities.CartItem) error {
	var err error

	tx := extractTx(ctx)
	if tx != nil {
		_, err = tx.ExecContext(ctx, qCreateCart, req.Quantity, req.UserId, req.PharmacyProductId)
	} else {
		_, err = r.db.ExecContext(ctx, qCreateCart, req.Quantity, req.UserId, req.PharmacyProductId)
	}

	if err != nil {
		return err
	}

	return nil
}

func (r *CartRepositoryPostgres) IncreaseCartQuantity(ctx context.Context, req entities.CartItem) error {
	var err error

	tx := extractTx(ctx)
	if tx != nil {
		_, err = tx.ExecContext(ctx, qIncreaseCartQuantity, req.Quantity, req.Id)
	} else {
		_, err = r.db.ExecContext(ctx, qIncreaseCartQuantity, req.Quantity, req.Id)
	}

	if err != nil {
		return err
	}

	return nil
}

func (r *CartRepositoryPostgres) DecreaseCartQuantity(ctx context.Context, req entities.CartItem) (*entities.CartItem, error) {
	var cart entities.CartItem
	var err error

	tx := extractTx(ctx)
	if tx != nil {
		err = tx.QueryRowContext(ctx, qDecreaseCartQuantity, req.Quantity, req.Id).Scan(&cart.Id, &cart.Quantity)
	} else {
		err = r.db.QueryRowContext(ctx, qDecreaseCartQuantity, req.Quantity, req.Id).Scan(&cart.Id, &cart.Quantity)
	}

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, custom_errors.NotFound(err)
		}
		return nil, err
	}

	return &cart, nil
}

func (r *CartRepositoryPostgres) DeleteCartItem(ctx context.Context, cartId int64) error {
	var err error

	tx := extractTx(ctx)
	if tx != nil {
		_, err = tx.ExecContext(ctx, qDeleteCartItem, cartId)
	} else {
		_, err = r.db.ExecContext(ctx, qDeleteCartItem, cartId)
	}

	if err != nil {
		return err
	}

	return nil
}

func (r *CartRepositoryPostgres) FindCartItem(ctx context.Context, req entities.CartItem) (*entities.CartItem, error) {
	c := entities.CartItem{}

	var err error

	tx := extractTx(ctx)
	if tx != nil {
		err = tx.QueryRowContext(ctx, qFindOneCartItem, req.UserId, req.PharmacyProductId).Scan(&c.Id, &c.Quantity)
	} else {
		err = r.db.QueryRowContext(ctx, qFindOneCartItem, req.UserId, req.PharmacyProductId).Scan(&c.Id, &c.Quantity)
	}

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, custom_errors.NotFound(err)
		}
		return nil, err
	}

	return &c, nil
}

func (r *CartRepositoryPostgres) FindAllUserCartItem(ctx context.Context, userId int64) ([]entities.CartItem, error) {
	carts := []entities.CartItem{}

	var err error
	var rows *sql.Rows

	tx := extractTx(ctx)
	if tx != nil {
		rows, err = tx.QueryContext(ctx, qFindCartItem, userId)
	} else {
		rows, err = r.db.QueryContext(ctx, qFindCartItem, userId)
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		cart := entities.CartItem{}
		err := rows.Scan(&cart.Id, &cart.UpdatedAt, &cart.ProductName, &cart.ProductPicture, &cart.SellingUnit, &cart.Price, &cart.Quantity, &cart.PharmacyName, &cart.PharmacyId, &cart.SlugId, &cart.TotalStock, &cart.IsAvailable, &cart.Weight)
		if err != nil {
			return nil, err
		}

		carts = append(carts, cart)
	}

	return carts, nil
}

func (r *CartRepositoryPostgres) FindPharmacyIdByCartId(ctx context.Context, id int64) (*entities.CartItem, error) {
	c := entities.CartItem{}

	var err error

	tx := extractTx(ctx)
	if tx != nil {
		err = tx.QueryRowContext(ctx, qFindPharmacyIdByCartId, id).Scan(&c.PharmacyId, &c.Quantity, &c.PharmacyProductId)
	} else {
		err = r.db.QueryRowContext(ctx, qFindPharmacyIdByCartId, id).Scan(&c.PharmacyId, &c.Quantity, &c.PharmacyProductId)
	}

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, custom_errors.NotFound(err)
		}
		return nil, err
	}

	return &c, nil
}

func (r *CartRepositoryPostgres) CartBulkDelete(ctx context.Context, id []int64) error {
	valueStrings := make([]string, 0, len(id))
	valueArgs := make([]interface{}, 0, len(id)*1)
	i := 0

	for _, id := range id {
		valueStrings = append(valueStrings, fmt.Sprintf("($%d)", i*1+1))
		valueArgs = append(valueArgs, strconv.Itoa(int(id)))
		i++
	}

	stmt := fmt.Sprintf(qCartItemsBulkDelete, strings.Join(valueStrings, ","))

	var err error

	tx := extractTx(ctx)
	if tx != nil {
		_, err = tx.ExecContext(ctx, stmt, valueArgs...)
	} else {
		_, err = r.db.ExecContext(ctx, stmt, valueArgs...)
	}

	if err != nil {
		return err
	}

	return nil
}
