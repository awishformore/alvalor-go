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

package node

import (
	"errors"

	"github.com/alvalor/alvalor-go/types"
)

type pathManager interface {
	Add(header *types.Header) error
	Longest() []types.Hash
}

// simplePath is a path manager using topological sort of all added headers to find the longest path to the root.
type simplePath struct {
	root     types.Hash
	headers  map[types.Hash]*types.Header
	children map[types.Hash][]types.Hash
}

// newsimplePath creates a new simple path manager using the given header as root.
func newSimplePaths(root *types.Header) *simplePath {
	sp := &simplePath{
		root:     root.Hash,
		headers:  make(map[types.Hash]*types.Header),
		children: make(map[types.Hash][]types.Hash),
	}
	sp.headers[root.Hash] = root
	return sp
}

// Add adds a new header to the graph.
func (sp *simplePath) Add(header *types.Header) error {
	_, ok := sp.headers[header.Hash]
	if ok {
		return errors.New("header already in graph")
	}
	_, ok = sp.headers[header.Parent]
	if !ok {
		return errors.New("header parent not in graph")
	}
	sp.children[header.Parent] = append(sp.children[header.Parent], header.Hash)
	sp.headers[header.Hash] = header
	return nil
}

// Longest returns the longest path of the graph.
func (sp *simplePath) Longest() []types.Hash {

	// create a topological sort of all headers starting at the root
	var hash types.Hash
	sorted := make([]types.Hash, 0, len(sp.headers))
	queue := []types.Hash{sp.root}
	queue = append(queue, sp.root)
	for len(queue) > 0 {
		hash, queue = queue[0], queue[1:]
		sorted = append(sorted, hash)
		queue = append(queue, sp.children[hash]...)
	}

	// find the maximum distance of each header from the root
	var max uint64
	var best types.Hash
	distances := make(map[types.Hash]uint64)
	for len(sorted) > 0 {
		hash, sorted = sorted[0], sorted[1:]
		for _, child := range sp.children[hash] {
			header := sp.headers[child]
			distance := distances[hash] + header.Diff
			if distances[child] >= distance {
				continue
			}
			distances[child] = distance
			if distance <= max {
				continue
			}
			max = distance
			best = child
		}
	}

	// if we have no distance, we are stuck at the root
	if max == 0 {
		return []types.Hash{sp.root}
	}

	// otherwise, iterate back to parents from best child
	var path []types.Hash
	header := sp.headers[best]
	for header.Parent != types.ZeroHash {
		path = append(path, header.Hash)
		header = sp.headers[header.Parent]
	}

	return path
}