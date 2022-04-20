package main

import (
	"log"
	"net/http"
	"net/rpc"
	"stringsvc/service"
	"stringsvc/service/transport/netrpc"
)

func main() {
	s := &service.ServerLogic{}
	server := netrpc.NewServiceBinding(s)
	err := rpc.DefaultServer.Register(server)
	if err != nil {
		log.Fatal(err)
	}

	rpc.HandleHTTP()
	err = http.ListenAndServe("127.0.0.1:8888", rpc.DefaultServer)
	if err != nil {
		log.Fatal(err)
	}
}
