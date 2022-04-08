package generator

import (
	"fmt"
	"go/importer"
	"go/token"
	"go/types"

	_ "github.com/nnnewb/jk/pkg/generator/contrib/service"
	_ "github.com/nnnewb/jk/pkg/generator/contrib/transports/grpc"
	"github.com/nnnewb/jk/pkg/generator/driver"
)

type Generator interface {
	Parse() error
	GenerateService(drv string, output string) error
	GenerateTransport(drv string, output string) error
}

type JKGenerator struct {
	fst      *token.FileSet
	pkgPath  string
	pkgTypes *types.Package
	svcName  string
	svcTypes *types.Interface
}

func NewJKGenerator(serviceName, packagePath string) *JKGenerator {
	return &JKGenerator{
		svcName: serviceName,
		pkgPath: packagePath,
	}
}

func (j *JKGenerator) Parse() error {
	fst := token.NewFileSet()
	j.fst = fst

	pkg, err := importer.ForCompiler(fst, "source", nil).Import(j.pkgPath)
	if err != nil {
		return err
	}
	j.pkgTypes = pkg

	// service type lookup
	svcTypeLookupResult := j.pkgTypes.Scope().Lookup(j.svcName)
	if types.IsInterface(svcTypeLookupResult.Type()) {
		j.svcTypes = svcTypeLookupResult.(*types.TypeName).Type().(*types.Named).Underlying().(*types.Interface)
	} else {
		return fmt.Errorf("%s: type not found", j.svcName)
	}

	return nil
}

func (j *JKGenerator) GenerateService(drv string, output string) error {
	d, ok := driver.ServiceGenDrivers[drv]
	if !ok {
		return fmt.Errorf("%s: driver not exists", drv)
	}

	req := driver.NewServiceGenerateRequest(j.fst, j.pkgTypes, j.svcName, j.svcTypes, output)

	err := d.GenerateService(req)
	if err != nil {
		return err
	}

	err = req.SaveGeneratedFiles()
	if err != nil {
		return err
	}

	return nil
}

func (j *JKGenerator) GenerateTransport(drv string, output string) error {
	d, ok := driver.TransportGenDrivers[drv]
	if !ok {
		return fmt.Errorf("%s: driver not exists", drv)
	}

	req := driver.NewServiceGenerateRequest(j.fst, j.pkgTypes, j.svcName, j.svcTypes, output)
	err := d.GenerateTransport(req)
	if err != nil {
		return err
	}

	err = req.SaveGeneratedFiles()
	if err != nil {
		return err
	}

	return nil
}
