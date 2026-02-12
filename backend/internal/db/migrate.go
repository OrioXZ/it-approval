package db

import (
	"fmt"
	"os"
	"path/filepath"

	"gorm.io/gorm"
)

func execSQLFile(gdb *gorm.DB, path string) error {
	b, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	if err := gdb.Exec(string(b)).Error; err != nil {
		return fmt.Errorf("exec %s failed: %w", path, err)
	}
	return nil
}

func MigrateAndSeed(gdb *gorm.DB, migrationsDir string) error {
	files := []string{
		filepath.Join(migrationsDir, "001_init.sql"),
		filepath.Join(migrationsDir, "002_seed.sql"),
	}
	for _, f := range files {
		if err := execSQLFile(gdb, f); err != nil {
			return err
		}
	}
	return nil
}
