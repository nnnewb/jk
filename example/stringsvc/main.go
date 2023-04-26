package main

import (
	"context"
	"github.com/nnnewb/jk/example/stringsvc/service"
	"log"
	"net/http"
)

type svc struct{}

func (s svc) Buy(ctx context.Context, req service.BuyRequest) (res service.BuyResponse, err error) {
	// TODO implement me
	panic("implement me")
}

func (s svc) Uppercase(ctx context.Context, req service.UppercaseRequest) (res service.UppercaseResponse, err error) {
	// TODO implement me
	panic("implement me")
}

func (s svc) Lowercase(ctx context.Context, req service.LowercaseRequest) (res service.LowercaseResponse, err error) {
	// TODO implement me
	panic("implement me")
}

func (s svc) Join(ctx context.Context, req service.JoinRequest) (res service.JoinResponse, err error) {
	// TODO implement me
	panic("implement me")
}

func (s svc) Join2(ctx context.Context, req service.Join2Request) (res service.Join2Response, err error) {
	// TODO implement me
	panic("implement me")
}

func (s svc) Join3(ctx context.Context, req service.Join3Request) (res service.Join3Response, err error) {
	// TODO implement me
	panic("implement me")
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	service.Register(svc{}, http.DefaultServeMux)
	log.Println("Server now listening at https://127.0.0.1:8888/")
	// generate self-signed certificate with openssl cli
	//
	// openssl req \
	//     -x509 \
	//     -newkey rsa:4096 \
	//     -keyout key.pem \
	//     -out cert.pem \
	//     -sha256 \
	//     -days 3650 \
	//     -nodes \
	//     -subj "/C=XX/ST=StateName/L=CityName/O=CompanyName/OU=CompanySectionName/CN=CommonNameOrHostname"
	err := http.ListenAndServeTLS("127.0.0.1:8888", "cert.pem", "key.pem", http.DefaultServeMux)
	if err != nil {
		log.Fatalf("Serve failed, error %+v", err)
	}
}
