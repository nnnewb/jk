package stringsvc

import "context"

type Good struct {
	Name  string `json:"name,omitempty"`
	Price string `json:"price,omitempty"`
}

type BuyRequest struct {
	Good Good `json:"good"`
}

type BuyResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	OrderID string `json:"order_id,omitempty"`
}

type UppercaseRequest struct {
	Name  string `json:"name,omitempty"`
	Name2 string `json:"name_2,omitempty"`
}

type UppercaseResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Text    string `json:"text,omitempty"`
}

type LowercaseRequest struct {
	Name  string `json:"name,omitempty"`
	Name2 string `json:"name_2,omitempty"`
}

type LowercaseResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Text    string `json:"text,omitempty"`
}

type JoinRequest struct {
	Parts []string `json:"parts,omitempty"`
}

type JoinResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Text    string `json:"text,omitempty"`
}

type Join2Request struct {
	Parts map[string]string `json:"parts,omitempty"`
}

type Join2Response struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Text    string `json:"text,omitempty"`
}

type Join3Request struct {
	Parts map[string]string `json:"parts,omitempty"`
}

type Join3Response struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Text    string `json:"text,omitempty"`
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
