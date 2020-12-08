package common

import (
	"fmt"
	"github.com/saintfish/chardet"
	"io/ioutil"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestEncConvert_map(t *testing.T) {

	Convey("test enc map", t, func() {
		_, err := ConvToUTF8("GB-18030", []byte("abc"))
		So(
			err,
			ShouldBeNil,
		)
	})
}
func TestEncConvert_convert(t *testing.T) {

	Convey("test topsec encoding convert", t, func() {
		b, err := ioutil.ReadFile("/Users/work/go/src/github.com/sky-cloud-tec/netd/x.txt")
		if err == nil {
			fmt.Println(err)
			return
		}
		So(
			err,
			ShouldBeNil,
		)

		d := chardet.NewTextDetector()
		dr, err := d.DetectBest(b)
		So(
			err,
			ShouldBeNil,
		)
		fmt.Println(dr)
		So(
			dr.Charset == "ISO-8859-1",
			ShouldBeTrue,
		)
		nb, err := ConvToUTF8(dr.Charset, b)
		So(
			err,
			ShouldBeNil,
		)
		dr, err = d.DetectBest(nb)
		So(
			err,
			ShouldBeNil,
		)
		So(
			dr.Charset == "UTF-8",
			ShouldBeTrue,
		)

	})
}
