// This file is generated by jk, DO NOT EDIT.

package endpoint

import (
	"context"
	endpoint "github.com/go-kit/kit/endpoint"
	service "stringsvc/service"
)

type ServiceEndpoints struct {
	buy       endpoint.Endpoint
	join      endpoint.Endpoint
	join2     endpoint.Endpoint
	lowercase endpoint.Endpoint
	uppercase endpoint.Endpoint
}

func (s ServiceEndpoints) Buy(ctx context.Context, good service.Good) (orderID string, err error) {
	request := &BuyRequest{Good: good}
	resp, err := s.buy(ctx, request)
	if err != nil {
		return "", err
	}
	response := resp.(*BuyResponse)
	orderID = response.OrderID
	return orderID, nil
}

func (s ServiceEndpoints) Join(ctx context.Context, parts []string) (text string, err error) {
	request := &JoinRequest{Parts: parts}
	resp, err := s.join(ctx, request)
	if err != nil {
		return "", err
	}
	response := resp.(*JoinResponse)
	text = response.Text
	return text, nil
}

func (s ServiceEndpoints) Join2(ctx context.Context, parts map[string]string) (text string, err error) {
	request := &Join2Request{Parts: parts}
	resp, err := s.join2(ctx, request)
	if err != nil {
		return "", err
	}
	response := resp.(*Join2Response)
	text = response.Text
	return text, nil
}

func (s ServiceEndpoints) Lowercase(ctx context.Context, name string, name2 string) (text string, err error) {
	request := &LowercaseRequest{
		Name:  name,
		Name2: name2,
	}
	resp, err := s.lowercase(ctx, request)
	if err != nil {
		return "", err
	}
	response := resp.(*LowercaseResponse)
	text = response.Text
	return text, nil
}

func (s ServiceEndpoints) Uppercase(ctx context.Context, name string, name2 string) (text string, err error) {
	request := &UppercaseRequest{
		Name:  name,
		Name2: name2,
	}
	resp, err := s.uppercase(ctx, request)
	if err != nil {
		return "", err
	}
	response := resp.(*UppercaseResponse)
	text = response.Text
	return text, nil
}