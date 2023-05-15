package order

import (
	"context"
)

//go:generate jk generate all -t OrderService --api-version v2
type OrderService interface {
	// CreateOrder 创建订单
	CreateOrder(ctx context.Context, req CreateOrderRequest) (CreateOrderResponse, error)
	// CancelOrder 取消订单
	CancelOrder(ctx context.Context, req CancelOrderRequest) (CancelOrderResponse, error)
	// GetOrderDetail 获取订单详情
	GetOrderDetail(ctx context.Context, req GetOrderDetailRequest) (GetOrderDetailResponse, error)
}

// CreateOrderRequest 创建订单请求结构体
type CreateOrderRequest struct {
	// 订单信息
	OrderInfo []OrderItem `json:"order_info"`
}

// OrderItem 订单项结构体
type OrderItem struct {
	// 商品ID
	ItemID string `json:"item_id"`
	// 商品数量
	Quantity int `json:"quantity"`
}

// CreateOrderResponse 创建订单响应结构体
type CreateOrderResponse struct {
	// 响应状态码
	Code int `json:"code"`
	// 响应消息
	Message string `json:"message"`
	// 订单号
	OrderID string `json:"order_id"`
}

// CancelOrderRequest 取消订单请求结构体
type CancelOrderRequest struct {
	// 订单号
	OrderID string `json:"order_id"`
}

// CancelOrderResponse 取消订单响应结构体
type CancelOrderResponse struct {
	// 响应状态码
	Code int `json:"code"`
	// 响应消息
	Message string `json:"message"`
}

// GetOrderDetailRequest 获取订单详情请求结构体
type GetOrderDetailRequest struct {
	// 订单号
	OrderID string `json:"order_id"`
}

// GetOrderDetailResponse 获取订单详情响应结构体
type GetOrderDetailResponse struct {
	// 响应状态码
	Code int `json:"code"`
	// 响应消息
	Message string `json:"message"`
	// 订单信息
	OrderInfo []OrderItem `json:"order_info"`
}
