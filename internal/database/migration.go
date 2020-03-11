package database

import (
	"fmt"
	migration "github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/mysql"
	_ "github.com/golang-migrate/migrate/source/file"
	log "github.com/sirupsen/logrus"
)

func (d *Database) Migrate(migrationDirectory string) error {
	driver, err := mysql.WithInstance(d.DB, &mysql.Config{})
	if err != nil {
		return err
	}

	migrate, err := migration.NewWithDatabaseInstance(fmt.Sprintf("file://%s", migrationDirectory), "mysql", driver)
	if err != nil {
		return err
	}

	if err := migrate.Up(); err != nil {
		if err == migration.ErrNoChange {
			log.Info("no new migration")
			return nil
		}
		return err
	}

	return nil
}