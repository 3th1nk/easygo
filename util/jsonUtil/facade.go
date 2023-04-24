package jsonUtil

import (
	jsonIter "github.com/json-iterator/go"
)

func Wrap(val interface{}) jsonIter.Any {
	return jsonIter.Wrap(val)
}

func Get(data []byte, path ...interface{}) jsonIter.Any {
	return defaultApi.Get(data, path...)
}

func GetString(str string, path ...interface{}) jsonIter.Any {
	return defaultApi.GetString(str, path...)
}

func Unmarshal(data []byte, v interface{}) error {
	return defaultApi.Unmarshal(data, v)
}

func UnmarshalFromString(str string, v interface{}) error {
	return defaultApi.UnmarshalFromString(str, v)
}

func UnmarshalFromObject(src, dest interface{}) error {
	return defaultApi.UnmarshalFromObject(src, dest)
}

func Marshal(v interface{}) ([]byte, error) {
	return defaultApi.Marshal(v)
}

func MustMarshal(v interface{}) []byte {
	return defaultApi.MustMarshal(v)
}

func MarshalIndent(v interface{}, indentAndPrefix ...string) ([]byte, error) {
	return defaultApi.MarshalIndent(v, indentAndPrefix...)
}

func MustMarshalIndent(v interface{}, indentAndPrefix ...string) []byte {
	return defaultApi.MustMarshalIndent(v, indentAndPrefix...)
}

func MarshalToString(v interface{}) (string, error) {
	return defaultApi.MarshalToString(v)
}

func MustMarshalToString(v interface{}) string {
	return defaultApi.MustMarshalToString(v)
}

func MarshalToStringIndent(v interface{}, indentAndPrefix ...string) (string, error) {
	return defaultApi.MarshalToStringIndent(v, indentAndPrefix...)
}

func MustMarshalToStringIndent(v interface{}, indentAndPrefix ...string) string {
	return defaultApi.MustMarshalToStringIndent(v, indentAndPrefix...)
}
