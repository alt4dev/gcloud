package log

import (
	"github.com/alt4dev/protobuff/proto"
	"github.com/google/go-cmp/cmp"
	"sort"
	"testing"
)

func sortClaims(claims []*proto.Claim) []*proto.Claim {
	sort.Slice(claims, func(i, j int) bool {
		return claims[i].Name < claims[j].Name
	})
	return claims
}

func TestParseString(t *testing.T) {
	alt4Claims := parseClaims(Claim{
		"name": "John Doe",
	})

	expectedClaims := []proto.Claim{
		{Name: "name", DataType: 4, Value: "John Doe"},
	}
	if !cmp.Equal(alt4Claims, expectedClaims) {
		t.Error("Claims are not equal\n", cmp.Diff(alt4Claims, expectedClaims))
	}
}

func TestParsIntegers(t *testing.T) {
	var ui uint = 1000000000
	var ui8 uint8 = 100
	var ui16 uint16 = 1000
	var ui32 uint32 = 10000
	var ui64 uint64 = 100000
	var i int = -1000000000
	var i8 int8 = -100
	var i16 int16 = -1000
	var i32 int32 = -10000
	var i64 int64 = -100000
	alt4Claims := sortClaims(parseClaims(Claim{
		"uint":   ui,
		"uint8":  ui8,
		"uint16": ui16,
		"uint32": ui32,
		"uint64": ui64,
		"int":    i,
		"int8":   i8,
		"int16":  i16,
		"int32":  i32,
		"int64":  i64,
	}))
	expectedClaims := sortClaims([]*proto.Claim{
		{Name: "uint", DataType: 1, Value: "1000000000"},
		{Name: "uint8", DataType: 1, Value: "100"},
		{Name: "uint16", DataType: 1, Value: "1000"},
		{Name: "uint32", DataType: 1, Value: "10000"},
		{Name: "uint64", DataType: 1, Value: "100000"},
		{Name: "int", DataType: 1, Value: "-1000000000"},
		{Name: "int8", DataType: 1, Value: "-100"},
		{Name: "int16", DataType: 1, Value: "-1000"},
		{Name: "int32", DataType: 1, Value: "-10000"},
		{Name: "int64", DataType: 1, Value: "-100000"},
	})

	if !cmp.Equal(alt4Claims, expectedClaims) {
		t.Error("Claims are not equal Expected:\n", expectedClaims, "Got:\n", alt4Claims, "Diff:\n", cmp.Diff(alt4Claims, expectedClaims))
	}
}

func TestParseFloat(t *testing.T) {
	var f32 float32 = 35.6
	var f64 float64 = 77.8
	alt4Claims := sortClaims(parseClaims(Claim{
		"float32": f32,
		"float64": f64,
	}))

	expectedClaims := sortClaims([]*proto.Claim{
		{Name: "float32", DataType: 2, Value: "35.6"},
		{Name: "float64", DataType: 2, Value: "77.8"},
	})
	if !cmp.Equal(alt4Claims, expectedClaims) {
		t.Error("Claims are not equal\n", cmp.Diff(alt4Claims, expectedClaims))
	}
}
