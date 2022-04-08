package generator

import (
	"fmt"
	"go/ast"
	"go/importer"
	"go/parser"
	"go/token"
	"go/types"
	"io/fs"
	"regexp"

	_ "github.com/nnnewb/jk/pkg/generator/contrib/service"
	_ "github.com/nnnewb/jk/pkg/generator/contrib/transports/grpc"
	"github.com/nnnewb/jk/pkg/generator/driver"
)

type Generator interface {
	Parse() error
	GenerateService(drv string) error
	GenerateTransport(drv string) error
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

	packages, err := parser.ParseDir(fst, j.pkgPath, func(fi fs.FileInfo) bool {
		return !regexp.MustCompile(`.*_test\.go`).MatchString(fi.Name())
	}, parser.ParseComments)
	if err != nil {
		return err
	}
	if len(packages) > 1 {
		return fmt.Errorf("more than one package found in path %s", j.pkgPath)
	}

	var pkg *ast.Package
	for _, p := range packages {
		pkg = p
		break
	}

	// type-check
	cfg := &types.Config{Importer: importer.Default()}
	files := []*ast.File{}
	for _, file := range pkg.Files {
		files = append(files, file)
	}

	j.pkgTypes, err = cfg.Check(j.pkgPath, fst, files, nil)
	if err != nil {
		return err
	}

	// service type lookup
	svcTypeLookupResult := j.pkgTypes.Scope().Lookup(j.svcName)
	if types.IsInterface(svcTypeLookupResult.Type()) {
		j.svcTypes = svcTypeLookupResult.(*types.TypeName).Type().(*types.Named).Underlying().(*types.Interface)
	} else {
		return fmt.Errorf("%s: type not found", j.svcName)
	}

	return nil
}

func (j *JKGenerator) GenerateService(drv string) error {
	d, ok := driver.ServiceGenDrivers[drv]
	if !ok {
		return fmt.Errorf("%s: driver not exists", drv)
	}

	req := driver.NewServiceGenerateRequest(j.fst, j.pkgTypes, j.svcName, j.svcTypes)

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

	req := driver.NewServiceGenerateRequest(j.fst, j.pkgTypes, j.svcName, j.svcTypes)
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
