package utils

import (
	"go/types"
	"strings"

	"github.com/dave/jennifer/jen"
	"github.com/iancoleman/strcase"
)

// GenCopyVar generate a assign statement for given selector.
// generated assign statement like `id1.v = id2.v`.
//
// only for serializable type, vars like pointer to pointer are not
// allowed.
func GenCopyVar(srcVar *types.Var, dst, src string) ([]*jen.Statement, error) {
	ret := make([]*jen.Statement, 0)
	paramName := strcase.ToCamel(srcVar.Name())
	switch typ := srcVar.Type().(type) {
	case *types.Basic, *types.Map, *types.Slice, *types.Array:
		// TODO: shallow copy map and slice, need reconsider it.
		ret = append(ret, jen.Id(dst).Dot(paramName).Op("=").Id(src).Dot(paramName))
		return ret, nil
	case *types.Pointer:
		// TODO: only allowed pointer to struct, need reconsider it.
		stmt, err := GenCopyVar(types.NewVar(srcVar.Pos(), srcVar.Pkg(), srcVar.Name(), typ.Elem()), dst, src)
		if err != nil {
			return nil, err
		}

		ret = append(ret, jen.If(jen.Id(src).Dot(paramName).Op("!=").Nil()).BlockFunc(func(g *jen.Group) {
			for _, v := range stmt {
				g.Add(v)
			}
		}))
		return ret, nil
	case *types.Named:
		stmt, err := GenCopyVar(types.NewVar(srcVar.Pos(), srcVar.Pkg(), srcVar.Name(), typ.Underlying()), dst, src)
		if err != nil {
			return nil, err
		}

		if IsStruct(typ.Underlying()) {
			// FIXME: typ.Obj().Name() is name of src field but not dst field, this line can generate wrong code.
			ret = append(ret, jen.Id(dst).Dot(paramName).Op("=").Op("&").Id(typ.Obj().Name()).Block())
			ret = append(ret, stmt...)
		} else {
			// FIXME: aliases should be qualified.
			ret = append(ret, jen.Id(src).Dot(paramName).Op("=").Id(typ.Obj().Name()).Params(jen.Id(dst).Dot(paramName)))
		}

		return ret, nil
	case *types.Struct:
		for i := 0; i < typ.NumFields(); i++ {
			field := typ.Field(i)
			stmt, err := GenCopyVar(field, strings.Join([]string{dst, paramName}, "."), strings.Join([]string{src, paramName}, "."))
			if err != nil {
				ret = append(ret, jen.Commentf("%s(%s) omitted, error %s", srcVar.Name(), srcVar.Type(), err.Error()))
			} else {
				ret = append(ret, stmt...)
			}
		}
		return ret, nil
	default:
		return ret, nil
	}
}
