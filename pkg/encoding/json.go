package encoding

import (
	"bytes"
	"encoding/json"
)

func init() {
	RegisterCodec("json", new(JsonCodec))
}

// JsonCodec uses json to encode and decode messages.
// It is the default codec.
//
// You can use it as a template for your own codecs.
// For example, you can use a different json implementation.
//
//	import (
	//		"github.com/aetherbus/aetherbus/core/encoding"
	//		jsoniter "github.com/json-iterator/go"
	//	)
	//
	//	func init() {
	//		encoding.RegisterCodec("json", new(JsonCodec))
	//	}
	//
	//	type JsonCodec struct{}
	//
	//	func (c *JsonCodec) Encode(v interface{}) ([]byte, error) {
	//		return jsoniter.Marshal(v)
	//	}
	//
	//	func (c *JsonCodec) Decode(data []byte, v interface{}) error {
	//		return jsoniter.Unmarshal(data, v)
	//	}

type JsonCodec struct{}

func (c *JsonCodec) Encode(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

func (c *JsonCodec) Decode(data []byte, v interface{}) error {
	// If v is a nil pointer, Decode returns an InvalidUnmarshalError.
	// We need to pass in a non-nil pointer to json.Unmarshal, because
	// otherwise it will quote the input and return an error.
	if v == nil {
		return nil
	}
	d := json.NewDecoder(bytes.NewReader(data))
	// d.UseNumber()
	return d.Decode(v)
}
