package encoding

import (
	"fmt"
)

var (
	ErrCodecNotRegistered = func(name string) error {
		return fmt.Errorf("codec '%s' is not registered", name)
	}
	ErrCannotEncodeNil = fmt.Errorf("cannot encode nil")
)

type Codec interface {
	Encode(interface{}) ([]byte, error)
	Decode([]byte, interface{}) error
}

var registeredCodecs = make(map[string]Codec)

func RegisterCodec(name string, codec Codec) {
	registeredCodecs[name] = codec
}

func GetCodec(name string) (Codec, error) {
	codec, ok := registeredCodecs[name]
	if !ok {
		return nil, ErrCodecNotRegistered(name)
	}
	return codec, nil
}
