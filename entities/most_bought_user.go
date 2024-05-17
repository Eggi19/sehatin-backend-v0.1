package entities

type MostBoughtUser struct {
	PharmacyProduct PharmacyProduct
}

type MostBoughtUserParams struct {
	SortBy    string
	Sort      string
	Limit     int
	Page      int
	Longitude float64
	Latitude  float64
	Radius    int
	ProductId int64
}
