package mysql

import (
	"github.com/jmoiron/sqlx"
)

func RunInTx(db *sqlx.DB, txFn func(tx *sqlx.Tx) error) error {
	// トランザクション開始
	tx, err := db.Beginx()
	if err != nil {
		return err
	}
	// トランザクション実行
	err = txFn(tx)
	if err != nil {
		// ロールバック
		if err := tx.Rollback(); err != nil {
			return err
		}
		return err
	}
	// コミット
	return tx.Commit()
}
