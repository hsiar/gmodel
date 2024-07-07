package corm

import "gorm.io/gorm"

// DBOpener opens a gorm Dialector.
//
// See:
//   - gorm.io/driver/mysql:    https://github.com/go-gorm/mysql/blob/f46a79cf94a9d67edcc7d5f6f2606e21bf6525fe/mysql.go#L52
//   - gorm.io/driver/postgres: https://github.com/go-gorm/postgres/blob/c2cfceb161687324cb399c9f60ec775428335957/postgres.go#L31
//   - gorm.io/driver/sqlite:   https://github.com/go-gorm/sqlite/blob/1d1e7723862758a6e6a860f90f3e7a3bea9cc94a/sqlite.go#L28
type DBOpener func(dsn string) gorm.Dialector

type DBDriver string
