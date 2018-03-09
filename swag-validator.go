package swag_validator

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"runtime"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/miketonks/swag/swagger"
	"github.com/xeipuuv/gojsonschema"
)

// RequestSchema ...
type RequestSchema struct {
	Title                string                    `json:"title"`
	Type                 string                    `json:"type"`
	Summary              string                    `json:"summary"`
	Properties           map[string]interface{}    `json:"properties"`
	Required             []string                  `json:"required"`
	Definitions          map[string]swagger.Object `json:"definitions"`
	AdditionalProperties bool                      `json:"additionalProperties"`
}

// RequestParameter ...
type RequestParameter struct {
	Name   string `json:"name,omitempty"`
	Type   string `json:"type,omitempty"`
	Format string `json:"format,omitempty"`
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
			if e != nil {
				schema := buildRequestSchema(e)
				schema.Definitions = api.Definitions
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
		document := map[string]interface{}{}

		for _, p := range c.Params {
			document[p.Key] = coerce(p.Value)
		}
		for k, v := range c.Request.URL.Query() {
			// TODO Consider if we need to support multiple param values
			document[k] = coerce(v[0])
		}

		// TODO Maybe c.ContentType(), but it's nice to assume json as default

		// read the response body to a variable
		if c.Request.ContentLength > 0 {
			var body interface{}
			b, err := ioutil.ReadAll(c.Request.Body)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Failed to read request body"})
				return
			}
			err = json.Unmarshal(b, &body)
			// TODO Consider different error cases: Empty Body, Invalid JSON, Form Data
			if err != nil {
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid json format"})
				return
			}
			document["body"] = body

			//reset the response body to the original unread state
			c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(b))
		}
		documentLoader := gojsonschema.NewGoLoader(document)

		result, err := gojsonschema.Validate(schemaLoader, documentLoader)
		if err != nil {
			fmt.Printf("ERROR: %s", err)
			c.Next()

		} else if result.Valid() {
			//fmt.Printf("The document is valid\n")
			c.Next()

		} else {
			//fmt.Printf("The document is not valid. see errors :\n")
			errors := []string{}
			for _, err := range result.Errors() {
				// Err implements the ResultError interface
				errors = append(errors, fmt.Sprintf("%s", err))
			}
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": errors})
		}
	}
}

func coerce(value string) interface{} {
	// TODO Add other types, float, bool.. etc
	i, err := strconv.Atoi(value)
	if err == nil {
		return i
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

			r.Properties[p.Name] = RequestParameter{
				Name:   p.Name,
				Type:   p.Type,
				Format: p.Format,
			}
		}
	}

	return &r
}
