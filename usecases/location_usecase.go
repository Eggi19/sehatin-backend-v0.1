package usecases

import (
	"context"
	"regexp"
	"strings"

	"github.com/tsanaativa/sehatin-backend-v0.1/constants"
	"github.com/tsanaativa/sehatin-backend-v0.1/custom_errors"
	"github.com/tsanaativa/sehatin-backend-v0.1/dtos"
	"github.com/tsanaativa/sehatin-backend-v0.1/entities"
	"github.com/tsanaativa/sehatin-backend-v0.1/repositories"
)

type LocationUsecaseOpts struct {
	LocationRepo repositories.LocationRepository
}

type LocationUsecase interface {
	GetAllProvinces(ctx context.Context) ([]dtos.ProvinceResponse, error)
	GetCitiesByProvinceId(ctx context.Context, provinceId int16) ([]dtos.CityResponse, error)
	GetDistrictsByCityId(ctx context.Context, cityId int16) ([]dtos.DistrictResponse, error)
	GetSubDistrictsByDistrictId(ctx context.Context, districtId int16) ([]dtos.SubDistrictRespone, error)
	ReverseCoordinate(ctx context.Context, latitude string, longitude string) (*dtos.GeoReverseResponse, error)
}

type LocationUsecaseImpl struct {
	LocationRepository repositories.LocationRepository
}

func NewLocationUsecaseImpl(lUseOpts *LocationUsecaseOpts) LocationUsecase {
	return &LocationUsecaseImpl{lUseOpts.LocationRepo}
}

func (u *LocationUsecaseImpl) GetAllProvinces(ctx context.Context) ([]dtos.ProvinceResponse, error) {
	provinces, err := u.LocationRepository.FindProvinces(ctx)
	if err != nil {
		return nil, custom_errors.InternalServerError(err)
	}

	return provinces, nil
}

func (u *LocationUsecaseImpl) GetCitiesByProvinceId(ctx context.Context, provinceId int16) ([]dtos.CityResponse, error) {
	cities, err := u.LocationRepository.FindCities(ctx, provinceId)
	if err != nil {
		return nil, custom_errors.InternalServerError(err)
	}

	return cities, nil
}

func (u *LocationUsecaseImpl) GetDistrictsByCityId(ctx context.Context, cityId int16) ([]dtos.DistrictResponse, error) {
	districts, err := u.LocationRepository.FindDistricts(ctx, cityId)
	if err != nil {
		return nil, custom_errors.InternalServerError(err)
	}

	return districts, nil
}

func (u *LocationUsecaseImpl) GetSubDistrictsByDistrictId(ctx context.Context, districtId int16) ([]dtos.SubDistrictRespone, error) {
	subDistricts, err := u.LocationRepository.FindSubDistricts(ctx, districtId)
	if err != nil {
		return nil, custom_errors.InternalServerError(err)
	}

	return subDistricts, nil
}

func (u *LocationUsecaseImpl) ReverseCoordinate(ctx context.Context, latitude string, longitude string) (*dtos.GeoReverseResponse, error) {
	longitudeIsNumber := regexp.MustCompile(`^-?[0-9]\d*(\.\d+)?$`).MatchString(longitude)
	latitudeIsNumber := regexp.MustCompile(`^-?[0-9]\d*(\.\d+)?$`).MatchString(latitude)

	if !(longitudeIsNumber && latitudeIsNumber) {
		return nil, custom_errors.BadRequest(custom_errors.ErrNonNumberCoordinate, constants.NonNumberCoordinateErrMsg)
	}

	geodata, err := u.LocationRepository.GetLocationByCoord(ctx, latitude, longitude)
	if err != nil {
		return nil, custom_errors.InternalServerError(err)
	}

	userAddress := entities.UserAddress{
		Address:     geodata.DisplayName,
		SubDistrict: geodata.Address.Village,
	}

	geoResponse := dtos.GeoReverseResponse{
		Address: geodata.DisplayName,
	}

	cityType := "Kota"
	if geodata.Address.City != nil {
		userAddress.City = *geodata.Address.City
		if *geodata.Address.City == "Daerah Khusus Jakarta" {
			userAddress.City = *geodata.Address.CityDistrict
		}
	}
	if geodata.Address.County != nil && userAddress.City == "" {
		userAddress.City = *geodata.Address.County
		cityType = "Kabupaten"
	}
	if geodata.Address.CityDistrict != nil && userAddress.City == "" {
		userAddress.City = *geodata.Address.CityDistrict
	}

	if strings.HasPrefix(userAddress.SubDistrict, "Desa ") || strings.HasPrefix(userAddress.SubDistrict, "Kelurahan ") {
		splittedAddr := strings.Split(userAddress.SubDistrict, " ")
		userAddress.SubDistrict = strings.Join(splittedAddr[1:len(splittedAddr)-1], " ")
	}

	cleanCity := regexp.MustCompile(`Kota |Kabupaten `)

	city := entities.City{Name: cleanCity.ReplaceAllString(userAddress.City, ""), Type: cityType}
	provinceAndCity, err := u.LocationRepository.FindProvinceAndCity(ctx, city)
	switch err.(type) {
	case *custom_errors.AppError:
		return &geoResponse, nil
	case error:
		return nil, custom_errors.InternalServerError(err)
	}

	geoResponse.ProvinceId = provinceAndCity.ProvinceId
	geoResponse.CityId = provinceAndCity.CityId
	districtAndSubDistrict, err := u.LocationRepository.FindDistrictAndSubDistrict(ctx, *geoResponse.CityId, userAddress.SubDistrict)
	switch err.(type) {
	case *custom_errors.AppError:
		return &geoResponse, nil
	case error:
		return nil, custom_errors.InternalServerError(err)
	}

	geoResponse.DistrictId = districtAndSubDistrict.DistrictId
	geoResponse.SubDistrictId = districtAndSubDistrict.SubDistrictId
	geoResponse.PostalCode = districtAndSubDistrict.PostalCode

	return &geoResponse, nil
}
