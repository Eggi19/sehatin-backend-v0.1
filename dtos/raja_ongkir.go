package dtos

type Query struct {
	Origin      string  `json:"origin"`
	Destination string  `json:"destination"`
	Weight      float64 `json:"weight"`
	Courier     string  `json:"courier"`
}

type Status struct {
	Code        int    `json:"code"`
	Description string `json:"description"`
}

type AddressDetails struct {
	CityId     string `json:"city_id"`
	ProvinceId string `json:"province_id"`
	Province   string `json:"province"`
	Type       string `json:"type"`
	CityName   string `json:"city_name"`
	PostalCode string `json:"postal_code"`
}

type Cost struct {
	Value float64 `json:"value"`
	Etd   string  `json:"etd"`
	Note  string  `json:"note"`
}

type Costs struct {
	Service     string `json:"service"`
	Description string `json:"description"`
	Cost        []Cost `json:"cost"`
}

type CostResponse struct {
	Code  string  `json:"code"`
	Name  string  `json:"name"`
	Costs []Costs `json:"costs"`
}

type RajaOngkirCostResponse struct {
	Query              Query          `json:"query"`
	Status             Status         `json:"status"`
	OriginDetails      AddressDetails `json:"origin_details"`
	DestinationDetails AddressDetails `json:"destination_details"`
	Results            []CostResponse `json:"results"`
}

type RajaOngkirResponse struct {
	RajaOngkir RajaOngkirCostResponse `json:"rajaongkir"`
}
