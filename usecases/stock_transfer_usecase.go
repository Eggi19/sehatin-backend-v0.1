package usecases

import (
	"context"
	"fmt"

	"github.com/tsanaativa/sehatin-backend-v0.1/constants"
	"github.com/tsanaativa/sehatin-backend-v0.1/custom_errors"
	"github.com/tsanaativa/sehatin-backend-v0.1/entities"
	"github.com/tsanaativa/sehatin-backend-v0.1/repositories"
)

type StockTransferUsecaseOpts struct {
	StockTransferRepo   repositories.StockTransferRepository
	PharmacyRepo        repositories.PharmacyRepository
	PharmacyProductRepo repositories.PharmacyProductRepository
	Transactor          repositories.Transactor
	StockHistoryRepo    repositories.StockHistoryRepository
}

type StockTransferUsecase interface {
	CreateStockRequest(ctx context.Context, st entities.StockTransfer) error
	GetAllStockTransfer(ctx context.Context, params entities.StockTransferParams) ([]entities.StockTransfer, *entities.PaginationInfo, error)
	UpdateStatusProcessed(ctx context.Context, stockTransferId, mutationStatusId int64) error
	UpdateStatusProcessedWithTransaction(ctx context.Context, stockTransferId, mutationStatusId int64) error
	UpdateStatusCanceled(ctx context.Context, stockTransferId, mutationStatusId int64) error
	UpdateStatusCanceledWithTransaction(ctx context.Context, stockTransferId, mutationStatusId int64) error
	UpdateStatusPending(ctx context.Context, stockTransferId, mutationStatusId int64) error
	UpdateStatusPendingWithTransaction(ctx context.Context, stockTransferId, mutationStatusId int64) error
}

type StockTransferUsecaseImpl struct {
	StockTransferRepository   repositories.StockTransferRepository
	PharmacyRepository        repositories.PharmacyRepository
	PharmacyProductRepository repositories.PharmacyProductRepository
	Transactor                repositories.Transactor
	StockHistoryRepository    repositories.StockHistoryRepository
}

func NewStockTransferUsecaseImpl(stuOpts *StockTransferUsecaseOpts) StockTransferUsecase {
	return &StockTransferUsecaseImpl{
		StockTransferRepository:   stuOpts.StockTransferRepo,
		PharmacyRepository:        stuOpts.PharmacyRepo,
		PharmacyProductRepository: stuOpts.PharmacyProductRepo,
		Transactor:                stuOpts.Transactor,
		StockHistoryRepository:    stuOpts.StockHistoryRepo,
	}
}

func (u *StockTransferUsecaseImpl) CreateStockRequest(ctx context.Context, st entities.StockTransfer) error {
	var err error

	senderPharmacy, err := u.PharmacyRepository.FindOneById(ctx, st.PharmacySender.Id)
	if err != nil {
		if senderPharmacy == nil {
			return custom_errors.BadRequest(err, "pharmacy not found")
		}
		return err
	}

	receiverPharmacy, err := u.PharmacyRepository.FindOneById(ctx, st.PharmacyReceiver.Id)
	if err != nil {
		if receiverPharmacy == nil {
			return custom_errors.BadRequest(err, "pharmacy not found")
		}
		return err
	}

	if senderPharmacy.PharmacyManager.Id != receiverPharmacy.PharmacyManager.Id {
		return custom_errors.BadRequest(err, constants.PharmacyManagerNotMatchErrMsg)
	}

	senderProduct, err := u.PharmacyProductRepository.FindOneByPharmacyAndProductId(ctx, senderPharmacy.Id, st.Product.Id)
	if err != nil {
		if senderProduct == nil {
			return custom_errors.BadRequest(err, "product not found in sender pharmacy")
		}
		return err
	}

	receiverProduct, err := u.PharmacyProductRepository.FindOneByPharmacyAndProductId(ctx, receiverPharmacy.Id, st.Product.Id)
	if err != nil {
		if receiverProduct == nil {
			return custom_errors.BadRequest(err, "product not found in receiver pharmacy")
		}
		return err
	}

	if senderProduct.TotalStock == 0 || senderProduct.TotalStock < st.Quantity || !senderProduct.IsAvailable {
		return custom_errors.BadRequest(err, "product out of stock or not available")
	}

	_, err = u.StockTransferRepository.CreateOne(ctx, st)
	if err != nil {
		return err
	}

	return nil
}

func (u *StockTransferUsecaseImpl) GetAllStockTransfer(ctx context.Context, params entities.StockTransferParams) ([]entities.StockTransfer, *entities.PaginationInfo, error) {
	stockTransfers, totalData, err := u.StockTransferRepository.FindAll(ctx, params)
	if err != nil {
		return nil, nil, err
	}

	totalPage := totalData / params.Limit
	if totalData%params.Limit > 0 {
		totalPage++
	}

	pagination := entities.PaginationInfo{
		Page:      params.Page,
		Limit:     params.Limit,
		TotalData: totalData,
		TotalPage: totalPage,
	}

	return stockTransfers, &pagination, nil
}

func (u *StockTransferUsecaseImpl) UpdateStatusProcessed(ctx context.Context, stockTransferId int64, mutationStatusId int64) error {
	st, err := u.StockTransferRepository.FindOneById(ctx, stockTransferId)
	if err != nil {
		return err
	}

	sender, _ := u.PharmacyProductRepository.FindOneByPharmacyAndProductId(ctx, st.PharmacySender.Id, st.Product.Id)
	receiver, _ := u.PharmacyProductRepository.FindOneByPharmacyAndProductId(ctx, st.PharmacyReceiver.Id, st.Product.Id)
	product := st.Product

	if sender == nil || receiver == nil {
		return custom_errors.BadRequest(err, "pharmacy not found")
	}

	if sender.TotalStock == 0 || sender.TotalStock < st.Quantity {
		return custom_errors.BadRequest(err, "product out of stock")
	}

	senderChangeStock := sender.TotalStock - st.Quantity
	receiverChangeStock := receiver.TotalStock + st.Quantity

	err = u.PharmacyProductRepository.UpdateTotalStock(ctx, sender.Pharmacy.Id, product.Id, senderChangeStock)
	if err != nil {
		return err
	}

	err = u.PharmacyProductRepository.UpdateTotalStock(ctx, receiver.Pharmacy.Id, product.Id, receiverChangeStock)
	if err != nil {
		return err
	}

	senderStockHistory := entities.StockHistory{
		PharmacyProduct: entities.PharmacyProduct{Id: sender.Id},
		Pharmacy:        entities.Pharmacy{Id: sender.Pharmacy.Id},
		Quantity:        -st.Quantity,
		Description:     fmt.Sprintf(`sent %d %s to %s`, st.Quantity, sender.Product.Name, receiver.Pharmacy.Name),
	}

	err = u.StockHistoryRepository.CreateOne(ctx, senderStockHistory)
	if err != nil {
		return err
	}

	receiverStockHistory := entities.StockHistory{
		PharmacyProduct: entities.PharmacyProduct{Id: receiver.Id},
		Pharmacy:        entities.Pharmacy{Id: receiver.Pharmacy.Id},
		Quantity:        +st.Quantity,
		Description:     fmt.Sprintf(`received %d %s from %s`, st.Quantity, sender.Product.Name, sender.Pharmacy.Name),
	}

	err = u.StockHistoryRepository.CreateOne(ctx, receiverStockHistory)
	if err != nil {
		return err
	}

	err = u.StockTransferRepository.UpdateMutationStatus(ctx, stockTransferId, mutationStatusId)
	if err != nil {
		return err
	}
	return nil
}

func (u *StockTransferUsecaseImpl) UpdateStatusProcessedWithTransaction(ctx context.Context, stockTransferId int64, mutationStatusId int64) error {
	_, err := u.Transactor.WithinTransaction(ctx, func(ctx context.Context) (interface{}, error) {
		err := u.UpdateStatusProcessed(ctx, stockTransferId, mutationStatusId)
		if err != nil {
			return nil, err
		}
		return nil, nil
	})
	if err != nil {
		return err
	}
	return nil
}

func (u *StockTransferUsecaseImpl) UpdateStatusCanceled(ctx context.Context, stockTransferId int64, mutationStatusId int64) error {
	err := u.StockTransferRepository.UpdateMutationStatus(ctx, stockTransferId, mutationStatusId)
	if err != nil {
		return err
	}
	return nil

}

func (u *StockTransferUsecaseImpl) UpdateStatusCanceledWithTransaction(ctx context.Context, stockTransferId int64, mutationStatusId int64) error {
	_, err := u.Transactor.WithinTransaction(ctx, func(ctx context.Context) (interface{}, error) {
		err := u.UpdateStatusCanceled(ctx, stockTransferId, mutationStatusId)
		if err != nil {
			return nil, err
		}
		return nil, nil
	})
	if err != nil {
		return err
	}
	return nil
}

func (u *StockTransferUsecaseImpl) UpdateStatusPending(ctx context.Context, stockTransferId int64, mutationStatusId int64) error {
	err := u.StockTransferRepository.UpdateMutationStatus(ctx, stockTransferId, mutationStatusId)
	if err != nil {
		return err
	}
	return nil
}

func (u *StockTransferUsecaseImpl) UpdateStatusPendingWithTransaction(ctx context.Context, stockTransferId int64, mutationStatusId int64) error {
	_, err := u.Transactor.WithinTransaction(ctx, func(ctx context.Context) (interface{}, error) {
		err := u.UpdateStatusPending(ctx, stockTransferId, mutationStatusId)
		if err != nil {
			return nil, err
		}
		return nil, nil
	})
	if err != nil {
		return err
	}
	return nil
}
