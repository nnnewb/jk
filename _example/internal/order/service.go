package order

import (
	"context"

	"example/pkg/order"
)

type OrderSvc struct{}

func (o OrderSvc) CreateOrder(ctx context.Context, req order.CreateOrderRequest) (order.CreateOrderResponse, error) {
	return order.CreateOrderResponse{
		Code:    -1,
		Message: "not implemented",
	}, nil
}

func (o OrderSvc) CancelOrder(ctx context.Context, req order.CancelOrderRequest) (order.CancelOrderResponse, error) {
	return order.CancelOrderResponse{
		Code:    -1,
		Message: "not implemented",
	}, nil
}

func (o OrderSvc) GetOrderDetail(ctx context.Context, req order.GetOrderDetailRequest) (order.GetOrderDetailResponse, error) {
	return order.GetOrderDetailResponse{
		Code:    -1,
		Message: "not implemented",
	}, nil
}
