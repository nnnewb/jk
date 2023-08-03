// Code generated by jk generate endpoints -t Service; DO NOT EDIT.

package order

import (
	"context"
	"github.com/go-kit/kit/endpoint"
)

func makeEndpointFromFunc[REQ, RESP any](f func(context.Context, REQ) (RESP, error)) endpoint.Endpoint {
	return func(ctx context.Context, r interface{}) (interface{}, error) {
		req := r.(REQ)
		return f(ctx, req)
	}
}

type EndpointSet struct {
	CancelOrderEndpoint endpoint.Endpoint
	CreateOrderEndpoint endpoint.Endpoint
	OrderDetailEndpoint endpoint.Endpoint
	UpdateEndpoint      endpoint.Endpoint
}

func (o EndpointSet) CancelOrder(ctx context.Context, req *CancelOrderRequest) (*CancelOrderResponse, error) {
	resp, err := o.CancelOrderEndpoint(ctx, req)

	if err != nil {
		return &CancelOrderResponse{}, err
	}
	return resp.(*CancelOrderResponse), nil
}

func (o EndpointSet) CreateOrder(ctx context.Context, req *CreateOrderRequest) (*CreateOrderResponse, error) {
	resp, err := o.CreateOrderEndpoint(ctx, req)

	if err != nil {
		return &CreateOrderResponse{}, err
	}
	return resp.(*CreateOrderResponse), nil
}

func (o EndpointSet) OrderDetail(ctx context.Context, req *GetOrderDetailRequest) (*GetOrderDetailResponse, error) {
	resp, err := o.OrderDetailEndpoint(ctx, req)

	if err != nil {
		return &GetOrderDetailResponse{}, err
	}
	return resp.(*GetOrderDetailResponse), nil
}

func (o EndpointSet) Update(ctx context.Context, req *UpdateOrderRequest) (*UpdateOrderResponse, error) {
	resp, err := o.UpdateEndpoint(ctx, req)

	if err != nil {
		return &UpdateOrderResponse{}, err
	}
	return resp.(*UpdateOrderResponse), nil
}

func NewEndpointSet(svc Service) EndpointSet {
	return EndpointSet{
		CancelOrderEndpoint: makeEndpointFromFunc(svc.CancelOrder),
		CreateOrderEndpoint: makeEndpointFromFunc(svc.CreateOrder),
		OrderDetailEndpoint: makeEndpointFromFunc(svc.OrderDetail),
		UpdateEndpoint:      makeEndpointFromFunc(svc.Update),
	}
}

func (s EndpointSet) With(outer endpoint.Middleware, others ...endpoint.Middleware) EndpointSet {
	return EndpointSet{
		CancelOrderEndpoint: endpoint.Chain(outer, others...)(s.CancelOrderEndpoint),
		CreateOrderEndpoint: endpoint.Chain(outer, others...)(s.CreateOrderEndpoint),
		OrderDetailEndpoint: endpoint.Chain(outer, others...)(s.OrderDetailEndpoint),
		UpdateEndpoint:      endpoint.Chain(outer, others...)(s.UpdateEndpoint),
	}
}
