package utils

import (
	"encoding/json"
	"fmt"
	"go/ast"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"github.com/hashicorp/go-multierror"
)

func parseCommentAnnotations(cg *ast.CommentGroup) map[string]string {
	var ret = make(map[string]string)

	if cg == nil {
		return ret
	}

	// e.g. @http-method get
	hasValueRegexp := regexp.MustCompile(`^@([[:word:]-]+)[[:space:]]+(.*)$`)
	// e.g. @swagger-deprecated
	noValueRegexp := regexp.MustCompile(`^@([[:word:]-]+)$`)

	for _, comment := range cg.List {
		line := strings.TrimSpace(strings.TrimPrefix(comment.Text, "// "))
		if hasValueRegexp.MatchString(line) {
			// parse annotation with argument
			match := hasValueRegexp.FindStringSubmatch(line)
			ret[match[1]] = match[2]
		} else if noValueRegexp.MatchString(line) {
			// parse annotation without argument
			match := noValueRegexp.FindStringSubmatch(line)
			ret[match[1]] = ""
		}
	}

	return ret
}

func unmarshalPrimitive(dest reflect.Value, value string) error {
	switch dest.Kind() {
	case reflect.Ptr:
		elem := dest.Elem()
		switch elem.Kind() {
		case reflect.Int:
			i, err := strconv.Atoi(value)
			if err != nil {
				return err
			}
			elem.SetInt(int64(i))
		case reflect.Int8:
			i, err := strconv.ParseInt(value, 10, 8)
			if err != nil {
				return err
			}
			elem.SetInt(i)
		case reflect.Int16:
			i, err := strconv.ParseInt(value, 10, 16)
			if err != nil {
				return err
			}
			elem.SetInt(i)
		case reflect.Int32:
			i, err := strconv.ParseInt(value, 10, 32)
			if err != nil {
				return err
			}
			elem.SetInt(i)
		case reflect.Int64:
			i, err := strconv.ParseInt(value, 10, 64)
			if err != nil {
				return err
			}
			elem.SetInt(i)
		case reflect.Uint:
			u, err := strconv.ParseUint(value, 10, 0)
			if err != nil {
				return err
			}
			elem.SetUint(u)
		case reflect.Uint8:
			u, err := strconv.ParseUint(value, 10, 8)
			if err != nil {
				return err
			}
			elem.SetUint(u)
		case reflect.Uint16:
			u, err := strconv.ParseUint(value, 10, 16)
			if err != nil {
				return err
			}
			elem.SetUint(u)
		case reflect.Uint32:
			u, err := strconv.ParseUint(value, 10, 32)
			if err != nil {
				return err
			}
			elem.SetUint(u)
		case reflect.Uint64:
			u, err := strconv.ParseUint(value, 10, 64)
			if err != nil {
				return err
			}
			elem.SetUint(u)
		case reflect.Float32:
			f, err := strconv.ParseFloat(value, 32)
			if err != nil {
				return err
			}
			elem.SetFloat(f)
		case reflect.Float64:
			f, err := strconv.ParseFloat(value, 64)
			if err != nil {
				return err
			}
			elem.SetFloat(f)
		case reflect.String:
			elem.SetString(value)
		case reflect.Bool:
			elem.SetBool(true)
		default:
			return fmt.Errorf("unsupported type: %v", elem.Type())
		}
	default:
		return fmt.Errorf("unsupported type: %v", dest.Type())
	}
	return nil
}

func unmarshalStruct(dest reflect.Value, value map[string]string) error {
	var ret error
	destType := dest.Type()

	for i := 0; i < dest.NumField(); i++ {
		field := destType.Field(i)
		tag := field.Tag.Get("jk")
		fieldName := tag
		if tag == "" {
			fieldName = field.Name
		}

		fieldValue := dest.Field(i)
		val, ok := value[fieldName]
		if !ok {
			continue
		}

		switch field.Type.Kind() {
		case reflect.Slice, reflect.Struct, reflect.Map:
			err := json.Unmarshal([]byte(val), dest.Field(i).Interface())
			if err != nil {
				ret = multierror.Append(ret, err)
			}
		case reflect.Uint, reflect.Int,
			reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
			reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
			reflect.String, reflect.Bool:
			err := unmarshalPrimitive(fieldValue.Addr(), val)
			if err != nil {
				ret = multierror.Append(ret, err)
			}
		default:
			ret = multierror.Append(ret, fmt.Errorf("unsupported type: %v (%s)", field.Type, field.Name))
		}
	}

	return ret
}

func unmarshal(dest interface{}, values map[string]string) error {
	destValue := reflect.ValueOf(dest)
	if destValue.Kind() != reflect.Ptr {
		return fmt.Errorf("dest must be a pointer")
	}

	if destValue.Elem().Type().Kind() != reflect.Struct {
		return fmt.Errorf("dest must be a pointer to struct")
	}

	destElem := destValue.Elem()
	return unmarshalStruct(destElem, values)
}

func UnmarshalAnnotations(cg *ast.CommentGroup, dest any) error {
	values := parseCommentAnnotations(cg)
	err := unmarshal(dest, values)
	if err != nil {
		return err
	}
	return nil
}
