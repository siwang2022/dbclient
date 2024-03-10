package main

import (
	"context"
	"database/sql"
	"errors"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
	"log"
	"time"
)

type DbClient struct {
	db    *sqlx.DB
	state int8 // 1 running, 0 stop
	url   string
}

func ConnectDB(url string) *DbClient {
	dbClient := DbClient{
		url:   url,
		db:    nil,
		state: 1,
	}

	dbClient.refreshSqlxDB()

	return &dbClient
}

func (db *DbClient) refreshSqlxDB() error {
	connConfig, err := pgx.ParseConfig(db.url)
	if err != nil {
		return err
	}

	connector := stdlib.GetConnector(*connConfig)

	sqlDB := sql.OpenDB(connector)
	sqlxDB := sqlx.NewDb(sqlDB, "pgx")

	db.db = sqlxDB

	return nil
}

func (db *DbClient) SwitchDB() error {
	log.Println("switch started")

	// The url could change.
	// In this test, as I only have one db instance so it does not change
	db.state = 0

	t1 := time.Now()
	db.db.Close()
	log.Println("Actual closing db spent", time.Now().Sub(t1).String())

	time.Sleep(3 * time.Second)

	db.refreshSqlxDB()

	db.state = 1
	log.Println("switch done")

	return nil
}

func (db DbClient) checkState() bool {
	return db.state == 1
}

func (db DbClient) MustBegin() *sqlx.Tx {
	if !db.checkState() {
		panic("Tx begin failed due to dbclient is in stop mode")
	}

	return db.db.MustBegin()
}

func (db DbClient) BeginTxx(ctx context.Context, opts *sql.TxOptions) (*sqlx.Tx, error) {
	if !db.checkState() {
		return nil, errors.New("Tx begin failed due to dbclient in stop mode")
	}

	return db.db.BeginTxx(ctx, opts)
}

func (db DbClient) QueryxContext(ctx context.Context, query string, args ...interface{}) (*sqlx.Rows, error) {
	if !db.checkState() {
		return nil, errors.New("QueryxContext failed due to dbclient in stop mode")
	}

	return db.db.QueryxContext(ctx, query, args...)
}
