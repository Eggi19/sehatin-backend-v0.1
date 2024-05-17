package dtos

import (
	"mime/multipart"
	"time"

	"github.com/tsanaativa/sehatin-backend-v0.1/entities"
	"github.com/shopspring/decimal"
)

type OrderRequest struct {
	CartItemId     []int64         `json:"cart_item_id"`
	UserAddressId  int64           `json:"user_address_id"`
	TotalPrice     decimal.Decimal `json:"total_price"`
	ShippingFee    decimal.Decimal `json:"shipping_fee"`
	ShippingMethod string          `json:"shipping_method"`
}

type OrderItem struct {
	Name              string          `json:"name"`
	SellingUnit       string          `json:"selling_unit"`
	Price             decimal.Decimal `json:"price"`
	Quantity          int             `json:"quantity"`
	ProductPicture    string          `json:"product_picture"`
	PharmacyProductId int64           `json:"pharmacy_product_id"`
}

type OrderResponse struct {
	Id              int64                   `json:"id"`
	UserName        string                  `json:"user_name,omitempty"`
	UserEmail       string                  `json:"user_email,omitempty"`
	OrderNumber     string                  `json:"order_number"`
	TotalPrice      decimal.Decimal         `json:"total_price"`
	PaymentProof    *string                 `json:"payment_proof"`
	PaymentDeadline time.Time               `json:"payment_deadline"`
	ShippingFee     decimal.Decimal         `json:"shipping_fee"`
	ShippingMethod  string                  `json:"shipping_method"`
	OrderStatus     string                  `json:"order_status"`
	PharmacyName    string                  `json:"pharmacy_name"`
	UserAddress     UserAddressResponse     `json:"user_address"`
	PharmacyAddress PharmacyAddressResponse `json:"pharmacy_address"`
	OrderItems      []OrderItem             `json:"order_items"`
}

type GetAllOrdersResponse struct {
	PaginationInfo PaginationResponse `json:"pagination_info"`
	Data           []OrderResponse    `json:"orders"`
}

type OrderStatusRequest struct {
	Id          int64  `json:"id"`
	OrderStatus string `json:"order_status"`
}

type UploadPaymentProofResponse struct {
	PaymentProof *multipart.FileHeader `form:"payment_proof"`
	OrderId      int64                 `form:"order_id" binding:"required"`
}

type CreateOrderResponse struct {
	PaymentDeadline time.Time `json:"payment_deadline"`
}

func ConvertToOrderResponse(req entities.Order) *OrderResponse {
	userAddress := ConvertToUserAddressResponse(&req.UserAddress)
	pharmacyAddress := ConvertToPharmacyAddressResponse(req.PharmacyAddress)
	return &OrderResponse{
		Id:              req.Id,
		UserName:        req.UserName,
		UserEmail:       req.UserEmail,
		OrderNumber:     req.OrderNumber,
		TotalPrice:      req.TotalPrice,
		PaymentProof:    req.PaymentProof,
		PaymentDeadline: req.PaymentDeadline,
		ShippingFee:     req.ShippingFee,
		ShippingMethod:  req.ShippingMethod,
		OrderStatus:     req.OrderStatus,
		PharmacyName:    req.PharmacyName,
		UserAddress:     *userAddress,
		PharmacyAddress: *pharmacyAddress,
	}
}

func ConvertToOrderResponses(req []entities.Order) []OrderResponse {
	result := []OrderResponse{}

	for i := 0; i < len(req); i++ {
		item := ConvertToOrderResponse(req[i])

		result = append(result, *item)
	}
	return result
}
