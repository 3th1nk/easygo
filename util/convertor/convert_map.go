package convertor

import (
	"fmt"
	"github.com/modern-go/reflect2"
	"reflect"
	"strings"
	"unicode"
)

var (
	strMapType    = reflect.TypeOf(map[string]string{})
	strObjMapType = reflect.TypeOf(map[string]interface{}{})
)

func ToStringMap(a interface{}, api ...JsonAPI) (map[string]string, error) {
	if reflect2.IsNil(a) {
		return nil, nil
	}

	rv := reflect.Indirect(reflect.ValueOf(a))
	if rt := rv.Type(); rt.Kind() == reflect.Map {
		if rt.ConvertibleTo(strMapType) {
			m := rv.Convert(strMapType).Interface().(map[string]string)
			b := make(map[string]string, len(m))
			for k, v := range m {
				b[k] = v
			}
			return b, nil
		} else if rt.ConvertibleTo(strObjMapType) {
			m := rv.Convert(strObjMapType).Interface().(map[string]interface{})
			b := make(map[string]string, len(m))
			var err error
			for k, v := range m {
				if b[k], err = ToString(v); err != nil {
					return nil, err
				}
			}
			return b, nil
		}
	}

	var rtn map[string]string
	if str, err := ToString(a); err != nil {
		return nil, err
	} else if str = strings.TrimLeftFunc(str, unicode.IsSpace); str != "" && str != "null" {
		if str[0] != '{' {
			return nil, fmt.Errorf(`cannot unmarshal array into map[string]string`)
		}

		var theApi JsonAPI
		if len(api) != 0 && api[0] != nil {
			theApi = api[0]
		} else {
			theApi = jsonApi
		}
		if err := theApi.UnmarshalFromString(str, &rtn); err != nil {
			return nil, err
		}
	}
	return rtn, nil
}

func ToStringMapNoError(a interface{}, api ...JsonAPI) map[string]string {
	val, _ := ToStringMap(a, api...)
	return val
}

func ToStringObjectMap(a interface{}, api ...JsonAPI) (map[string]interface{}, error) {
	if reflect2.IsNil(a) {
		return nil, nil
	}

	rv := reflect.Indirect(reflect.ValueOf(a))
	if rt := rv.Type(); rt.Kind() == reflect.Map {
		if rt.ConvertibleTo(strMapType) {
			m := rv.Convert(strMapType).Interface().(map[string]string)
			b := make(map[string]interface{}, len(m))
			for k, v := range m {
				b[k] = v
			}
			return b, nil
		} else if rt.ConvertibleTo(strObjMapType) {
			m := rv.Convert(strObjMapType).Interface().(map[string]interface{})
			b := make(map[string]interface{}, len(m))
			for k, v := range m {
				b[k] = v
			}
			return b, nil
		}
	}

	var rtn map[string]interface{}
	if str, err := ToString(a); err != nil {
		return nil, err
	} else if str = strings.TrimLeftFunc(str, unicode.IsSpace); str != "" && str != "null" {
		if str[0] != '{' {
			return nil, fmt.Errorf(`cannot unmarshal array into map[string]interface{}`)
		}

		var theApi JsonAPI
		if len(api) != 0 && api[0] != nil {
			theApi = api[0]
		} else {
			theApi = jsonApi
		}
		if err := theApi.UnmarshalFromString(str, &rtn); err != nil {
			return nil, err
		}
	}
	return rtn, nil
}

func ToStringObjectMapNoError(a interface{}, api ...JsonAPI) map[string]interface{} {
	val, _ := ToStringObjectMap(a, api...)
	return val
}

func StrToStringMap(str string) (map[string]string, error) {
	m := map[string]string{}
	if err := jsonApi.UnmarshalFromString(str, &m); err != nil {
		return nil, err
	}
	return m, nil
}

func MustStrToStringMap(str string) map[string]string {
	m, _ := StrToStringMap(str)
	return m
}

func StrToStringObjectMap(str string) (map[string]interface{}, error) {
	m := map[string]interface{}{}
	if err := jsonApi.UnmarshalFromString(str, &m); err != nil {
		return nil, err
	}
	return m, nil
}
