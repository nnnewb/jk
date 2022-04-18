package main

import (
	"flag"
	"fmt"
	"log"
	"reflect"

	"github.com/nnnewb/jk/pkg/generator"
)

func requireCliOption(name string, option interface{}) {
	if reflect.ValueOf(option).IsZero() {
		flag.Usage()
		panic(fmt.Errorf("option %s must be set", name))
	}
}

func main() {
	var (
		serviceName   string
		packagePath   string
		transportName string
	)
	flag.StringVar(&packagePath, "package", "", "go package path. e.g. github.com/nnnewb/jk")
	flag.StringVar(&serviceName, "service", "", "service interface name")
	flag.StringVar(&transportName, "transport", "", "transport name")
	flag.Parse()
	requireCliOption("package", packagePath)
	requireCliOption("service", serviceName)

	g := generator.NewJKGenerator(serviceName, packagePath)
	err := g.Parse()
	if err != nil {
		log.Fatal(err)
	}

	err = g.GenerateService("")
	if err != nil {
		log.Fatal(err)
	}

	if transportName != "" {
		err = g.GenerateTransport(transportName)
		if err != nil {
			log.Fatal(err)
		}
	}
}
