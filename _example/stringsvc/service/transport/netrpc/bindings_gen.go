package netrpc

import (
	"context"
	kitendpoint "github.com/go-kit/kit/endpoint"
	service "stringsvc/service"
	endpoint "stringsvc/service/endpoint"
)

type ServiceBinding struct {
	BuyEndpoint       kitendpoint.Endpoint
	JoinEndpoint      kitendpoint.Endpoint
	Join2Endpoint     kitendpoint.Endpoint
	LowercaseEndpoint kitendpoint.Endpoint
	UppercaseEndpoint kitendpoint.Endpoint
}

func NewServiceBinding(svc service.Service) *ServiceBinding {
	return &ServiceBinding{
		BuyEndpoint:       endpoint.MakeBuyEndpoint(svc),
		Join2Endpoint:     endpoint.MakeJoin2Endpoint(svc),
		JoinEndpoint:      endpoint.MakeJoinEndpoint(svc),
		LowercaseEndpoint: endpoint.MakeLowercaseEndpoint(svc),
		UppercaseEndpoint: endpoint.MakeUppercaseEndpoint(svc),
	}
}

func (b ServiceBinding) Buy(req endpoint.BuyRequest, response *endpoint.BuyResponse) error {
	resp, err := b.BuyEndpoint(context.Background(), req)
	if err != nil {
		return err
	}
	*response = *resp.(*endpoint.BuyResponse)
	return nil
}

func (b ServiceBinding) Join(req endpoint.JoinRequest, response *endpoint.JoinResponse) error {
	resp, err := b.JoinEndpoint(context.Background(), req)
	if err != nil {
		return err
	}
	*response = *resp.(*endpoint.JoinResponse)
	return nil
}

func (b ServiceBinding) Join2(req endpoint.Join2Request, response *endpoint.Join2Response) error {
	resp, err := b.Join2Endpoint(context.Background(), req)
	if err != nil {
		return err
	}
	*response = *resp.(*endpoint.Join2Response)
	return nil
}

func (b ServiceBinding) Lowercase(req endpoint.LowercaseRequest, response *endpoint.LowercaseResponse) error {
	resp, err := b.LowercaseEndpoint(context.Background(), req)
	if err != nil {
		return err
	}
	*response = *resp.(*endpoint.LowercaseResponse)
	return nil
}

func (b ServiceBinding) Uppercase(req endpoint.UppercaseRequest, response *endpoint.UppercaseResponse) error {
	resp, err := b.UppercaseEndpoint(context.Background(), req)
	if err != nil {
		return err
	}
	*response = *resp.(*endpoint.UppercaseResponse)
	return nil
}
