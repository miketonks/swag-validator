package swagvalidator_test

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo"
	swag "github.com/miketonks/swag"
	"github.com/miketonks/swag/endpoint"
	"github.com/miketonks/swag/swagger"
	"github.com/stretchr/testify/assert"

	sv "github.com/miketonks/swag-validator"
)

func createEngineEcho(api *swagger.API) (r *echo.Echo) {
	r = echo.New()
	r.Use(sv.SwaggerValidatorEcho(api))
	api.Walk(func(path string, endpoint *swagger.Endpoint) {
		h := endpoint.Handler.(func(echo.Context) error)
		path = swag.ColonPath(path)

		r.Router().Add(endpoint.Method, path, h)
	})
	return
}

func handler(c echo.Context) error {
	return nil
}

func TestQueryEcho(t *testing.T) {
	testTable := []struct {
		description      string
		query            string
		expectedStatus   int
		expectedResponse map[string]interface{}
	}{
		{
			description:    "Non-int value in an int query param",
			query:          "int_param=abc",
			expectedStatus: 400,
			expectedResponse: map[string]interface{}{
				"int_param": "Invalid type. Expected: integer, given: string",
			},
		},
		{
			description:      "Int value in an int query param",
			query:            "int_param=10",
			expectedStatus:   200,
			expectedResponse: nil,
		},
		{
			description:    "Non-UUID value in an uuid query param",
			query:          "uuid_param=abc",
			expectedStatus: 400,
			expectedResponse: map[string]interface{}{
				"uuid_param": "Field does not match format 'uuid'",
			},
		},
		{
			description:      "UUID value in an int query param",
			query:            "uuid_param=" + testUUID,
			expectedStatus:   200,
			expectedResponse: nil,
		},
		{
			description:    "Not allowed enum value in enum query param",
			query:          "enum_param=baz",
			expectedStatus: 400,
			expectedResponse: map[string]interface{}{
				"enum_param": "Must be one of the following: \"foo\", \"bar\"",
			},
		},
	}

	api := swag.New(swag.Endpoints(endpoint.New("GET", "/validate-test", "Test query params",
		endpoint.Handler(handler),
		endpoint.QueryMap(map[string]swagger.Parameter{
			"int_param": {
				Type: "integer",
			},
			"uuid_param": {
				Type:   "string",
				Format: "uuid",
			},
			"enum_param": {
				Type: "string",
				Enum: []string{"foo", "bar"},
			},
		}),
	)))

	r := createEngineEcho(api)

	for _, tt := range testTable {
		t.Run(tt.description, func(t *testing.T) {

			w := httptest.NewRecorder()

			url := fmt.Sprintf("/validate-test?%s", tt.query)

			req, err := http.NewRequest("GET", url, nil)
			if err != nil {
				log.Fatalf("Error preparing request: %s", err)
			}

			r.ServeHTTP(w, req)

			var body map[string]interface{}

			if w.Body != nil && w.Body.String() != "" {
				err := json.Unmarshal(w.Body.Bytes(), &body)
				if err != nil {
					panic(fmt.Sprintf("Failed to unmarshal body while running test: %q. Error: %s", tt.description, err))
				}

				assert.Equal(t, tt.expectedResponse, body["details"])
			}

			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}

func TestPathEcho(t *testing.T) {
	testTable := []struct {
		url      string
		urlWParm string
		path     endpoint.Option
		cases    []pathCase
	}{
		{
			url:      "/int-test",
			urlWParm: "/int-test/{int_id}",
			path:     endpoint.Path("int_id", "integer", "integer", ""),
			cases: []pathCase{
				{
					description:    "non-int path param",
					pathParam:      "abc",
					expectedStatus: 400,
					expectedResponse: map[string]interface{}{
						"int_id": "Invalid type. Expected: integer, given: string",
					},
				},
				{
					description:      "int path param",
					pathParam:        "10",
					expectedStatus:   200,
					expectedResponse: nil,
				},
			},
		},
		{
			url:      "/uuid-test",
			urlWParm: "/uuid-test/{uuid_id}",
			path:     endpoint.Path("uuid_id", "string", "uuid", ""),
			cases: []pathCase{
				{
					description:    "non-uuid path param",
					pathParam:      "10",
					expectedStatus: 400,
					expectedResponse: map[string]interface{}{
						"uuid_id": "Field does not match format 'uuid'",
					},
				},
				{
					description:      "uuid path param",
					pathParam:        testUUID,
					expectedStatus:   200,
					expectedResponse: nil,
				},
			}},
	}

	// Can't bind multiple endpoints to the same handler
	// Even if the handler is a lambda, it still does not work.
	// Therefore have to create a new api for each endpoint iteratively
	for _, testCase := range testTable {
		api := swag.New(
			swag.Endpoints(endpoint.New("GET", "/validate-test"+testCase.urlWParm, "Test the validator",
				endpoint.Handler(func(echo.Context) error { return nil }),
				testCase.path,
			)))

		r := createEngineEcho(api)

		for _, tt := range testCase.cases {
			t.Run(tt.description, func(t *testing.T) {

				w := httptest.NewRecorder()

				url := fmt.Sprintf("/validate-test%s/%s", testCase.url, tt.pathParam)

				req, err := http.NewRequest("GET", url, nil)
				if err != nil {
					log.Fatalf("Error preparing request: %s", err)
				}

				r.ServeHTTP(w, req)

				var body map[string]interface{}

				if w.Body != nil && w.Body.String() != "" {
					err := json.Unmarshal(w.Body.Bytes(), &body)
					if err != nil {
						panic(fmt.Sprintf("Failed to unmarshal body while running test: %q. Error: %s", tt.description, err))
					}

					assert.Equal(t, tt.expectedResponse, body["details"])
				}

				assert.Equal(t, tt.expectedStatus, w.Code)
			})
		}
	}
}

func TestPayloadEcho(t *testing.T) {

	testTable := []struct {
		description      string
		in               payload
		expectedStatus   int
		expectedResponse map[string]interface{}
	}{
		{
			description:    "Scalar uuid tag with non-uuid value",
			in:             payload{FormatString: "not-a-uuid"},
			expectedStatus: 400,
			expectedResponse: map[string]interface{}{
				"format_str": "Field does not match format 'uuid'",
			},
		},
		{
			description:      "Scalar uuid tag with uuid value",
			in:               payload{FormatString: testUUID},
			expectedStatus:   200,
			expectedResponse: nil,
		},
		{
			description:    "Non-UUID string in a UUID array",
			in:             payload{FormatStringArr: []string{"not-a-uuid"}},
			expectedStatus: 400,
			expectedResponse: map[string]interface{}{
				"format_str_arr.0": "Field does not match format 'uuid'",
			},
		},
		{
			description:      "UUID strings in a UUID array",
			in:               payload{FormatStringArr: []string{testUUID}},
			expectedStatus:   200,
			expectedResponse: nil,
		},
		{
			description:    "String shorter than minimum required",
			in:             payload{MinLenString: "1234"},
			expectedStatus: 400,
			expectedResponse: map[string]interface{}{
				"min_len_str": "String length must be greater than or equal to 5",
			},
		},
		{
			description:      "String longer than minimum required",
			in:               payload{MinLenString: "123456"},
			expectedStatus:   200,
			expectedResponse: nil,
		},
		{
			description:    "String in an array shorter than minimum required",
			in:             payload{MinLenStringArr: []string{"1234"}},
			expectedStatus: 400,
			expectedResponse: map[string]interface{}{
				"min_len_str_arr.0": "String length must be greater than or equal to 5",
			},
		},
		{
			description:      "Strings in an array longer than minimum required",
			in:               payload{MinLenStringArr: []string{"12345", "123456"}},
			expectedStatus:   200,
			expectedResponse: nil,
		},
		{
			description:    "String longer than maximum allowed",
			in:             payload{MaxLenString: "12345678"},
			expectedStatus: 400,
			expectedResponse: map[string]interface{}{
				"max_len_str": "String length must be less than or equal to 7",
			},
		},
		{
			description:      "String shoter or equal to maximum allowed",
			in:               payload{MaxLenString: "123456"},
			expectedStatus:   200,
			expectedResponse: nil,
		},
		{
			description:    `String in an array longer than maximum allowed`,
			in:             payload{MaxLenStringArr: []string{"12345678"}},
			expectedStatus: 400,
			expectedResponse: map[string]interface{}{
				"max_len_str_arr.0": "String length must be less than or equal to 7",
			},
		},
		{
			description:      "Strings in an array shorter than or euqal to maximum allowed",
			in:               payload{MaxLenStringArr: []string{"123456", "1234567"}},
			expectedStatus:   200,
			expectedResponse: nil,
		},
		{
			description:    "String does not match enumaration",
			in:             payload{EnumString: "test"},
			expectedStatus: 400,
			expectedResponse: map[string]interface{}{
				"enum_str": "Must be one of the following: \"Foo\", \"Bar\"",
			},
		},
		{
			description:      "String matches enumeration",
			in:               payload{EnumString: "Foo"},
			expectedStatus:   200,
			expectedResponse: nil,
		},
		{
			description:    "String in an array does not match enumeration",
			in:             payload{EnumStringArr: []string{"test"}},
			expectedStatus: 400,
			expectedResponse: map[string]interface{}{
				"enum_str_arr.0": "Must be one of the following: \"Bar\"",
			},
		},
		{
			description:      `Strings in an arrya match enumeration`,
			in:               payload{EnumStringArr: []string{"Bar"}},
			expectedStatus:   200,
			expectedResponse: nil,
		},
		{
			description:    "Number is smaller than minimum required",
			in:             payload{Minimum: 4},
			expectedStatus: 400,
			expectedResponse: map[string]interface{}{
				"minimum": "Must be greater than or equal to 5",
			},
		},
		{
			description:      "Number is gte to minimum required",
			in:               payload{Minimum: 5},
			expectedStatus:   200,
			expectedResponse: nil,
		},
		{
			description:    "Number is greater than allowed",
			in:             payload{Maximum: 2},
			expectedStatus: 400,
			expectedResponse: map[string]interface{}{
				"maximum": "Must be less than or equal to 1",
			},
		},
		{
			description:      "Number is lte to maximum allowed",
			in:               payload{Maximum: 1},
			expectedStatus:   200,
			expectedResponse: nil,
		},
		{
			description:    "Number is smaller than minimum required",
			in:             payload{Minimum: 4},
			expectedStatus: 400,
			expectedResponse: map[string]interface{}{
				"minimum": "Must be greater than or equal to 5",
			},
		},
		{
			description:    "Number is gte to excl minimum required",
			in:             payload{ExclMinimum: 5},
			expectedStatus: 400,
			expectedResponse: map[string]interface{}{
				"excl_minimum": "Must be greater than 5",
			},
		},
		{
			description:    "Number is lte to excl maximum allowed",
			in:             payload{ExclMaximum: 1},
			expectedStatus: 400,
			expectedResponse: map[string]interface{}{
				"excl_maximum": "Must be less than 1",
			},
		},
		{
			description:    "Nested struct field is missing",
			in:             payload{Nested: &nested{}},
			expectedStatus: 400,
			expectedResponse: map[string]interface{}{
				"nested.foo": "foo is required",
			},
		},
		{
			description:      "Nested struct field is present",
			in:               payload{Nested: &nested{Foo: "bar"}},
			expectedStatus:   200,
			expectedResponse: nil,
		},
		{
			description:    "Array contains more items than maximum allowed",
			in:             payload{MaxItemsArr: []string{"1", "2", "3", "4"}},
			expectedStatus: 400,
			expectedResponse: map[string]interface{}{
				"max_items_arr": "Array must have at most 3 items",
			},
		},
		{
			description:      "Array contains less items than maximum allowed",
			in:               payload{MaxItemsArr: []string{"1", "2"}},
			expectedStatus:   200,
			expectedResponse: nil,
		},
		{
			description:    "Array contains less items than minimum allowed",
			in:             payload{MinItemsArr: []string{"1"}},
			expectedStatus: 400,
			expectedResponse: map[string]interface{}{
				"min_items_arr": "Array must have at least 2 items",
			},
		},
		{
			description:      "Array contains more items than minimum required",
			in:               payload{MinItemsArr: []string{"1", "2", "3"}},
			expectedStatus:   200,
			expectedResponse: nil,
		},
		{
			description:    "Unique array contains non-unique items",
			in:             payload{UniqueItemsAarr: []string{"foo", "foo"}},
			expectedStatus: 400,
			expectedResponse: map[string]interface{}{
				"unique_items_arr": "array items[0,1] must be unique",
			},
		},
		{
			description:      "Unique array contains unique items",
			in:               payload{UniqueItemsAarr: []string{"foo", "bar"}},
			expectedStatus:   200,
			expectedResponse: nil,
		},
	}

	api := swag.New(swag.Endpoints(endpoint.New("POST", "/validate-test", "Test the validator",
		endpoint.Handler(func(echo.Context) error { return nil }),
		endpoint.Body(payload{}, "Validation body", true),
	)))

	r := createEngineEcho(api)

	for _, tt := range testTable {
		t.Run(tt.description, func(t *testing.T) {

			w := httptest.NewRecorder()
			req := preparePostRequest("/validate-test", tt.in)
			r.ServeHTTP(w, req)

			var body map[string]interface{}

			if w.Body != nil && w.Body.String() != "" {
				err := json.Unmarshal(w.Body.Bytes(), &body)
				if err != nil {
					panic(fmt.Sprintf("Failed to unmarshal body while running test: %q. Error: %s", tt.description, err))
				}

				assert.Equal(t, tt.expectedResponse, body["details"])
			}

			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}
