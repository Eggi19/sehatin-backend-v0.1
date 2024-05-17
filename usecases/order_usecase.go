package usecases

import (
	"context"
	"strings"
	"time"

	"github.com/tsanaativa/sehatin-backend-v0.1/constants"
	"github.com/tsanaativa/sehatin-backend-v0.1/custom_errors"
	"github.com/tsanaativa/sehatin-backend-v0.1/dtos"
	"github.com/tsanaativa/sehatin-backend-v0.1/entities"
	"github.com/tsanaativa/sehatin-backend-v0.1/repositories"
	"github.com/tsanaativa/sehatin-backend-v0.1/utils"
	"github.com/google/uuid"
)

type OrderUsecaseOpts struct {
	OrderRepository           repositories.OrderRepository
	CartRepository            repositories.CartRepository
	PharmacyProductRepository repositories.PharmacyProductRepository
	StockHistoryRepository    repositories.StockHistoryRepository
	Transactor                repositories.Transactor
	UploadFile                utils.FileUploader
}

type OrderUsecase interface {
	CreateNewOrder(ctx context.Context, req dtos.OrderRequest) (*entities.Order, error)
	CreateOrderWithTransaction(ctx context.Context, req dtos.OrderRequest) (*entities.Order, error)
	GetAllOrderByUser(ctx context.Context, userId int64, params entities.OrderParams) ([]dtos.OrderResponse, int, error)
	GetAllOrderByPharmacyManager(ctx context.Context, pharmacyManagerId int64, params entities.OrderParams) ([]dtos.OrderResponse, int, error)
	GetAllOrderByAdmin(ctx context.Context, params entities.OrderParams) ([]dtos.OrderResponse, int, error)
	UpdateOrderStatusToProcessing(ctx context.Context, orderId int64) error
	UpdateOrderStatusToShipped(ctx context.Context, orderId int64, pharmacyManagerId int64) error
	UpdateOrderStatusToCompleted(ctx context.Context, orderId int64, userId int64) error
	UploadPaymentProof(ctx context.Context, req dtos.UploadPaymentProofResponse, userId int64) error
	UpdateOrderStatusToCanceled(ctx context.Context, orderId int64, userId int64) error
	CancelOrderByAdmin(ctx context.Context, orderId int64) error
	CancelOrderByPharmacyManager(ctx context.Context, orderId int64, pharmacyManagerId int64) error
	GetOrderDetail(ctx context.Context, orderId int64) (*dtos.OrderResponse, error)
}

type OrderUsecaseImpl struct {
	OrderRepository           repositories.OrderRepository
	CartRepository            repositories.CartRepository
	PharmacyProductRepository repositories.PharmacyProductRepository
	StockHistoryRepository    repositories.StockHistoryRepository
	Transactor                repositories.Transactor
	UploadFile                utils.FileUploader
}

func NewOrderUsecaseImpl(oUseOpts *OrderUsecaseOpts) OrderUsecase {
	return &OrderUsecaseImpl{
		OrderRepository:           oUseOpts.OrderRepository,
		CartRepository:            oUseOpts.CartRepository,
		PharmacyProductRepository: oUseOpts.PharmacyProductRepository,
		StockHistoryRepository:    oUseOpts.StockHistoryRepository,
		Transactor:                oUseOpts.Transactor,
		UploadFile:                oUseOpts.UploadFile,
	}
}

func (u *OrderUsecaseImpl) CreateNewOrder(ctx context.Context, req dtos.OrderRequest) (*entities.Order, error) {
	timeNow := time.Now()
	cart, err := u.CartRepository.FindPharmacyIdByCartId(ctx, req.CartItemId[0])
	if err != nil {
		return nil, err
	}

	uuid, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}

	order := entities.Order{
		OrderNumber:     uuid.String(),
		TotalPrice:      req.TotalPrice,
		PaymentDeadline: timeNow.Add(time.Hour * 1),
		ShippingFee:     req.ShippingFee,
		ShippingMethod:  req.ShippingMethod,
		UserAddressId:   req.UserAddressId,
		OrderStatus:     constants.Pending,
		PharmacyId:      cart.PharmacyId,
	}
	newOrder, err := u.OrderRepository.CreateOrder(ctx, order)
	if err != nil {
		return nil, err
	}

	err = u.OrderRepository.CreateOrderItems(ctx, req.CartItemId, newOrder.Id)
	if err != nil {
		return nil, err
	}

	err = u.CartRepository.CartBulkDelete(ctx, req.CartItemId)
	if err != nil {
		return nil, err
	}

	return newOrder, nil
}

func (u *OrderUsecaseImpl) CreateOrderWithTransaction(ctx context.Context, req dtos.OrderRequest) (*entities.Order, error) {
	var order *entities.Order
	var err error

	_, err = u.Transactor.WithinTransaction(ctx, func(txCtx context.Context) (interface{}, error) {
		order, err = u.CreateNewOrder(txCtx, req)
		if err != nil {
			return nil, err
		}
		return nil, nil
	})
	if err != nil {
		return nil, err
	}

	orderItems, err := u.OrderRepository.GetOrderItems(ctx, order.Id)
	if err != nil {
		return nil, err
	}

	for _, item := range orderItems {
		_, err := u.Transactor.WithinTransaction(ctx, func(txCtx context.Context) (interface{}, error) {
			err := u.decreaseStock(txCtx, item.PharmacyProductId, item.Quantity)
			if err != nil {
				return nil, err
			}
			return nil, nil
		})

		if err != nil {
			return nil, err
		}
	}

	return order, nil
}

func (u *OrderUsecaseImpl) GetAllOrderByUser(ctx context.Context, userId int64, params entities.OrderParams) ([]dtos.OrderResponse, int, error) {
	var total int
	
	orders, err := u.OrderRepository.FindAllOrdersByUserId(ctx, userId, params)
	if err != nil {
		return nil, 0, err
	}
	if len(orders) > 0 {
		total = orders[0].Total
	}

	result := dtos.ConvertToOrderResponses(orders)

	for i, res := range result {
		orderItems, err := u.OrderRepository.GetOrderItems(ctx, res.Id)
		if err != nil {
			return nil, 0, err
		}

		result[i].OrderItems = append(res.OrderItems, orderItems...)
	}

	return result, total, nil
}

func (u *OrderUsecaseImpl) GetAllOrderByPharmacyManager(ctx context.Context, pharmacyManagerId int64, params entities.OrderParams) ([]dtos.OrderResponse, int, error) {
	var total int

	orders, err := u.OrderRepository.FindAllOrdersByPharmacyManagerId(ctx, pharmacyManagerId, params)
	if err != nil {
		return nil, 0, err
	}
	if len(orders) > 0 {
		total = orders[0].Total
	}

	result := dtos.ConvertToOrderResponses(orders)

	for i, res := range result {
		orderItems, err := u.OrderRepository.GetOrderItems(ctx, res.Id)
		if err != nil {
			return nil, 0, err
		}

		result[i].OrderItems = append(res.OrderItems, orderItems...)
	}

	return result, total, nil
}

func (u *OrderUsecaseImpl) GetAllOrderByAdmin(ctx context.Context, params entities.OrderParams) ([]dtos.OrderResponse, int, error) {
	var total int
	
	orders, err := u.OrderRepository.FindAllOrdersByAdmin(ctx, params)
	if err != nil {
		return nil, 0, err
	}
	if len(orders) > 0 {
		total = orders[0].Total
	}

	result := dtos.ConvertToOrderResponses(orders)

	for i, res := range result {
		orderItems, err := u.OrderRepository.GetOrderItems(ctx, res.Id)
		if err != nil {
			return nil, 0, err
		}

		result[i].OrderItems = append(res.OrderItems, orderItems...)
	}

	return result, total, nil
}

func (u *OrderUsecaseImpl) decreaseStock(ctx context.Context, pharmacyProductId int64, quantity int) error {
	err := u.PharmacyProductRepository.LockRow(ctx, pharmacyProductId)
	if err != nil {
		return err
	}

	pharmacyProduct, err := u.PharmacyProductRepository.GetOnePharmacyProduct(ctx, pharmacyProductId)
	if err != nil {
		return err
	}
	if quantity > pharmacyProduct.TotalStock {
		return custom_errors.NotEnoughStock()
	}

	err = u.PharmacyProductRepository.DecreaseStock(ctx, quantity, pharmacyProductId)
	if err != nil {
		return err
	}

	err = u.StockHistoryRepository.CreateOne(ctx, entities.StockHistory{
		Quantity:        0 - quantity,
		PharmacyProduct: entities.PharmacyProduct{Id: pharmacyProductId},
		Pharmacy:        entities.Pharmacy{Id: pharmacyProduct.Pharmacy.Id},
		Description:     "selling",
	})
	if err != nil {
		return err
	}

	return nil
}

func (u *OrderUsecaseImpl) increaseStock(ctx context.Context, pharmacyProductId int64, quantity int, description string) error {
	err := u.PharmacyProductRepository.LockRow(ctx, pharmacyProductId)
	if err != nil {
		return err
	}

	pharmacyProduct, err := u.PharmacyProductRepository.GetOnePharmacyProduct(ctx, pharmacyProductId)
	if err != nil {
		return err
	}

	err = u.PharmacyProductRepository.IncreaseStock(ctx, quantity, pharmacyProductId)
	if err != nil {
		return err
	}

	err = u.StockHistoryRepository.CreateOne(ctx, entities.StockHistory{
		Quantity:        quantity,
		PharmacyProduct: entities.PharmacyProduct{Id: pharmacyProductId},
		Pharmacy:        entities.Pharmacy{Id: pharmacyProduct.Pharmacy.Id},
		Description:     description,
	})
	if err != nil {
		return err
	}

	return nil
}

func (u *OrderUsecaseImpl) UpdateOrderStatusToProcessing(ctx context.Context, orderId int64) error {
	order, err := u.OrderRepository.GetOrder(ctx, orderId)
	if err != nil {
		return err
	}

	if order.OrderStatus != constants.Pending {
		return custom_errors.BadRequest(nil, constants.OrderStatusNotPendingErrMsg)
	}

	req := entities.UpdateOrderStatus{
		OrderId:           orderId,
		OrderStatus:       constants.Processing,
		PharmacyManagerId: 0,
		UserId:            0,
	}

	err = u.OrderRepository.UpdateOrderStatus(ctx, req)
	if err != nil {
		return err
	}

	return nil
}

func (u *OrderUsecaseImpl) UpdateOrderStatusToShipped(ctx context.Context, orderId int64, pharmacyManagerId int64) error {
	order, err := u.OrderRepository.GetOrder(ctx, orderId)
	if err != nil {
		return err
	}

	if order.OrderStatus != constants.Processing {
		return custom_errors.BadRequest(nil, constants.OrderStatusNotProcessingErrMsg)
	}

	req := entities.UpdateOrderStatus{
		OrderId:           orderId,
		OrderStatus:       constants.Shipped,
		PharmacyManagerId: pharmacyManagerId,
		UserId:            0,
	}

	err = u.OrderRepository.UpdateOrderStatus(ctx, req)
	if err != nil {
		return err
	}

	return nil
}

func (u *OrderUsecaseImpl) UpdateOrderStatusToCompleted(ctx context.Context, orderId int64, userId int64) error {
	order, err := u.OrderRepository.GetOrder(ctx, orderId)
	if err != nil {
		return err
	}

	if order.OrderStatus != constants.Shipped {
		return custom_errors.BadRequest(nil, constants.OrderStatusNotShippedErrMsg)
	}

	req := entities.UpdateOrderStatus{
		OrderId:           orderId,
		OrderStatus:       constants.Completed,
		PharmacyManagerId: 0,
		UserId:            userId,
	}

	err = u.OrderRepository.UpdateOrderStatus(ctx, req)
	if err != nil {
		return err
	}

	return nil
}

func (u *OrderUsecaseImpl) UploadPaymentProof(ctx context.Context, req dtos.UploadPaymentProofResponse, userId int64) error {
	convertedReq := entities.UploadPaymentProof{
		OrderId: req.OrderId,
		UserId:  userId,
	}

	order, err := u.OrderRepository.GetOrder(ctx, convertedReq.OrderId)
	if err != nil {
		return err
	}

	if order.OrderStatus != constants.Pending {
		return custom_errors.BadRequest(nil, constants.CannotUploadPaymentProofErrMsg)
	}

	file, _ := req.PaymentProof.Open()
	if req.PaymentProof.Size > 500000 {
		return custom_errors.FileTooLarge()
	}

	if strings.Split(req.PaymentProof.Filename, ".")[1] == "png" || strings.Split(req.PaymentProof.Filename, ".")[1] == "jpg" || strings.Split(req.PaymentProof.Filename, ".")[1] == "jpeg" {
		fileUrl, err := u.UploadFile.UploadFile(ctx, file)
		if err != nil {
			return custom_errors.UploadFile()
		}
		convertedReq.PaymentProof = fileUrl
	} else {
		return custom_errors.FileNotImage()
	}

	err = u.OrderRepository.UploadPaymentProof(ctx, convertedReq)
	if err != nil {
		return err
	}

	return nil
}

func (u *OrderUsecaseImpl) UpdateOrderStatusToCanceled(ctx context.Context, orderId int64, userId int64) error {
	order, err := u.OrderRepository.GetOrder(ctx, orderId)
	if err != nil {
		return err
	}

	if order.PaymentProof != nil {
		return custom_errors.BadRequest(nil, constants.CannotCancelOrderErrMsg)
	}

	orderItems, err := u.OrderRepository.GetOrderItems(ctx, orderId)
	if err != nil {
		return err
	}

	for _, item := range orderItems {
		_, err := u.Transactor.WithinTransaction(ctx, func(txCtx context.Context) (interface{}, error) {
			err := u.increaseStock(txCtx, item.PharmacyProductId, item.Quantity, "cancel order")
			if err != nil {
				return nil, err
			}
			return nil, nil
		})

		if err != nil {
			return err
		}
	}

	req := entities.UpdateOrderStatus{
		OrderId:           orderId,
		OrderStatus:       constants.Canceled,
		PharmacyManagerId: 0,
		UserId:            userId,
	}

	err = u.OrderRepository.UpdateOrderStatus(ctx, req)
	if err != nil {
		return err
	}

	return nil
}

func (u *OrderUsecaseImpl) CancelOrderByAdmin(ctx context.Context, orderId int64) error {
	order, err := u.OrderRepository.GetOrder(ctx, orderId)
	if err != nil {
		return err
	}

	if order.OrderStatus != constants.Pending {
		return custom_errors.BadRequest(nil, constants.CannotCancelOrderErrMsg)
	}

	orderItems, err := u.OrderRepository.GetOrderItems(ctx, orderId)
	if err != nil {
		return err
	}

	for _, item := range orderItems {
		_, err := u.Transactor.WithinTransaction(ctx, func(txCtx context.Context) (interface{}, error) {
			err := u.increaseStock(txCtx, item.PharmacyProductId, item.Quantity, "cancel order")
			if err != nil {
				return nil, err
			}
			return nil, nil
		})

		if err != nil {
			return err
		}
	}

	req := entities.UpdateOrderStatus{
		OrderId:           orderId,
		OrderStatus:       constants.Canceled,
		PharmacyManagerId: 0,
		UserId:            0,
	}

	err = u.OrderRepository.UpdateOrderStatus(ctx, req)
	if err != nil {
		return err
	}

	return nil
}

func (u *OrderUsecaseImpl) CancelOrderByPharmacyManager(ctx context.Context, orderId int64, pharmacyManagerId int64) error {
	order, err := u.OrderRepository.GetOrder(ctx, orderId)
	if err != nil {
		return err
	}

	if order.OrderStatus != constants.Processing {
		return custom_errors.BadRequest(nil, constants.CannotCancelOrderErrMsg)
	}

	orderItems, err := u.OrderRepository.GetOrderItems(ctx, orderId)
	if err != nil {
		return err
	}

	for _, item := range orderItems {
		_, err := u.Transactor.WithinTransaction(ctx, func(txCtx context.Context) (interface{}, error) {
			err := u.increaseStock(txCtx, item.PharmacyProductId, item.Quantity, "cancel order")
			if err != nil {
				return nil, err
			}
			return nil, nil
		})

		if err != nil {
			return err
		}
	}

	req := entities.UpdateOrderStatus{
		OrderId:           orderId,
		OrderStatus:       constants.Canceled,
		PharmacyManagerId: pharmacyManagerId,
		UserId:            0,
	}

	err = u.OrderRepository.UpdateOrderStatus(ctx, req)
	if err != nil {
		return err
	}

	return nil
}

func (u *OrderUsecaseImpl) GetOrderDetail(ctx context.Context, orderId int64) (*dtos.OrderResponse, error) {
	order, err := u.OrderRepository.FindOrderDetail(ctx, orderId)
	if err != nil {
		return nil, err
	}

	result := dtos.ConvertToOrderResponse(*order)

	orderItems, err := u.OrderRepository.GetOrderItems(ctx, result.Id)
	if err != nil {
		return nil, err
	}

	result.OrderItems = append(result.OrderItems, orderItems...)

	return result, nil
}
