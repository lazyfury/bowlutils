package openapi

import "strings"

func (d *Document) WithInfo(title, version string) *Document {
    d.Info = NewInfo(title, version)
    return d
}

func (d *Document) AddServer(s Server) *Document {
    d.Servers = append(d.Servers, s)
    return d
}

func (d *Document) EnsurePaths() Paths {
    if d.Paths == nil {
        d.Paths = make(Paths)
    }
    return d.Paths
}

func (d *Document) AddPathItem(path string, item PathItem) *Document {
    d.EnsurePaths()
    d.Paths[path] = item
    return d
}

func (d *Document) AddOperation(path string, method string, op Operation) *Document {
    d.EnsurePaths()
    pi := d.Paths[path]
    switch strings.ToLower(method) {
    case "get":
        pi.Get = &op
    case "put":
        pi.Put = &op
    case "post":
        pi.Post = &op
    case "delete":
        pi.Delete = &op
    case "options":
        pi.Options = &op
    case "head":
        pi.Head = &op
    case "patch":
        pi.Patch = &op
    case "trace":
        pi.Trace = &op
    }
    d.Paths[path] = pi
    return d
}

func (d *Document) AddGet(path string, op Operation) *Document  { return d.AddOperation(path, "get", op) }
func (d *Document) AddPost(path string, op Operation) *Document { return d.AddOperation(path, "post", op) }
func (d *Document) AddPut(path string, op Operation) *Document  { return d.AddOperation(path, "put", op) }
func (d *Document) AddDel(path string, op Operation) *Document  { return d.AddOperation(path, "delete", op) }

func (d *Document) EnsureComponents() *Components {
    if d.Components == nil {
        d.Components = &Components{}
    }
    if d.Components.Schemas == nil {
        d.Components.Schemas = make(map[string]*Schema)
    }
    if d.Components.Responses == nil {
        d.Components.Responses = make(map[string]Response)
    }
    if d.Components.RequestBodies == nil {
        d.Components.RequestBodies = make(map[string]RequestBody)
    }
    return d.Components
}

func (d *Document) RegisterSchema(name string, v any) *Schema {
    d.EnsureComponents()
    if s, ok := d.Components.Schemas[name]; ok {
        return s
    }
    s := NewSchemaFrom(v)
    d.Components.Schemas[name] = s
    return s
}

func (d *Document) BodyRef(name string, contentType string) map[string]MediaType {
    ct := contentType
    if ct == "" {
        ct = "application/json"
    }
    return map[string]MediaType{ct: {Schema: &Schema{Ref: "#/components/schemas/" + name}}}
}

func (d *Document) NewRequestBodyRef(name string, v any, required bool, contentType string) *RequestBody {
    d.RegisterSchema(name, v)
    rb := &RequestBody{Content: d.BodyRef(name, contentType), Required: required}
    d.EnsureComponents()
    d.Components.RequestBodies[name] = *rb
    return rb
}

func (d *Document) NewResponseRef(name string, v any, description string, contentType string) Response {
    d.RegisterSchema(name, v)
    desc := description
    if desc == "" {
        desc = "OK"
    }
    r := Response{Description: desc, Content: d.BodyRef(name, contentType)}
    d.EnsureComponents()
    d.Components.Responses[name] = r
    return r
}

