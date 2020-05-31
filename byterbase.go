package byter

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"math"
	"reflect"

	"github.com/eaciit/toolkit"
)

type ByterBase struct {
	encoder EncoderFunction
	decoder DecoderFunction
}

const (
	KeyReferenceObj string = "HttpReferenceObj"
)

func (b *ByterBase) Encode(data interface{}) ([]byte, error) {
	if b.encoder != nil {
		return b.encoder(data)
	}

	switch data.(type) {
	case string, *string:
		return []byte(data.(string)), nil

	case int, int8, int16, int32, int64, float32, float64:
		bits := math.Float64bits(toolkit.ToFloat64(data, 12, toolkit.RoundingAuto))
		bs := make([]byte, 8)
		binary.LittleEndian.PutUint64(bs, bits)
		return bs, nil

	default:
		bs, e := json.Marshal(data)
		if e != nil {
			return nil, fmt.Errorf("error: %s", e.Error())
		}
		return bs, nil
	}
}

func (b *ByterBase) Decode(bits []byte, target interface{}, config toolkit.M) (interface{}, error) {
	if b.decoder != nil {
		return b.decoder(bits, target, config)
	}

	var dest interface{}
	targetIsPtr := false
	v := reflect.ValueOf(target)
	if v.Kind() == reflect.Ptr {
		targetIsPtr = true
		dest = v.Elem().Interface()
	} else {
		dest = target
	}
	switch dest.(type) {
	case string:
		return string(bits), nil

	case int, int8, int16, int32, int64, float32, float64:
		bits := binary.LittleEndian.Uint64(bits)
		f := math.Float64frombits(bits)

		switch dest.(type) {
		case int, int8, int16, int32, int64:
			return int(f), nil
		case float32:
			return float32(f), nil
		case float64:
			return f, nil
		default:
			return 0, fmt.Errorf("invalid type")
		}

	default:
		var targetPtr interface{}
		if !targetIsPtr {
			targetPtr = reflect.New(reflect.TypeOf(target)).Interface()
		} else {
			targetPtr = target
		}
		if err := toolkit.FromBytes(bits, "json", targetPtr); err != nil {
			return nil, fmt.Errorf("unable to serialize return object. %s", err.Error())
		}
		if targetIsPtr {
			return targetPtr, nil
		}
		return reflect.ValueOf(targetPtr).Elem().Interface(), nil
	}
}

func (b *ByterBase) DecodeTo(bits []byte, dest interface{}, config toolkit.M) error {
	if config == nil {
		config = toolkit.M{}
	}
	config.Set(KeyReferenceObj, dest)
	result, err := b.Decode(bits, dest, config)

	vdest := reflect.ValueOf(dest)
	vres := reflect.ValueOf(result)
	if vdest.Kind() == reflect.Ptr {
		if vres.Kind() == reflect.Ptr {
			vdest.Elem().Set(vres.Elem())
		} else {
			vdest.Elem().Set(vres)
		}
	} else {
		if vres.Kind() == reflect.Ptr {
			vdest.Set(vres.Elem())
		} else {
			vdest.Set(vres)
		}
	}
	return err
}

func (b *ByterBase) SetEncoder(encoder func(interface{}) ([]byte, error)) {
	b.encoder = encoder
}

func (b *ByterBase) SetDecoder(decoder func([]byte, interface{}, toolkit.M) (interface{}, error)) {
	b.decoder = decoder
}
