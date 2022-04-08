package grpc

import (
	"fmt"
	"go/types"
	"log"
	"strconv"
	"strings"

	"github.com/iancoleman/strcase"
	"github.com/nnnewb/jk/pkg/generator/driver"
	"github.com/nnnewb/jk/pkg/generator/utils"
)

func init() {
	driver.RegisterTransportGenDriver("grpc", defaultGRPCTransportGenerator{})
}

type defaultGRPCTransportGenerator struct {
	req *driver.GenerateRequest
}

func (d defaultGRPCTransportGenerator) GenerateTransport(req *driver.GenerateRequest) error {
	d.req = req
	gf := req.GenFile("transport/grpc/rpc.proto")
	gf.P("// This file is generated jk, DO NOT EDIT.\n")
	gf.P("syntax = \"proto3\";\n")
	gf.P("option go_package=\"", req.Pkg.Path(), "\";\n")

	// parameter types may depending on other types, we need define thoese types in Protocol Buffer for serialization
	messages := make(map[string]string)

	for i := 0; i < req.Svc.NumMethods(); i++ {
		method := req.Svc.Method(i)
		if !method.Exported() {
			continue
		}
		if !utils.CheckMethodSignature(method) {
			continue
		}

		methodName := method.Name()
		methodType := method.Type().(*types.Signature)

		gf.P("message ", methodName, "Request {\n")
		for i := 1; i < methodType.Params().Len(); i++ {
			param := methodType.Params().At(i)
			gf.P("    ", typeToProto(param.Type()), " ", param.Name(), " = ", strconv.FormatInt(int64(i), 10), ";\n")
			namedTypeToProtoMessage(param.Type(), messages)
		}
		gf.P("}\n\n")

		gf.P("message ", methodName, "Response {\n")
		for i := 0; i < methodType.Results().Len()-1; i++ {
			result := methodType.Results().At(i)
			gf.P("    ", typeToProto(result.Type()), " ", result.Name(), " = ", strconv.FormatInt(int64(i+1), 10), ";\n")
			namedTypeToProtoMessage(result.Type(), messages)
		}
		gf.P("}\n\n")
	}

	for _, message := range messages {
		gf.P(message)
	}

	return nil
}

func namedTypeToProtoMessageName(tp *types.Named) string {
	sb := strings.Builder{}
	if tp.Obj().Pkg() != nil {
		sb.WriteString(tp.Obj().Pkg().Name())
		sb.WriteString("_")
		sb.WriteString(tp.Obj().Name())
	} else {
		sb.WriteString(tp.Obj().Name())
	}
	return sb.String()
}

// namedTypeToProtoMessage transform types.Named to corresponding Protocol Buffer message.
//
// Only works for struct type, passing type alias will return nil.
func namedTypeToProtoMessage(tp types.Type, messages map[string]string) error {
	switch t := tp.(type) {
	case *types.Named:
		if t.Obj().IsAlias() {
			return fmt.Errorf("type %s is type alias, which is not allowed", t.Obj().Name())
		}

		if !utils.IsNamedStruct(tp) {
			return fmt.Errorf("type %s is not a named struct", t.Obj().Name())
		}

		if _, ok := messages[namedTypeToProtoMessageName(t)]; ok {
			return nil
		}

		sb := strings.Builder{}
		s := tp.Underlying().(*types.Struct)
		sb.WriteString("message ")
		sb.WriteString(namedTypeToProtoMessageName(t))
		sb.WriteString(" {\n")
		for i := 0; i < s.NumFields(); i++ {
			field := s.Field(i)
			sb.WriteString("    ")
			sb.WriteString(typeToProto(field.Type()))
			sb.WriteString(" ")
			sb.WriteString(strcase.ToSnake(field.Name()))
			sb.WriteString(" = ")
			sb.WriteString(strconv.FormatInt(int64(i), 10))
			sb.WriteString(";\n")

			namedTypeToProtoMessage(field.Type(), messages)
		}
		sb.WriteString("}\n\n")
		messages[namedTypeToProtoMessageName(t)] = sb.String()
		return nil
	case *types.Array:
		return namedTypeToProtoMessage(t.Elem(), messages)
	case *types.Slice:
		return namedTypeToProtoMessage(t.Elem(), messages)
	case *types.Pointer:
		return namedTypeToProtoMessage(t.Elem(), messages)
	case *types.Map:
		// map key must be integral or string, skip
		return namedTypeToProtoMessage(t.Elem(), messages)
	default:
		return nil
	}
}

// typeToProto transform given type to corresponding Protocol Buffer type.
func typeToProto(tp types.Type) string {
	sb := strings.Builder{}
	switch t := tp.(type) {
	case *types.Pointer:
		sb.WriteString(typeToProto(t.Elem()))
		return sb.String()
	case *types.Basic:
		switch t.Kind() {
		case types.Int, types.Int8, types.Int16, types.Int32:
			sb.WriteString("int32")
		case types.Int64:
			sb.WriteString("int64")
		case types.Uint, types.Uint8, types.Uint16, types.Uint32:
			sb.WriteString("uint32")
		case types.Uint64:
			sb.WriteString("uint64")
		case types.Bool:
			sb.WriteString("bool")
		case types.String:
			sb.WriteString("string")
		default:
			panic(fmt.Errorf("unexpected kind %d of basic type", t.Kind()))
		}
		return sb.String()
	case *types.Named:
		if t.Obj().IsAlias() {
			return typeToProto(t.Underlying())
		}

		return namedTypeToProtoMessageName(t)
	case *types.Array, *types.Slice:
		_, err := sb.WriteString("repeated ")
		if err != nil {
			log.Fatal(err)
		}

		var elem types.Type
		if arr, ok := t.(*types.Array); ok {
			elem = arr.Elem()
		} else if sl, ok := t.(*types.Slice); ok {
			elem = sl.Elem()
		}

		_, err = sb.WriteString(typeToProto(elem))
		if err != nil {
			log.Fatal(err)
		}
		return sb.String()
	case *types.Map:
		if (utils.IsIntegral(t.Key()) || utils.IsString(t.Key())) && !utils.IsMap(t.Elem()) && !utils.IsSlice(t.Elem()) {
			sb.WriteString("map<")
			sb.WriteString(typeToProto(t.Key()))
			sb.WriteString(",")
			sb.WriteString(typeToProto(t.Elem()))
			sb.WriteString(">")
			return sb.String()
		}
		panic(fmt.Errorf("proto map key can only be any integral or string type except floating point types and bytes"))
	default:
		panic(fmt.Errorf("unexpected type %s", tp))
	}
}
