package entities

type SalesReportCategory struct {
	Category  Category
	TotalSold int
	Month     int
	Year      string
}

type SalesReportCategoryParams struct {
	SortBy     string
	Sort       string
	Limit      int
	Page       int
	Keyword    string
	CategoryId int64
}
