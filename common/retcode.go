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

const (
	// OK ..
	OK = 0
	// [1, 1000] preserved range

	// [1001, 2000] for cli handler and cli module

	// ErrNoOpFound no operator match
	ErrNoOpFound = 1001
	// ErrAcquireConn acquire cli conn error
	ErrAcquireConn = 1002
	// ErrCliExec execute cli command error
	ErrCliExec = 1003
	// ErrTimeout timeout error
	ErrTimeout = 1005
	// ErrNoMode mode not specified error
	ErrNoMode = 1006
)
