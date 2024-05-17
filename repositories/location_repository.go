package repositories

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/tsanaativa/sehatin-backend-v0.1/custom_errors"
	"github.com/tsanaativa/sehatin-backend-v0.1/dtos"
	"github.com/tsanaativa/sehatin-backend-v0.1/entities"
)

type LocationRepoOpt struct {
	Db *sql.DB
}

type LocationRepository interface {
	FindProvinces(ctx context.Context) ([]dtos.ProvinceResponse, error)
	FindCities(ctx context.Context, provinceId int16) ([]dtos.CityResponse, error)
	FindDistricts(ctx context.Context, cityId int16) ([]dtos.DistrictResponse, error)
	FindSubDistricts(ctx context.Context, districtId int16) ([]dtos.SubDistrictRespone, error)
	GetLocationByCoord(ctx context.Context, latitude string, longitude string) (*dtos.GeodataResponse, error)
	FindProvinceAndCity(ctx context.Context, c entities.City) (*dtos.GeoReverseResponse, error)
	FindDistrictAndSubDistrict(ctx context.Context, cityId int16, subDistrict string) (*dtos.GeoReverseResponse, error)
}

type LocationRepositoryPostgres struct {
	db *sql.DB
}

func NewLocationRepositoryPostgres(lOpt *LocationRepoOpt) LocationRepository {
	return &LocationRepositoryPostgres{lOpt.Db}
}

func (r *LocationRepositoryPostgres) FindProvinces(ctx context.Context) ([]dtos.ProvinceResponse, error) {
	ps := []dtos.ProvinceResponse{}

	var err error
	var rows *sql.Rows

	tx := extractTx(ctx)
	if tx != nil {
		rows, err = tx.QueryContext(ctx, qFindProvinces)
	} else {
		rows, err = r.db.QueryContext(ctx, qFindProvinces)
	}
	defer rows.Close()

	for rows.Next() {
		p := dtos.ProvinceResponse{}
		rows.Scan(&p.Id, &p.Name)
		ps = append(ps, p)
	}

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, custom_errors.NotFound(err)
		}
		return nil, err
	}

	return ps, nil
}

func (r *LocationRepositoryPostgres) FindCities(ctx context.Context, provinceId int16) ([]dtos.CityResponse, error) {
	cs := []dtos.CityResponse{}

	var err error
	var rows *sql.Rows

	tx := extractTx(ctx)
	if tx != nil {
		rows, err = tx.QueryContext(ctx, qFindCities, provinceId)
	} else {
		rows, err = r.db.QueryContext(ctx, qFindCities, provinceId)
	}
	defer rows.Close()

	for rows.Next() {
		c := dtos.CityResponse{}
		rows.Scan(&c.Id, &c.Name, &c.Type)
		cs = append(cs, c)
	}

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, custom_errors.NotFound(err)
		}
		return nil, err
	}

	return cs, nil
}

func (r *LocationRepositoryPostgres) FindDistricts(ctx context.Context, cityId int16) ([]dtos.DistrictResponse, error) {
	ds := []dtos.DistrictResponse{}

	var err error
	var rows *sql.Rows

	tx := extractTx(ctx)
	if tx != nil {
		rows, err = tx.QueryContext(ctx, qFindDistricts, cityId)
	} else {
		rows, err = r.db.QueryContext(ctx, qFindDistricts, cityId)
	}
	defer rows.Close()

	for rows.Next() {
		d := dtos.DistrictResponse{}
		rows.Scan(&d.Id, &d.Name)
		ds = append(ds, d)
	}

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, custom_errors.NotFound(err)
		}
		return nil, err
	}

	return ds, nil
}

func (r *LocationRepositoryPostgres) FindSubDistricts(ctx context.Context, districtId int16) ([]dtos.SubDistrictRespone, error) {
	sds := []dtos.SubDistrictRespone{}

	var err error
	var rows *sql.Rows

	tx := extractTx(ctx)
	if tx != nil {
		rows, err = tx.QueryContext(ctx, qFindSubDistricts, districtId)
	} else {
		rows, err = r.db.QueryContext(ctx, qFindSubDistricts, districtId)
	}
	defer rows.Close()

	for rows.Next() {
		sd := dtos.SubDistrictRespone{}
		rows.Scan(&sd.Id, &sd.Name, &sd.PostalCode, &sd.Coordinate)
		if strings.Split(*sd.Coordinate, " ")[1] == "0)" {
			sd.Coordinate = nil
		}
		sds = append(sds, sd)
	}

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, custom_errors.NotFound(err)
		}
		return nil, err
	}

	return sds, nil
}

func (r *LocationRepositoryPostgres) GetLocationByCoord(ctx context.Context, latitude string, longitude string) (*dtos.GeodataResponse, error) {
	res, err := http.Get("https://nominatim.openstreetmap.org/reverse?format=jsonv2&zoom=13&lat=" + latitude + "&lon=" + longitude)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var geodata dtos.GeodataResponse
	err = json.NewDecoder(res.Body).Decode(&geodata)
	if err != nil {
		return nil, err
	}

	return &geodata, nil
}

func (r *LocationRepositoryPostgres) FindProvinceAndCity(ctx context.Context, c entities.City) (*dtos.GeoReverseResponse, error) {
	g := dtos.GeoReverseResponse{}

	var err error
	values := []interface{}{c.Name, c.Type}

	tx := extractTx(ctx)
	if tx != nil {
		err = tx.QueryRowContext(ctx, qFindProvinceAndCity, values...).Scan(&g.ProvinceId, &g.CityId)
	} else {
		err = r.db.QueryRowContext(ctx, qFindProvinceAndCity, values...).Scan(&g.ProvinceId, &g.CityId)
	}

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, custom_errors.NotFound(err)
		}
		return nil, err
	}

	return &g, nil
}

func (r *LocationRepositoryPostgres) FindDistrictAndSubDistrict(ctx context.Context, cityId int16, subDistrict string) (*dtos.GeoReverseResponse, error) {
	g := dtos.GeoReverseResponse{}

	var err error
	subDistrictTokens := strings.Split(subDistrict, " ")
	subDistrictConcat := ""
	subDistrictIlike := ""

	values := []interface{}{cityId}
	for i, subDistrictToken := range subDistrictTokens {
		values = append(values, subDistrictToken)
		subDistrictConcat += "$" + strconv.Itoa(i+2) + "::text"
		if i < len(subDistrictTokens)-1 {
			subDistrictConcat += ", ' ', "
		}

		subDistrictIlike += "sd.name ILIKE '%' || $" + strconv.Itoa(i+2) + " || '%'"
		if i < len(subDistrictTokens)-1 {
			subDistrictIlike += "OR "
		}
	}

	q := dFindDistrictAndSubDistrict
	q = strings.ReplaceAll(q, "--sub-district-concat--", subDistrictConcat)
	q = strings.ReplaceAll(q, "--sub-district-ilike--", subDistrictIlike)

	tx := extractTx(ctx)
	if tx != nil {
		err = tx.QueryRowContext(ctx, q, values...).Scan(&g.DistrictId, &g.SubDistrictId, &g.PostalCode)
	} else {
		err = r.db.QueryRowContext(ctx, q, values...).Scan(&g.DistrictId, &g.SubDistrictId, &g.PostalCode)
	}

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, custom_errors.NotFound(err)
		}
		return nil, err
	}

	return &g, nil
}
