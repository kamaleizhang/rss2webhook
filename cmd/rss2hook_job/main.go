package main

import (
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"

	"r2w"
)

func main() {
	db, err := startMigrate()
	if err != nil {
		panic(err)
	}
	err = startJob(db)
	if err != nil {
		panic(err)
	}
}

func startJob(db *sqlx.DB) error {
	configManager := r2w.NewConfigManager(nil, "source-config.json")
	configs, err := configManager.GetConfigs()
	if err != nil {
		return err
	}
	httpClient := &http.Client{
		Timeout: 10 * time.Second,
	}
	updater := r2w.NewRssUpdaterImpl(db)
	for _, config := range configs {
		requests, err := updater.SyncRss(config)
		if err != nil {
			return err
		}
		caller := r2w.NewHookCallerImpl(httpClient)
		for _, request := range requests {
			err = caller.CallHooks(config, request)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func startMigrate() (*sqlx.DB, error) {
	db, err := sqlx.Open("sqlite3", "test.db")
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	defer db.Close()
	driver, err := sqlite3.WithInstance(db.DB, &sqlite3.Config{})
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	m, err := migrate.NewWithDatabaseInstance(
		"file:///migrations",
		"sqlite3",
		driver)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Fatal(err)
		return nil, err
	}
	return db, nil
}
