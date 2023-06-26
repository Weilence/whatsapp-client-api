package model

import "gorm.io/gorm"

type Pagination struct {
	Page int `json:"page,omitempty" query:"page,omitempty"`
	Size int `json:"size,omitempty" query:"size,omitempty"`
}

func Paginate(p Pagination) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Limit(p.Size).
			Offset((p.Page - 1) * p.Size)
	}
}

func WhereIf(b bool, query interface{}, args ...interface{}) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if b {
			return db.Where(query, args)
		} else {
			return db
		}
	}
}
