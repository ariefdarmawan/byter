package byter

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"math"
	"reflect"

	"github.com/sebarcode/codekit"
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

	case int, int8, int16, int32, int64:
		bits := math.Float64bits(codekit.ToFloat64(data, 8, codekit.RoundingAuto))
		bs := make([]byte, 8)
		binary.LittleEndian.PutUint64(bs, bits)
		return bs, nil

	case float32, float64:
		bits := math.Float64bits(codekit.ToFloat64(data, 8, codekit.RoundingAuto))
		bs := make([]byte, 8)
		binary.LittleEndian.PutUint64(bs, bits)
		return bs, nil

	default:
		if data == nil {
			return []byte{}, nil
		}

		bs, e := json.Marshal(data)
		if e != nil {
			return nil, fmt.Errorf("error: %s", e.Error())
		}
		return bs, nil
	}
}

func (b *ByterBase) Decode(bits []byte, typeref interface{}, config codekit.M) (interface{}, error) {
	if b.decoder != nil {
		return b.decoder(bits, typeref, config)
	}

	//-- get indirect type and prepare a result object based on indirect type
	var res interface{}
	var resType reflect.Type
	targetIsPtr := false
	v := reflect.ValueOf(typeref)
	if v.Kind() == reflect.Ptr {
		targetIsPtr = true
		res = v.Elem().Interface()
		resType = v.Elem().Type()
	} else {
		res = typeref
		resType = v.Type()
	}

	//-- decode based on indirect type and store the result into res object
	switch res.(type) {
	case string:
		res = string(bits)

	case int, int8, int16, int32, int64, float32, float64:
		bits := binary.LittleEndian.Uint64(bits)
		f := math.Float64frombits(bits)

		switch res.(type) {
		case int, int8, int16, int32, int64:
			res = int(f)
		case float32:
			res = float32(f)
		case float64:
			res = f
		default:
			return 0, fmt.Errorf("invalid type")
		}

	default:
		if len(bits) == 0 {
			return nil, nil
		}

		var targetPtr interface{}
		targetPtr = reflect.New(resType).Interface()
		if err := codekit.FromBytes(bits, "json", targetPtr); err != nil {
			return nil, fmt.Errorf("unable to serialize return object, type: %s, error: %s", resType, err.Error())
		}
		if targetIsPtr {
			return targetPtr, nil
		}
		return reflect.ValueOf(targetPtr).Elem().Interface(), nil
	}

	// if target is pointer, we need to return pointer as well
	if targetIsPtr {
		vres := reflect.New(resType)
		vres.Elem().Set(reflect.ValueOf(res))
		return vres.Interface(), nil
	}
	return res, nil
}

func (b *ByterBase) DecodeTo(bits []byte, dest interface{}, config codekit.M) error {
	if config == nil {
		config = codekit.M{}
	}

	vdest := reflect.ValueOf(dest)
	if vdest.Kind() != reflect.Ptr {
		return fmt.Errorf("decode need a pointer reference as the target")
	}

	result, err := b.Decode(bits, dest, config)
	if err != nil {
		return err
	}

	vres := reflect.ValueOf(result)
	if vres.Kind() == reflect.Ptr {
		vdest.Elem().Set(vres.Elem())
	} else {
		vdest.Elem().Set(vres)
	}
	return nil
}

func (b *ByterBase) SetEncoder(encoder func(interface{}) ([]byte, error)) {
	b.encoder = encoder
}

func (b *ByterBase) SetDecoder(decoder func([]byte, interface{}, codekit.M) (interface{}, error)) {
	b.decoder = decoder
}
