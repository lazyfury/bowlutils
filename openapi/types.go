package openapi

import (
	"reflect"
	"strconv"
	"strings"
)

type Document struct {
	OpenAPI      string                `json:"openapi"`
	Info         Info                  `json:"info"`
	Servers      []Server              `json:"servers,omitempty"`
	Paths        Paths                 `json:"paths"`
	Components   *Components           `json:"components,omitempty"`
	Security     []SecurityRequirement `json:"security,omitempty"`
	Tags         []Tag                 `json:"tags,omitempty"`
	ExternalDocs *ExternalDoc          `json:"externalDocs,omitempty"`
}

type Info struct {
	Title          string   `json:"title"`
	Summary        string   `json:"summary,omitempty"`
	Description    string   `json:"description,omitempty"`
	TermsOfService string   `json:"termsOfService,omitempty"`
	Contact        *Contact `json:"contact,omitempty"`
	License        *License `json:"license,omitempty"`
	Version        string   `json:"version"`
}

type Contact struct {
	Name  string `json:"name,omitempty"`
	URL   string `json:"url,omitempty"`
	Email string `json:"email,omitempty"`
}

type License struct {
	Name       string `json:"name"`
	Identifier string `json:"identifier,omitempty"`
	URL        string `json:"url,omitempty"`
}

type Server struct {
	URL         string                    `json:"url"`
	Description string                    `json:"description,omitempty"`
	Variables   map[string]ServerVariable `json:"variables,omitempty"`
}

type ServerVariable struct {
	Enum        []string `json:"enum,omitempty"`
	Default     string   `json:"default"`
	Description string   `json:"description,omitempty"`
}

type Paths map[string]PathItem

type PathItem struct {
	Ref         string      `json:"$ref,omitempty"`
	Summary     string      `json:"summary,omitempty"`
	Description string      `json:"description,omitempty"`
	Get         *Operation  `json:"get,omitempty"`
	Put         *Operation  `json:"put,omitempty"`
	Post        *Operation  `json:"post,omitempty"`
	Delete      *Operation  `json:"delete,omitempty"`
	Options     *Operation  `json:"options,omitempty"`
	Head        *Operation  `json:"head,omitempty"`
	Patch       *Operation  `json:"patch,omitempty"`
	Trace       *Operation  `json:"trace,omitempty"`
	Servers     []Server    `json:"servers,omitempty"`
	Parameters  []Parameter `json:"parameters,omitempty"`
}

type Operation struct {
	Tags         []string              `json:"tags,omitempty"`
	Summary      string                `json:"summary,omitempty"`
	Description  string                `json:"description,omitempty"`
	ExternalDocs *ExternalDoc          `json:"externalDocs,omitempty"`
	OperationID  string                `json:"operationId,omitempty"`
	Parameters   []Parameter           `json:"parameters,omitempty"`
	RequestBody  *RequestBody          `json:"requestBody,omitempty"`
	Responses    Responses             `json:"responses"`
	Callbacks    map[string]Callback   `json:"callbacks,omitempty"`
	Deprecated   bool                  `json:"deprecated,omitempty"`
	Security     []SecurityRequirement `json:"security,omitempty"`
	Servers      []Server              `json:"servers,omitempty"`
}

var (
	InQuery  = "query"
	InPath   = "path"
	InHeader = "header"
	InCookie = "cookie"
)

type Parameter struct {
	Ref             string               `json:"$ref,omitempty"`
	Name            string               `json:"name,omitempty"`
	In              string               `json:"in,omitempty"`
	Description     string               `json:"description,omitempty"`
	Required        bool                 `json:"required,omitempty"`
	Deprecated      bool                 `json:"deprecated,omitempty"`
	AllowEmptyValue bool                 `json:"allowEmptyValue,omitempty"`
	Style           string               `json:"style,omitempty"`
	Explode         *bool                `json:"explode,omitempty"`
	AllowReserved   bool                 `json:"allowReserved,omitempty"`
	Schema          *Schema              `json:"schema,omitempty"`
	Example         any                  `json:"example,omitempty"`
	Examples        map[string]Example   `json:"examples,omitempty"`
	Content         map[string]MediaType `json:"content,omitempty"`
}

type RequestBody struct {
	Ref         string               `json:"$ref,omitempty"`
	Description string               `json:"description,omitempty"`
	Content     map[string]MediaType `json:"content"`
	Required    bool                 `json:"required,omitempty"`
}

type Responses map[string]Response

type Response struct {
	Ref         string               `json:"$ref,omitempty"`
	Description string               `json:"description"`
	Headers     map[string]Header    `json:"headers,omitempty"`
	Content     map[string]MediaType `json:"content,omitempty"`
	Links       map[string]Link      `json:"links,omitempty"`
}

type MediaType struct {
	Schema   *Schema             `json:"schema,omitempty"`
	Example  any                 `json:"example,omitempty"`
	Examples map[string]Example  `json:"examples,omitempty"`
	Encoding map[string]Encoding `json:"encoding,omitempty"`
}

type Encoding struct {
	ContentType   string            `json:"contentType,omitempty"`
	Headers       map[string]Header `json:"headers,omitempty"`
	Style         string            `json:"style,omitempty"`
	Explode       *bool             `json:"explode,omitempty"`
	AllowReserved bool              `json:"allowReserved,omitempty"`
}

type Example struct {
	Ref           string `json:"$ref,omitempty"`
	Summary       string `json:"summary,omitempty"`
	Description   string `json:"description,omitempty"`
	Value         any    `json:"value,omitempty"`
	ExternalValue string `json:"externalValue,omitempty"`
}

type Header struct {
	Ref         string               `json:"$ref,omitempty"`
	Description string               `json:"description,omitempty"`
	Required    bool                 `json:"required,omitempty"`
	Deprecated  bool                 `json:"deprecated,omitempty"`
	Style       string               `json:"style,omitempty"`
	Explode     *bool                `json:"explode,omitempty"`
	Schema      *Schema              `json:"schema,omitempty"`
	Example     any                  `json:"example,omitempty"`
	Examples    map[string]Example   `json:"examples,omitempty"`
	Content     map[string]MediaType `json:"content,omitempty"`
}

type Components struct {
	Schemas         map[string]*Schema        `json:"schemas,omitempty"`
	Responses       map[string]Response       `json:"responses,omitempty"`
	Parameters      map[string]Parameter      `json:"parameters,omitempty"`
	Examples        map[string]Example        `json:"examples,omitempty"`
	RequestBodies   map[string]RequestBody    `json:"requestBodies,omitempty"`
	Headers         map[string]Header         `json:"headers,omitempty"`
	SecuritySchemes map[string]SecurityScheme `json:"securitySchemes,omitempty"`
	Links           map[string]Link           `json:"links,omitempty"`
	Callbacks       map[string]Callback       `json:"callbacks,omitempty"`
	PathItems       map[string]PathItem       `json:"pathItems,omitempty"`
}

type SecurityScheme struct {
	Type             string      `json:"type"`
	Description      string      `json:"description,omitempty"`
	Name             string      `json:"name,omitempty"`
	In               string      `json:"in,omitempty"`
	Scheme           string      `json:"scheme,omitempty"`
	BearerFormat     string      `json:"bearerFormat,omitempty"`
	Flows            *OAuthFlows `json:"flows,omitempty"`
	OpenIdConnectUrl string      `json:"openIdConnectUrl,omitempty"`
}

type OAuthFlows struct {
	Implicit          *OAuthFlow `json:"implicit,omitempty"`
	Password          *OAuthFlow `json:"password,omitempty"`
	ClientCredentials *OAuthFlow `json:"clientCredentials,omitempty"`
	AuthorizationCode *OAuthFlow `json:"authorizationCode,omitempty"`
}

type OAuthFlow struct {
	AuthorizationUrl string            `json:"authorizationUrl,omitempty"`
	TokenUrl         string            `json:"tokenUrl,omitempty"`
	RefreshUrl       string            `json:"refreshUrl,omitempty"`
	Scopes           map[string]string `json:"scopes,omitempty"`
}

type SecurityRequirement map[string][]string

type Tag struct {
	Name         string       `json:"name"`
	Description  string       `json:"description,omitempty"`
	ExternalDocs *ExternalDoc `json:"externalDocs,omitempty"`
}

type ExternalDoc struct {
	Description string `json:"description,omitempty"`
	URL         string `json:"url"`
}

type XML struct {
	Name      string `json:"name,omitempty"`
	Namespace string `json:"namespace,omitempty"`
	Prefix    string `json:"prefix,omitempty"`
	Attribute bool   `json:"attribute,omitempty"`
	Wrapped   bool   `json:"wrapped,omitempty"`
}

type Callback map[string]PathItem

type Link struct {
	OperationRef string         `json:"operationRef,omitempty"`
	OperationId  string         `json:"operationId,omitempty"`
	Parameters   map[string]any `json:"parameters,omitempty"`
	RequestBody  any            `json:"requestBody,omitempty"`
	Description  string         `json:"description,omitempty"`
	Server       *Server        `json:"server,omitempty"`
}

type Schema struct {
	Ref                  string             `json:"$ref,omitempty"`
	Title                string             `json:"title,omitempty"`
	MultipleOf           *float64           `json:"multipleOf,omitempty"`
	Maximum              *float64           `json:"maximum,omitempty"`
	ExclusiveMaximum     *bool              `json:"exclusiveMaximum,omitempty"`
	Minimum              *float64           `json:"minimum,omitempty"`
	ExclusiveMinimum     *bool              `json:"exclusiveMinimum,omitempty"`
	MaxLength            *int               `json:"maxLength,omitempty"`
	MinLength            *int               `json:"minLength,omitempty"`
	Pattern              string             `json:"pattern,omitempty"`
	MaxItems             *int               `json:"maxItems,omitempty"`
	MinItems             *int               `json:"minItems,omitempty"`
	UniqueItems          *bool              `json:"uniqueItems,omitempty"`
	MaxProperties        *int               `json:"maxProperties,omitempty"`
	MinProperties        *int               `json:"minProperties,omitempty"`
	Required             []string           `json:"required,omitempty"`
	Type                 string             `json:"type,omitempty"`
	Not                  *Schema            `json:"not,omitempty"`
	AllOf                []*Schema          `json:"allOf,omitempty"`
	OneOf                []*Schema          `json:"oneOf,omitempty"`
	AnyOf                []*Schema          `json:"anyOf,omitempty"`
	Items                *Schema            `json:"items,omitempty"`
	Properties           map[string]*Schema `json:"properties,omitempty"`
	AdditionalProperties any                `json:"additionalProperties,omitempty"`
	Description          string             `json:"description,omitempty"`
	Format               string             `json:"format,omitempty"`
	Default              any                `json:"default,omitempty"`
	Nullable             bool               `json:"nullable,omitempty"`
	Discriminator        *Discriminator     `json:"discriminator,omitempty"`
	ReadOnly             bool               `json:"readOnly,omitempty"`
	WriteOnly            bool               `json:"writeOnly,omitempty"`
	XML                  *XML               `json:"xml,omitempty"`
	ExternalDocs         *ExternalDoc       `json:"externalDocs,omitempty"`
	Example              any                `json:"example,omitempty"`
	Deprecated           bool               `json:"deprecated,omitempty"`
}

type Discriminator struct {
	PropertyName string            `json:"propertyName"`
	Mapping      map[string]string `json:"mapping,omitempty"`
}

func NewDocument(version string, info Info) *Document {
	return &Document{OpenAPI: version, Info: info, Paths: make(Paths)}
}

func NewInfo(title, version string) Info {
	return Info{Title: title, Version: version}
}

func NewPaths() Paths { return make(Paths) }

func NewPathItem() PathItem { return PathItem{} }

func NewSchemaFrom(v any) *Schema {
	if v == nil {
		return &Schema{Type: "object"}
	}
	t := reflect.TypeOf(v)
	nullable := false
	for t.Kind() == reflect.Ptr {
		nullable = true
		t = t.Elem()
	}
	s := &Schema{}
	if nullable {
		s.Nullable = true
	}
	switch t.Kind() {
	case reflect.Bool:
		s.Type = "boolean"
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32:
		s.Type = "integer"
		s.Format = "int32"
	case reflect.Int64:
		s.Type = "integer"
		s.Format = "int64"
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32:
		s.Type = "integer"
		s.Format = "int32"
	case reflect.Uint64:
		s.Type = "integer"
		s.Format = "int64"
	case reflect.Float32:
		s.Type = "number"
		s.Format = "float"
	case reflect.Float64:
		s.Type = "number"
		s.Format = "double"
	case reflect.String:
		s.Type = "string"
	case reflect.Slice, reflect.Array:
		s.Type = "array"
		itemZero := reflect.Zero(t.Elem()).Interface()
		s.Items = NewSchemaFrom(itemZero)
	case reflect.Map:
		s.Type = "object"
		if t.Key().Kind() == reflect.String {
			valZero := reflect.Zero(t.Elem()).Interface()
			s.AdditionalProperties = NewSchemaFrom(valZero)
		}
	case reflect.Struct:
		s.Type = "object"
		s.Properties = make(map[string]*Schema)
		var required []string
		for i := 0; i < t.NumField(); i++ {
			f := t.Field(i)

			name, req, ignored := parseJSONTag(f)
			if f.PkgPath != "" {
				continue
			}
			if f.Anonymous {
				// if is ptr
				if f.Type.Kind() == reflect.Ptr {
					continue
				}
				// if struct
				if f.Type.Kind() == reflect.Struct {
					explode := f.Tag.Get("explode") == "1"
					if explode {
						numFields := f.Type.NumField()
						for j := 0; j < numFields; j++ {
							f2 := f.Type.Field(j)
							if f2.PkgPath != "" {
								continue
							}
							if f2.Anonymous {
								continue
							}
							name2, req2, ignored2 := parseJSONTag(f2)
							if ignored2 {
								continue
							}
							if req2 {
								required = append(required, name2)
							}
							zero2 := reflect.Zero(f2.Type).Interface()
							s.Properties[name2] = NewSchemaFrom(zero2)
						}
					}
				}
				continue
			}

			if ignored {
				continue
			}
			if req {
				required = append(required, name)
			}
			zero := reflect.Zero(f.Type).Interface()
			s.Properties[name] = NewSchemaFrom(zero)
		}
		if len(required) > 0 {
			s.Required = required
		}
	default:
		s.Type = "object"
	}

	s.Example = v
	return s
}

// NewParameter
func NewParameter(name string, in string, typ any, required bool) Parameter {
	return Parameter{Name: name, In: in, Required: required, Schema: NewSchemaFrom(typ)}
}

func NewMediaTypeFrom(v any) MediaType {
	return MediaType{Schema: NewSchemaFrom(v)}
}

func NewMediaBody(v any, contentType string) map[string]MediaType {
	ct := contentType
	if ct == "" {
		ct = "application/json"
	}
	return map[string]MediaType{ct: NewMediaTypeFrom(v)}
}

func NewRequestBody(v any, required bool, contentType string) *RequestBody {
	return &RequestBody{Content: NewMediaBody(v, contentType), Required: required}
}

func NewResponseFrom(v any, description string, contentType string) Response {
	d := description
	if d == "" {
		d = "OK"
	}
	return Response{Description: d, Content: NewMediaBody(v, contentType)}
}

func NewResponses(status int, v any, description string, contentType string) Responses {
	code := strconv.Itoa(status)
	return Responses{code: NewResponseFrom(v, description, contentType)}
}

type jsonTagInfo struct {
	Name     string
	Required bool
	Ignored  bool
}

func parseJSONTag(f reflect.StructField) (string, bool, bool) {
	required := false
	if _, ok := f.Tag.Lookup("validate"); ok {
		if strings.Contains(f.Tag.Get("validate"), "required") {
			required = true
		}
	}
	name := f.Name
	tag := f.Tag.Get("json")
	if tag == "" {
		return name, required, false
	}
	parts := strings.Split(tag, ",")
	if len(parts) > 0 {
		if parts[0] == "-" {
			return name, required, true
		}
		if parts[0] != "" {
			name = parts[0]
		}
	}
	return name, required, false
}
