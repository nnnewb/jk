package main

import (
	"flag"
	"fmt"
	"go/importer"
	"go/token"
	"go/types"
	"log"
	"os"
	"path/filepath"
	"reflect"

	"github.com/dave/jennifer/jen"
	"github.com/nnnewb/jk/pkg/gen/genreq"
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

	req := &genreq.GenRequest{
		Fst:      fst,
		Importer: importer,
		Svc:      svc,
	}

	if err := generateService(req); err != nil {
		log.Fatal(err)
	}

	f := jen.NewFile("netrpc")
	err = genrpc.GenBindings(f, req)
	if err != nil {
		log.Fatal(err)
	}

	err = f.Render(os.Stdout)
	if err != nil {
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

func generateService(req *genreq.GenRequest) error {
	f := jen.NewFile("endpoint")

	err := gensvc.GenParamsResultsStruct(f, req)
	if err != nil {
		log.Fatal(err)
	}

	err = gensvc.GenEndpointMaker(f, req)
	if err != nil {
		log.Fatal(err)
	}

	dst := filepath.Join(req.GetServiceLocalPath(), "endpoint")
	info, err := os.Stat(dst)
	if err != nil {
		if os.IsNotExist(err) {
			os.MkdirAll(dst, 0o755)
		}
		return err
	}
	if !info.IsDir() {
		return fmt.Errorf("%s is not dir", dst)
	}

	file, err := os.OpenFile(filepath.Join(dst, "endpoint.go"), os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o644)
	if err != nil {
		return err
	}
	defer file.Close()

	err = f.Render(file)
	if err != nil {
		log.Fatal(err)
	}

	return nil
}
