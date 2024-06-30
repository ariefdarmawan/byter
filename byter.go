package byter

import (
	"fmt"

	"github.com/sebarcode/codekit"
)

var (
	byters = map[string]func() Byter{}
)

type EncoderFunction func(interface{}) ([]byte, error)
type DecoderFunction func([]byte, interface{}, codekit.M) (interface{}, error)

type Byter interface {
	Encode(data interface{}) ([]byte, error)
	Decode(bits []byte, target interface{}, config codekit.M) (interface{}, error)
	DecodeTo(bits []byte, dest interface{}, config codekit.M) error
	SetEncoder(func(data interface{}) ([]byte, error))
	SetDecoder(func(bits []byte, target interface{}, config codekit.M) (interface{}, error))
}

func NewByter(name string) Byter {
	fn, ok := byters[name]
	if !ok {
		return new(ByterBase)
	}
	return fn()
}

func Cast(b Byter, source interface{}, destTo interface{}, config codekit.M) error {
	if config == nil {
		config = codekit.M{}
	}

	bs, err := b.Encode(source)
	if err != nil {
		return fmt.Errorf("encode: %s", err.Error())
	}

	func() {
		defer func() {
			if r := recover(); r != nil {
				if err != nil {
					return
				}
				err = fmt.Errorf("panic: %v", r)
			}
		}()

		err = b.DecodeTo(bs, destTo, config)
		if err != nil {
			return
		}
	}()

	return err
}
