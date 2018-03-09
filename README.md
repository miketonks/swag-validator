# Swagger Validation Middleware

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

Generate Swagger Documentation automatically:

```
$ curl http://localhost:8089/swagger | jq
  % Total    % Received % Xferd  Average Speed   Time    Time     Time  Current
                                 Dload  Upload   Total   Spent    Left  Speed
100  1950  100  1950    0     0  1386k      0 --:--:-- --:--:-- --:--:-- 1904k
{
  "swagger": "2.0",
  "info": {
    "title": "Your API Title",
    "description": "Describe your API",
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

and Validate your API requests automatically based on this definition.

# Sample

See /sample for fully working example and test cases.
