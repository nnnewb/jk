package gensvc

import (
	"fmt"
	"go/types"
	"log"

	"github.com/dave/jennifer/jen"
	"github.com/iancoleman/strcase"
	"github.com/nnnewb/jk/pkg/gen/genreq"
	"github.com/nnnewb/jk/pkg/gen/genutils"
	"github.com/nnnewb/jk/pkg/typecheck"
)

func GenParamsResultsStruct(f *jen.File, req *genreq.GenRequest) error {
outer:
	for i := 0; i < req.GetServiceInterface().NumMethods(); i++ {
		method := req.GetServiceInterface().Method(i)
		signature := method.Type().(*types.Signature)
		if !method.Exported() {
			return nil
		}

		fields := make([]jen.Code, 0, signature.Params().Len())
		for i := 0; i < signature.Params().Len(); i++ {
			param := signature.Params().At(i)
			if !typecheck.CanUse(param.Type()) {
				if !types.Identical(param.Type(), req.GetContextType()) && !types.Identical(param.Type(), req.GetErrorType()) {
					log.Printf("omit method %s, bad parameter %s %s", method.Name(), param.Name(), param.Type())
					continue outer
				}
				// do not add context.Context and error to struct
				continue
			}

			fields = append(fields, jen.Id(strcase.ToCamel(param.Name())).Add(genutils.Qual(param.Type())))
		}

		f.Type().Id(fmt.Sprintf("%sRequest", method.Name())).Struct(fields...).Line()

		fields = make([]jen.Code, 0, signature.Results().Len())
		for i := 0; i < signature.Results().Len(); i++ {
			param := signature.Results().At(i)
			if !typecheck.CanUse(param.Type()) {
				if !types.Identical(param.Type(), req.GetContextType()) && !types.Identical(param.Type(), req.GetErrorType()) {
					log.Printf("omit method %s, bad parameter %s %s", method.Name(), param.Name(), param.Type())
					continue outer
				}
				// do not add context.Context and error to struct
				continue
			}

			fields = append(fields, jen.Id(strcase.ToCamel(param.Name())).Add(genutils.Qual(param.Type())))
		}

		f.Type().Id(fmt.Sprintf("%sResponse", method.Name())).Struct(fields...).Line()
	}

	return nil
}
