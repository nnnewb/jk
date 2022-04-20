package endpoint

import (
	"context"
	endpoint "github.com/go-kit/kit/endpoint"
	service "stringsvc/service"
)

func MakeBuyEndpoint(svc service.Service) endpoint.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		request := req.(*BuyRequest)
		orderID, err := svc.Buy(ctx, request.Good)
		if err != nil {
			return nil, err
		}
		resp := &BuyResponse{OrderID: orderID}
		return resp, nil
	}
}

func MakeJoinEndpoint(svc service.Service) endpoint.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		request := req.(*JoinRequest)
		text, err := svc.Join(ctx, request.Parts)
		if err != nil {
			return nil, err
		}
		resp := &JoinResponse{Text: text}
		return resp, nil
	}
}

func MakeJoin2Endpoint(svc service.Service) endpoint.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		request := req.(*Join2Request)
		text, err := svc.Join2(ctx, request.Parts)
		if err != nil {
			return nil, err
		}
		resp := &Join2Response{Text: text}
		return resp, nil
	}
}

func MakeLowercaseEndpoint(svc service.Service) endpoint.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		request := req.(*LowercaseRequest)
		text, err := svc.Lowercase(ctx, request.Name, request.Name2)
		if err != nil {
			return nil, err
		}
		resp := &LowercaseResponse{Text: text}
		return resp, nil
	}
}

func MakeUppercaseEndpoint(svc service.Service) endpoint.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		request := req.(*UppercaseRequest)
		text, err := svc.Uppercase(ctx, request.Name, request.Name2)
		if err != nil {
			return nil, err
		}
		resp := &UppercaseResponse{Text: text}
		return resp, nil
	}
}
