package openapi

import (
	"reflect"
	"testing"
)

func TestNewDocument(t *testing.T) {
	doc := NewDocument("3.0.0", NewInfo("Test API", "1.0.0"))

	if doc.OpenAPI != "3.0.0" {
		t.Errorf("OpenAPI = %v, want '3.0.0'", doc.OpenAPI)
	}
	if doc.Info.Title != "Test API" {
		t.Errorf("Info.Title = %v, want 'Test API'", doc.Info.Title)
	}
	if doc.Info.Version != "1.0.0" {
		t.Errorf("Info.Version = %v, want '1.0.0'", doc.Info.Version)
	}
	if doc.Paths == nil {
		t.Fatal("Paths should not be nil")
	}
}

func TestNewInfo(t *testing.T) {
	info := NewInfo("Test", "1.0.0")
	if info.Title != "Test" {
		t.Errorf("Title = %v, want 'Test'", info.Title)
	}
	if info.Version != "1.0.0" {
		t.Errorf("Version = %v, want '1.0.0'", info.Version)
	}
}

func TestDocument_WithInfo(t *testing.T) {
	doc := NewDocument("3.0.0", Info{})
	doc.WithInfo("New Title", "2.0.0")

	if doc.Info.Title != "New Title" {
		t.Errorf("Title = %v, want 'New Title'", doc.Info.Title)
	}
	if doc.Info.Version != "2.0.0" {
		t.Errorf("Version = %v, want '2.0.0'", doc.Info.Version)
	}
}

func TestDocument_AddServer(t *testing.T) {
	doc := NewDocument("3.0.0", NewInfo("Test", "1.0.0"))
	doc.AddServer(Server{URL: "https://api.example.com"})

	if len(doc.Servers) != 1 {
		t.Fatalf("Servers length = %d, want 1", len(doc.Servers))
	}
	if doc.Servers[0].URL != "https://api.example.com" {
		t.Errorf("Server URL = %v, want 'https://api.example.com'", doc.Servers[0].URL)
	}
}

func TestDocument_EnsurePaths(t *testing.T) {
	doc := NewDocument("3.0.0", NewInfo("Test", "1.0.0"))

	paths := doc.EnsurePaths()
	if paths == nil {
		t.Fatal("EnsurePaths() should not return nil")
	}
	if doc.Paths == nil {
		t.Fatal("Paths should not be nil after EnsurePaths()")
	}
}

func TestDocument_AddPathItem(t *testing.T) {
	doc := NewDocument("3.0.0", NewInfo("Test", "1.0.0"))

	pathItem := NewPathItem()
	pathItem.Summary = "Test Path"

	doc.AddPathItem("/test", pathItem)

	if doc.Paths["/test"].Summary != "Test Path" {
		t.Errorf("PathItem Summary = %v, want 'Test Path'", doc.Paths["/test"].Summary)
	}
}

func TestDocument_AddOperation(t *testing.T) {
	doc := NewDocument("3.0.0", NewInfo("Test", "1.0.0"))

	op := Operation{
		Summary: "Test Operation",
		Tags:    []string{"test"},
	}

	doc.AddOperation("/test", "get", op)

	if doc.Paths["/test"].Get == nil {
		t.Fatal("Get operation should not be nil")
	}
	if doc.Paths["/test"].Get.Summary != "Test Operation" {
		t.Errorf("Operation Summary = %v, want 'Test Operation'", doc.Paths["/test"].Get.Summary)
	}
}

func TestDocument_AddOperation_Methods(t *testing.T) {
	doc := NewDocument("3.0.0", NewInfo("Test", "1.0.0"))

	op := Operation{Summary: "Test"}

	methods := []string{"get", "post", "put", "delete", "patch", "head", "options", "trace"}
	for _, method := range methods {
		doc.AddOperation("/test", method, op)
	}

	pi := doc.Paths["/test"]
	if pi.Get == nil || pi.Post == nil || pi.Put == nil || pi.Delete == nil {
		t.Fatal("Operations should be set")
	}
}

func TestDocument_AddGet(t *testing.T) {
	doc := NewDocument("3.0.0", NewInfo("Test", "1.0.0"))
	doc.AddGet("/test", Operation{Summary: "GET test"})

	if doc.Paths["/test"].Get == nil {
		t.Fatal("Get operation should not be nil")
	}
}

func TestDocument_AddPost(t *testing.T) {
	doc := NewDocument("3.0.0", NewInfo("Test", "1.0.0"))
	doc.AddPost("/test", Operation{Summary: "POST test"})

	if doc.Paths["/test"].Post == nil {
		t.Fatal("Post operation should not be nil")
	}
}

func TestDocument_EnsureComponents(t *testing.T) {
	doc := NewDocument("3.0.0", NewInfo("Test", "1.0.0"))

	components := doc.EnsureComponents()
	if components == nil {
		t.Fatal("EnsureComponents() should not return nil")
	}
	if doc.Components == nil {
		t.Fatal("Components should not be nil")
	}
	if components.Schemas == nil {
		t.Fatal("Schemas should not be nil")
	}
}

func TestDocument_RegisterSchema(t *testing.T) {
	doc := NewDocument("3.0.0", NewInfo("Test", "1.0.0"))

	type TestStruct struct {
		Name string
		Age  int
	}

	schema := doc.RegisterSchema("TestStruct", TestStruct{})
	if schema == nil {
		t.Fatal("RegisterSchema() should not return nil")
	}
	if doc.Components.Schemas["TestStruct"] == nil {
		t.Fatal("Schema should be registered")
	}

	// 测试重复注册（应该返回已存在的 schema）
	schema2 := doc.RegisterSchema("TestStruct", TestStruct{})
	if schema != schema2 {
		t.Fatal("RegisterSchema() should return existing schema for duplicate name")
	}
}

func TestNewSchemaFrom(t *testing.T) {
	tests := []struct {
		name     string
		input    any
		expected string
	}{
		{"bool", true, "boolean"},
		{"int", 1, "integer"},
		{"int64", int64(1), "integer"},
		{"float32", float32(1.0), "number"},
		{"float64", 1.0, "number"},
		{"string", "test", "string"},
		{"slice", []int{1, 2}, "array"},
		{"map", map[string]int{"a": 1}, "object"},
		{"struct", struct{ Name string }{"test"}, "object"},
		{"nil", nil, "object"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			schema := NewSchemaFrom(tt.input)
			if schema.Type != tt.expected {
				t.Errorf("NewSchemaFrom(%v).Type = %v, want %v", tt.input, schema.Type, tt.expected)
			}
		})
	}
}

func TestNewSchemaFrom_Struct(t *testing.T) {
	type Person struct {
		Name  string `json:"name"`
		Age   int    `json:"age"`
		Email string `json:"email,omitempty"`
	}

	schema := NewSchemaFrom(Person{})
	if schema.Type != "object" {
		t.Errorf("Type = %v, want 'object'", schema.Type)
	}
	if schema.Properties == nil {
		t.Fatal("Properties should not be nil")
	}
	if schema.Properties["name"] == nil {
		t.Fatal("name property should exist")
	}
	if schema.Properties["age"] == nil {
		t.Fatal("age property should exist")
	}
}

func TestNewSchemaFrom_Pointer(t *testing.T) {
	type Person struct {
		Name string
	}

	schema := NewSchemaFrom((*Person)(nil))
	if !schema.Nullable {
		t.Error("Schema should be nullable for pointer type")
	}
}

func TestNewSchemaFrom_Array(t *testing.T) {
	schema := NewSchemaFrom([]string{})
	if schema.Type != "array" {
		t.Errorf("Type = %v, want 'array'", schema.Type)
	}
	if schema.Items == nil {
		t.Fatal("Items should not be nil")
	}
	if schema.Items.Type != "string" {
		t.Errorf("Items.Type = %v, want 'string'", schema.Items.Type)
	}
}

func TestNewParameter(t *testing.T) {
	param := NewParameter("id", InPath, 1, true)

	if param.Name != "id" {
		t.Errorf("Name = %v, want 'id'", param.Name)
	}
	if param.In != InPath {
		t.Errorf("In = %v, want '%v'", param.In, InPath)
	}
	if !param.Required {
		t.Error("Required should be true")
	}
	if param.Schema == nil {
		t.Fatal("Schema should not be nil")
	}
}

func TestNewMediaTypeFrom(t *testing.T) {
	mt := NewMediaTypeFrom(struct {
		Name string
	}{})

	if mt.Schema == nil {
		t.Fatal("Schema should not be nil")
	}
}

func TestNewMediaBody(t *testing.T) {
	body := NewMediaBody(struct {
		Name string
	}{}, "application/json")

	if body["application/json"].Schema == nil {
		t.Fatal("Schema should not be nil")
	}

	// 测试默认 content type
	body2 := NewMediaBody(struct{}{}, "")
	if _, exists := body2["application/json"]; !exists {
		t.Fatal("Should use default content type 'application/json'")
	}
}

func TestNewRequestBody(t *testing.T) {
	rb := NewRequestBody(struct {
		Name string
	}{}, true, "application/json")

	if !rb.Required {
		t.Error("Required should be true")
	}
	if rb.Content == nil {
		t.Fatal("Content should not be nil")
	}
}

func TestNewResponseFrom(t *testing.T) {
	resp := NewResponseFrom(struct {
		Message string
	}{}, "Success", "application/json")

	if resp.Description != "Success" {
		t.Errorf("Description = %v, want 'Success'", resp.Description)
	}
	if resp.Content == nil {
		t.Fatal("Content should not be nil")
	}

	// 测试默认 description
	resp2 := NewResponseFrom(struct{}{}, "", "")
	if resp2.Description != "OK" {
		t.Errorf("Description = %v, want 'OK'", resp2.Description)
	}
}

func TestNewResponses(t *testing.T) {
	responses := NewResponses(200, struct {
		Message string
	}{}, "Success", "application/json")

	if responses["200"].Description != "Success" {
		t.Errorf("Description = %v, want 'Success'", responses["200"].Description)
	}
}

func TestParseJSONTag(t *testing.T) {
	type TestStruct struct {
		Field1 string `json:"field1"`
		Field2 string `json:"field2,omitempty"`
		Field3 string `json:"-"`
		Field4 string `json:",omitempty"`
		Field5 string
		Field6 string `json:"field6" validate:"required"`
	}

	structType := reflect.TypeOf(TestStruct{})

	tests := []struct {
		fieldName string
		expected  string
		required  bool
		ignored   bool
	}{
		{"Field1", "field1", false, false},
		{"Field2", "field2", false, false},
		{"Field3", "Field3", false, true},
		{"Field4", "Field4", false, false},
		{"Field5", "Field5", false, false},
		{"Field6", "field6", true, false},
	}

	for _, tt := range tests {
		field, _ := structType.FieldByName(tt.fieldName)
		name, required, ignored := parseJSONTag(field)

		if name != tt.expected {
			t.Errorf("parseJSONTag(%s).name = %v, want %v", tt.fieldName, name, tt.expected)
		}
		if required != tt.required {
			t.Errorf("parseJSONTag(%s).required = %v, want %v", tt.fieldName, required, tt.required)
		}
		if ignored != tt.ignored {
			t.Errorf("parseJSONTag(%s).ignored = %v, want %v", tt.fieldName, ignored, tt.ignored)
		}
	}
}

func TestNewPaths(t *testing.T) {
	paths := NewPaths()
	if paths == nil {
		t.Fatal("NewPaths() should not return nil")
	}
}

func TestNewPathItem(t *testing.T) {
	item := NewPathItem()
	if item.Get != nil {
		t.Error("New PathItem should have nil operations")
	}
}
