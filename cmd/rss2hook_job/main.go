package main

import (
	"errors"
	"flag"
	"log"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"

	"r2w"
)

var (
	dbFile       = flag.String("db", "test.db", "Database file path")
	migrationDir = flag.String("migrate", "/Users/seven/dev/rss2webhook/migrations", "Migration directory path")
)

func main() {
	flag.Parse()
	db, err := sqlx.Open("sqlite3", *dbFile)
	if err != nil {
		log.Fatalf("err opening db:%v\n%s", err, debug.Stack())
		return
	}
	err = startMigrate(db, *migrationDir)
	if err != nil {
		log.Fatalf("err startMigrate:%v\n%s", err, debug.Stack())
		return
	}
	err = startJob(db)
	if err != nil {
		log.Fatalf("err startJob:%v\n%s", err, debug.Stack())
		return
	}
	err = db.Close()
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

func startMigrate(db *sqlx.DB, migrationDir string) error {
	driver, err := sqlite3.WithInstance(db.DB, &sqlite3.Config{})
	if err != nil {
		log.Fatal(err)
		return err
	}
	m, err := migrate.NewWithDatabaseInstance(
		"file://"+migrationDir,
		"sqlite3",
		driver)
	if err != nil {
		log.Fatal(err)
		return err
	}
	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Fatal(err)
		return err
	}
	return nil
}
