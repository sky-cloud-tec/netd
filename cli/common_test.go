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

package cli

import (
	"regexp"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestMoreOK(t *testing.T) {
	Convey("test more", t, func() {

		str1 := "--More--"
		So(
			IsSymmetricalMore(str1),
			ShouldBeTrue,
		)

		str2 := "<--- More --->"
		So(
			IsSymmetricalMore(str2),
			ShouldBeTrue,
		)
	})
}

func TestPromptMatch(t *testing.T) {
	Convey("test tsos login prompt", t, func() {
		So(
			Match(regexp.MustCompile("^[[:alnum:]]{1,}> $"), "\nUSG> "),
			ShouldBeTrue,
		)
	})
}
