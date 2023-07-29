package domain

import (
	"go/ast"
	"go/types"

	"emperror.dev/errors"
	"github.com/dave/jennifer/jen"
)

type MethodAnnotations struct {
	HTTPMethod string `jk:"http-method"`
	HTTPPath   string `jk:"http-path"`
}

type Method struct {
	parent      *Service
	Func        *types.Func
	Field       *ast.Field
	Annotations *MethodAnnotations
}

func (m *Method) RequestType() types.Type {
	signature := m.Func.Type().(*types.Signature)
	if signature.Params().Len() != 2 {
		panic(errors.Errorf("unexpected function signature: params type mismatch, got %s", signature))
	}

	return signature.Params().At(1).Type()
}

func (m *Method) RequestTypeName() string {
	signature := m.Func.Type().(*types.Signature)
	if signature.Params().Len() != 2 {
		panic(errors.Errorf("unexpected function signature: params type mismatch, got %s", signature))
	}

	ptrType := signature.Params().At(1).Type().(*types.Pointer)
	named := ptrType.Elem().(*types.Named)
	return named.Obj().Name()
}

func (m *Method) ResponseType() types.Type {
	signature := m.Func.Type().(*types.Signature)
	if signature.Results().Len() != 2 {
		panic(errors.Errorf("unexpected function signature: results type mismatch, got %s", signature))
	}

	return signature.Results().At(0).Type()
}

func (m *Method) ResponseTypeName() string {
	signature := m.Func.Type().(*types.Signature)
	if signature.Params().Len() != 2 {
		panic(errors.Errorf("unexpected function signature: params type mismatch, got %s", signature))
	}

	ptrType := signature.Results().At(0).Type().(*types.Pointer)
	named := ptrType.Elem().(*types.Named)
	return named.Obj().Name()
}

func (m *Method) RequestTypeCodeJen() *jen.Statement {
	// 一般来说这个类型是 *T，也就是 ptr->named->struct
	ptr := m.RequestType().(*types.Pointer)
	named := ptr.Elem().(*types.Named)
	return jen.Qual(named.Obj().Pkg().Path(), named.Obj().Name())
}

func (m *Method) ResponseTypeCodeJen() *jen.Statement {
	// 一般来说这个类型是 *T，也就是 ptr->named->struct
	ptr := m.ResponseType().(*types.Pointer)
	named := ptr.Elem().(*types.Named)
	return jen.Qual(named.Obj().Pkg().Path(), named.Obj().Name())
}
