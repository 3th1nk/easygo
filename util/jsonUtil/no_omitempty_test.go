package jsonUtil

import (
	"fmt"
	"testing"
)

func TestMarshalNoOmitempty(t *testing.T) {
	type tp2 struct {
		Id   int    `json:"id,omitempty"`
		Name string `json:"name,omitempty"`
	}
	type tp1 struct {
		Id    int                    `json:"id,omitempty"`
		Name  string                 `json:"name,omitempty"`
		Info  map[string]interface{} `json:"info,omitempty"`
		Info2 map[string]interface{} `json:"info2,omitempty"`
		Data  *tp2                   `json:"data,omitempty"`
	}

	obj := &tp1{
		Id: 1,
		Info: map[string]interface{}{
			"id": 1,
			"data": &tp2{
				Id:   123,
				Name: "",
			},
		},
		Data: &tp2{
			Id:   123,
			Name: "",
		},
	}
	fmt.Println(defaultApi.MustMarshalToString(obj))
	fmt.Println(defaultApi.NoOmitemptyApi().MustMarshalToString(obj))
}
