package dtos

import (
	"github.com/tsanaativa/sehatin-backend-v0.1/entities"
	"github.com/shopspring/decimal"
)

type StockHistoryResponse struct {
	Id          int64           `json:"id"`
	ProductName string          `json:"product_name"`
	Price       decimal.Decimal `json:"price"`
	Quantity    int             `json:"quantity"`
	Description string          `json:"description"`
}

type StockHistoryResponses struct {
	Pagination     PaginationResponse     `json:"pagination_info"`
	StockHistories []StockHistoryResponse `json:"stock_histories"`
}

func ConvertToStockHistoryResponse(stockHistory *entities.StockHistory) *StockHistoryResponse {
	return &StockHistoryResponse{
		Id:          stockHistory.Id,
		ProductName: stockHistory.PharmacyProduct.Product.Name,
		Price:       stockHistory.PharmacyProduct.Price,
		Quantity:    stockHistory.Quantity,
		Description: stockHistory.Description,
	}
}

func ConvertToStockHistoryResponses(stockHistories []entities.StockHistory, pagination entities.PaginationInfo) *StockHistoryResponses {
	stockHistoriesResponses := []StockHistoryResponse{}

	for _, stockHistory := range stockHistories {
		stockHistoriesResponses = append(stockHistoriesResponses, *ConvertToStockHistoryResponse(&stockHistory))
	}

	return &StockHistoryResponses{
		Pagination:     *ConvertToPaginationResponse(pagination),
		StockHistories: stockHistoriesResponses,
	}
}
