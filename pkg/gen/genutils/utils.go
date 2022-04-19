package genutils

import (
	"go/types"

	"github.com/dave/jennifer/jen"
	"github.com/nnnewb/jk/pkg/typecheck"
)

// IsOmitMethod return code generation should omit this function or not.
func IsOmitMethod(f *types.Func, ctxType, errType types.Type) bool {
	signature := f.Type().(*types.Signature)
	if !f.Exported() {
		return true
	}

	for i := 0; i < signature.Params().Len(); i++ {
		param := signature.Params().At(i)
		if !typecheck.CanUse(param.Type()) && !types.Identical(param.Type(), ctxType) && !types.Identical(param.Type(), errType) {
			return true
		}
	}

	return false
}

// Qual if given type is types.Named and not builtin type, return a jen.Qual for it.
func Qual(tp types.Type) jen.Code {
	switch t := tp.(type) {
	case *types.Named:
		if t.Obj().Pkg() != nil {
			return jen.Qual(t.Obj().Pkg().Path(), t.Obj().Name())
		}
		return jen.Id(t.Obj().Name())
	default:
		return jen.Id(tp.String())
	}
}
