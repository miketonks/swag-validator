# Swagger Validation Middleware

Swagger generation and validation for gin and echo server.

[![Build Status](https://travis-ci.com/miketonks/swag-validator.svg?branch=master)](https://travis-ci.com/miketonks/swag-validator)
[![GoDoc](http://img.shields.io/badge/go-documentation-blue.svg)](http://godoc.org/github.com/miketonks/swag-validator)

## Definition

Define API Specification using code:

```
get := endpoint.New("get", "/pet/{petId}", "Find pet by ID",
  endpoint.Handler(GetPet),
  endpoint.Path("petId", "integer", "", "ID of pet to return", true),
  endpoint.Query("foo", "integer", "Some foo", false),
  endpoint.Response(http.StatusOK, Pet{}, "successful operation"),
)

post := endpoint.New("post", "/pet", "Add a new pet to the store",
  endpoint.Handler(PostPet),
  endpoint.Description("Additional information on adding a pet to the store"),
  endpoint.Body(Pet{}, "Pet object that needs to be added to the store", true),
  endpoint.Response(http.StatusOK, Pet{}, "Successfully added pet"),
)
```

## Middleware

Add the middleware to your server:

```
r.GET("/swagger", gin.WrapH(api.Handler(enableCors)))

r.Use(swag_validator.SwaggerValidator(api))
```

## Swagger Docs

Generates Swagger Documentation automatically:

```
$ curl http://localhost:8089/swagger | jq
{
  "swagger": "2.0",
  "info": {
    "title": "Petstore",
    "description": "Sample Petstore API",
  },
  "basePath": "/",
  "paths": {
    "/pet/{petId}": {
      "get": {
        "summary": "Find pet by ID",
        "parameters": [
          {
            "in": "path",
            "name": "petId",
            "description": "ID of pet to return",
            "required": true,
            "type": "integer"
          },
...
```

## Validation

Validate your API requests automatically based on this definition.

```
curl http://localhost:8089/pet/123
{"id":0,"uuid":"00000000-0000-0000-0000-000000000000","category":{"category":0,"name":""},"name":"Ollie","photoUrls":null,"tags":null,"age":0,"grumpy":false,"dob":"00

curl http://localhost:8089/pet/foo
{"error":["petId: Invalid type. Expected: integer, given: string"]}
```

# Sample

See /sample for working example and test cases.

# Credits and Thanks

Thanks to savaki, xeipuuv, and of course the gin team, for their awesome libraries that make this possible!

swag: https://github.com/savaki/swag

gojsonschema: https://github.com/xeipuuv/gojsonschema

gin: https://github.com/gin-gonic/gin
