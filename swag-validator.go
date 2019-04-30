package swagvalidator

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"runtime"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/miketonks/swag/swagger"
	"github.com/xeipuuv/gojsonschema"
)

// MaxMemory ...
const MaxMemory = 1 * 1024 * 1024

// RequestSchema ...
type RequestSchema struct {
	Title                string                      `json:"title"`
	Type                 string                      `json:"type"`
	Summary              string                      `json:"summary"`
	Properties           map[string]interface{}      `json:"properties"`
	Required             []string                    `json:"required"`
	Definitions          map[string]SchemaDefinition `json:"definitions"`
	AdditionalProperties bool                        `json:"additionalProperties"`
}

// RequestParameter ...
type RequestParameter struct {
	Name                 string         `json:"name,omitempty"`
	Type                 string         `json:"type,omitempty"`
	Format               string         `json:"format,omitempty"`
	Items                *swagger.Items `json:"items,omitempty"`
	Nullable             bool           `json:"nullable,omitempty"`
	Enum                 []string       `json:"enum,omitempty"`
	Pattern              string         `json:"pattern,omitempty"`
	MaxLength            int            `json:"maxLength,omitempty"`
	MinLength            int            `json:"minLength,omitempty"`
	Minimum              *int64         `json:"minimum,omitempty"`
	Maximum              *int64         `json:"maximum,omitempty"`
	ExclusiveMinimum     bool           `json:"exclusiveMinimum,omitempty"`
	ExclusiveMaximum     bool           `json:"exclusiveMaximum,omitempty"`
	AdditionalProperties interface{}    `json:"additionalProperties,omitempty"`
}

// SchemaDefinition ...
type SchemaDefinition struct {
	Name                 string                    `json:"-"`
	Type                 string                    `json:"type"`
	Format               string                    `json:"format,omitempty"`
	Required             []string                  `json:"required,omitempty"`
	Properties           map[string]SchemaProperty `json:"properties,omitempty"`
	Enum                 []string                  `json:"enum,omitempty"`
	Pattern              string                    `json:"pattern,omitempty"`
	MinLength            int                       `json:"minLength,omitempty"`
	MaxLength            int                       `json:"maxLength,omitempty"`
	Minimum              *int64                    `json:"minimum,omitempty"`
	Maximum              *int64                    `json:"maximum,omitempty"`
	ExclusiveMinimum     bool                      `json:"exclusiveMinimum,omitempty"`
	ExclusiveMaximum     bool                      `json:"exclusiveMaximum,omitempty"`
	AdditionalProperties bool                      `json:"additionalProperties"`
}

// SchemaProperty ...
type SchemaProperty struct {
	Type                 []string       `json:"type,omitempty"`
	Description          string         `json:"description,omitempty"`
	Enum                 []string       `json:"enum,omitempty"`
	Format               string         `json:"format,omitempty"`
	Ref                  string         `json:"$ref,omitempty"`
	Example              string         `json:"example,omitempty"`
	Items                *swagger.Items `json:"items,omitempty"`
	Pattern              string         `json:"pattern,omitempty"`
	MinLength            int            `json:"minLength,omitempty"`
	MaxLength            int            `json:"maxLength,omitempty"`
	Minimum              *int64         `json:"minimum,omitempty"`
	Maximum              *int64         `json:"maximum,omitempty"`
	ExclusiveMinimum     bool           `json:"exclusiveMinimum,omitempty"`
	ExclusiveMaximum     bool           `json:"exclusiveMaximum,omitempty"`
	AdditionalProperties interface{}    `json:"additionalProperties,omitempty"`
}

func loadValueForKey(properties map[string]interface{}, key string, values []string) interface{} {
	valueType := ""
	valueFormat := ""
	elemType := ""
	elemFormat := ""
	propI, found := properties[key]
	if found {
		prop := propI.(map[string]interface{})
		t, found := prop["type"]
		if found {
			valueType = t.(string)
		}
		f, found := prop["format"]
		if found {
			valueFormat = f.(string)
		}
		if items, ok := prop["items"]; ok {
			if t, ok := items.(map[string]interface{})["type"]; ok {
				elemType = t.(string)
			}
			if f, ok := items.(map[string]interface{})["format"]; ok {
				elemFormat = f.(string)
			}
		}
	}

	// if parameter isn't an array and we didn't receive multiple values, pass it as a normal value
	if len(values) == 1 && valueType != "array" {
		return coerce(values[0], valueType, valueFormat)
	}

	// if we received multiple values, use them as the elements; otherwise, split the value we got
	var items []string
	if len(values) > 1 {
		items = values
	} else {
		items = strings.Split(values[0], ",")
	}

	result := []interface{}{}
	for _, item := range items {
		result = append(result, coerce(strings.TrimSpace(item), elemType, elemFormat))
	}
	return result
}

// SwaggerValidator middleware
func SwaggerValidator(api *swagger.API) gin.HandlerFunc {

	apiMap := map[string]gojsonschema.JSONLoader{}
	for _, p := range api.Paths {
		for _, e := range []*swagger.Endpoint{
			p.Delete,
			p.Get,
			p.Post,
			p.Put,
			p.Patch,
			p.Head,
			p.Options,
			p.Trace,
			p.Connect} {
			if e != nil && e.Handler != nil {
				schema := buildRequestSchema(e)
				schema.Definitions = buildSchemaDefinitions(api)
				schemaLoader := gojsonschema.NewGoLoader(schema)

				handler := nameOfFunction(e.Handler)
				apiMap[handler] = schemaLoader
			}
		}
	}

	// This part runs at runtime, with context for individual request
	return func(c *gin.Context) {
		schemaLoader, found := apiMap[c.HandlerName()]
		if !found {
			c.Next()
			return
		}
		ref, _ := schemaLoader.LoadJSON()
		properties, _ := ref.(map[string]interface{})["properties"].(map[string]interface{})

		document := map[string]interface{}{}

		for _, p := range c.Params {
			document[p.Key] = loadValueForKey(properties, p.Key, []string{p.Value})
		}
		for k, v := range c.Request.URL.Query() {
			document[k] = loadValueForKey(properties, k, v)
		}

		// For muiltipart form, handle params and file uploads
		if c.ContentType() == "multipart/form-data" {
			r := c.Request
			r.ParseMultipartForm(MaxMemory)

			for k, v := range c.Request.PostForm {
				document[k] = coerce(v[0], "", "")
			}
			if r.MultipartForm != nil && r.MultipartForm.File != nil {
				for k := range r.MultipartForm.File {
					document[k] = "x"
				}
			}
		} else if c.ContentType() == "application/x-www-form-urlencoded" {
			r := c.Request
			r.ParseForm()

			body := map[string]interface{}{}
			for k, v := range c.Request.PostForm {
				body[k] = coerce(v[0], "", "")
			}
			document["body"] = body
		} else if c.Request.ContentLength > 0 {
			// For all other types parse body as json, if possible

			// read the response body to a variable
			var body interface{}
			b, err := ioutil.ReadAll(c.Request.Body)
			if err != nil {
				c.AbortWithStatusJSON(
					http.StatusBadRequest,
					gin.H{
						"message": "Validation error",
						"details": map[string]string{
							"body": "Failed to read request body",
						},
					},
				)
				return
			}
			err = json.Unmarshal(b, &body)
			// TODO Consider different error cases: Empty Body, Invalid JSON, Form Data
			if err != nil {
				c.AbortWithStatusJSON(
					http.StatusBadRequest,
					gin.H{
						"message": "Validation error",
						"details": map[string]string{
							"body": "Invalid JSON format",
						},
					},
				)
				return
			}
			document["body"] = body

			//reset the response body to the original unread state
			c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(b))
		}

		gojsonschema.Locale = CustomLocale{}

		documentLoader := gojsonschema.NewGoLoader(document)
		result, err := gojsonschema.Validate(schemaLoader, documentLoader)

		if err != nil {
			// fmt.Printf("ERROR: %s\n", err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"message": "swagger document " + err.Error(),
			})

		} else if result.Valid() {
			// fmt.Printf("The document is valid\n")
			c.Next()
		} else {
			// fmt.Printf("The document is not valid. see errors :\n")
			errors := map[string]string{}
			for _, err := range result.Errors() {
				description := err.Description()
				details := err.Details()

				if val, ok := details["property"]; ok {
					field := val.(string)
					errors[field] = description
				} else {
					field := details["field"].(string)
					field = strings.TrimPrefix(field, "body.")
					errors[field] = description
				}
			}
			// fmt.Printf("The document is not valid. see errors : %+v\n", errors)
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"message": "Validation error",
				"details": errors,
			})
		}
	}
}

// Data types are defined here: https://swagger.io/specification/#dataTypes
func coerce(value string, valueType string, valueFormat string) interface{} {
	switch valueType {
	case "integer":
		bitSize := 32
		if valueFormat == "int64" {
			bitSize = 64
		}
		v, err := strconv.ParseInt(value, 10, bitSize)
		if err == nil {
			return v
		}
	case "number":
		bitSize := 32
		if valueFormat == "double" {
			bitSize = 64
		}
		v, err := strconv.ParseFloat(value, bitSize)
		if err == nil {
			return v
		}
	case "string":
		if valueFormat == "byte" {
			return []byte(value)
		}
		return value
	case "boolean":
		v, err := strconv.ParseBool(value)
		if err == nil {
			return v
		}
	default:
		return value
	}

	return value
}

func nameOfFunction(f interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()
}

func buildRequestSchema(e *swagger.Endpoint) *RequestSchema {
	r := RequestSchema{
		Title:      fmt.Sprintf("%s %s", e.Method, e.Path),
		Type:       "object",
		Properties: make(map[string]interface{}),
		Required:   []string{},
	}

	if len(e.Parameters) == 0 {
		return &r
	}

	for _, p := range e.Parameters {
		if p.Required {
			r.Required = append(r.Required, p.Name)
		}

		if p.Name != "" && p.Schema != nil {
			r.Properties[p.Name] = p.Schema

			// TODO Consider if we should use Ref to optimise definitions
			// if p.Schema.Ref != "" {
			// 	parts := strings.Split(p.Schema.Ref, "/")
			// 	last := parts[len(parts)-1]
			//
			// 	fmt.Printf("DEF: %+v", defs[last])
			// 	//r.Properties[p.Name] = defs[last]
			// }

		} else if p.Name != "" {

			param := RequestParameter{
				Name:                 p.Name,
				Type:                 p.Type,
				Format:               p.Format,
				Nullable:             p.Nullable,
				Items:                p.Items,
				Enum:                 p.Enum,
				MinLength:            p.MinLength,
				MaxLength:            p.MaxLength,
				Minimum:              p.Minimum,
				Maximum:              p.Maximum,
				ExclusiveMinimum:     p.ExclusiveMinimum,
				ExclusiveMaximum:     p.ExclusiveMaximum,
				AdditionalProperties: p.AdditionalProperties,
			}
			// for validation purposes, file can be treated as string type
			if p.Type == "file" {
				param.Type = "string"
			}

			r.Properties[p.Name] = param
		}
	}

	return &r
}

func buildSchemaDefinitions(api *swagger.API) map[string]SchemaDefinition {
	defs := map[string]SchemaDefinition{}
	for _, d := range api.Definitions {
		schemaDef := SchemaDefinition{
			Name:       d.Name,
			Type:       d.Type,
			Format:     d.Format,
			Required:   d.Required,
			Properties: map[string]SchemaProperty{},
		}
		for k, p := range d.Properties {
			sp := SchemaProperty{
				Description:          p.Description,
				Enum:                 p.Enum,
				Format:               p.Format,
				Ref:                  p.Ref,
				Example:              p.Example,
				Items:                p.Items,
				MinLength:            p.MinLength,
				MaxLength:            p.MaxLength,
				Minimum:              p.Minimum,
				Maximum:              p.Maximum,
				ExclusiveMinimum:     p.ExclusiveMinimum,
				ExclusiveMaximum:     p.ExclusiveMaximum,
				AdditionalProperties: p.AdditionalProperties,
			}
			if p.Type != "" {
				sp.Type = strings.Split(p.Type, ",")
			}
			if p.Nullable {
				sp.Type = append(sp.Type, "null")
			}

			// for json.RawMessage
			if p.GoType.PkgPath() == "encoding/json" && p.GoType.Name() == "RawMessage" {
				sp.Type = []string{"raw_message"}
			}

			schemaDef.Properties[k] = sp
		}
		defs[d.Name] = schemaDef
	}
	return defs
}
