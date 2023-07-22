package utils

import (
	"go/ast"
	"reflect"
	"testing"
)

// 测试parseCommentAnnotations函数
func TestParseCommentAnnotations(t *testing.T) {
	// 测试解析多个注释的情况
	cg1 := &ast.CommentGroup{
		List: []*ast.Comment{
			{
				Text: "// @http-method get",
			},
			{
				Text: "// @swagger-deprecated",
			},
			{
				Text: "// @custom-annotation value",
			},
		},
	}
	result1 := parseCommentAnnotations(cg1)
	if len(result1) != 3 {
		t.Errorf("Expected 3, but got %v", len(result1))
	}
	if result1["http-method"] != "get" {
		t.Errorf("Expected 'get', but got %v", result1["http-method"])
	}
	if result1["swagger-deprecated"] != "" {
		t.Errorf("Expected '', but got %v", result1["swagger-deprecated"])
	}
	if result1["custom-annotation"] != "value" {
		t.Errorf("Expected 'value', but got %v", result1["custom-annotation"])
	}
	// 测试解析空注释的情况
	cg2 := &ast.CommentGroup{
		List: []*ast.Comment{
			{
				Text: "//",
			},
		},
	}
	result2 := parseCommentAnnotations(cg2)
	if len(result2) != 0 {
		t.Errorf("Expected 0, but got %v", len(result2))
	}
	// 测试解析无注释的情况
	cg3 := &ast.CommentGroup{
		List: []*ast.Comment{},
	}
	result3 := parseCommentAnnotations(cg3)
	if len(result3) != 0 {
		t.Errorf("Expected 0, but got %v", len(result3))
	}
	// 测试解析注释中有多个空格的情况
	cg4 := &ast.CommentGroup{
		List: []*ast.Comment{
			{
				Text: "// @http-method     post",
			},
		},
	}
	result4 := parseCommentAnnotations(cg4)
	if len(result4) != 1 {
		t.Errorf("Expected 1, but got %v", len(result4))
	}
	if result4["http-method"] != "post" {
		t.Errorf("Expected 'post', but got %v", result4["http-method"])
	}
}

func TestUnmarshalInt(t *testing.T) {
	testCases := []struct {
		name      string
		value     string
		expected  int
		expectErr bool
	}{
		{
			name:      "Test positive integer",
			value:     "123",
			expected:  123,
			expectErr: false,
		},
		{
			name:      "Test negative integer",
			value:     "-456",
			expected:  -456,
			expectErr: false,
		},
		{
			name:      "Test zero",
			value:     "0",
			expected:  0,
			expectErr: false,
		},
		{
			name:      "Test invalid integer",
			value:     "abc",
			expected:  0, // Expecting default value for invalid input
			expectErr: true,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var result int
			err := unmarshalPrimitive(reflect.ValueOf(&result), tc.value)
			if err != nil && !tc.expectErr {
				t.Errorf("Unexpected error: %v", err)
			}

			if result != tc.expected {
				t.Errorf("Expected %d, but got %d", tc.expected, result)
			}
		})
	}
}

func TestUnmarshalFloat(t *testing.T) {
	testCases := []struct {
		name      string
		value     string
		expected  float64
		expectErr bool
	}{
		{
			name:      "Test positive float",
			value:     "3.14",
			expected:  3.14,
			expectErr: false,
		},
		{
			name:      "Test negative float",
			value:     "-2.5",
			expected:  -2.5,
			expectErr: false,
		},
		{
			name:      "Test zero",
			value:     "0.0",
			expected:  0.0,
			expectErr: false,
		},
		{
			name:      "Test invalid float",
			value:     "abc",
			expected:  0.0, // Expecting default value for invalid input
			expectErr: true,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var result float64
			err := unmarshalPrimitive(reflect.ValueOf(&result), tc.value)
			if err != nil && !tc.expectErr {
				t.Errorf("Unexpected error: %v", err)
			}

			if result != tc.expected {
				t.Errorf("Expected %f, but got %f", tc.expected, result)
			}
		})
	}
}

func TestUnmarshalString(t *testing.T) {
	testCases := []struct {
		name     string
		value    string
		expected string
	}{
		{
			name:     "Test non-empty string",
			value:    "hello",
			expected: "hello",
		},
		{
			name:     "Test empty string",
			value:    "",
			expected: "",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var result string
			err := unmarshalPrimitive(reflect.ValueOf(&result), tc.value)
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
			if result != tc.expected {
				t.Errorf("Expected '%s', but got '%s'", tc.expected, result)
			}
		})
	}
}

// 测试unmarshal函数
func TestUnmarshal(t *testing.T) {
	type TestStruct struct {
		Name string `jk:"name"`
		Age  int    `jk:"age"`
	}
	var ts TestStruct
	err := unmarshal(&ts, map[string]string{
		"name": "John",
		"age":  "25",
	})
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if ts.Name != "John" {
		t.Errorf("Expected 'John', but got %v", ts.Name)
	}
	if ts.Age != 25 {
		t.Errorf("Expected 25, but got %v", ts.Age)
	}

	type TestStruct2 struct {
		Name string `jk:"name"`
		Age  int
	}
	var ts2 TestStruct2
	err = unmarshal(&ts2, map[string]string{
		"name": "John",
		"Age":  "25",
	})
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if ts2.Name != "John" {
		t.Errorf("Expected 'John', but got %v", ts2.Name)
	}
	if ts2.Age != 25 {
		t.Errorf("Expected 25, but got %v", ts2.Age)
	}

	type TestStruct3 struct {
		Name string `jk:"name"`
		Age  int    `jk:"age"`
	}
	var ts3 TestStruct3
	err = unmarshal(&ts3, map[string]string{
		"name": "John",
		"Age":  "25",
	})
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if ts3.Name != "John" {
		t.Errorf("Expected 'John', but got %v", ts3.Name)
	}
	if ts3.Age != 0 {
		t.Errorf("Expected 0, but got %v", ts3.Age)
	}
}

// 测试UnmarshalAnnotations函数
func TestUnmarshalAnnotations(t *testing.T) {
	type TestStruct struct {
		HttpMethod string `jk:"http-method"`
		Deprecated bool   `jk:"swagger-deprecated"`
	}
	cg := &ast.CommentGroup{
		List: []*ast.Comment{
			{
				Text: "// @http-method get",
			},
			{
				Text: "// @swagger-deprecated",
			},
		},
	}
	var ts TestStruct
	err := UnmarshalAnnotations(cg, &ts)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if ts.HttpMethod != "get" {
		t.Errorf("Expected 'get', but got %v", ts.HttpMethod)
	}
	if ts.Deprecated != true {
		t.Errorf("Expected true, but got %v", ts.Deprecated)
	}
}
