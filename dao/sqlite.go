package dao

import (
	"os"
	"path/filepath"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	SqliteDB *gorm.DB
)

func InitSqlite() (err error) {
	absDir := func() string {
		absPath, err := filepath.Abs(os.Args[0])
		if err != nil {
			panic(err)
		}
		return filepath.Dir(absPath)
	}()

	sqliteFile := filepath.Join(absDir, "sqlite.db")
	SqliteDB, err = gorm.Open(sqlite.Open(sqliteFile), &gorm.Config{})
	return err
}
