// NetD makes network device operations easy.
// Copyright (C) 2019  sky-cloud.net
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

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
