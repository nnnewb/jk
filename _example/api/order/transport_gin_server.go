// Code generated by jk generate transport -t OrderService --protocol http --server --language go --framework gin --swagger; DO NOT EDIT.

package order

import (
	"fmt"
	gin "github.com/gin-gonic/gin"
	"github.com/go-kit/kit/endpoint"
)

type RequestDecoder func(c *gin.Context, req any) error

func QueryStringDecoder(c *gin.Context, req any) error {
	return c.BindQuery(req)
}

func JSONBodyDecoder(c *gin.Context, req any) error {
	return c.BindJSON(req)
}

type ResponseEncoder func(c *gin.Context, resp any)

func JSONBodyEncoder(c *gin.Context, resp any) {
	c.JSON(200, resp)
}

func Handler[Request any](ep endpoint.Endpoint, decoder RequestDecoder, encoder ResponseEncoder) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req = new(Request)
		err := decoder(c, req)
		if err != nil {
			c.AbortWithStatusJSON(400, gin.H{
				"code":    -1,
				"message": fmt.Sprintf("unable to parse request payload, error %v", err),
			})
			return
		}

		resp, err := ep(c.Request.Context(), req)
		encoder(c, resp)
	}
}

type GinServerSet struct {
	CancelOrderHandler gin.HandlerFunc
	CreateOrderHandler gin.HandlerFunc
	OrderDetailHandler gin.HandlerFunc
	UpdateHandler      gin.HandlerFunc
}

func NewGinServerSet(eps EndpointSet) *GinServerSet {
	return &GinServerSet{
		CancelOrderHandler: Handler[CancelOrderRequest](eps.CancelOrderEndpoint, JSONBodyDecoder, JSONBodyEncoder),
		CreateOrderHandler: Handler[CreateOrderRequest](eps.CreateOrderEndpoint, JSONBodyDecoder, JSONBodyEncoder),
		OrderDetailHandler: Handler[GetOrderDetailRequest](eps.OrderDetailEndpoint, QueryStringDecoder, JSONBodyEncoder),
		UpdateHandler:      Handler[UpdateOrderRequest](eps.UpdateEndpoint, JSONBodyDecoder, JSONBodyEncoder),
	}
}

func (s *GinServerSet) Register(router gin.IRouter) {
	router.POST("/api/v1/order-service/order/cancel", s.CancelOrderHandler)
	router.POST("/api/v1/order-service/order", s.CreateOrderHandler)
	router.GET("/api/v1/order-service/order/detail", s.OrderDetailHandler)
	router.PUT("/api/v1/order-service/order", s.UpdateHandler)
}