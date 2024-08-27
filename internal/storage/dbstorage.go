package storage

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type DBStorage struct {
	dsn     string
	timeout time.Duration
}

func (dbs *DBStorage) New(p ...string) error {
	if len(p) < 2 {
		return fmt.Errorf("specify dsn and connection timeout")
	}

	dbs.dsn = p[0]
	var err error
	var i int
	if i, err = strconv.Atoi(p[1]); err != nil {
		return fmt.Errorf("timeout value %s convertion error %w ", p[1], err)
	}
	dbs.timeout = time.Duration(i)

	// db, err := sql.Open("pgx", dbs.dsn)
	// if err != nil {
	// 	return err
	// }
	// defer db.Close()

	// ctx, cancel := context.WithTimeout(context.Background(), dbs.timeout*time.Second)
	// defer cancel()

	// if err = db.PingContext(ctx); err != nil {
	// 	return err
	// }
	return nil
}

func (dbs *DBStorage) CheckStorage() error {
	db, err := sql.Open("pgx", dbs.dsn)
	if err != nil {
		return err
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), dbs.timeout*time.Second)
	defer cancel()

	if err = db.PingContext(ctx); err != nil {
		return err
	}
	return nil
}
