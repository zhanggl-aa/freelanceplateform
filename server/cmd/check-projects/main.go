package main

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/spf13/viper"
)

func main() {
	viper.SetConfigFile("config.yaml")
	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("Error reading config: %v\n", err)
		return
	}

	dbURL := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s",
		viper.GetString("database.user"),
		viper.GetString("database.password"),
		viper.GetString("database.host"),
		viper.GetInt("database.port"),
		viper.GetString("database.dbname"),
		viper.GetString("database.sslmode"),
	)

	ctx := context.Background()
	pool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		fmt.Printf("Error connecting: %v\n", err)
		return
	}
	defer pool.Close()

	// Total count
	var total int
	err = pool.QueryRow(ctx, "SELECT COUNT(*) FROM projects").Scan(&total)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("Total Projects: %d\n", total)

	// Status breakdown
	rows, err := pool.Query(ctx, "SELECT status, COUNT(*) FROM projects GROUP BY status ORDER BY status")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Println("\nStatus Breakdown:")
	for rows.Next() {
		var status string
		var count int
		rows.Scan(&status, &count)
		fmt.Printf("  %-15s %d\n", status, count)
	}
	rows.Close()

	// Sample projects
	fmt.Println("\nSample Projects (first 5):")
	rows, err = pool.Query(ctx, "SELECT id, title, status, budget_min, budget_max FROM projects ORDER BY created_at LIMIT 5")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	for rows.Next() {
		var id, title, status string
		var min, max float64
		rows.Scan(&id, &title, &status, &min, &max)
		fmt.Printf("  %s... - %s (%s) - %.0f-%.0f\n", id[:8], title, status, min, max)
	}
	rows.Close()

	fmt.Println("\n✓ Done!")
}
