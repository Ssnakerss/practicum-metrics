package storage

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"time"

	"github.com/Ssnakerss/practicum-metrics/internal/metric"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type DBStorage struct {
	dsn     string
	timeout time.Duration
	DB      *sql.DB
}

// TODO добавить создание таблицы
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
	//Открываем коннект
	dbs.DB, err = sql.Open("pgx", dbs.dsn)
	if err != nil {
		return err
	}
	//Проверям коннект
	if err = dbs.CheckStorage(); err != nil {
		dbs.Close()
		return fmt.Errorf("db connection failure: %w", err)
	}

	//Создаем таблицу в базе если ее еще нет
	sql := `CREATE TABLE IF NOT EXISTS public.metrics
		(
			name text COLLATE pg_catalog."default" NOT NULL,
			type text COLLATE pg_catalog."default" NOT NULL,
			gauge double precision,
			counter bigint,
			CONSTRAINT metrics_pkey PRIMARY KEY (name, type)
		)
		TABLESPACE pg_default;

		ALTER TABLE IF EXISTS public.metrics
			OWNER to postgres;
	`
	ctx, cancel := context.WithTimeout(context.Background(), dbs.timeout*time.Second)
	defer cancel()

	_, err = dbs.DB.ExecContext(ctx, sql)
	if err != nil {
		dbs.Close()
		return fmt.Errorf("db table creation failure: %w", err)
	}

	return nil
}

// Закрываем соединение с базой
func (dbs *DBStorage) Close() {
	dbs.DB.Close()
}

// Проверяем соединение
func (dbs *DBStorage) CheckStorage() error {
	ctx, cancel := context.WithTimeout(context.Background(), dbs.timeout*time.Second)
	defer cancel()

	if err := dbs.DB.PingContext(ctx); err != nil {
		return err
	}
	return nil
}

func (dbs *DBStorage) Write(m *metric.Metric) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbs.timeout*time.Second)
	defer cancel()

	sql := `insert into metrics 
	(name, type, gauge, counter)
	values
	($1,$2,$3,$4)
	on conflict(name,type) do update
		set gauge = excluded.gauge,
			counter = metrics.counter + excluded.counter`

	_, err := dbs.DB.ExecContext(ctx,
		sql,
		m.Name,
		m.Type,
		m.Gauge,
		m.Counter,
	)

	if err != nil {
		return fmt.Errorf("db write error: %W", err)
	}
	return nil
}

// Сохраняем в базу [] метрик
func (dbs *DBStorage) WriteAll(mm *([]metric.Metric)) (int, error) {
	cnt := 0
	ctx, cancel := context.WithTimeout(context.Background(), dbs.timeout*time.Second)
	defer cancel()

	sql := `insert into metrics 
	(name, type, gauge, counter)
	values
	($1,$2,$3,$4)
	on conflict(name,type) do update
	set gauge = excluded.gauge,
	counter = metrics.counter + excluded.counter`

	tx, err := dbs.DB.BeginTx(ctx, nil)
	if err != nil {
		return 0, fmt.Errorf("db tx open error: %w", err)
	}
	// defer tx.Rollback()

	for _, m := range *mm {
		_, err := tx.Exec(sql,
			m.Name,
			m.Type,
			m.Gauge,
			m.Counter,
		)
		if err != nil {
			tx.Rollback()
			return 0, fmt.Errorf("db insert error: %W", err)
		}
		cnt++
	}
	err = tx.Commit()
	if err != nil {
		return 0, fmt.Errorf("db tx commit error: %w", err)
	}

	return cnt, nil

}

func (dbs *DBStorage) Read(m *metric.Metric) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbs.timeout*time.Second)
	defer cancel()

	sql := `select 
		name
		, type
		, gauge 
		, counter
	from metrics 
	where  name = $1 and type = $2`
	row := dbs.DB.QueryRowContext(ctx, sql, m.Name, m.Type)

	if err := row.Scan(&m.Name, &m.Type, &m.Gauge, &m.Counter); err != nil {
		return fmt.Errorf("db read error: %w", err)
	}
	return nil
}

//Читаем из базы массив метрик

func (dbs *DBStorage) ReadAll(mm *([]metric.Metric)) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbs.timeout*time.Second)
	defer cancel()

	rows, err := dbs.DB.QueryContext(ctx, `select 
		name
		, type
		, gauge
		, counter 
	from metrics `)
	if err != nil {
		return 0, fmt.Errorf("db select error: %w", err)
	}
	defer rows.Close()

	cnt := 0

	for rows.Next() {
		m := metric.Metric{}
		if err := rows.Scan(&m.Name, &m.Type, &m.Gauge, &m.Counter); err != nil {
			return cnt, fmt.Errorf("db scan row error: %w", err)
		} else {
			*mm = append(*mm, m)
		}
		cnt++
	}
	//Проверяем на ошибки в процессе чтения
	err = rows.Err()
	if err != nil {
		return cnt, fmt.Errorf("error during reading data: %w", err)
	}
	return cnt, nil
}

func (dbs *DBStorage) Truncate() error {
	ctx, cancel := context.WithTimeout(context.Background(), dbs.timeout*time.Second)
	defer cancel()

	sql := `truncate table metrics`
	_, err := dbs.DB.ExecContext(ctx, sql)
	if err != nil {
		return fmt.Errorf("db truncate error: %w", err)
	}
	return nil
}
