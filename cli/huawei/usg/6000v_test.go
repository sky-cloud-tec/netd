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

package usg

import (
	"testing"

	"github.com/sky-cloud-tec/netd/cli"
	. "github.com/smartystreets/goconvey/convey"
)

func TestAsaOp(t *testing.T) {

	Convey("usg op", t, func() {
		op := createopUsg6000V()
		So(
			cli.AnyMatch(op.GetPrompts("login"), "HRP_M<BH-MgtFw-1>"),
			ShouldBeTrue,
		)
		So(
			cli.AnyMatch(op.GetPrompts("login"), "HRP_A<NF-MgtFW-1-(227.33-A-B02)>"),
			ShouldBeTrue,
		)
		So(
			cli.AnyMatch(op.GetPrompts("system_View"), "HRP_A[NF-MgtFW-1-(227.33-A-B02)]"),
			ShouldBeTrue,
		)
		So(
			cli.AnyMatch(op.GetPrompts("system_View"), `[~!@#$%^&*()+-={}[]'";:/\|.<>]`),
			ShouldBeTrue,
		)
	})
}
