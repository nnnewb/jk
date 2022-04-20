package main

import (
	"flag"
	"fmt"
	"go/importer"
	"go/token"
	"go/types"
	"log"
	"reflect"

	"github.com/nnnewb/jk/pkg/gen/gencore"
	"github.com/nnnewb/jk/pkg/gen/gensvc"
	"github.com/nnnewb/jk/pkg/gen/transports/genrpc"
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

	fst, importer, svc, err := findService(packagePath, serviceName)
	if err != nil {
		log.Fatal(err)
	}

	req := &gencore.GenRequest{
		Fst:      fst,
		Importer: importer,
		Svc:      svc,
	}

	data := gencore.NewPluginData(req)

	if err := gensvc.GenerateEndpoint(data); err != nil {
		log.Fatal(err)
	}

	if err := genrpc.GenerateBindings(data); err != nil {
		log.Fatal(err)
	}

	if err := data.WriteToDisk(); err != nil {
		log.Fatal(err)
	}
}

func findService(pkgPath, serviceName string) (*token.FileSet, types.Importer, *types.TypeName, error) {
	fst := token.NewFileSet()

	i := importer.ForCompiler(fst, "source", nil)
	pkg, err := i.Import(pkgPath)
	if err != nil {
		return nil, nil, nil, err
	}

	// service type lookup
	svcTypeLookupResult := pkg.Scope().Lookup(serviceName)
	svcTypeName, ok := svcTypeLookupResult.(*types.TypeName)
	if !ok {
		return nil, nil, nil, fmt.Errorf("%v: not a type name", svcTypeLookupResult)
	}

	return fst, i, svcTypeName, nil
}
