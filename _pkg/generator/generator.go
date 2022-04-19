package generator

import (
	"fmt"
	"go/importer"
	"go/token"
	"go/types"
	"path/filepath"

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
	fst           *token.FileSet
	pkgPath       string // package import path. e.g. github.com/nnnewb/jk
	pkgLocalPath  string // package local path. e.g. ./
	pkgTypes      *types.Package
	svcName       string
	svcType       *types.Named
	typesImporter types.Importer
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
	svcTypeName, ok := svcTypeLookupResult.(*types.TypeName)
	if !ok {
		return fmt.Errorf("%v: not a type name", svcTypeLookupResult)
	}

	svcNamedType, ok := svcTypeName.Type().(*types.Named)
	if !ok {
		return fmt.Errorf("%v: not a named type", svcTypeName)
	}

	// svc type
	j.svcType = svcNamedType

	// find local path of package
	j.pkgLocalPath = filepath.Dir(fst.Position(svcTypeLookupResult.Pos()).Filename)

	return nil
}

func (j *JKGenerator) GenerateService(drv string) error {
	d, ok := driver.ServiceGenDrivers[drv]
	if !ok {
		return fmt.Errorf("%s: driver not exists", drv)
	}

	req := driver.NewServiceGenerateRequest(j.fst, j.pkgTypes, j.svcName, j.svcType, j.pkgLocalPath, j.typesImporter)

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

func (j *JKGenerator) GenerateTransport(drv string) error {
	d, ok := driver.TransportGenDrivers[drv]
	if !ok {
		return fmt.Errorf("%s: driver not exists", drv)
	}

	req := driver.NewServiceGenerateRequest(j.fst, j.pkgTypes, j.svcName, j.svcType, j.pkgLocalPath, j.typesImporter)
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
