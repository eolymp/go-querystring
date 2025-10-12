package querystring

import (
	"net/url"
	"reflect"
	"testing"
)

func TestUnmarshal_BasicTypes(t *testing.T) {
	type TestStruct struct {
		StringField string
		IntField    int
		Int32Field  int32
		UintField   uint
		Uint32Field uint32
		FloatField  float64
		BoolField   bool
	}

	query := url.Values{
		"StringField": {"hello world"},
		"IntField":    {"-123"},
		"Int32Field":  {"456"},
		"UintField":   {"789"},
		"Uint32Field": {"999"},
		"FloatField":  {"3.14"},
		"BoolField":   {"true"},
	}

	var got TestStruct
	if err := Unmarshal(query, &got); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	want := TestStruct{
		StringField: "hello world",
		IntField:    -123,
		Int32Field:  456,
		UintField:   789,
		Uint32Field: 999,
		FloatField:  3.14,
		BoolField:   true,
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("Expected %+v, got %+v", want, got)
	}
}

func TestUnmarshal_MissingFields(t *testing.T) {
	type TestStruct struct {
		Present string
		Missing string
	}

	query := url.Values{
		"Present": {"value"},
	}

	var got TestStruct
	if err := Unmarshal(query, &got); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if got.Present != "value" {
		t.Errorf("Expected Present='value', got %q", got.Present)
	}
	if got.Missing != "" {
		t.Errorf("Expected Missing to be zero value, got %q", got.Missing)
	}
}

func TestUnmarshal_Slices(t *testing.T) {
	type TestStruct struct {
		StringSlice []string
		IntSlice    []int
		UintSlice   []uint32
		FloatSlice  []float64
		BoolSlice   []bool
	}

	query := url.Values{
		"StringSlice": {"one", "two", "three"},
		"IntSlice":    {"1", "-2", "3"},
		"UintSlice":   {"10", "20", "30"},
		"FloatSlice":  {"1.1", "2.2", "3.3"},
		"BoolSlice":   {"true", "false", "true"},
	}

	var got TestStruct
	if err := Unmarshal(query, &got); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	want := TestStruct{
		StringSlice: []string{"one", "two", "three"},
		IntSlice:    []int{1, -2, 3},
		UintSlice:   []uint32{10, 20, 30},
		FloatSlice:  []float64{1.1, 2.2, 3.3},
		BoolSlice:   []bool{true, false, true},
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("Expected %+v, got %+v", want, got)
	}
}

func TestUnmarshal_EmptySlices(t *testing.T) {
	type TestStruct struct {
		StringSlice []string
		IntSlice    []int
	}

	query := url.Values{
		"StringSlice": {},
		"IntSlice":    {},
	}

	var got TestStruct
	if err := Unmarshal(query, &got); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if len(got.StringSlice) != 0 {
		t.Errorf("Expected empty slice, got %v", got.StringSlice)
	}

	if len(got.IntSlice) != 0 {
		t.Errorf("Expected empty slice, got %v", got.IntSlice)
	}
}

func TestUnmarshal_ComplexTypes(t *testing.T) {
	type NestedStruct struct {
		Name  string `json:"name"`
		Value int    `json:"value"`
	}

	type TestStruct struct {
		NestedField NestedStruct
		MapField    map[string]string
		SliceJSON   []NestedStruct
	}

	query := url.Values{
		"NestedField": {`{"name":"test","value":42}`},
		"MapField":    {`{"key1":"value1","key2":"value2"}`},
		"SliceJSON":   {`{"name":"first","value":1}`, `{"name":"second","value":2}`},
	}

	var got TestStruct
	if err := Unmarshal(query, &got); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	want := TestStruct{
		NestedField: NestedStruct{Name: "test", Value: 42},
		MapField:    map[string]string{"key1": "value1", "key2": "value2"},
		SliceJSON:   []NestedStruct{{Name: "first", Value: 1}, {Name: "second", Value: 2}},
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("Expected %+v, got %+v", want, got)
	}
}

func TestUnmarshal_PointerTypes(t *testing.T) {
	type NestedStruct struct {
		Value string `json:"value"`
	}

	type TestStruct struct {
		PtrField *NestedStruct
	}

	query := url.Values{
		"PtrField": {`{"value":"test"}`},
	}

	var got TestStruct
	if err := Unmarshal(query, &got); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if got.PtrField == nil {
		t.Fatal("Expected PtrField to be non-nil")
	}

	if got.PtrField.Value != "test" {
		t.Errorf("Expected PtrField.Value='test', got %q", got.PtrField.Value)
	}
}

func TestUnmarshal_StructTags(t *testing.T) {
	type TestStruct struct {
		DefaultName  string `json:"default_name"`
		CustomName   string `json:"custom"`
		IgnoredField string `json:"-"`
		WithOptions  string `json:"with_opts,omitempty"`
	}

	query := url.Values{
		"default_name": {"default_value"},
		"custom":       {"custom_value"},
		"IgnoredField": {"ignored_value"},
		"with_opts":    {"options_value"},
	}

	var got TestStruct
	if err := Unmarshal(query, &got); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if got.DefaultName != "default_value" {
		t.Errorf("Expected DefaultName='default_value', got %q", got.DefaultName)
	}

	if got.CustomName != "custom_value" {
		t.Errorf("Expected CustomName='custom_value', got %q", got.CustomName)
	}

	if got.IgnoredField != "" {
		t.Errorf("Expected IgnoredField to be empty (ignored), got %q", got.IgnoredField)
	}

	if got.WithOptions != "options_value" {
		t.Errorf("Expected WithOptions='options_value', got %q", got.WithOptions)
	}
}

func TestUnmarshal_ErrorCases(t *testing.T) {
	tests := []struct {
		name  string
		data  any
		error bool
	}{
		{
			name:  "not a pointer",
			data:  struct{}{},
			error: true,
		},
		{
			name:  "nil pointer",
			data:  (*struct{})(nil),
			error: true,
		},
		{
			name:  "pointer to non-struct",
			data:  new(string),
			error: true,
		},
	}

	query := url.Values{"test": {"value"}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Unmarshal(query, tt.data)
			if tt.error && err == nil {
				t.Fatal("Expected error but got none")
			}

			if !tt.error && err != nil {
				t.Fatalf("Expected no error but got: %v", err)
			}
		})
	}
}

func TestUnmarshal_ParseErrors(t *testing.T) {
	type TestStruct struct {
		IntField   int
		UintField  uint
		FloatField float64
		BoolField  bool
	}

	tests := []struct {
		name  string
		query url.Values
	}{
		{
			name:  "invalid int",
			query: url.Values{"IntField": {"not-a-number"}},
		},
		{
			name:  "invalid uint",
			query: url.Values{"UintField": {"-123"}},
		},
		{
			name:  "invalid float",
			query: url.Values{"FloatField": {"not-a-float"}},
		},
		{
			name:  "invalid bool",
			query: url.Values{"BoolField": {"maybe"}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var result TestStruct
			if err := Unmarshal(tt.query, &result); err == nil {
				t.Error("Expected parse error but got none")
			}
		})
	}
}

func TestUnmarshal_SliceParseErrors(t *testing.T) {
	type TestStruct struct {
		IntSlice []int
	}

	query := url.Values{
		"IntSlice": {"1", "not-a-number", "3"},
	}

	var got TestStruct
	if err := Unmarshal(query, &got); err == nil {
		t.Error("Expected parse error for slice element but got none")
	}
}

func TestUnmarshal_JSONParseErrors(t *testing.T) {
	type TestStruct struct {
		MapField map[string]string
	}

	query := url.Values{
		"MapField": {"invalid-json"},
	}

	var got TestStruct
	if err := Unmarshal(query, &got); err == nil {
		t.Error("Expected JSON parse error but got none")
	}
}

func TestUnmarshaler_CustomStructTag(t *testing.T) {
	type TestStruct struct {
		Field1 string `query:"field_1"`
		Field2 string `query:"field_2"`
		Field3 string `json:"field_3"`
	}

	query := url.Values{
		"field_1": {"value1"},
		"field_2": {"value2"},
		"field_3": {"value3"},
	}

	var got TestStruct
	if err := (&UnmarshalOptions{StructTag: "query"}).Unmarshal(query, &got); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if got.Field1 != "value1" {
		t.Errorf("Expected Field1='value1', got %q", got.Field1)
	}

	if got.Field2 != "value2" {
		t.Errorf("Expected Field2='value2', got %q", got.Field2)
	}

	if got.Field3 != "" {
		t.Errorf("Expected Field3 to be empty (no query StructTag), got %q", got.Field3)
	}
}
