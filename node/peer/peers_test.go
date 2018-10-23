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
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with Alvalor.  If not, see <http://www.gnu.org/licenses/>.

package peer

import (
	"github.com/alvalor/alvalor-go/types"
	"github.com/stretchr/testify/mock"
)

type PeersMock struct {
	mock.Mock
}

func (p *PeersMock) Active(address string) {
	p.Called(address)
}

func (p *PeersMock) Inactive(address string) {
	p.Called(address)
}

func (p *PeersMock) Actives() []string {
	args := p.Called()
	var active []string
	if args.Get(0) != nil {
		active = args.Get(0).([]string)
	}
	return active
}

func (p *PeersMock) Tag(address string, hash types.Hash) {
	p.Called(address, hash)
}

func (p *PeersMock) Tags(hash types.Hash) []string {
	args := p.Called(hash)
	var seen []string
	if args.Get(0) != nil {
		seen = args.Get(0).([]string)
	}
	return seen
}

func (p *PeersMock) Pending(address string) (uint, error) {
	args := p.Called(address)
	return uint(args.Int(0)), args.Error(1)
}

func (p *PeersMock) Find(filters ...FilterFunc) []string {
	args := p.Called(filters)
	var peers []string
	if args.Get(0) != nil {
		peers = args.Get(0).([]string)
	}
	return peers
}

// func TestNewPeers(t *testing.T) {
// 	peers := newPeers()
// 	assert.NotNil(t, peers.peers)
// }
//
// func TestPeersActive(t *testing.T) {
// 	address1 := "192.0.2.100:1337"
// 	address2 := "192.0.2.200:1337"
// 	peers := &simplePeers{actives: make(map[string]bool)}
//
// 	peers.Active(address1)
// 	if assert.Len(t, peers.actives, 1) {
// 		assert.Contains(t, peers.actives, address1)
// 	}
//
// 	peers.Active(address1)
// 	assert.Len(t, peers.actives, 1)
//
// 	peers.Active(address2)
// 	if assert.Len(t, peers.actives, 2) {
// 		assert.Contains(t, peers.actives, address2)
// 	}
// }
//
// func TestPeersInactive(t *testing.T) {
// 	address1 := "192.0.2.100:1337"
// 	address2 := "192.0.2.200:1337"
// 	peers := &simplePeers{actives: make(map[string]bool)}
//
// 	peers.actives[address1] = true
// 	peers.actives[address2] = true
// 	peers.Inactive(address1)
// 	if assert.Len(t, peers.actives, 1) {
// 		assert.NotContains(t, peers.actives, address1)
// 	}
//
// 	peers.Inactive(address1)
// 	assert.Len(t, peers.actives, 1)
//
// 	peers.Inactive(address2)
// 	assert.Len(t, peers.actives, 0)
// }
//
// func TestPeersActives(t *testing.T) {
// 	address1 := "192.0.2.100:1337"
// 	address2 := "192.0.2.200:1337"
// 	peers := &simplePeers{actives: make(map[string]bool)}
//
// 	actives := peers.Actives()
// 	assert.Empty(t, actives)
//
// 	peers.actives[address1] = true
// 	peers.actives[address2] = true
// 	actives = peers.Actives()
// 	assert.ElementsMatch(t, []string{address1, address2}, actives)
// }
//
// func TestPeersTag(t *testing.T) {
// 	id1 := [32]byte{1, 2, 3, 4}
// 	id2 := [32]byte{5, 6, 7, 8}
// 	address1 := "192.0.2.100:1337"
// 	address2 := "192.0.2.200:1337"
// 	peers := &simplePeers{tags: make(map[types.Hash][]string)}
//
// 	peers.Tag(address1, id1)
// 	if assert.Len(t, peers.tags[id1], 1) {
// 		assert.Contains(t, peers.tags[id1], address1)
// 	}
//
// 	assert.Empty(t, peers.tags[id2])
//
// 	peers.Tag(address1, id2)
// 	if assert.Len(t, peers.tags[id2], 1) {
// 		assert.Contains(t, peers.tags[id2], address1)
// 	}
//
// 	peers.Tag(address2, id1)
// 	if assert.Len(t, peers.tags[id1], 2) {
// 		assert.Contains(t, peers.tags[id1], address2)
// 	}
// }
//
// func TestPeersTags(t *testing.T) {
// 	id := [32]byte{1, 2, 3, 4}
// 	address1 := "192.0.2.100:1337"
// 	address2 := "192.0.2.200:1337"
// 	peers := &simplePeers{tags: make(map[types.Hash][]string)}
//
// 	peers.tags[id] = []string{address1, address2}
// 	tags := peers.Tags(id)
// 	assert.ElementsMatch(t, []string{address1, address2}, tags)
// }
