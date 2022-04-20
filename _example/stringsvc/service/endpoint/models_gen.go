package endpoint

import service "stringsvc/service"

type BuyRequest struct {
	Good service.Good
}

type BuyResponse struct {
	OrderID string
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

type LowercaseRequest struct {
	Name  string
	Name2 string
}

type LowercaseResponse struct {
	Text string
}

type UppercaseRequest struct {
	Name  string
	Name2 string
}

type UppercaseResponse struct {
	Text string
}
