package order

import (
	"context"
)

// OrderService 订单服务定义
//
// @swagger-info-version v1
// @swagger-info-title 订单服务
// @http-base-path /api/v1/order-service/
//
//go:generate jk generate endpoints -t OrderService
//go:generate jk generate transport -t OrderService --protocol http --server --language go --framework gin --swagger
//go:generate jk generate transport -t OrderService --protocol http --client --language go --framework http
//go:generate jk generate transport -t OrderService --protocol http --client --language ts --framework fetch
//go:generate prettier -w client.ts
//go:generate tsc
type OrderService interface {
	// CreateOrder 创建订单
	// @http-method post
	// @http-path /api/v1/order-service/order
	CreateOrder(ctx context.Context, req *CreateOrderRequest) (*CreateOrderResponse, error)
	// CancelOrder 取消订单
	// @http-method post
	// @http-path /api/v1/order-service/order/cancel
	CancelOrder(ctx context.Context, req *CancelOrderRequest) (*CancelOrderResponse, error)
	// OrderDetail 获取订单详情
	// @http-method get
	// @http-path /api/v1/order-service/order/detail
	OrderDetail(ctx context.Context, req *GetOrderDetailRequest) (*GetOrderDetailResponse, error)
	// Update 更新订单
	// @http-method put
	// @http-path /api/v1/order-service/order
	Update(ctx context.Context, req *UpdateOrderRequest) (*UpdateOrderResponse, error)
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

// UpdateOrderRequest 更新订单请求
type UpdateOrderRequest struct {
	OrderInfo []OrderItem `json:"order_info"`
}

// UpdateOrderResponse 更新订单响应
type UpdateOrderResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
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
