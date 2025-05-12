package dbtest

import (
	"fmt"
	"io"
	"os"

	"gorm.io/gorm"
)

func MigrateFromFile(db *gorm.DB, fileNames ...string) error {
	for _, fileName := range fileNames {
		fh, err := os.Open(fileName)
		if err != nil {
			return fmt.Errorf("os.Open: %w", err)
		}

		fileBytes, err := io.ReadAll(fh)
		if err != nil {
			return fmt.Errorf("io.ReadAll: %w", err)
		}

		if err = fh.Close(); err != nil {
			return fmt.Errorf("fh.Close: %w", err)
		}

		if err = db.Exec(string(fileBytes)).Error; err != nil {
			return fmt.Errorf("db.Exec: %w", err)
		}
	}

	return nil
}
