package main

import (
	"fmt"
	"log"
	"strconv"

	"github.com/Ssnakerss/practicum-metrics/internal/metric"
	"github.com/Ssnakerss/practicum-metrics/internal/storage"
)

// Модуль для тестирования и отладки
func main() {

	//`postgres://postgres:postgres@postgres:5432/praktikum?sslmode=disable`
	dsn := `postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable`

	var db storage.DBStorage
	if err := db.New(dsn, "1"); err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err := db.CheckStorage(); err != nil {
		log.Fatal(err)
	}
	fmt.Println("connection ok")

	//читаем метрику
	m := metric.Metric{}
	row := db.Db.QueryRow(`select 
			name
			, type
			, gauge 
			, counter
		from metrics 
		where  name = 'GCSys' and type = 'gauge'`)

	if err := row.Scan(&m.Name, &m.Type, &m.Gauge, &m.Counter); err != nil {
		log.Fatalf("could not scan row: %v", err)
	}

	fmt.Printf("found row %v", m)

	for i := 0; i < 10; i++ {
		sql := `insert into metrics 
		(name, type, gauge, counter)
		values
		($1,$2,$3,$4)
		on conflict(name,type) do update
			set gauge = excluded.gauge,
				counter = metrics.counter + excluded.counter`

		result, err := db.Db.Exec(
			sql,
			"TCounter"+strconv.Itoa(i+100),
			"counter",
			i,
			i,
		)

		if err != nil {
			log.Fatalf("could not insert row: %v", err)
		}
		rowsAffected, err := result.RowsAffected()
		if err != nil {
			log.Fatalf("could not get affected rows: %v", err)
		}
		// we can log how many rows were inserted
		fmt.Println("inserted", rowsAffected, "rows")
	}

	mm := []metric.Metric{}

	rows, err := db.Db.Query(`select 
		name
		, type
		, gauge
		, counter 
	from metrics `)
	if err != nil {
		log.Fatal("error query db")
	}
	for rows.Next() {
		if err := rows.Scan(&m.Name, &m.Type, &m.Gauge, &m.Counter); err != nil {
			fmt.Printf("could not scan row: %v\n", err)
		} else {
			mm = append(mm, m)
		}
	}
	fmt.Printf("collected metrcics %v", mm)

}
