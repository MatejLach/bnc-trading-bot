package bot

import "database/sql"

// TODO
func initDb(db *sql.DB) error {
	stmt, err := db.Prepare(`
    
      CREATE TABLE IF NOT EXISTS trading_bot(
      crypto TEXT PRIMARY KEY, currency TEXT NOT NULL,
      original_price REAL NOT NULL, original_price_timestamp TEXT NOT NULL,
      current_price REAL NOT NULL, current_price_timestamp TEXT NOT NULL)`)

	if err != nil {
		return err
	}

	_, err = stmt.Exec()
	if err != nil {
		return err
	}

	return nil
}
