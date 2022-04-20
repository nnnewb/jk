package gencore

import (
	"go/token"
	"go/types"
	"log"
	"path/filepath"
)

type GenRequest struct {
	Fst      *token.FileSet
	Importer types.Importer
	Svc      *types.TypeName
}

func (g *GenRequest) GetServiceLocalPath() string {
	return filepath.Dir(g.Fst.Position(g.Svc.Pos()).Filename)
}

func (g *GenRequest) GetServiceInterface() *types.Interface {
	return g.Svc.Type().(*types.Named).Underlying().(*types.Interface)
}

func (g *GenRequest) GetContextType() types.Type {
	ctxPkg, err := g.Importer.Import("context")
	if err != nil {
		log.Fatal(err)
	}

	ctxTypeNameObj := ctxPkg.Scope().Lookup("Context")
	return ctxTypeNameObj.(*types.TypeName).Type().(*types.Named)
}

func (g *GenRequest) GetErrorType() types.Type {
	errTypeNameObj := types.Universe.Lookup("error")
	return errTypeNameObj.(*types.TypeName).Type().(*types.Named)
}
