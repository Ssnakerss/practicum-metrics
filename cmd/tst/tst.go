package main

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

// Модуль для тестирования и отладки
func main() {
	// ps := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
	// 	`localhost`, `5432`, `postgres`, `postgres`, `postgres`)

	//`postgres://postgres:postgres@postgres:5432/praktikum?sslmode=disable`
	dsn := `postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable`

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	if err = db.PingContext(ctx); err != nil {
		panic(err)
	}
	fmt.Println("PING OK")
}
