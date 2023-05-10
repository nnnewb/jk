package order

import (
	"context"
	"example/pkg/order"
)

type OrderSvc struct{}

func (o OrderSvc) CreateOrder(ctx context.Context, req order.CreateOrderRequest) (order.CreateOrderResponse, error) {
	// TODO implement me
	panic("implement me")
}

func (o OrderSvc) CancelOrder(ctx context.Context, req order.CancelOrderRequest) (order.CancelOrderResponse, error) {
	// TODO implement me
	panic("implement me")
}

func (o OrderSvc) GetOrderDetail(ctx context.Context, req order.GetOrderDetailRequest) (order.GetOrderDetailResponse, error) {
	// TODO implement me
	panic("implement me")
}
