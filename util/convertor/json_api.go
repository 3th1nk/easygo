package convertor

import (
	jsonIter "github.com/json-iterator/go"
	"github.com/modern-go/reflect2"
)

type JsonAPI interface {
	MarshalToString(v interface{}) (string, error)
	UnmarshalFromString(str string, v interface{}) error
}

var (
	jsonApi JsonAPI = defaultJsonApi

	defaultJsonApi = jsonIter.Config{
		EscapeHTML:              false,
		MarshalFloatWith6Digits: true,
		SortMapKeys:             true,
	}.Froze()
)

func SetDefaultJsonApi(api JsonAPI) {
	if !reflect2.IsNil(api) {
		jsonApi = api
	} else {
		jsonApi = defaultJsonApi
	}
}
