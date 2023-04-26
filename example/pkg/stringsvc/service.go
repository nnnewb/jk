package stringsvc

import "context"

type Good struct {
	Name  string
	Price string
}

type BuyRequest struct {
	Good Good
}

type BuyResponse struct {
	OrderID string
}

type UppercaseRequest struct {
	Name  string
	Name2 string
}

type UppercaseResponse struct {
	Text string
}

type LowercaseRequest struct {
	Name  string
	Name2 string
}

type LowercaseResponse struct {
	Text string
}

type JoinRequest struct {
	Parts []string
}

type JoinResponse struct {
	Text string
}

type Join2Request struct {
	Parts map[string]string
}

type Join2Response struct {
	Text string
}

type Join3Request struct {
	Parts map[string]string
}

type Join3Response struct {
	Text string
}

//go:generate jk generate all -t Service
type Service interface {
	Buy(ctx context.Context, req BuyRequest) (res BuyResponse, err error)
	Uppercase(ctx context.Context, req UppercaseRequest) (res UppercaseResponse, err error)
	Lowercase(ctx context.Context, req LowercaseRequest) (res LowercaseResponse, err error)
	Join(ctx context.Context, req JoinRequest) (res JoinResponse, err error)
	Join2(ctx context.Context, req Join2Request) (res Join2Response, err error)
	Join3(ctx context.Context, req Join3Request) (res Join3Response, err error)
}
