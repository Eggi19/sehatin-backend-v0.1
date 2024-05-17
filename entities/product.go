package entities

import (
	"database/sql"
	"time"
)

type Product struct {
	Id                    int64
	Name                  string
	GenericName           string
	Content               string
	Description           string
	UnitInPack            string
	SellingUnit           string
	Weight                int
	Height                int
	Length                int
	Width                 int
	ProductPicture        string
	SlugId                string
	ProductForm           ProductForm
	ProductClassification ProductClassification
	Manufacture           Manufacture
	Categories            []Category
	CreatedAt             time.Time
	UpdatedAt             time.Time
	DeletedAt             sql.NullTime
}
