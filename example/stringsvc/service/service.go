package service

import "context"

type Good struct {
	Name  string
	Price string
}

type Service interface {
	Buy(ctx context.Context, good Good) (orderID string, err error)
	Uppercase(ctx context.Context, name string, name2 string) (text string, err error)
	Lowercase(ctx context.Context, name string, name2 string) (text string, err error)
	Join(ctx context.Context, parts []string) (text string, err error)
	Join2(ctx context.Context, parts map[string]string) (text string, err error)
	Join3(ctx context.Context, parts *map[string]string) (text string, err error)
}
