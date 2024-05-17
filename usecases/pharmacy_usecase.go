package usecases

import (
	"context"

	"github.com/tsanaativa/sehatin-backend-v0.1/entities"
	"github.com/tsanaativa/sehatin-backend-v0.1/repositories"
)

type PharmacyUsecaseOpts struct {
	PharmacyRepo        repositories.PharmacyRepository
	PharmacyAddressRepo repositories.PharmacyAddressRepository
	ShippingMethodRepo  repositories.ShippingMethodRepository
	PharmacyProductRepo repositories.PharmacyProductRepository
	StockHistoryRepo    repositories.StockHistoryRepository
}

type PharmacyUsecase interface {
	CreatePharmacy(ctx context.Context, pharmacy entities.Pharmacy) error
	GetPharmacyById(ctx context.Context, pharmacyId int64) (*entities.Pharmacy, error)
	UpdatePharmacy(ctx context.Context, pharmacy entities.Pharmacy) error
	DeletePharmacyById(ctx context.Context, pharmacyId int64) error
	GetAllPharmacyByPharmacyManagerId(ctx context.Context, pharmacyManagerId int64, params entities.PharmacyParams) ([]entities.Pharmacy, *entities.PaginationInfo, error)
}

type PharmacyUsecaseImpl struct {
	PharmacyRepository        repositories.PharmacyRepository
	PharmacyAddressRepository repositories.PharmacyAddressRepository
	ShippingMethodRepository  repositories.ShippingMethodRepository
	PharmacyProductRepository repositories.PharmacyProductRepository
	StockHistoryRepository    repositories.StockHistoryRepository
}

func NewPharmacyUsecaseImpl(puOpts *PharmacyUsecaseOpts) PharmacyUsecase {
	return &PharmacyUsecaseImpl{
		PharmacyRepository:        puOpts.PharmacyRepo,
		PharmacyAddressRepository: puOpts.PharmacyAddressRepo,
		ShippingMethodRepository:  puOpts.ShippingMethodRepo,
		PharmacyProductRepository: puOpts.PharmacyProductRepo,
		StockHistoryRepository:    puOpts.StockHistoryRepo,
	}
}

func (u *PharmacyUsecaseImpl) CreatePharmacy(ctx context.Context, pharmacy entities.Pharmacy) error {
	createdPharmacy, err := u.PharmacyRepository.CreateOne(ctx, pharmacy)
	if err != nil {
		return err
	}

	for i := 0; i < len(pharmacy.OfficialShippingMethod); i++ {
		_, err := u.ShippingMethodRepository.CreatePharmacyOfficialShippingMethod(ctx, entities.PharmacyShippingMethod{
			PharmacyId:         createdPharmacy.Id,
			OfficialShippingId: pharmacy.OfficialShippingMethod[i].Id,
		})
		if err != nil {
			return err
		}
	}

	for i := 0; i < len(pharmacy.NonOfficialShippingMethod); i++ {
		_, err := u.ShippingMethodRepository.CreatePharmacyNonOfficialShippingMethod(ctx, entities.PharmacyShippingMethod{
			PharmacyId:            createdPharmacy.Id,
			NonOfficialShippingId: pharmacy.NonOfficialShippingMethod[i].Id,
		})
		if err != nil {
			return err
		}
	}

	pharmacyAddress := entities.PharmacyAddress{
		PharmacyId:  createdPharmacy.Id,
		City:        pharmacy.PharmacyAddress.City,
		CityId:      pharmacy.PharmacyAddress.CityId,
		Province:    pharmacy.PharmacyAddress.Province,
		Address:     pharmacy.PharmacyAddress.Address,
		District:    pharmacy.PharmacyAddress.District,
		SubDistrict: pharmacy.PharmacyAddress.SubDistrict,
		PostalCode:  pharmacy.PharmacyAddress.PostalCode,
		Longitude:   pharmacy.PharmacyAddress.Longitude,
		Latitude:    pharmacy.PharmacyAddress.Latitude,
	}

	err = u.PharmacyAddressRepository.CreateOne(ctx, pharmacyAddress)
	if err != nil {
		return err
	}

	return nil
}

func (u *PharmacyUsecaseImpl) GetPharmacyById(ctx context.Context, pharmacyId int64) (*entities.Pharmacy, error) {
	pharmacy, err := u.PharmacyRepository.FindOneById(ctx, pharmacyId)
	if err != nil {
		return nil, err
	}

	officialShipping, err := u.ShippingMethodRepository.GetOfficialShippingMethod(ctx, pharmacyId)
	if err != nil {
		return nil, err
	}

	nonOfficialShipping, _ := u.ShippingMethodRepository.GetNonOfficialShippingMethod(ctx, pharmacyId)

	if nonOfficialShipping != nil {
		pharmacy.NonOfficialShippingMethod = nonOfficialShipping
	} else {
		pharmacy.NonOfficialShippingMethod = nil
	}

	pharmacy.OfficialShippingMethod = officialShipping

	return pharmacy, nil
}

func (u *PharmacyUsecaseImpl) UpdatePharmacy(ctx context.Context, pharmacy entities.Pharmacy) error {
	foundedPharmacy, err := u.GetPharmacyById(ctx, pharmacy.Id)
	if err != nil {
		return err
	}

	err = u.ShippingMethodRepository.DeleteShippingMethodByPharmacyId(ctx, pharmacy.Id)
	if err != nil {
		return err
	}

	for i := 0; i < len(pharmacy.OfficialShippingMethod); i++ {
		_, err := u.ShippingMethodRepository.CreatePharmacyOfficialShippingMethod(ctx, entities.PharmacyShippingMethod{
			PharmacyId:         pharmacy.Id,
			OfficialShippingId: pharmacy.OfficialShippingMethod[i].Id,
		})
		if err != nil {
			return err
		}
	}

	for i := 0; i < len(pharmacy.NonOfficialShippingMethod); i++ {
		_, err := u.ShippingMethodRepository.CreatePharmacyNonOfficialShippingMethod(ctx, entities.PharmacyShippingMethod{
			PharmacyId:            pharmacy.Id,
			NonOfficialShippingId: pharmacy.NonOfficialShippingMethod[i].Id,
		})
		if err != nil {
			return err
		}
	}

	err = u.PharmacyRepository.UpdateOne(ctx, pharmacy)
	if err != nil {
		return err
	}

	pharmacyAddress := entities.PharmacyAddress{
		Id:          foundedPharmacy.PharmacyAddress.Id,
		PharmacyId:  pharmacy.Id,
		City:        pharmacy.PharmacyAddress.City,
		Province:    pharmacy.PharmacyAddress.Province,
		Address:     pharmacy.PharmacyAddress.Address,
		District:    pharmacy.PharmacyAddress.District,
		SubDistrict: pharmacy.PharmacyAddress.SubDistrict,
		PostalCode:  pharmacy.PharmacyAddress.PostalCode,
		Longitude:   pharmacy.PharmacyAddress.Longitude,
		Latitude:    pharmacy.PharmacyAddress.Latitude,
	}

	err = u.PharmacyAddressRepository.UpdateOne(ctx, pharmacyAddress)
	if err != nil {
		return err
	}

	return nil
}

func (u *PharmacyUsecaseImpl) DeletePharmacyById(ctx context.Context, pharmacyId int64) error {
	_, err := u.GetPharmacyById(ctx, pharmacyId)
	if err != nil {
		return err
	}

	err = u.PharmacyAddressRepository.DeleteOne(ctx, pharmacyId)
	if err != nil {
		return err
	}

	err = u.ShippingMethodRepository.DeleteAllOfficialShippingMethod(ctx, pharmacyId)
	if err != nil {
		return err
	}

	err = u.PharmacyProductRepository.DeletedAllPharmacyProduct(ctx, pharmacyId)
	if err != nil {
		return err
	}

	err = u.PharmacyRepository.DeleteById(ctx, pharmacyId)
	if err != nil {
		return err
	}
	return nil
}

func (u *PharmacyUsecaseImpl) GetAllPharmacyByPharmacyManagerId(ctx context.Context, pharmacyManagerId int64, params entities.PharmacyParams) ([]entities.Pharmacy, *entities.PaginationInfo, error) {
	pharmacies, totalData, err := u.PharmacyRepository.FindAllByPharmacyManagerId(ctx, pharmacyManagerId, params)
	if err != nil {
		return nil, nil, err
	}

	for i := 0; i < len(pharmacies); i++ {
		officialShipping, err := u.ShippingMethodRepository.GetOfficialShippingMethod(ctx, pharmacies[i].Id)
		if err != nil {
			return nil, nil, err
		}
		nonOfficialShipping, _ := u.ShippingMethodRepository.GetNonOfficialShippingMethod(ctx, pharmacies[i].Id)
		if nonOfficialShipping != nil {
			pharmacies[i].NonOfficialShippingMethod = nonOfficialShipping
		} else {
			pharmacies[i].NonOfficialShippingMethod = nil
		}
		pharmacies[i].OfficialShippingMethod = officialShipping
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

	return pharmacies, &pagination, nil
}
