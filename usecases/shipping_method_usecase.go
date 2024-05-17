package usecases

import (
	"context"
	"math"

	"github.com/tsanaativa/sehatin-backend-v0.1/dtos"
	"github.com/tsanaativa/sehatin-backend-v0.1/repositories"
)

type ShippingMethodOpts struct {
	ShippingMethodRepository  repositories.ShippingMethodRepository
	UserAddressRepository     repositories.UserAddressRepository
	PharmacyAddressRepository repositories.PharmacyAddressRepository
}

type ShippingMethodUsecase interface {
	GetShippingMethod(ctx context.Context, pharmacyId int64) (*dtos.ShippingMethod, error)
	GetOfficialFee(ctx context.Context, req dtos.OfficialShippingFeeRequest) (float64, error)
	GetNonOfficialFee(ctx context.Context, req dtos.NonOfficialShippingFeeRequest, userId int) (float64, error)
}

type ShippingMethodUsecaseImpl struct {
	ShippingMethodRepository  repositories.ShippingMethodRepository
	UserAddressRepository     repositories.UserAddressRepository
	PharmacyAddressRepository repositories.PharmacyAddressRepository
}

func NewShippingMethodUsecaseImpl(smOpts *ShippingMethodOpts) ShippingMethodUsecase {
	return &ShippingMethodUsecaseImpl{
		ShippingMethodRepository:  smOpts.ShippingMethodRepository,
		UserAddressRepository:     smOpts.UserAddressRepository,
		PharmacyAddressRepository: smOpts.PharmacyAddressRepository,
	}
}

func (u *ShippingMethodUsecaseImpl) GetShippingMethod(ctx context.Context, pharmacyId int64) (*dtos.ShippingMethod, error) {
	official, err := u.ShippingMethodRepository.GetOfficialShippingMethod(ctx, pharmacyId)
	if err != nil {
		return nil, err
	}

	nonOfficial, err := u.ShippingMethodRepository.GetNonOfficialShippingMethod(ctx, pharmacyId)
	if err != nil {
		return nil, err
	}

	shippingMethods := dtos.ConvertToShippingMethodDto(official, nonOfficial)

	return &shippingMethods, nil
}

func (u *ShippingMethodUsecaseImpl) GetOfficialFee(ctx context.Context, req dtos.OfficialShippingFeeRequest) (float64, error) {
	officialShipping, err := u.ShippingMethodRepository.GetOfficialShippingFee(ctx, req.OfficialShippingMethodId)
	if err != nil {
		return 0, err
	}

	distance, err := u.UserAddressRepository.GetDistanceFromPharmacy(ctx, req.UserAddressId, req.PharmacyId)
	if err != nil {
		return 0, err
	}

	shippingFee, _ := officialShipping.Fee.Float64()
	cost := math.Round(shippingFee * (*distance))
	if cost < shippingFee {
		cost = shippingFee
	}

	return cost, nil
}

func (u *ShippingMethodUsecaseImpl) GetNonOfficialFee(ctx context.Context, req dtos.NonOfficialShippingFeeRequest, userId int) (float64, error) {
	userAddress, err := u.UserAddressRepository.FindById(ctx, req.UserAddressId, int64(userId))
	if err != nil {
		return 0, err
	}

	pharmacyAddress, err := u.PharmacyAddressRepository.FindById(ctx, req.PharmacyId)
	if err != nil {
		return 0, err
	}

	NonOfficialShipping, err := u.ShippingMethodRepository.GetNonOfficialShippingService(ctx, req.NonOfficialShippingMethodId)
	if err != nil {
		return 0, err
	}

	response, err := u.ShippingMethodRepository.GetNonOfficialFee(userAddress.CityId, pharmacyAddress.CityId, req.TotalWeight, NonOfficialShipping.Courier)
	if err != nil {
		return 0, err
	}

	for _, res := range response.RajaOngkir.Results{
		for _, courier := range res.Costs {
			if courier.Service == NonOfficialShipping.Service {
				return courier.Cost[0].Value, nil
			}
		}
	}

	return 0, nil
}
