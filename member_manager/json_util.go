package member_manager

import "github.com/json-iterator/go"

//unmarshal
var jsonCompatible = jsoniter.Config{
	EscapeHTML:             true,
	SortMapKeys:            true,
	ValidateJsonRawMessage: true,
	UseNumber:              true,
}.Froze()

//marshal
var jsonFastest = jsoniter.Config{
	EscapeHTML:                    false,
	MarshalFloatWith6Digits:       true,
	ObjectFieldMustBeSimpleString: true,
	UseNumber:                     true,
}.Froze()
