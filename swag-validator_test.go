package swagvalidator_test

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
)

var testUUID = "00000000-0000-0000-0000-000000000000"

type pathCase struct {
	description      string
	pathParam        string
	expectedStatus   int
	expectedResponse map[string]interface{}
}

type nested struct {
	Foo string `json:"foo,omitempty" binding:"required"`
}

type payload struct {
	FormatString    string   `json:"format_str,omitempty" format:"uuid"`
	FormatStringArr []string `json:"format_str_arr,omitempty" format:"uuid"`
	MinLenString    string   `json:"min_len_str,omitempty" min_length:"5"`
	MinLenStringArr []string `json:"min_len_str_arr,omitempty" min_length:"5"`
	MaxLenString    string   `json:"max_len_str,omitempty" max_length:"7"`
	MaxLenStringArr []string `json:"max_len_str_arr,omitempty" max_length:"7"`
	EnumString      string   `json:"enum_str,omitempty" enum:"Foo,Bar"`
	EnumStringArr   []string `json:"enum_str_arr,omitempty" enum:"Foo,Bar"`
	PatternString   string   `json:"pattern_str,omitempty" pattern:"^test$"`
	Minimum         int      `json:"minimum,omitempty" minimum:"5"`
	Maximum         int      `json:"maximum,omitempty" maximum:"1"`
	ExclMinimum     int      `json:"excl_minimum,omitempty" minimum:"5" exclusive_minimum:"true"`
	ExclMaximum     int      `json:"excl_maximum,omitempty" maximum:"1" exclusive_maximum:"true"`
	Nested          *nested  `json:"nested,omitempty"`
	MaxItemsArr     []string `json:"max_items_arr,omitempty" max_items:"3"`
	MinItemsArr     []string `json:"min_items_arr,omitempty" min_items:"2"`
	UniqueItemsAarr []string `json:"unique_items_arr,omitempty" unique_items:"true"`
}

func preparePostRequest(url string, body payload) *http.Request {
	buff, err := json.Marshal(body)
	if err != nil {
		log.Fatalf("Failed to marshal the body: %s", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(buff))
	if err != nil {
		log.Fatalf("Error preparing request: %s", err)
	}

	return req
}
