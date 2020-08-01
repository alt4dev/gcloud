package logging

import (
	"github.com/google/go-cmp/cmp"
	"sort"
	"testing"
)

func sortClaims (claims []protoClaim) []protoClaim{
	sort.Slice(claims, func(i, j int) bool {
		return claims[i].Name < claims[j].Name
	})
	return claims
}

func TestParseString(t *testing.T) {
	alt4Claims := parseClaims(Claim{
		"name": "John Doe",
	})

	expectedClaims := []protoClaim{
		{
			Name:     "name",
			DataType: 4,
			Value:    "John Doe",
		},
	}
	if !cmp.Equal(alt4Claims, expectedClaims) {
		t.Error("Claims are not equal\n", cmp.Diff(alt4Claims, expectedClaims))
	}
}

func TestParsIntegers(t *testing.T) {
	var ui uint = 18
	var ui8 uint8 = 18
	var ui16 uint16 = 18
	var ui32 uint32 = 18
	var ui64 uint64 = 18
	var i int = 18
	var i8 int8 = 18
	var i16 int16 = 18
	var i32 int32 = 18
	var i64 int64 = 18
	alt4Claims := sortClaims(parseClaims(Claim{
		"uint": ui,
		"uint8": ui8,
		"uint16": ui16,
		"uint32": ui32,
		"uint64": ui64,
		"int": i,
		"int8": i8,
		"int16": i16,
		"int32": i32,
		"int64": i64,
	}))
	expectedClaims := sortClaims([]protoClaim{
		{ Name: "uint", DataType: 1, Value: "18"},
		{ Name: "uint8", DataType: 1, Value: "18"},
		{ Name: "uint16", DataType: 1, Value: "18"},
		{ Name: "uint32", DataType: 1, Value: "18"},
		{ Name: "uint64", DataType: 1, Value: "18"},
		{ Name: "int", DataType: 1, Value: "18"},
		{ Name: "int8", DataType: 1, Value: "18"},
		{ Name: "int16", DataType: 1, Value: "18"},
		{ Name: "int32", DataType: 1, Value: "18"},
		{ Name: "int64", DataType: 1, Value: "18"},
	})

	if !cmp.Equal(alt4Claims, expectedClaims) {
		t.Error("Claims are not equal Expected:\n", expectedClaims, "Got:\n", alt4Claims, "Diff:\n", cmp.Diff(alt4Claims, expectedClaims))
	}
}
