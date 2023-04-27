package main

import (
	"github.com/julienschmidt/httprouter"
	stringsvc1 "github.com/nnnewb/jk/example/internal/stringsvc"
	"github.com/nnnewb/jk/example/pkg/stringsvc"
	"log"
	"net/http"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	router := httprouter.New()
	stringsvc.Register(stringsvc1.Svc{}, router)
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
	err := http.ListenAndServeTLS("127.0.0.1:8888", "cert.pem", "key.pem", router)
	if err != nil {
		log.Fatalf("Serve failed, error %+v", err)
	}
}
