// Copyright (c) Kopexa GmbH
// SPDX-License-Identifier: BUSL-1.1

package convert_test

import (
	"encoding/json"
	"testing"

	"github.com/kopexa-grc/x/convert"
)

// Old JSON-based implementation for comparison
func jsonToDictOld(v any) (map[string]any, error) {
	res := make(map[string]any)

	data, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal([]byte(data), &res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func jsonToDictSliceOld(v any) ([]any, error) {
	res := []any{}

	data, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal([]byte(data), &res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

type BenchmarkStruct struct {
	ID       int                    `json:"id"`
	Name     string                 `json:"name"`
	Email    string                 `json:"email"`
	Age      int                    `json:"age"`
	Active   bool                   `json:"active"`
	Score    float64                `json:"score"`
	Tags     []string               `json:"tags"`
	Metadata map[string]interface{} `json:"metadata"`
}

type NestedBenchmarkStruct struct {
	User    BenchmarkStruct   `json:"user"`
	Config  BenchmarkStruct   `json:"config"`
	Items   []BenchmarkStruct `json:"items"`
	Version string            `json:"version"`
}

func createBenchmarkData() BenchmarkStruct {
	return BenchmarkStruct{
		ID:     12345,
		Name:   "John Doe",
		Email:  "john.doe@example.com",
		Age:    30,
		Active: true,
		Score:  95.7,
		Tags:   []string{"developer", "golang", "backend", "api"},
		Metadata: map[string]interface{}{
			"department": "engineering",
			"level":      "senior",
			"projects":   3,
		},
	}
}

func createNestedBenchmarkData() NestedBenchmarkStruct {
	base := createBenchmarkData()
	return NestedBenchmarkStruct{
		User:   base,
		Config: base,
		Items: []BenchmarkStruct{
			base, base, base, base, base,
		},
		Version: "1.0.0",
	}
}

// Benchmark simple struct conversion
func BenchmarkJSONToDict_New_Simple(b *testing.B) {
	data := createBenchmarkData()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := convert.JSONToDict(data)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkJSONToDict_Old_Simple(b *testing.B) {
	data := createBenchmarkData()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := jsonToDictOld(data)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// Benchmark nested struct conversion
func BenchmarkJSONToDict_New_Nested(b *testing.B) {
	data := createNestedBenchmarkData()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := convert.JSONToDict(data)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkJSONToDict_Old_Nested(b *testing.B) {
	data := createNestedBenchmarkData()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := jsonToDictOld(data)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// Benchmark slice conversion
func BenchmarkJSONToDictSlice_New(b *testing.B) {
	data := []BenchmarkStruct{
		createBenchmarkData(),
		createBenchmarkData(),
		createBenchmarkData(),
		createBenchmarkData(),
		createBenchmarkData(),
	}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := convert.JSONToDictSlice(data)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkJSONToDictSlice_Old(b *testing.B) {
	data := []BenchmarkStruct{
		createBenchmarkData(),
		createBenchmarkData(),
		createBenchmarkData(),
		createBenchmarkData(),
		createBenchmarkData(),
	}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := jsonToDictSliceOld(data)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// Benchmark map input (should be fast for both)
func BenchmarkJSONToDict_New_Map(b *testing.B) {
	data := map[string]interface{}{
		"key1": "value1",
		"key2": 42,
		"key3": true,
		"key4": 3.14,
		"key5": []string{"a", "b", "c"},
	}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := convert.JSONToDict(data)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkJSONToDict_Old_Map(b *testing.B) {
	data := map[string]interface{}{
		"key1": "value1",
		"key2": 42,
		"key3": true,
		"key4": 3.14,
		"key5": []string{"a", "b", "c"},
	}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := jsonToDictOld(data)
		if err != nil {
			b.Fatal(err)
		}
	}
}
