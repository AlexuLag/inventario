package repository

import (
	"database/sql"
	"fmt"
	"time"
)

type MySQLBaseRepository struct {
	db *sql.DB
}

func NewMySQLBaseRepository(db *sql.DB) *MySQLBaseRepository {
	return &MySQLBaseRepository{
		db: db,
	}
}

func (r *MySQLBaseRepository) GetDB() *sql.DB {
	return r.db
}

func (r *MySQLBaseRepository) BeginTx() (*sql.Tx, error) {
	return r.db.Begin()
}

func (r *MySQLBaseRepository) CommitTx(tx *sql.Tx) error {
	return tx.Commit()
}

func (r *MySQLBaseRepository) RollbackTx(tx *sql.Tx) error {
	return tx.Rollback()
}

func (r *MySQLBaseRepository) GetCurrentTimestamp() time.Time {
	return time.Now().UTC()
}

func (r *MySQLBaseRepository) FormatTimestamp(t time.Time) string {
	return t.Format("2006-01-02 15:04:05")
}

func (r *MySQLBaseRepository) ParseTimestamp(timestamp string) (time.Time, error) {
	return time.Parse("2006-01-02 15:04:05", timestamp)
}

func (r *MySQLBaseRepository) Close() error {
	return r.db.Close()
}

func (r *MySQLBaseRepository) Ping() error {
	return r.db.Ping()
}

func (r *MySQLBaseRepository) GetLastInsertID(result sql.Result) (int64, error) {
	return result.LastInsertId()
}

func (r *MySQLBaseRepository) GetRowsAffected(result sql.Result) (int64, error) {
	return result.RowsAffected()
}

func (r *MySQLBaseRepository) IsNotFound(err error) bool {
	return err == sql.ErrNoRows
}

func (r *MySQLBaseRepository) IsDuplicateEntry(err error) bool {
	if err == nil {
		return false
	}
	return fmt.Sprintf("%v", err) == "Error 1062: Duplicate entry"
}
