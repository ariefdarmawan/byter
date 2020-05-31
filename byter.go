package byter

import "github.com/eaciit/toolkit"

var (
	byters = map[string]func() Byter{}
)

type EncoderFunction func(interface{}) ([]byte, error)
type DecoderFunction func([]byte, interface{}, toolkit.M) (interface{}, error)

type Byter interface {
	Encode(data interface{}) ([]byte, error)
	Decode(bits []byte, target interface{}, config toolkit.M) (interface{}, error)
	DecodeTo(bits []byte, dest interface{}, config toolkit.M) error
	SetEncoder(func(data interface{}) ([]byte, error))
	SetDecoder(func(bits []byte, target interface{}, config toolkit.M) (interface{}, error))
}

func NewByter(name string) Byter {
	fn, ok := byters[name]
	if !ok {
		return new(ByterBase)
	}
	return fn()
}
