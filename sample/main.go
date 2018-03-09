package main

import "net/http"

import "github.com/miketonks/swag-validator/sample/server"

func main() {
	api := server.SetupAPI()
	router := server.SetupRouter(api)

	http.ListenAndServe(":8089", router)
}
