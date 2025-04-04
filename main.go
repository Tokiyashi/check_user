package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func checkSubscription(ctx context.Context, pool *pgxpool.Pool, username string) bool {
	var isPaid bool
	var expirationDate time.Time

	err := pool.QueryRow(ctx, `
		SELECT is_paid, expiration_date FROM subscriptions
		WHERE username = $1
	`, username).Scan(&isPaid, &expirationDate)

	if err != nil {
		if err.Error() == "no rows in result set" {
			// Нет подписки
			return false
		}
		fmt.Println("Error querying database:", err)
		return false
	}

	return isPaid && expirationDate.After(time.Now())
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: ./check_user <username>")
		os.Exit(1)
	}

	username := os.Args[1]

	// 🔐 Настрой подключение к базе
	dsn := "postgres://your_db_user:your_db_password@localhost:5432/your_db_name?sslmode=disable"

	ctx := context.Background()

	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		fmt.Println("Error connecting to database:", err)
		os.Exit(1)
	}
	defer pool.Close()

	if checkSubscription(ctx, pool, username) {
		os.Exit(0) // Успех
	} else {
		os.Exit(1) // Подписки нет или истекла
	}
}
