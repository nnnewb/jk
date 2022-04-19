package genrpc

import (
	"fmt"
	"go/types"

	"github.com/dave/jennifer/jen"
	"github.com/nnnewb/jk/pkg/gen/genreq"
	"github.com/nnnewb/jk/pkg/gen/genutils"
)

func genBindStruct(t types.TypeName) jen.Code {
	return jen.Type().Id(t.Name()).Struct()
}

func genBindFunc() {
}

func GenBindings(f *jen.File, req *genreq.GenRequest) error {
	// type XxxBinding struct { svc Xxx }
	f.Type().Id(fmt.Sprintf("%sBinding", req.Svc.Name())).Struct(jen.Id("svc").Add(genutils.Qual(req.Svc.Type()))).Line()

	for i := 0; i < req.GetServiceInterface().NumMethods(); i++ {
		method := req.GetServiceInterface().Method(i)
		if genutils.IsOmitMethod(method, req.GetContextType(), req.GetErrorType()) {
			continue
		}

		f.Func().
			Params(jen.Id("b").Id(fmt.Sprintf("%sBinding", req.Svc.Name()))).
			Id(method.Name()).
			Params(
				jen.Id("req").Qual(req.Svc.Pkg().Path()+"/endpoint", fmt.Sprintf("%sRequest", method.Name())),
				jen.Id("response").Op("*").Qual(req.Svc.Pkg().Path()+"/endpoint", fmt.Sprintf("%sResponse", method.Name())),
			).
			Params(jen.Error()).
			BlockFunc(func(g *jen.Group) {
				// resp, err := b.svc.XxxEndpoint(context.Background(), req)
				g.List(jen.Id("resp"), jen.Id("err")).Op(":=").Add(jen.Id("b").Dot("svc").Dot(method.Name()+"Endpoint")).Call(jen.Qual("context", "Background").Call(), jen.Id("req"))
				// if err != nil { return err }
				g.If(jen.Id("err").Op("!=").Nil()).Block(jen.Return(jen.Id("err")))
				// (*resp) = *resp.(*XxxResponse)
				g.Op("*").Id("response").Op("=").Op("*").Id("resp").Assert(jen.Op("*").Qual(req.Svc.Pkg().Path()+"/endpoint", fmt.Sprintf("%sResponse", method.Name())))
				// return nil
				g.Return(jen.Nil())
			}).Line()
	}
	return nil
}
