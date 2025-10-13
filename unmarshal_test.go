package querystring_test

import (
	"net/url"
	"testing"
	"time"

	"github.com/eolymp/go-querystring"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestUnmarshal_ProtoMessage(t *testing.T) {
	tests := []struct {
		name  string
		query url.Values
		want  *TestData
	}{
		{
			name: "string value",
			query: url.Values{
				"string_value": {"hello world"},
			},
			want: &TestData{
				StringValue: "hello world",
			},
		},
		{
			name: "int32 value",
			query: url.Values{
				"int32_value": {"42"},
			},
			want: &TestData{
				Int32Value: 42,
			},
		},
		{
			name: "uint32 value",
			query: url.Values{
				"uint32_value": {"100"},
			},
			want: &TestData{
				Uint32Value: 100,
			},
		},
		{
			name: "bool value true",
			query: url.Values{
				"bool_value": {"true"},
			},
			want: &TestData{
				BoolValue: true,
			},
		},
		{
			name: "bool value false",
			query: url.Values{
				"bool_value": {"false"},
			},
			want: &TestData{
				BoolValue: false,
			},
		},
		{
			name: "float value",
			query: url.Values{
				"float_value": {"3.14"},
			},
			want: &TestData{
				FloatValue: 3.14,
			},
		},
		{
			name: "timestamp value",
			query: url.Values{
				"timestamp_value": {"2023-12-25T10:30:45Z"},
			},
			want: &TestData{
				TimestampValue: timestamppb.New(time.Date(2023, 12, 25, 10, 30, 45, 0, time.UTC)),
			},
		},
		{
			name: "timestamp value with different time",
			query: url.Values{
				"timestamp_value": {"2024-01-15T08:45:30Z"},
			},
			want: &TestData{
				TimestampValue: timestamppb.New(time.Date(2024, 1, 15, 8, 45, 30, 0, time.UTC)),
			},
		},
		{
			name: "multiple basic fields",
			query: url.Values{
				"string_value": {"test"},
				"int32_value":  {"123"},
				"bool_value":   {"true"},
			},
			want: &TestData{
				StringValue: "test",
				Int32Value:  123,
				BoolValue:   true,
			},
		},
		{
			name: "string values",
			query: url.Values{
				"string_values": {"one", "two", "three"},
			},
			want: &TestData{
				StringValues: []string{"one", "two", "three"},
			},
		},
		{
			name: "int32 values",
			query: url.Values{
				"int32_values": {"1", "2", "3"},
			},
			want: &TestData{
				Int32Values: []int32{1, 2, 3},
			},
		},
		{
			name: "uint32 values",
			query: url.Values{
				"uint32_values": {"10", "20", "30"},
			},
			want: &TestData{
				Uint32Values: []uint32{10, 20, 30},
			},
		},
		{
			name: "bool values",
			query: url.Values{
				"bool_values": {"true", "false", "true"},
			},
			want: &TestData{
				BoolValues: []bool{true, false, true},
			},
		},
		{
			name: "float values",
			query: url.Values{
				"float_values": {"1.1", "2.2", "3.3"},
			},
			want: &TestData{
				FloatValues: []float32{1.1, 2.2, 3.3},
			},
		},
		{
			name: "timestamp values",
			query: url.Values{
				"timestamp_values": {"2023-01-01T00:00:00Z", "2023-12-31T23:59:59Z", "2024-06-15T12:00:00Z"},
			},
			want: &TestData{
				TimestampValues: []*timestamppb.Timestamp{
					timestamppb.New(time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)),
					timestamppb.New(time.Date(2023, 12, 31, 23, 59, 59, 0, time.UTC)),
					timestamppb.New(time.Date(2024, 6, 15, 12, 0, 0, 0, time.UTC)),
				},
			},
		},
		{
			name: "mixed repeated fields",
			query: url.Values{
				"string_values": {"hello", "world"},
				"int32_values":  {"100", "200"},
			},
			want: &TestData{
				StringValues: []string{"hello", "world"},
				Int32Values:  []int32{100, 200},
			},
		},
		{
			name: "enum by name",
			query: url.Values{
				"enum_value": {"TWO"},
			},
			want: &TestData{
				EnumValue: TestData_TWO,
			},
		},
		{
			name: "enum by number",
			query: url.Values{
				"enum_value": {"1"},
			},
			want: &TestData{
				EnumValue: TestData_TWO,
			},
		},
		{
			name: "enum values by name",
			query: url.Values{
				"enum_values": {"ONE", "THREE"},
			},
			want: &TestData{
				EnumValues: []TestData_Enumeration{TestData_ONE, TestData_THREE},
			},
		},
		{
			name: "enum values by number",
			query: url.Values{
				"enum_values": {"0", "2"},
			},
			want: &TestData{
				EnumValues: []TestData_Enumeration{TestData_ONE, TestData_THREE},
			},
		},
		{
			name: "mixed enum formats",
			query: url.Values{
				"enum_values": {"ONE", "1", "THREE"},
			},
			want: &TestData{
				EnumValues: []TestData_Enumeration{TestData_ONE, TestData_TWO, TestData_THREE},
			},
		},
		{
			name: "flat struct",
			query: url.Values{
				"flat_struct": {`{"string_value":"nested test","enum_value":"THREE"}`},
			},
			want: &TestData{
				FlatStruct: &TestData_FlatStruct{
					StringValue: "nested test",
					EnumValue:   TestData_THREE,
				},
			},
		},
		{
			name: "complex struct",
			query: url.Values{
				"complex_struct": {`{"string_value":"complex test","enum_value":"TWO","flat_struct":{"string_value":"inner"}}`},
			},
			want: &TestData{
				ComplexStruct: &TestData_ComplexStruct{
					StringValue: "complex test",
					EnumValue:   TestData_TWO,
					FlatStruct: &TestData_FlatStruct{
						StringValue: "inner",
					},
				},
			},
		},
		{
			name: "flat structs array",
			query: url.Values{
				"flat_structs": {
					`{"string_value":"first","enum_value":"ONE"}`,
					`{"string_value":"second","enum_value":"TWO"}`,
				},
			},
			want: &TestData{
				FlatStructs: []*TestData_FlatStruct{
					{StringValue: "first", EnumValue: TestData_ONE},
					{StringValue: "second", EnumValue: TestData_TWO},
				},
			},
		},
		{
			name: "complex structs array",
			query: url.Values{
				"complex_structs": {
					`{"string_value":"first complex"}`,
					`{"string_value":"second complex","enum_value":"THREE"}`,
				},
			},
			want: &TestData{
				ComplexStructs: []*TestData_ComplexStruct{
					{StringValue: "first complex"},
					{StringValue: "second complex", EnumValue: TestData_THREE},
				},
			},
		},
		{
			name:  "empty query",
			query: url.Values{},
			want:  &TestData{},
		},
		{
			name: "empty string",
			query: url.Values{
				"string_value": {""},
			},
			want: &TestData{
				StringValue: "",
			},
		},
		{
			name: "zero values",
			query: url.Values{
				"int32_value":  {"0"},
				"uint32_value": {"0"},
				"bool_value":   {"false"},
				"float_value":  {"0"},
			},
			want: &TestData{
				Int32Value:  0,
				Uint32Value: 0,
				BoolValue:   false,
				FloatValue:  0,
			},
		},
		{
			name: "empty arrays",
			query: url.Values{
				"string_values": {},
				"int32_values":  {},
			},
			want: &TestData{},
		},
		{
			name: "flat struct with outer field",
			query: url.Values{
				"string_value": {"outer"},
				"flat_struct":  {`{"string_value":"inner","enum_value":"ONE"}`},
			},
			want: &TestData{
				StringValue: "outer",
				FlatStruct: &TestData_FlatStruct{
					StringValue: "inner",
					EnumValue:   TestData_ONE,
				},
			},
		},
		{
			name: "flat struct with arrays",
			query: url.Values{
				"flat_struct": {`{"string_values":["a","b"],"enum_values":["ONE","TWO"]}`},
			},
			want: &TestData{
				FlatStruct: &TestData_FlatStruct{
					StringValues: []string{"a", "b"},
					EnumValues:   []TestData_Enumeration{TestData_ONE, TestData_TWO},
				},
			},
		},
		{
			name: "complex struct with nested flat struct",
			query: url.Values{
				"complex_struct": {`{"string_value":"complex","flat_struct":{"string_value":"nested","enum_value":"THREE"}}`},
			},
			want: &TestData{
				ComplexStruct: &TestData_ComplexStruct{
					StringValue: "complex",
					FlatStruct: &TestData_FlatStruct{
						StringValue: "nested",
						EnumValue:   TestData_THREE,
					},
				},
			},
		},
		{
			name: "comprehensive test",
			query: url.Values{
				"string_value":     {"comprehensive"},
				"int32_value":      {"42"},
				"bool_value":       {"true"},
				"enum_value":       {"TWO"},
				"timestamp_value":  {"2023-07-15T14:30:00Z"},
				"string_values":    {"hello", "world"},
				"enum_values":      {"ONE", "THREE"},
				"timestamp_values": {"2023-01-01T00:00:00Z", "2023-12-31T23:59:59Z"},
				"flat_struct":      {`{"string_value":"nested","enum_value":"ONE"}`},
			},
			want: &TestData{
				StringValue:    "comprehensive",
				Int32Value:     42,
				BoolValue:      true,
				EnumValue:      TestData_TWO,
				TimestampValue: timestamppb.New(time.Date(2023, 7, 15, 14, 30, 0, 0, time.UTC)),
				StringValues:   []string{"hello", "world"},
				EnumValues:     []TestData_Enumeration{TestData_ONE, TestData_THREE},
				TimestampValues: []*timestamppb.Timestamp{
					timestamppb.New(time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)),
					timestamppb.New(time.Date(2023, 12, 31, 23, 59, 59, 0, time.UTC)),
				},
				FlatStruct: &TestData_FlatStruct{
					StringValue: "nested",
					EnumValue:   TestData_ONE,
				},
			},
		},
		{
			name: "all array types",
			query: url.Values{
				"string_values": {"one", "two"},
				"int32_values":  {"1", "2"},
				"uint32_values": {"10", "20"},
				"bool_values":   {"true", "false"},
				"float_values":  {"1.5", "2.5"},
				"enum_values":   {"ONE", "TWO"},
			},
			want: &TestData{
				StringValues: []string{"one", "two"},
				Int32Values:  []int32{1, 2},
				Uint32Values: []uint32{10, 20},
				BoolValues:   []bool{true, false},
				FloatValues:  []float32{1.5, 2.5},
				EnumValues:   []TestData_Enumeration{TestData_ONE, TestData_TWO},
			},
		},
		{
			name: "single values with arrays",
			query: url.Values{
				"string_value":  {"single"},
				"int32_value":   {"100"},
				"string_values": {"array1", "array2"},
				"int32_values":  {"200", "300"},
				"enum_value":    {"THREE"},
				"enum_values":   {"ONE", "TWO"},
			},
			want: &TestData{
				StringValue:  "single",
				Int32Value:   100,
				StringValues: []string{"array1", "array2"},
				Int32Values:  []int32{200, 300},
				EnumValue:    TestData_THREE,
				EnumValues:   []TestData_Enumeration{TestData_ONE, TestData_TWO},
			},
		},
		// todo: add support for maps
		//{
		//	name: "map value",
		//	query: url.Values{
		//		"map_value": {`{"key1":"value1","key2":"value2"}`},
		//	},
		//	want: &TestData{
		//		MapValue: map[string]string{
		//			"key1": "value1",
		//			"key2": "value2",
		//		},
		//	},
		//},
		//{
		//	name: "map value with single pair",
		//	query: url.Values{
		//		"map_value": {`{"single":"value"}`},
		//	},
		//	want: &TestData{
		//		MapValue: map[string]string{
		//			"single": "value",
		//		},
		//	},
		//},
		//{
		//	name: "empty map value",
		//	query: url.Values{
		//		"map_value": {`{}`},
		//	},
		//	want: &TestData{
		//		MapValue: map[string]string{},
		//	},
		//},
		{
			name: "complex struct with map",
			query: url.Values{
				"complex_struct": {`{"string_value":"complex","map_value":{"nested_key":"nested_value"}}`},
			},
			want: &TestData{
				ComplexStruct: &TestData_ComplexStruct{
					StringValue: "complex",
					MapValue: map[string]string{
						"nested_key": "nested_value",
					},
				},
			},
		},
		//{
		//	name: "mixed fields with maps",
		//	query: url.Values{
		//		"string_value":   {"test"},
		//		"int32_value":    {"42"},
		//		"map_value":      {`{"top_key":"top_value"}`},
		//		"complex_struct": {`{"string_value":"nested","enum_value":"TWO","map_value":{"nested_key":"nested_value"}}`},
		//	},
		//	want: &TestData{
		//		StringValue: "test",
		//		Int32Value:  42,
		//		MapValue: map[string]string{
		//			"top_key": "top_value",
		//		},
		//		ComplexStruct: &TestData_ComplexStruct{
		//			StringValue: "nested",
		//			EnumValue:   TestData_TWO,
		//			MapValue: map[string]string{
		//				"nested_key": "nested_value",
		//			},
		//		},
		//	},
		//},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := &TestData{}
			if err := querystring.Unmarshal(tt.query, got); err != nil {
				t.Errorf("Unmarshal() error = %v", err)
				return
			}

			if !proto.Equal(got, tt.want) {
				t.Errorf("Unmarshal() = %v, want %v", got, tt.want)
			}
		})
	}
}

// Backward compatibility tests using the 'q' parameter
func TestUnmarshal_ProtoMessage_BackwardCompatibility(t *testing.T) {
	tests := []struct {
		name  string
		query url.Values
		want  *TestData
	}{
		{
			name: "q parameter simple",
			query: url.Values{
				"q": {`{"string_value":"from q parameter","int32_value":123}`},
			},
			want: &TestData{
				StringValue: "from q parameter",
				Int32Value:  123,
			},
		},
		{
			name: "q parameter with enum",
			query: url.Values{
				"q": {`{"string_value":"test","enum_value":"TWO","bool_value":true}`},
			},
			want: &TestData{
				StringValue: "test",
				EnumValue:   TestData_TWO,
				BoolValue:   true,
			},
		},
		{
			name: "q parameter with arrays",
			query: url.Values{
				"q": {`{"string_values":["hello","world"],"int32_values":[1,2,3]}`},
			},
			want: &TestData{
				StringValues: []string{"hello", "world"},
				Int32Values:  []int32{1, 2, 3},
			},
		},
		{
			name: "q parameter with nested",
			query: url.Values{
				"q": {`{"string_value":"parent","flat_struct":{"string_value":"child","enum_value":"THREE"}}`},
			},
			want: &TestData{
				StringValue: "parent",
				FlatStruct: &TestData_FlatStruct{
					StringValue: "child",
					EnumValue:   TestData_THREE,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := &TestData{}
			if err := querystring.Unmarshal(tt.query, got); err != nil {
				t.Errorf("Unmarshal() error = %v", err)
				return
			}

			if !proto.Equal(got, tt.want) {
				t.Errorf("Unmarshal() = %v, want %v", got, tt.want)
			}
		})
	}
}
