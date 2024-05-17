package dtos

type ProvinceResponse struct {
	Id   int16  `json:"id"`
	Name string `json:"name"`
}

type CityResponse struct {
	Id   int16  `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
}

type DistrictResponse struct {
	Id   int16  `json:"id"`
	Name string `json:"name"`
}

type SubDistrictRespone struct {
	Id         int32   `json:"id"`
	Name       string  `json:"name"`
	PostalCode int32   `json:"postal_code"`
	Coordinate *string `json:"coordinate"`
}

type GeoAddress struct {
	Village      string  `json:"village,omitempty"`
	CityDistrict *string `json:"city_district,omitempty"`
	City         *string `json:"city,omitempty"`
	County       *string `json:"county,omitempty"`
	State        string  `json:"state"`
	Postcode     string  `json:"postcode"`
}

type GeodataResponse struct {
	Address     GeoAddress `json:"address"`
	DisplayName string     `json:"display_name"`
}

type GeoReverseResponse struct {
	ProvinceId    *int16 `json:"province_id"`
	CityId        *int16 `json:"city_id"`
	DistrictId    *int16 `json:"district_id"`
	SubDistrictId *int32 `json:"sub_district_id"`
	PostalCode    int32  `json:"postal_code"`
	Address       string `json:"address"`
}
