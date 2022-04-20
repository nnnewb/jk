package main

import (
	"flag"
	"fmt"
	"go/importer"
	"go/token"
	"go/types"
	"log"

	"github.com/nnnewb/jk/pkg/gen/gencore"
	_ "github.com/nnnewb/jk/pkg/gen/plugins"
)

type StringSliceFlag struct {
	flags []string
}

func (s *StringSliceFlag) String() string {
	return ""
}

func (s *StringSliceFlag) Set(v string) error {
	s.flags = append(s.flags, v)
	return nil
}

var (
	pkgPath, svcName string
	genTargets       StringSliceFlag
)

func init() {
	flag.Var(&genTargets, "gen", "generate target.")
	flag.StringVar(&pkgPath, "package", "", "go package path. e.g. github.com/nnnewb/jk")
	flag.StringVar(&svcName, "service", "", "service interface name")
}

func main() {
	flag.Parse()

	if pkgPath == "" {
		flag.Usage()
		log.Fatal("-package must be set")
	}

	if svcName == "" {
		flag.Usage()
		log.Fatal("-service must be set")
	}

	if len(genTargets.flags) == 0 {
		flag.Usage()
		log.Fatal("-gen must be set")
	}

	fst, importer, svc, err := findService(pkgPath, svcName)
	if err != nil {
		log.Fatal(err)
	}

	req := &gencore.GenRequest{
		Fst:      fst,
		Importer: importer,
		Svc:      svc,
	}

	data := gencore.NewPluginData(req)

	for _, name := range genTargets.flags {
		plugin := gencore.GetPlugin(name)
		if plugin == nil {
			log.Fatalf("%s: no such plugin", name)
		}

		if err = plugin.Generate(data); err != nil {
			log.Fatal(err)
		}
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
