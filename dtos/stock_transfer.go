package dtos

import (
	"github.com/tsanaativa/sehatin-backend-v0.1/entities"
)

type StockTransferCreateRequest struct {
	PharmacySenderId   int64 `json:"pharmacy_sender_id" binding:"required"`
	PharmacyReceiverId int64 `json:"phramacy_receiver_id" binding:"required"`
	ProductId          int64 `json:"product_id" binding:"required"`
	Quantity           int   `json:"quantity" binding:"required"`
}

type MutationStatusIdRequest struct {
	StockTransferId  int64 `json:"stock_transfer_id" binding:"required"`
	MutationStatusId int64 `json:"mutation_status_id" binding:"required"`
}

type StockTransferResponse struct {
	Id                   int64  `json:"id"`
	PharmacySenderId     int64  `json:"pharmacy_sender_id"`
	PharmacySenderName   string `json:"pharmacy_sender_name"`
	PharmacyReceiverId   int64  `json:"pharmacy_receiver_id"`
	PharmacyReceiverName string `json:"pharmacy_receiver_name"`
	ProductId            int64  `json:"product_id"`
	ProductName          string `json:"product_name"`
	Quantity             int    `json:"quantity"`
	MutationStatusId     int64  `json:"mutation_status_id"`
	MutationStatusName   string `json:"mutation_status_name"`
}

type StockTransferResponses struct {
	Pagination     PaginationResponse      `json:"pagination_info"`
	StockTransfers []StockTransferResponse `json:"stock_transfers"`
}

func ConvertToStockTransferResponse(stockTransfer *entities.StockTransfer) *StockTransferResponse {

	return &StockTransferResponse{
		Id:                   stockTransfer.Id,
		PharmacySenderId:     stockTransfer.PharmacySender.Id,
		PharmacySenderName:   stockTransfer.PharmacySender.Name,
		PharmacyReceiverId:   stockTransfer.PharmacyReceiver.Id,
		PharmacyReceiverName: stockTransfer.PharmacyReceiver.Name,
		ProductId:            stockTransfer.Product.Id,
		ProductName:          stockTransfer.Product.Name,
		Quantity:             stockTransfer.Quantity,
		MutationStatusId:     stockTransfer.MutationStatus.Id,
		MutationStatusName:   stockTransfer.MutationStatus.Name,
	}
}

func ConvertToStockTransferResponses(stockTransfers []entities.StockTransfer, pagination entities.PaginationInfo) *StockTransferResponses {
	stockTransferResponses := []StockTransferResponse{}

	for _, stocktransfer := range stockTransfers {
		stockTransferResponses = append(stockTransferResponses, *ConvertToStockTransferResponse(&stocktransfer))
	}

	return &StockTransferResponses{
		Pagination:     *ConvertToPaginationResponse(pagination),
		StockTransfers: stockTransferResponses,
	}
}
