package stringsvc

import (
	"context"
	service2 "github.com/nnnewb/jk/example/pkg/stringsvc"
)

type Svc struct{}

func (s Svc) Buy(ctx context.Context, req service2.BuyRequest) (res service2.BuyResponse, err error) {
	// TODO implement me
	panic("implement me")
}

func (s Svc) Uppercase(ctx context.Context, req service2.UppercaseRequest) (res service2.UppercaseResponse, err error) {
	// TODO implement me
	panic("implement me")
}

func (s Svc) Lowercase(ctx context.Context, req service2.LowercaseRequest) (res service2.LowercaseResponse, err error) {
	// TODO implement me
	panic("implement me")
}

func (s Svc) Join(ctx context.Context, req service2.JoinRequest) (res service2.JoinResponse, err error) {
	// TODO implement me
	panic("implement me")
}

func (s Svc) Join2(ctx context.Context, req service2.Join2Request) (res service2.Join2Response, err error) {
	// TODO implement me
	panic("implement me")
}

func (s Svc) Join3(ctx context.Context, req service2.Join3Request) (res service2.Join3Response, err error) {
	// TODO implement me
	panic("implement me")
}
