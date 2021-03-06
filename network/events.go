// Copyright (c) 2017 The Alvalor Authors
//
// This file is part of Alvalor.
//
// Alvalor is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// Alvalor is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more detailb.
//
// You should have received a copy of the GNU Affero General Public License
// along with Alvalor.  If not, see <http://www.gnu.org/licenses/>.

package network

import "time"

// Connected represents a connection event.
type Connected struct {
	Address   string
	Timestamp time.Time
}

// Disconnected represents a disconnection event.
type Disconnected struct {
	Address   string
	Timestamp time.Time
}

// Received represents a message received event.
type Received struct {
	Address   string
	Timestamp time.Time
	Message   interface{}
}
