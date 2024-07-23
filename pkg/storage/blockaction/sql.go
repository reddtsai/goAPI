package blockaction

import (
	"context"
	dbsql "database/sql"
	"time"

	_ "github.com/go-sql-driver/mysql"
	gormmysqldriver "gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func ConnMySQL(ctx context.Context, dsn string, maxOpenConn, maxIdleConn, maxConnLifetimeMinutes int) (*dbsql.DB, error) {
	db, err := dbsql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	db.SetConnMaxLifetime(time.Duration(maxConnLifetimeMinutes) * time.Minute)
	db.SetMaxOpenConns(maxOpenConn)
	db.SetMaxIdleConns(maxIdleConn)

	return db, nil
}

func ConnGormMySQL(ctx context.Context, dsn string, maxOpenConn, maxIdleConn, maxConnLifetimeMinutes int) (*gorm.DB, error) {
	db, err := ConnMySQL(ctx, dsn, maxOpenConn, maxIdleConn, maxConnLifetimeMinutes)
	if err != nil {
		return nil, err
	}

	cfg := gormmysqldriver.Config{
		Conn: db,
	}
	conn, err := gorm.Open(gormmysqldriver.New(cfg), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return conn, nil
}
