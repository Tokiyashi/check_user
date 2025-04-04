package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"
)

const (
	dbURL = "postgres://youruser:yourpass@localhost:5432/yourdb" // 🔁 Заменить на свои данные
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	username, err := reader.ReadString('\n')
	if err != nil {
		os.Exit(1)
	}
	password, err := reader.ReadString('\n')
	if err != nil {
		os.Exit(1)
	}

	username = strings.TrimSpace(username)
	password = strings.TrimSpace(password)

	ok := checkSubscription(username, password)
	if ok {
		os.Exit(0)
	} else {
		os.Exit(1)
	}
}

func checkSubscription(username, password string) bool {
	ctx := context.Background()
	conn, err := pgx.Connect(ctx, dbURL)
	if err != nil {
		fmt.Fprintln(os.Stderr, "DB connection error:", err)
		return false
	}
	defer conn.Close(ctx)

	var passwordHash string
	var subscriptionEnd time.Time
	var isActive bool

	err = conn.QueryRow(ctx, `
		SELECT password_hash, subscription_end, is_active
		FROM vpn_users
		WHERE username = $1
	`, username).Scan(&passwordHash, &subscriptionEnd, &isActive)

	if err != nil {
		fmt.Fprintln(os.Stderr, "Query error:", err)
		return false
	}

	if !isActive {
		return false
	}

	if err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(password)); err != nil {
		return false
	}

	if time.Now().After(subscriptionEnd) {
		return false
	}

	return true
}
