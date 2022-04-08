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
		serviceName     string
		packagePath     string
		transportName   string
		serviceOutDir   string
		transportOutDir string
	)
	flag.StringVar(&packagePath, "package", "", "input go package path")
	flag.StringVar(&serviceName, "service", "", "service interface name")
	flag.StringVar(&transportName, "transport", "", "transport name")
	flag.StringVar(&serviceOutDir, "service-outdir", "", "output folder path")
	flag.StringVar(&transportOutDir, "transport-outdir", "", "transport output folder path")
	flag.Parse()
	requireCliOption("package", packagePath)
	requireCliOption("service", serviceName)
	requireCliOption("service-outdir", serviceOutDir)

	if transportName != "" && transportOutDir == "" {
		flag.Usage()
		panic(fmt.Errorf("-transport-outdir is required when generating transport code"))
	}

	g := generator.NewJKGenerator(serviceName, packagePath)
	err := g.Parse()
	if err != nil {
		log.Fatal(err)
	}

	err = g.GenerateService("", serviceOutDir)
	if err != nil {
		log.Fatal(err)
	}

	if transportName != "" {
		err = g.GenerateTransport(transportName, transportOutDir)
		if err != nil {
			log.Fatal(err)
		}
	}
}
