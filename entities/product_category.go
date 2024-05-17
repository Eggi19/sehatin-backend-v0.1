package entities

type ProductCategory struct {
	Id         int64
	ProductId  int64
	CategoryId int64
}

type ProductCategoryParams struct {
	SortBy  string
	Sort    string
	Limit   int
	Page    int
	Keyword string
}
