// Copyright 2018 Vulcanize
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package tend

import (
	"fmt"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/shared"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/shared/constants"
)

type TendRepository struct {
	db *postgres.DB
}

func (repository TendRepository) Create(headerID int64, models []interface{}) error {
	tx, err := repository.db.Begin()
	if err != nil {
		return err
	}

	tic, err := shared.GetTicInTx(headerID, tx)
	if err != nil {
		return err
	}

	for _, model := range models {
		tend, ok := model.(TendModel)
		if !ok {
			tx.Rollback()
			return fmt.Errorf("model of type %T, not %T", model, TendModel{})
		}

		_, err = tx.Exec(
			`INSERT into maker.tend (header_id, bid_id, lot, bid, guy, tic, log_idx, tx_idx, raw_log)
			VALUES($1, $2, $3::NUMERIC, $4::NUMERIC, $5, $6, $7, $8, $9)`,
			headerID, tend.BidId, tend.Lot, tend.Bid, tend.Guy, tic, tend.LogIndex, tend.TransactionIndex, tend.Raw,
		)

		if err != nil {
			tx.Rollback()
			return err
		}
	}

	err = shared.MarkHeaderCheckedInTransaction(headerID, tx, constants.TendChecked)
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func (repository TendRepository) MarkHeaderChecked(headerId int64) error {
	return shared.MarkHeaderChecked(headerId, repository.db, constants.TendChecked)
}

func (repository *TendRepository) SetDB(db *postgres.DB) {
	repository.db = db
}