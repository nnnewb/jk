package service

import "context"

type ServerLogic struct {
}

func (s *ServerLogic) Buy(ctx context.Context, good Good) (orderID string, err error) {
	panic("not implemented") // TODO: Implement
}

func (s *ServerLogic) Uppercase(ctx context.Context, name string, name2 string) (text string, err error) {
	panic("not implemented") // TODO: Implement
}

func (s *ServerLogic) Lowercase(ctx context.Context, name string, name2 string) (text string, err error) {
	panic("not implemented") // TODO: Implement
}

func (s *ServerLogic) Join(ctx context.Context, parts []string) (text string, err error) {
	panic("not implemented") // TODO: Implement
}

func (s *ServerLogic) Join2(ctx context.Context, parts map[string]string) (text string, err error) {
	panic("not implemented") // TODO: Implement
}

// omit invalid signature
func (s *ServerLogic) Join3(ctx context.Context, parts *map[string]string) (text string, err error) {
	panic("not implemented") // TODO: Implement
}
