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
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
)

type Repository interface {
	Create(headerID int64, models []CatFilePitVowModel) error
	MarkHeaderChecked(headerID int64) error
	MissingHeaders(startingBlockNumber, endingBlockNumber int64) ([]core.Header, error)
}

type CatFilePitVowRepository struct {
	db *postgres.DB
}

func NewCatFilePitVowRepository(db *postgres.DB) CatFilePitVowRepository {
	return CatFilePitVowRepository{db: db}
}

func (repository CatFilePitVowRepository) Create(headerID int64, models []CatFilePitVowModel) error {
	tx, err := repository.db.Begin()
	if err != nil {
		return err
	}
	for _, model := range models {
		_, err = repository.db.Exec(
			`INSERT into maker.cat_file_pit_vow (header_id, what, data, tx_idx, raw_log)
        VALUES($1, $2, $3, $4, $5)`,
			headerID, model.What, model.Data, model.TransactionIndex, model.Raw,
		)
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	_, err = tx.Exec(`INSERT INTO public.checked_headers (header_id, cat_file_pit_vow_checked)
		VALUES ($1, $2) 
	ON CONFLICT (header_id) DO
		UPDATE SET cat_file_pit_vow_checked = $2`, headerID, true)
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func (repository CatFilePitVowRepository) MarkHeaderChecked(headerID int64) error {
	_, err := repository.db.Exec(`INSERT INTO public.checked_headers (header_id, cat_file_pit_vow_checked)
		VALUES ($1, $2) 
	ON CONFLICT (header_id) DO
		UPDATE SET cat_file_pit_vow_checked = $2`, headerID, true)
	return err
}

func (repository CatFilePitVowRepository) MissingHeaders(startingBlockNumber, endingBlockNumber int64) ([]core.Header, error) {
	var result []core.Header
	err := repository.db.Select(
		&result,
		`SELECT headers.id, headers.block_number FROM headers
               LEFT JOIN checked_headers on headers.id = header_id
               WHERE (header_id ISNULL OR cat_file_pit_vow_checked IS FALSE)
               AND headers.block_number >= $1
               AND headers.block_number <= $2
               AND headers.eth_node_fingerprint = $3`,
		startingBlockNumber,
		endingBlockNumber,
		repository.db.Node.ID,
	)
	return result, err
}