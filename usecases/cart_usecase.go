package usecases

import (
	"context"

	"github.com/tsanaativa/sehatin-backend-v0.1/constants"
	"github.com/tsanaativa/sehatin-backend-v0.1/custom_errors"
	"github.com/tsanaativa/sehatin-backend-v0.1/dtos"
	"github.com/tsanaativa/sehatin-backend-v0.1/entities"
	"github.com/tsanaativa/sehatin-backend-v0.1/repositories"
)

type CartUsecaseOpts struct {
	CartRepo                  repositories.CartRepository
	PharmacyProductRepository repositories.PharmacyProductRepository
	ShippingMethodUsecase     ShippingMethodUsecase
}

type CartUsecase interface {
	CreateCartItem(ctx context.Context, req entities.CartItem) error
	IncreaseCartItemQuantity(ctx context.Context, req entities.CartItem) error
	DecreaseCartItemQuantity(ctx context.Context, req entities.CartItem) error
	DeleteCartItem(ctx context.Context, cartId int64) error
	GetUserCartItems(ctx context.Context, userId int64) ([]dtos.CartItemResponse, error)
}

type CartUsecaseImpl struct {
	CartRepository            repositories.CartRepository
	PharmacyProductRepository repositories.PharmacyProductRepository
	ShippingMethodUsecase     ShippingMethodUsecase
}

func NewCartUsecaseImpl(cuOpts *CartUsecaseOpts) CartUsecase {
	return &CartUsecaseImpl{
		CartRepository:            cuOpts.CartRepo,
		PharmacyProductRepository: cuOpts.PharmacyProductRepository,
		ShippingMethodUsecase:     cuOpts.ShippingMethodUsecase,
	}
}

func (u *CartUsecaseImpl) CreateCartItem(ctx context.Context, req entities.CartItem) error {
	cart, err := u.CartRepository.FindCartItem(ctx, req)
	if err != nil && err.Error() != constants.ResponseMsgErrorNotFound {
		return err
	}

	if err != nil && err.Error() == constants.ResponseMsgErrorNotFound {
		isAvailable, err := u.checkStockAvailability(ctx, req.PharmacyProductId, 0, req.Quantity)
		if err != nil {
			return err
		}
		if !*isAvailable {
			return custom_errors.NotEnoughStock()
		}

		err = u.CartRepository.CreateOneCartItem(ctx, req)
		if err != nil {
			return err
		}

		return nil
	}

	req.Id = cart.Id

	err = u.IncreaseCartItemQuantity(ctx, req)
	if err != nil {
		return err
	}

	return nil
}

func (u *CartUsecaseImpl) IncreaseCartItemQuantity(ctx context.Context, req entities.CartItem) error {
	cart, err := u.CartRepository.FindPharmacyIdByCartId(ctx, req.Id)
	if err != nil {
		return err
	}

	isAvailable, err := u.checkStockAvailability(ctx, cart.PharmacyProductId, cart.Quantity, req.Quantity)
	if err != nil {
		return err
	}
	if !*isAvailable {
		return custom_errors.NotEnoughStock()
	}

	err = u.CartRepository.IncreaseCartQuantity(ctx, req)
	if err != nil {
		return err
	}

	return nil
}

func (u *CartUsecaseImpl) DecreaseCartItemQuantity(ctx context.Context, req entities.CartItem) error {
	cart, err := u.CartRepository.DecreaseCartQuantity(ctx, req)
	if err != nil {
		return err
	}

	if cart.Quantity <= 0 {
		err = u.CartRepository.DeleteCartItem(ctx, cart.Id)
		if err != nil {
			return err
		}
	}

	return nil
}

func (u *CartUsecaseImpl) DeleteCartItem(ctx context.Context, cartId int64) error {
	err := u.CartRepository.DeleteCartItem(ctx, cartId)
	if err != nil {
		return err
	}

	return nil
}

func (u *CartUsecaseImpl) GetUserCartItems(ctx context.Context, userId int64) ([]dtos.CartItemResponse, error) {
	response, err := u.CartRepository.FindAllUserCartItem(ctx, userId)
	if err != nil {
		return nil, err
	}

	result := dtos.ConvertToCartItemResponses(response)

	for key, res := range result {
		shippingMethodRes, err := u.ShippingMethodUsecase.GetShippingMethod(ctx, res.PharmacyId)
		if err != nil {
			return nil, err
		}

		result[key].ShippingMethods = *shippingMethodRes
	}

	return result, nil
}

func (u *CartUsecaseImpl) checkStockAvailability(ctx context.Context, PharmacyProductId int64, cartQuantity int, reqQuantity int) (*bool, error) {
	var result bool
	result = true

	pharmacyProduct, err := u.PharmacyProductRepository.FindOnePharmacyProduct(ctx, PharmacyProductId)
	if err != nil {
		return nil, err
	}

	if pharmacyProduct.TotalStock < cartQuantity+reqQuantity {
		result = false
		return &result, nil
	}

	return &result, nil
}
