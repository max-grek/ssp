package connector

import (
	"database/sql"
	"fmt"
	"log"
	"test-assignment-cookie-sync/config"
	"test-assignment-cookie-sync/connector/internal/cookie"
	"time"

	_ "github.com/lib/pq"
)

const q0 = "cookie"

type DB interface {
	Cookie() *cookie.QTX

	impl(q string) interface{}
}

type database map[string]interface{}

func new(conn *sql.DB) DB {
	return database{
		q0: cookie.New(conn),
	}
}

func (db database) Cookie() *cookie.QTX {
	v, ok := db.impl(q0).(*cookie.QTX)
	if ok {
		return v
	}
	return nil
}

func (db database) impl(q string) interface{} {
	if v, ok := db[q]; ok {
		return v
	}
	return nil
}

func Initialize(cfg *config.StorageConfig) (DB, error) {
	db, err := openConnection(cfg.Driver, resolveDSN(cfg))
	if err != nil {
		return nil, fmt.Errorf("open connection: %v", err)
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	log.Println("successfully established connection with database")
	return new(db), nil
}

func openConnection(driver, dsn string) (*sql.DB, error) {
	db, err := sql.Open(driver, dsn)
	if err != nil {
		return nil, fmt.Errorf("%v: %s, %s", err, driver, dsn)
	}

	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(1)
	db.SetConnMaxLifetime(time.Minute)
	return db, nil
}

func resolveDSN(cfg *config.StorageConfig) string {
	if cfg.IsDsnEnabled {
		return cfg.Dsn
	}
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s search_path=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Database, cfg.Sslmode, cfg.Schema)
}
