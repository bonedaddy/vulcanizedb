// VulcanizeDB
// Copyright © 2019 Vulcanize

// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.

// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.

// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package mocks

import (
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/event"
	"github.com/makerdao/vulcanizedb/pkg/core"
	"github.com/makerdao/vulcanizedb/pkg/datastore/postgres"
)

type MockConverter struct {
	ToModelsError           error
	ContractAbi             string
	LogsToConvert           []core.HeaderSyncLog
	PassedContractAddresses []string
	SetDBCalled             bool
	ToModelsCalledCounter   int
}

func (converter *MockConverter) ToModels(abi string, logs []core.HeaderSyncLog) ([]event.InsertionModel, error) {
	converter.LogsToConvert = logs
	converter.ContractAbi = abi
	converter.ToModelsCalledCounter = converter.ToModelsCalledCounter + 1
	return nil, converter.ToModelsError
}

func (converter *MockConverter) SetDB(db *postgres.DB) {
	converter.SetDBCalled = true
}
