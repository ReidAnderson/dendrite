// Copyright 2020 The Matrix.org Foundation C.I.C.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package sqlite3

import (
	"context"
	"database/sql"

	"github.com/matrix-org/dendrite/internal/sqlutil"
	"github.com/matrix-org/dendrite/roomserver/storage/shared"
	"github.com/matrix-org/dendrite/roomserver/storage/tables"
)

const expirySchema = `
-- Stores what time an event will expire and should subsequently be redacted
CREATE TABLE IF NOT EXISTS roomserver_event_expiry (
    event_nid INTEGER PRIMARY KEY,
	expiry_ts INTEGER NOT NULL
);
`

const insertExpirySQL = "" +
	"INSERT OR IGNORE INTO roomserver_event_expiry (event_nid, expiry_ts)" +
	" VALUES ($1, $2)"

// const selectRedactionInfoByRedactionEventIDSQL = "" +
// 	"SELECT redaction_event_id, redacts_event_id, validated FROM roomserver_redactions" +
// 	" WHERE redaction_event_id = $1"

// const selectRedactionInfoByEventBeingRedactedSQL = "" +
// 	"SELECT redaction_event_id, redacts_event_id, validated FROM roomserver_redactions" +
// 	" WHERE redacts_event_id = $1"

// const markRedactionValidatedSQL = "" +
// 	" UPDATE roomserver_redactions SET validated = $2 WHERE redaction_event_id = $1"

type expiryStatements struct {
	db               *sql.DB
	insertExpiryStmt *sql.Stmt
	// selectRedactionInfoByRedactionEventIDStmt   *sql.Stmt
	// selectRedactionInfoByEventBeingRedactedStmt *sql.Stmt
	// markRedactionValidatedStmt                  *sql.Stmt
}

// NewSqliteExpiryTable New table
func NewSqliteExpiryTable(db *sql.DB) (tables.Expiry, error) {
	s := &expiryStatements{
		db: db,
	}
	_, err := db.Exec(expirySchema)
	if err != nil {
		return nil, err
	}

	return s, shared.StatementList{
		{&s.insertExpiryStmt, insertExpirySQL},
		// {&s.selectRedactionInfoByRedactionEventIDStmt, selectRedactionInfoByRedactionEventIDSQL},
		// {&s.selectRedactionInfoByEventBeingRedactedStmt, selectRedactionInfoByEventBeingRedactedSQL},
		// {&s.markRedactionValidatedStmt, markRedactionValidatedSQL},
	}.Prepare(db)
}

func (s *expiryStatements) InsertExpiry(
	ctx context.Context, txn *sql.Tx, info tables.ExpiryInfo,
) error {
	stmt := sqlutil.TxStmt(txn, s.insertExpiryStmt)
	_, err := stmt.ExecContext(ctx, info.EventNID, info.ExpiryTimestamp)
	return err
}

// func (s *redactionStatements) SelectRedactionInfoByRedactionEventID(
// 	ctx context.Context, txn *sql.Tx, redactionEventID string,
// ) (info *tables.RedactionInfo, err error) {
// 	info = &tables.RedactionInfo{}
// 	stmt := sqlutil.TxStmt(txn, s.selectRedactionInfoByRedactionEventIDStmt)
// 	err = stmt.QueryRowContext(ctx, redactionEventID).Scan(
// 		&info.RedactionEventID, &info.RedactsEventID, &info.Validated,
// 	)
// 	if err == sql.ErrNoRows {
// 		info = nil
// 		err = nil
// 	}
// 	return
// }

// func (s *redactionStatements) SelectRedactionInfoByEventBeingRedacted(
// 	ctx context.Context, txn *sql.Tx, eventID string,
// ) (info *tables.RedactionInfo, err error) {
// 	info = &tables.RedactionInfo{}
// 	stmt := sqlutil.TxStmt(txn, s.selectRedactionInfoByEventBeingRedactedStmt)
// 	err = stmt.QueryRowContext(ctx, eventID).Scan(
// 		&info.RedactionEventID, &info.RedactsEventID, &info.Validated,
// 	)
// 	if err == sql.ErrNoRows {
// 		info = nil
// 		err = nil
// 	}
// 	return
// }

// func (s *redactionStatements) MarkRedactionValidated(
// 	ctx context.Context, txn *sql.Tx, redactionEventID string, validated bool,
// ) error {
// 	stmt := sqlutil.TxStmt(txn, s.markRedactionValidatedStmt)
// 	_, err := stmt.ExecContext(ctx, redactionEventID, validated)
// 	return err
// }
