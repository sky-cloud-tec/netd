package common

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestConnDial(t *testing.T) {

	Convey("test enc map", t, func() {
		_, err := ConvToUTF8("GB-18030", []byte("abc"))
		So(
			err,
			ShouldBeNil,
		)
	})
}
