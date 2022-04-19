package driver

import "fmt"

type TransportGenDriver interface {
	GenerateTransport(req *GenerateRequest) error
}

type ServiceGenDriver interface {
	GenerateService(req *GenerateRequest) error
}

var (
	TransportGenDrivers map[string]TransportGenDriver
	ServiceGenDrivers   map[string]ServiceGenDriver
)

func RegisterTransportGenDriver(name string, driver TransportGenDriver) {
	if TransportGenDrivers == nil {
		TransportGenDrivers = make(map[string]TransportGenDriver)
	}

	if _, ok := TransportGenDrivers[name]; ok {
		panic(fmt.Errorf("transport gen driver %s has already registered", name))
	}

	TransportGenDrivers[name] = driver
}

func RegisterServiceGenDriver(name string, driver ServiceGenDriver) {
	if ServiceGenDrivers == nil {
		ServiceGenDrivers = make(map[string]ServiceGenDriver)
	}

	if _, ok := ServiceGenDrivers[name]; ok {
		panic(fmt.Errorf("service gen driver %s has already registered", name))
	}

	ServiceGenDrivers[name] = driver
}
