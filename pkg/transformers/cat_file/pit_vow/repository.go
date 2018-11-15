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

package pit_vow

import (
	"fmt"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/shared"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/shared/constants"
)

type CatFilePitVowRepository struct {
	db *postgres.DB
}

func (repository CatFilePitVowRepository) Create(headerID int64, models []interface{}) error {
	tx, err := repository.db.Begin()
	if err != nil {
		return err
	}
	for _, model := range models {
		vow, ok := model.(CatFilePitVowModel)
		if !ok {
			tx.Rollback()
			return fmt.Errorf("model of type %T, not %T", model, CatFilePitVowModel{})
		}

		err = shared.ValidateHeaderConsistency(headerID, vow.Raw, repository.db)
		if err != nil {
			tx.Rollback()
			return err
		}

		_, err = repository.db.Exec(
			`INSERT into maker.cat_file_pit_vow (header_id, what, data, tx_idx, log_idx, raw_log)
			VALUES($1, $2, $3, $4, $5, $6)`,
			headerID, vow.What, vow.Data, vow.TransactionIndex, vow.LogIndex, vow.Raw,
		)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	err = shared.MarkHeaderCheckedInTransaction(headerID, tx, constants.CatFilePitVowChecked)
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func (repository CatFilePitVowRepository) MarkHeaderChecked(headerID int64) error {
	return shared.MarkHeaderChecked(headerID, repository.db, constants.CatFilePitVowChecked)
}

func (repository CatFilePitVowRepository) MissingHeaders(startingBlockNumber, endingBlockNumber int64) ([]core.Header, error) {
	return shared.MissingHeaders(startingBlockNumber, endingBlockNumber, repository.db, constants.CatFilePitVowChecked)
}

func (repository *CatFilePitVowRepository) SetDB(db *postgres.DB) {
	repository.db = db
}