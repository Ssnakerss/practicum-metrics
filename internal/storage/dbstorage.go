package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/Ssnakerss/practicum-metrics/internal/metric"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type DBStorage struct {
	dsn     string
	timeout time.Duration
	DB      *sql.DB
	ictx    context.Context    //Внутренний контекст на основе родительского
	cancel  context.CancelFunc //Объект для отмены таймаута
}

// Выбираем тип результирующей ошибки ошибки для методов пакета DB Storage
func errSelect(ctx context.Context, method string, err error) error {
	if err != nil {
		var pgerr *pgconn.PgError
		if errors.As(err, &pgerr) {
			//Если ишибка соединения - формируем специальную ошибку - будем использовать при повторных попытках соединения
			if pgerrcode.IsConnectionException(pgerr.SQLState()) {
				return NewStorageError("db", method, ConnectionError, err)
			}
		}
		return NewStorageError("db", method, 99, err)
	}
	//Проверяем была ли операция отменена по таймауту
	if errors.Is(ctx.Err(), context.DeadlineExceeded) {
		return NewStorageError("db", method, 13, ctx.Err())
	}
	return nil
}

func (dbs *DBStorage) New(ctx context.Context, p ...string) error {
	if len(p) < 2 {
		return errSelect(context.TODO(), "init", fmt.Errorf("specify dsn and connection timeout"))
	}

	dbs.dsn = p[0]
	var err error
	var i int

	if i, err = strconv.Atoi(p[1]); err != nil {
		return errSelect(context.TODO(), "init", fmt.Errorf("timeout value %s convertion error %w ", p[1], err))
	}

	dbs.timeout = time.Duration(i) * time.Second
	//Cоздаем экземпляр внутреннего контекста на основе внешнего родительского
	dbs.ictx, dbs.cancel = context.WithCancel(ctx)
	//Открываем коннект
	dbs.DB, err = sql.Open("pgx", dbs.dsn)

	if err != nil {
		return errSelect(context.TODO(), "conn", fmt.Errorf("connection open error: %w", err))
	}

	//Проверям коннект
	if err = dbs.CheckStorage(); err != nil {
		dbs.Close()
		return errSelect(context.TODO(), "conn", fmt.Errorf("connection checking failure: %w", err))
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
		
	`
	ctx, cancel := context.WithTimeout(dbs.ictx, dbs.timeout*time.Second)
	defer cancel()

	_, err = dbs.DB.ExecContext(ctx, sql)
	if err != nil {
		dbs.Close()
	}

	return errSelect(ctx, "tbcreate", err)
}

// Закрываем соединение с базой
func (dbs *DBStorage) Close() {
	//отменяем внутренний контекс
	dbs.cancel()
	dbs.DB.Close()
}

// Проверяем соединение
func (dbs *DBStorage) CheckStorage() error {
	ctx, cancel := context.WithTimeout(dbs.ictx, dbs.timeout*time.Second)
	defer cancel()

	err := dbs.DB.PingContext(ctx)

	return errSelect(ctx, "ping", err)
}

func (dbs *DBStorage) Write(m *metric.Metric) error {
	ctx, cancel := context.WithTimeout(dbs.ictx, dbs.timeout*time.Second)
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

	return errSelect(ctx, "WRITE", err)
}

// Сохраняем в базу [] метрик
func (dbs *DBStorage) WriteAll(mm *([]metric.Metric)) (int, error) {
	cnt := 0
	ctx, cancel := context.WithTimeout(dbs.ictx, dbs.timeout*time.Second)
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
		return 0, errSelect(ctx, "begintx", fmt.Errorf("db tx open error: %w", err))
	}

	for _, m := range *mm {
		_, err = tx.Exec(sql,
			m.Name,
			m.Type,
			m.Gauge,
			m.Counter,
		)
		if err != nil {
			tx.Rollback()
			return 0, errSelect(ctx, "insert", fmt.Errorf("db insert error: %W", err))
		}
		cnt++
	}
	err = tx.Commit()

	return cnt, errSelect(ctx, "WRITEALL", err)
}

func (dbs *DBStorage) Read(m *metric.Metric) error {
	ctx, cancel := context.WithTimeout(dbs.ictx, dbs.timeout*time.Second)
	defer cancel()

	sql := `select 
		name
		, type
		, gauge 
		, counter
	from metrics 
	where  name = $1 and type = $2`
	row := dbs.DB.QueryRowContext(ctx, sql, m.Name, m.Type)
	err := row.Scan(&m.Name, &m.Type, &m.Gauge, &m.Counter)

	return errSelect(ctx, "READ", err)
}

//Читаем из базы массив метрик

func (dbs *DBStorage) ReadAll(mm *([]metric.Metric)) (int, error) {
	ctx, cancel := context.WithTimeout(dbs.ictx, dbs.timeout*time.Second)
	defer cancel()

	rows, err := dbs.DB.QueryContext(ctx, `select 
		name
		, type
		, gauge
		, counter 
	from metrics `)
	if err != nil {
		return 0, errSelect(ctx, "select", fmt.Errorf("db select error: %w", err))
	}
	defer rows.Close()

	cnt := 0

	for rows.Next() {
		m := metric.Metric{}
		if err = rows.Scan(&m.Name, &m.Type, &m.Gauge, &m.Counter); err != nil {
			return cnt, errSelect(ctx, "scan", fmt.Errorf("db scan row error: %w", err))
		} else {
			*mm = append(*mm, m)
		}
		cnt++
	}
	//Проверяем на ошибки в процессе чтения
	err = rows.Err()

	return cnt, errSelect(ctx, "READALL", err)
}

func (dbs *DBStorage) Truncate() error {
	ctx, cancel := context.WithTimeout(dbs.ictx, dbs.timeout*time.Second)
	defer cancel()

	sql := `truncate table metrics`
	_, err := dbs.DB.ExecContext(ctx, sql)

	return errSelect(ctx, "TRUNCATE", err)
}
