package server

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/miketonks/swag"
	"github.com/miketonks/swag-validator"
	"github.com/miketonks/swag/endpoint"
	"github.com/miketonks/swag/swagger"
)

// Category example from the swagger pet store
type Category struct {
	ID   int64  `json:"category"`
	Name string `json:"name"`
}

// Pet example from the swagger pet store
type Pet struct {
	ID          int64        `json:"id"`
	UUID        swagger.UUID `json:"uuid"`
	Category    Category     `json:"category"`
	Name        string       `json:"name" binding:"required"`
	PhotoUrls   []string     `json:"photoUrls"`
	Tags        []string     `json:"tags"`
	Age         *float64     `json:"age"`
	Grumpy      bool         `json:"grumpy"`
	DateOfBirth time.Time    `json:"dob"`
	Tm          swagger.Time `json:"tm"`
	Dt          swagger.Date `json:"dt"`
}

// GetPet Handler
func GetPet(c *gin.Context) {
	c.JSON(http.StatusOK, Pet{Name: "Ollie"})
}

// PostPet Handler
func PostPet(c *gin.Context) {
	var pet Pet
	if err := c.ShouldBindJSON(&pet); err == nil {
		c.JSON(http.StatusOK, pet)
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
}

// SetupAPI ...
func SetupAPI() *swagger.API {
	post := endpoint.New("post", "/pet", "Add a new pet to the store",
		endpoint.Handler(PostPet),
		endpoint.Description("Additional information on adding a pet to the store"),
		endpoint.Body(Pet{}, "Pet object that needs to be added to the store", true),
		endpoint.FormData("upfile", "file", "", "file to upload", false),
		endpoint.Response(http.StatusOK, Pet{}, "Successfully added pet"),
		endpoint.Tags("petstore", "pet"),
	)
	get := endpoint.New("get", "/pet/{petId}", "Find pet by ID",
		endpoint.Handler(GetPet),
		endpoint.Path("petId", "integer", "", "ID of pet to return"),
		endpoint.Query("foo", "integer", "", "Some foo", false),
		endpoint.Response(http.StatusOK, Pet{}, "successful operation", endpoint.Header(
			"x-custom-header", "string", "integer", "custom number")),
	)

	api := swag.New(
		swag.Endpoints(post, get),
	)
	return api
}

// SetupRouter ...
func SetupRouter(api *swagger.API) *gin.Engine {
	router := gin.New()
	enableCors := true
	router.GET("/swagger", gin.WrapH(api.Handler(enableCors)))

	router.Use(swag_validator.SwaggerValidator(api))

	api.Walk(func(path string, endpoint *swagger.Endpoint) {
		h := endpoint.Handler.(func(c *gin.Context))
		path = swag.ColonPath(path)

		router.Handle(endpoint.Method, path, h)
	})
	return router
}
