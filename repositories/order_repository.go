package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/tsanaativa/sehatin-backend-v0.1/custom_errors"
	"github.com/tsanaativa/sehatin-backend-v0.1/dtos"
	"github.com/tsanaativa/sehatin-backend-v0.1/entities"
)

type OrderRepoOpts struct {
	Db *sql.DB
}

type OrderRepository interface {
	CreateOrder(ctx context.Context, req entities.Order) (*entities.Order, error)
	CreateOrderItems(ctx context.Context, cartId []int64, orderId int64) error
	FindAllOrdersByUserId(ctx context.Context, userId int64, params entities.OrderParams) ([]entities.Order, error)
	GetOrderItems(ctx context.Context, orderId int64) ([]dtos.OrderItem, error)
	FindAllOrdersByPharmacyManagerId(ctx context.Context, pharmacyManagerId int64, params entities.OrderParams) ([]entities.Order, error)
	FindAllOrdersByAdmin(ctx context.Context, params entities.OrderParams) ([]entities.Order, error)
	UpdateOrderStatus(ctx context.Context, req entities.UpdateOrderStatus) error
	GetOrder(ctx context.Context, orderId int64) (*entities.Order, error)
	UploadPaymentProof(ctx context.Context, req entities.UploadPaymentProof) error
	FindOrderDetail(ctx context.Context, orderId int64) (*entities.Order, error)
}

type OrderRepositoryPostgres struct {
	db *sql.DB
}

func NewOrderRepositoryPostgres(oOpt *OrderRepoOpts) OrderRepository {
	return &OrderRepositoryPostgres{
		db: oOpt.Db,
	}
}

func (r *OrderRepositoryPostgres) CreateOrder(ctx context.Context, req entities.Order) (*entities.Order, error) {
	o := entities.Order{}

	var err error

	tx := extractTx(ctx)
	if tx != nil {
		err = tx.QueryRowContext(ctx, qCreateOrder, req.OrderNumber, req.TotalPrice, req.PaymentDeadline, req.ShippingFee, req.ShippingMethod, req.UserAddressId, req.OrderStatus, req.PharmacyId).Scan(&o.Id, &o.PaymentDeadline)
	} else {
		err = r.db.QueryRowContext(ctx, qCreateOrder, req.OrderNumber, req.TotalPrice, req.PaymentDeadline, req.ShippingFee, req.ShippingMethod, req.UserAddressId, req.OrderStatus, req.PharmacyId).Scan(&o.Id, &o.PaymentDeadline)
	}

	if err != nil {
		return nil, err
	}

	return &o, nil
}

func (r *OrderRepositoryPostgres) CreateOrderItems(ctx context.Context, cartId []int64, orderId int64) error {
	valueStrings := make([]string, 0, len(cartId))
	valueArgs := make([]interface{}, 0, len(cartId)*3)
	i := 0

	for _, id := range cartId {
		valueStrings = append(valueStrings, fmt.Sprintf("((select quantity from cart_items where id = $%d), $%d, (select pharmacy_product_id from cart_items where id = $%d))", i*3+1, i*3+2, i*3+3))
		valueArgs = append(valueArgs, id)
		valueArgs = append(valueArgs, orderId)
		valueArgs = append(valueArgs, id)
		i++
	}
	stmt := fmt.Sprintf(qCreateOrderItem, strings.Join(valueStrings, ","))

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

func (r *OrderRepositoryPostgres) FindAllOrdersByUserId(ctx context.Context, userId int64, params entities.OrderParams) ([]entities.Order, error) {
	orders := []entities.Order{}

	var err error
	var rows *sql.Rows

	page := params.Limit * (params.Page - 1)

	tx := extractTx(ctx)
	if tx != nil {
		rows, err = tx.QueryContext(ctx, qFindUserOrder, userId, params.Status, page, params.Limit)
	} else {
		rows, err = r.db.QueryContext(ctx, qFindUserOrder, userId, params.Status, page, params.Limit)
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		order := entities.Order{}
		err := rows.Scan(&order.Id, &order.OrderNumber, &order.TotalPrice, &order.PaymentProof, &order.PaymentDeadline, &order.ShippingFee, &order.ShippingMethod, &order.OrderStatus, &order.PharmacyName, &order.UserAddress.City, &order.UserAddress.Province, &order.UserAddress.Address, &order.UserAddress.District, &order.UserAddress.SubDistrict, &order.UserAddress.PostalCode, &order.Total)
		if err != nil {
			return nil, err
		}

		orders = append(orders, order)
	}

	return orders, nil
}

func (r *OrderRepositoryPostgres) GetOrderItems(ctx context.Context, orderId int64) ([]dtos.OrderItem, error) {
	orderItems := []dtos.OrderItem{}

	var err error
	var rows *sql.Rows

	tx := extractTx(ctx)
	if tx != nil {
		rows, err = tx.QueryContext(ctx, qFindOrderItems, orderId)
	} else {
		rows, err = r.db.QueryContext(ctx, qFindOrderItems, orderId)
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		orderItem := dtos.OrderItem{}
		err := rows.Scan(&orderItem.Name, &orderItem.SellingUnit, &orderItem.Price, &orderItem.Quantity, &orderItem.ProductPicture, &orderItem.PharmacyProductId)
		if err != nil {
			return nil, err
		}

		orderItems = append(orderItems, orderItem)
	}

	return orderItems, nil
}

func (r *OrderRepositoryPostgres) FindAllOrdersByPharmacyManagerId(ctx context.Context, pharmacyManagerId int64, params entities.OrderParams) ([]entities.Order, error) {
	orders := []entities.Order{}

	var err error
	var rows *sql.Rows
	var q string

	page := params.Limit * (params.Page - 1)

	if params.PharmacyId == 0 {
		q = fmt.Sprintf(qFindUserOrderByPharmacyManager, "")
	} else {
		stmnt := fmt.Sprintf("AND p.id = %d", params.PharmacyId)
		q = fmt.Sprintf(qFindUserOrderByPharmacyManager, stmnt)
	}

	tx := extractTx(ctx)
	if tx != nil {
		rows, err = tx.QueryContext(ctx, q, pharmacyManagerId, params.Status, page, params.Limit)
	} else {
		rows, err = r.db.QueryContext(ctx, q, pharmacyManagerId, params.Status, page, params.Limit)
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		order := entities.Order{}
		err := rows.Scan(&order.Id, &order.OrderNumber, &order.TotalPrice, &order.PaymentProof, &order.PaymentDeadline, &order.ShippingFee, &order.ShippingMethod, &order.OrderStatus, &order.PharmacyName, &order.UserAddress.City, &order.UserAddress.Province, &order.UserAddress.Address, &order.UserAddress.District, &order.UserAddress.SubDistrict, &order.UserAddress.PostalCode, &order.PharmacyAddress.City, &order.PharmacyAddress.Province, &order.PharmacyAddress.Address, &order.PharmacyAddress.District, &order.PharmacyAddress.SubDistrict, &order.PharmacyAddress.PostalCode, &order.Total, &order.UserName, &order.UserEmail)
		if err != nil {
			return nil, err
		}

		orders = append(orders, order)
	}

	return orders, nil
}

func (r *OrderRepositoryPostgres) FindAllOrdersByAdmin(ctx context.Context, params entities.OrderParams) ([]entities.Order, error) {
	orders := []entities.Order{}

	var err error
	var rows *sql.Rows
	var q string

	page := params.Limit * (params.Page - 1)

	if params.PharmacyId == 0 {
		q = fmt.Sprintf(qFindUserOrderByAdmin, "")
	} else {
		stmnt := fmt.Sprintf("AND p.id = %d", params.PharmacyId)
		q = fmt.Sprintf(qFindUserOrderByAdmin, stmnt)
	}

	tx := extractTx(ctx)
	if tx != nil {
		rows, err = tx.QueryContext(ctx, q, params.Status, page, params.Limit)
	} else {
		rows, err = r.db.QueryContext(ctx, q, params.Status, page, params.Limit)
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		order := entities.Order{}
		err := rows.Scan(&order.Id, &order.OrderNumber, &order.TotalPrice, &order.PaymentProof, &order.PaymentDeadline, &order.ShippingFee, &order.ShippingMethod, &order.OrderStatus, &order.PharmacyName, &order.UserAddress.City, &order.UserAddress.Province, &order.UserAddress.Address, &order.UserAddress.District, &order.UserAddress.SubDistrict, &order.UserAddress.PostalCode, &order.PharmacyAddress.City, &order.PharmacyAddress.Province, &order.PharmacyAddress.Address, &order.PharmacyAddress.District, &order.PharmacyAddress.SubDistrict, &order.PharmacyAddress.PostalCode, &order.Total, &order.UserName, &order.UserEmail)
		if err != nil {
			return nil, err
		}

		orders = append(orders, order)
	}

	return orders, nil
}

func (r *OrderRepositoryPostgres) UpdateOrderStatus(ctx context.Context, req entities.UpdateOrderStatus) error {
	var err error
	var res sql.Result
	var q string

	if req.PharmacyManagerId != 0 {
		stmnt := fmt.Sprintf("AND pharmacies.id = orders.pharmacy_id AND pharmacies.pharmacy_manager_id = %d", req.PharmacyManagerId)
		q = fmt.Sprintf(qUpdateOrderStatus, stmnt)
	} else if req.UserId != 0 {
		stmnt := fmt.Sprintf("AND user_addresses.user_id = orders.user_address_id AND user_addresses.user_id = %d", req.UserId)
		q = fmt.Sprintf(qUpdateOrderStatus, stmnt)
	} else {
		q = fmt.Sprintf(qUpdateOrderStatus, "")
	}

	tx := extractTx(ctx)
	if tx != nil {
		res, err = tx.ExecContext(ctx, q, req.OrderStatus, req.OrderId)
	} else {
		res, err = r.db.ExecContext(ctx, q, req.OrderStatus, req.OrderId)
	}
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return custom_errors.NotFound(nil)
	}

	return nil
}

func (r *OrderRepositoryPostgres) GetOrder(ctx context.Context, orderId int64) (*entities.Order, error) {
	o := entities.Order{}

	var err error

	tx := extractTx(ctx)
	if tx != nil {
		err = tx.QueryRowContext(ctx, qFindOrderStatus, orderId).Scan(&o.OrderStatus, &o.PaymentProof, &o.PharmacyId)
	} else {
		err = r.db.QueryRowContext(ctx, qFindOrderStatus, orderId).Scan(&o.OrderStatus, &o.PaymentProof, &o.PharmacyId)
	}

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, custom_errors.NotFound(err)
		}
		return nil, err
	}

	return &o, nil
}

func (r *OrderRepositoryPostgres) UploadPaymentProof(ctx context.Context, req entities.UploadPaymentProof) error {
	var err error
	var res sql.Result

	tx := extractTx(ctx)
	if tx != nil {
		res, err = tx.ExecContext(ctx, qUploadPaymentProof, req.PaymentProof, req.OrderId, req.UserId)
	} else {
		res, err = r.db.ExecContext(ctx, qUploadPaymentProof, req.PaymentProof, req.OrderId, req.UserId)
	}
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return custom_errors.NotFound(nil)
	}

	return nil
}

func (r *OrderRepositoryPostgres) FindOrderDetail(ctx context.Context, orderId int64) (*entities.Order, error) {
	order := entities.Order{}

	var err error

	tx := extractTx(ctx)
	if tx != nil {
		err = tx.QueryRowContext(ctx, qFindOrderDetail, orderId).Scan(&order.Id, &order.OrderNumber, &order.TotalPrice, &order.PaymentProof, &order.PaymentDeadline, &order.ShippingFee, &order.ShippingMethod, &order.OrderStatus, &order.PharmacyName, &order.UserAddress.City, &order.UserAddress.Province, &order.UserAddress.Address, &order.UserAddress.District, &order.UserAddress.SubDistrict, &order.UserAddress.PostalCode)
	} else {
		err = r.db.QueryRowContext(ctx, qFindOrderDetail, orderId).Scan(&order.Id, &order.OrderNumber, &order.TotalPrice, &order.PaymentProof, &order.PaymentDeadline, &order.ShippingFee, &order.ShippingMethod, &order.OrderStatus, &order.PharmacyName, &order.UserAddress.City, &order.UserAddress.Province, &order.UserAddress.Address, &order.UserAddress.District, &order.UserAddress.SubDistrict, &order.UserAddress.PostalCode)
	}

	if err != nil {
		return nil, err
	}

	return &order, nil
}
