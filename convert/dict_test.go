// Copyright (c) Kopexa GmbH
// SPDX-License-Identifier: BUSL-1.1

package convert_test

import (
	"testing"

	"github.com/kopexa-grc/x/convert"
	"github.com/stretchr/testify/assert"
)

// Test JSONToDict function
func TestJSONToDict(t *testing.T) {
	// Test struct input
	type TestStruct struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}

	input := TestStruct{Name: "John", Age: 30}
	result, err := convert.JSONToDict(input)

	assert.NoError(t, err)
	assert.Equal(t, "John", result["name"])
	assert.Equal(t, 30, result["age"]) // Reflection preserves original int type
}

func TestJSONToDict_Map(t *testing.T) {
	input := map[string]interface{}{
		"key1": "value1",
		"key2": 42,
	}

	result, err := convert.JSONToDict(input)
	assert.NoError(t, err)
	assert.Equal(t, "value1", result["key1"])
	assert.Equal(t, 42, result["key2"]) // Reflection preserves original int type
}

func TestJSONToDict_Nil(t *testing.T) {
	result, err := convert.JSONToDict(nil)
	assert.NoError(t, err)
	assert.Nil(t, result)
}

// Test JSONToDictSlice function
func TestJSONToDictSlice(t *testing.T) {
	type TestStruct struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}

	input := []TestStruct{
		{Name: "John", Age: 30},
		{Name: "Jane", Age: 25},
	}

	result, err := convert.JSONToDictSlice(input)
	assert.NoError(t, err)
	assert.Len(t, result, 2)

	first := result[0].(map[string]interface{})
	assert.Equal(t, "John", first["name"])
	assert.Equal(t, 30, first["age"])

	second := result[1].(map[string]interface{})
	assert.Equal(t, "Jane", second["name"])
	assert.Equal(t, 25, second["age"])
}

func TestJSONToDictSlice_EmptySlice(t *testing.T) {
	input := []interface{}{}
	result, err := convert.JSONToDictSlice(input)
	assert.NoError(t, err)
	assert.Empty(t, result)
}

func TestJSONToDictSlice_NonSlice(t *testing.T) {
	input := "not a slice"
	result, err := convert.JSONToDictSlice(input)
	assert.Error(t, err) // Should error because string can't be unmarshaled as []any
	assert.Nil(t, result)
}

func TestDictToMapStr_Int(t *testing.T) {
	input := map[string]interface{}{
		"one": 1,
		"two": 2,
		"bad": "not-an-int",
	}
	expected := map[string]int{
		"one": 1,
		"two": 2,
	}

	output := convert.DictToTypedMap[int](input)
	assert.Equal(t, expected, output)
}

func TestDictToMapStr_String(t *testing.T) {
	input := map[string]interface{}{
		"a": "hello",
		"b": "world",
		"c": 123, // Should be ignored
	}
	expected := map[string]string{
		"a": "hello",
		"b": "world",
	}

	output := convert.DictToTypedMap[string](input)
	assert.Equal(t, expected, output)
}

func TestDictToMapStr_EmptyInput(t *testing.T) {
	output := convert.DictToTypedMap[int](nil)
	assert.Empty(t, output)
}

func TestDictToMapStr_WrongType(t *testing.T) {
	input := map[string]interface{}{
		"x": []int{1, 2, 3},
		"y": map[string]interface{}{"nested": "value"},
	}
	output := convert.DictToTypedMap[string](input)
	assert.Empty(t, output)
}

func TestDictToMapStr_NonMapInput(t *testing.T) {
	output := convert.DictToTypedMap[int]("not a map")
	assert.Empty(t, output)
}
