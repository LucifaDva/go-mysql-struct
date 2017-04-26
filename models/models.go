package models

import (
	"github.com/LucifaDva/go-mysql-struct/db"
)

type DBModel struct {
	DbPool *db.DBPool
}

func (m *DBModel) Init(pool *db.DBPool) {
	m.DbPool = pool
}
