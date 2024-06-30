package byter_test

import (
	"testing"
	"time"

	"github.com/ariefdarmawan/byter"
	"github.com/sebarcode/codekit"
	"github.com/smartystreets/goconvey/convey"
)

func TestByter(t *testing.T) {
	var (
		bs []byte
		e  error
	)
	b := byter.NewByter("")

	convey.Convey("String", t, func() {
		convey.Convey("encode", func() {
			str := "test data"
			bs, e = b.Encode(str)
			convey.So(e, convey.ShouldBeNil)
			convey.Convey("decode", func() {
				ret, e := b.Decode(bs, "", nil)
				convey.So(e, convey.ShouldBeNil)
				convey.So(ret, convey.ShouldEqual, str)
			})
		})
	})

	convey.Convey("Int", t, func() {
		convey.Convey("encode", func() {
			data := codekit.RandInt(1000)
			bs, e = b.Encode(data)
			convey.So(e, convey.ShouldBeNil)
			convey.Convey("decode", func() {
				var dest int
				e := b.DecodeTo(bs, &dest, nil)
				convey.So(e, convey.ShouldBeNil)
				convey.So(dest, convey.ShouldEqual, data)
			})
		})
	})

	convey.Convey("Float", t, func() {
		convey.Convey("encode", func() {
			data := float32(100.20)
			bs, e = b.Encode(data)
			convey.So(e, convey.ShouldBeNil)
			convey.Convey("decode", func() {
				dest, e := b.Decode(bs, float32(0), nil)
				convey.So(e, convey.ShouldBeNil)
				convey.So(dest.(float32), convey.ShouldEqual, data)
			})
		})
	})

	convey.Convey("Date DecodeTo with pointer result", t, func() {
		convey.Convey("encode", func() {
			data := time.Now()
			bs, e = b.Encode(data)
			convey.So(e, convey.ShouldBeNil)
			convey.Convey("decode", func() {
				var dest time.Time
				e := b.DecodeTo(bs, &dest, nil)
				//convey.Printf("\nOriginal: %v Result: %v\n", data, dest)
				convey.So(e, convey.ShouldBeNil)
				convey.So(dest.Unix(), convey.ShouldEqual, data.Unix())
			})
		})
	})

	convey.Convey("Date Decode without pointer result", t, func() {
		convey.Convey("encode", func() {
			data := time.Now()
			bs, e = b.Encode(data)
			convey.So(e, convey.ShouldBeNil)
			convey.Convey("decode", func() {
				dest, e := b.Decode(bs, time.Time{}, nil)
				//convey.Printf("\nOriginal: %v Result: %v\n", data, dest)
				convey.So(e, convey.ShouldBeNil)
				convey.So(dest.(time.Time).Unix(), convey.ShouldEqual, data.Unix())
			})
		})
	})

	convey.Convey("Date Decode with pointer result", t, func() {
		convey.Convey("encode", func() {
			data := time.Now()
			bs, e = b.Encode(data)
			convey.So(e, convey.ShouldBeNil)
			convey.Convey("decode", func() {
				dest, e := b.Decode(bs, &time.Time{}, nil)
				//convey.Printf("\nOriginal: %v Result: %v\n", data, dest)
				convey.So(e, convey.ShouldBeNil)
				convey.So(dest.(*time.Time).Unix(), convey.ShouldEqual, data.Unix())
			})
		})
	})
}

func TestCast(t *testing.T) {
	type tmpType struct {
		Name       string
		ValueInt   int
		ValueFloat float64
		Date       *time.Time
	}

	convey.Convey("cast", t, func() {
		var dest *tmpType

		b := byter.NewByter("")
		source := codekit.M{}.Set("Name", "Name 1").Set("ValueFloat", 120.05).Set("Date", time.Now())

		convey.Convey("cast to nil", func() {
			err := byter.Cast(b, source, dest, nil)
			convey.So(err, convey.ShouldNotBeNil)
			convey.Printf(" %s", err.Error())
		})

		convey.Convey("cast to not nil", func() {
			dest = new(tmpType)
			err := byter.Cast(b, source, dest, nil)
			convey.So(err, convey.ShouldBeNil)
			convey.Printf(" %s", codekit.JsonString(dest))
		})
	})
}
