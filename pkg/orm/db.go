package orm

import "gorm.io/gorm"

type DB struct {
	*gorm.DB
}

func (db *DB) Model(value interface{}) *DB {
	db.DB = db.DB.Model(value)
	return db
}

func (db *DB) WhereIf(condition bool, query interface{}, args ...interface{}) *DB {
	if condition {
		db.DB = db.DB.Where(query, args...)
	}
	return db
}
