package order

import (
	"context"

	"example/api/order"
)

type OrderSvc struct{}

func (o OrderSvc) CreateOrder(ctx context.Context, req *order.CreateOrderRequest) (*order.CreateOrderResponse, error) {
	return &order.CreateOrderResponse{
		Code:    -1,
		Message: "not implemented",
	}, nil
}

func (o OrderSvc) CancelOrder(ctx context.Context, req *order.CancelOrderRequest) (*order.CancelOrderResponse, error) {
	return &order.CancelOrderResponse{
		Code:    -1,
		Message: "not implemented",
	}, nil
}

func (o OrderSvc) OrderDetail(ctx context.Context, req *order.GetOrderDetailRequest) (*order.GetOrderDetailResponse, error) {
	return &order.GetOrderDetailResponse{
		Code:    -1,
		Message: "not implemented",
	}, nil
}

func (o OrderSvc) Update(ctx context.Context, req *order.UpdateOrderRequest) (*order.UpdateOrderResponse, error) {
	return &order.UpdateOrderResponse{
		Code:    -1,
		Message: "not implemented",
	}, nil
}
