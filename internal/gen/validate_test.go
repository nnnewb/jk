package gen

import (
	"go/token"
	"go/types"
	"testing"
)

// Unit test for isSerializable
func TestIsSerializable(t *testing.T) {
	// Test basic types
	if !IsSerializable(types.Typ[types.Int]) {
		t.Errorf("Expected Int to be serializable")
	}
	if !IsSerializable(types.Typ[types.Float32]) {
		t.Errorf("Expected Float32 to be serializable")
	}
	if !IsSerializable(types.Typ[types.String]) {
		t.Errorf("Expected String to be serializable")
	}
	if !IsSerializable(types.Typ[types.Bool]) {
		t.Errorf("Expected Bool to be serializable")
	}

	// Test struct types
	if !IsSerializable(types.NewStruct(
		[]*types.Var{
			types.NewVar(token.NoPos, nil, "Field1", types.Typ[types.Int]),
			types.NewVar(token.NoPos, nil, "Field2", types.Typ[types.String]),
		},
		nil,
	)) {
		t.Errorf("Expected TestStruct to be serializable")
	}

	// Test pointer types
	if !IsSerializable(types.NewPointer(types.Typ[types.Int])) {
		t.Errorf("Expected pointer to Int to be serializable")
	}
	if !IsSerializable(types.NewPointer(types.NewSlice(types.Typ[types.String]))) {
		t.Errorf("Expected pointer to slice of String to be serializable")
	}

	// Test slice types
	if !IsSerializable(types.NewSlice(types.Typ[types.Int])) {
		t.Errorf("Expected slice of Int to be serializable")
	}
	if !IsSerializable(types.NewSlice(types.NewStruct(
		[]*types.Var{
			types.NewVar(token.NoPos, nil, "Field1", types.Typ[types.Int]),
			types.NewVar(token.NoPos, nil, "Field2", types.Typ[types.String]),
		},
		[]string{},
	))) {
		t.Errorf("Expected slice of TestStruct to be serializable")
	}

	// Test map types
	if !IsSerializable(types.NewMap(types.Typ[types.String], types.Typ[types.Int])) {
		t.Errorf("Expected map with String key and Int value to be serializable")
	}
	if !IsSerializable(types.NewMap(types.Typ[types.String], types.NewSlice(types.Typ[types.String]))) {
		t.Errorf("Expected map with String key and slice of String value to be serializable")
	}
}

// Unit test for isSerializableStructureType
func TestIsSerializableStructureType(t *testing.T) {
	// Test case for a serializable struct
	if !isSerializableStructureType(types.NewStruct([]*types.Var{
		types.NewVar(token.NoPos, nil, "Field1", types.Typ[types.Int]),
		types.NewVar(token.NoPos, nil, "Field2", types.Typ[types.String]),
	}, []string{})) {
		t.Errorf("Expected SerializableStruct to be serializable")
	}

	// Test case for a struct with a non-serializable field
	if isSerializableStructureType(types.NewStruct([]*types.Var{
		types.NewVar(token.NoPos, nil, "Field1", types.Typ[types.Int]),
		types.NewVar(token.NoPos, nil, "Field2", types.NewChan(types.SendRecv, types.Typ[types.Int])),
	}, []string{})) {
		t.Errorf("Expected NonSerializableStruct to not be serializable")
	}
}

func TestIsSerializablePointerType(t *testing.T) {
	basicType := types.Universe.Lookup("int").Type()
	pointerType := types.NewPointer(basicType)
	if !isSerializablePointerType(pointerType) {
		t.Errorf("Expected pointer to basic type to be serializable")
	}

	sliceType := types.NewSlice(pointerType)
	pointerType = types.NewPointer(sliceType)
	if !isSerializablePointerType(pointerType) {
		t.Errorf("Expected pointer to slice of basic type to be serializable")
	}

	mapType := types.NewMap(basicType, pointerType)
	pointerType = types.NewPointer(mapType)
	if !isSerializablePointerType(pointerType) {
		t.Errorf("Expected pointer to map with basic key and pointer value to be serializable")
	}

	structType := types.NewStruct([]*types.Var{types.NewVar(token.NoPos, nil, "x", pointerType)}, nil)
	pointerType = types.NewPointer(structType)
	if !isSerializablePointerType(pointerType) {
		t.Errorf("Expected pointer to struct with pointer field to be serializable")
	}

	nonSerializableType := types.NewChan(types.SendRecv, basicType)
	pointerType = types.NewPointer(nonSerializableType)
	if isSerializablePointerType(pointerType) {
		t.Errorf("Expected pointer to non-serializable type to be non-serializable")
	}
}

func TestIsSerializableSliceType(t *testing.T) {
	// Test for a slice of basic types
	basicSlice := types.NewSlice(types.Typ[types.Int])
	if !isSerializableSliceType(basicSlice) {
		t.Errorf("Expected slice of basic type to be serializable")
	}

	// Test for a slice of struct types
	structType := types.NewStruct(nil, nil)
	structSlice := types.NewSlice(structType)
	if !isSerializableSliceType(structSlice) {
		t.Errorf("Expected slice of struct type to be serializable")
	}

	// Test for a slice of pointer types
	pointerType := types.NewPointer(structType)
	pointerSlice := types.NewSlice(pointerType)
	if !isSerializableSliceType(pointerSlice) {
		t.Errorf("Expected slice of pointer type to be serializable")
	}

	// Test for a slice of map types
	mapType := types.NewMap(types.Typ[types.String], structType)
	mapSlice := types.NewSlice(mapType)
	if !isSerializableSliceType(mapSlice) {
		t.Errorf("Expected slice of map type to be serializable")
	}

	// Test for a slice of non-serializable types
	nonSerializableType := types.NewChan(types.SendRecv, types.Typ[types.Int32])
	nonSerializableSlice := types.NewSlice(nonSerializableType)
	if isSerializableSliceType(nonSerializableSlice) {
		t.Errorf("Expected slice of non-serializable type to not be serializable")
	}
}

// Unit test for isSerializableMapType
func TestIsSerializableMapType(t *testing.T) {
	// Test case for a serializable map type
	serializableMap := types.NewMap(types.Typ[types.String], types.Typ[types.Int])
	if !isSerializableMapType(serializableMap) {
		t.Errorf("Expected serializable map type, but got non-serializable map type")
	}

	// Test case for a non-serializable map type with non-serializable key type
	nonSerializableMapKey := types.NewMap(types.Typ[types.Complex64], types.Typ[types.String])
	if isSerializableMapType(nonSerializableMapKey) {
		t.Errorf("Expected non-serializable map type with non-serializable key type, but got serializable map type")
	}

	// Test case for a non-serializable map type with non-serializable value type
	nonSerializableMapValue := types.NewMap(types.Typ[types.String], types.Typ[types.Complex64])
	if isSerializableMapType(nonSerializableMapValue) {
		t.Errorf("Expected non-serializable map type with non-serializable value type, but got serializable map type")
	}
}

// Unit test for isBasicSerializableType
func TestIsBasicSerializableType(t *testing.T) {
	// Test for basic serializable types
	if !isBasicSerializableType(types.Typ[types.Int]) {
		t.Errorf("Expected true for Int type")
	}
	if !isBasicSerializableType(types.Typ[types.Float32]) {
		t.Errorf("Expected true for Float32 type")
	}
	if !isBasicSerializableType(types.Typ[types.String]) {
		t.Errorf("Expected true for String type")
	}
	if !isBasicSerializableType(types.Typ[types.Bool]) {
		t.Errorf("Expected true for Bool type")
	}

	// Test for non-serializable types
	if isBasicSerializableType(types.Typ[types.Complex64]) {
		t.Errorf("Expected false for Complex64 type")
	}
	if isBasicSerializableType(types.Typ[types.Complex128]) {
		t.Errorf("Expected false for Complex128 type")
	}
	if isBasicSerializableType(types.Typ[types.Uintptr]) {
		t.Errorf("Expected false for Uintptr type")
	}
}
