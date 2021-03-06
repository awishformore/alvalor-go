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

package inventories

import (
	"github.com/pkg/errors"

	"github.com/alvalor/alvalor-go/types"
)

// Repo is a simple implementation of the inventory store.
type Repo struct {
	inventories map[types.Hash]*types.Inventory
}

// NewRepo creates a new store for block inventories.
func NewRepo() *Repo {
	repo := &Repo{
		inventories: make(map[types.Hash]*types.Inventory),
	}
	return repo
}

// Add stores a new inventory.
func (repo *Repo) Add(inv *types.Inventory) error {
	_, ok := repo.inventories[inv.Hash]
	if ok {
		return errors.Wrap(ErrExist, "inventory already exists")
	}
	repo.inventories[inv.Hash] = inv
	return nil
}

// Has checks if a given inventory is known.
func (repo *Repo) Has(hash types.Hash) bool {
	_, ok := repo.inventories[hash]
	return ok
}

// Get retrieves the inventory with the given block hash.
func (repo *Repo) Get(hash types.Hash) (*types.Inventory, error) {
	inv, ok := repo.inventories[hash]
	if !ok {
		return nil, errors.Wrap(ErrNotExist, "inventory does not exist")
	}
	return inv, nil
}
