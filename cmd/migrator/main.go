package migrator

import (
	"errors"
	"flag"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
)

func main() {
	var storagePath, migratorPath, migratorTable string

	flag.StringVar(&storagePath, "storage-path", "", "path to storage")
	flag.StringVar(&migratorPath, "migrator-path", "", "path to migrator")
	flag.StringVar(&migratorTable, "migrator-table", "migrations", "migrator table")
	flag.Parse()

	if storagePath == "" {
		panic("storage-path is empty")
	}
	if migratorPath == "" {
		panic("migrator-path is empty")
	}
	m, err := migrate.New(
		"file://"+migratorPath,
		fmt.Sprintf("sqlite3://%s?x-migration-table=%s", storagePath, migratorTable),
	)
	if err != nil {
		panic(err)
	}
	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			fmt.Println("no changes")
			return
		}
		panic(err)
	}

	fmt.Println()
}
