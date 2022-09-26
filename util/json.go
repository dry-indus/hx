package util

import (
	jsoniter "github.com/json-iterator/go"
	"github.com/json-iterator/go/extra"
)

func init() {
	extra.RegisterFuzzyDecoders()
}

var JSON = jsoniter.ConfigCompatibleWithStandardLibrary

func MustMarshalToString(v interface{}) string {
	s, _ := JSON.MarshalToString(v)
	return s
}
