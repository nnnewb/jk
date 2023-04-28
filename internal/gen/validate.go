package gen

import "go/types"

// IsSerializable check is given type serializable
func IsSerializable(t types.Type) bool {
	switch t.(type) {
	case *types.Basic:
		return isBasicSerializableType(t)
	case *types.Pointer:
		return isSerializablePointerType(t)
	case *types.Slice:
		return isSerializableSliceType(t)
	case *types.Map:
		return isSerializableMapType(t)
	case *types.Struct:
		return isSerializableStructureType(t)
	case *types.Named:
		return IsSerializable(t.Underlying())
	default:
		return false
	}
}

func isSerializableStructureType(t types.Type) bool {
	if s, ok := t.(*types.Struct); ok {
		for i := 0; i < s.NumFields(); i++ {
			field := s.Field(i)
			if !IsSerializable(field.Type()) {
				return false
			}
		}
		return true
	}
	return false
}

func isSerializablePointerType(t types.Type) bool {
	if p, ok := t.(*types.Pointer); ok {
		switch p.Elem().(type) {
		case *types.Pointer:
			return false
		default:
			return IsSerializable(p.Elem())
		}
	}
	return false
}

func isSerializableSliceType(t types.Type) bool {
	if s, ok := t.(*types.Slice); ok {
		return IsSerializable(s.Elem())
	}
	return false
}

func isSerializableMapType(t types.Type) bool {
	if m, ok := t.(*types.Map); ok {
		if !isBasicSerializableType(m.Key()) {
			return false
		}

		return IsSerializable(m.Elem())
	}
	return false
}

func isBasicSerializableType(t types.Type) bool {
	if b, ok := t.(*types.Basic); ok {
		switch b.Kind() {
		case types.Int, types.Uint,
			types.Int8, types.Int16, types.Int32, types.Int64,
			types.Uint8, types.Uint16, types.Uint32, types.Uint64,
			types.Float32, types.Float64,
			types.String, types.Bool:
			return true
		default:
			return false
		}
	}
	return false
}
