package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/freelanceplatform/server/internal/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	clean := flag.Bool("clean", false, "truncate test data before seeding")
	flag.Parse()

	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load config: %v\n", err)
		os.Exit(1)
	}

	ctx := context.Background()
	pool, err := pgxpool.New(ctx, cfg.Database.DSN())
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to connect database: %v\n", err)
		os.Exit(1)
	}
	defer pool.Close()

	if err := pool.Ping(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "failed to ping database: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Database connected")

	if *clean {
		fmt.Println("Cleaning test data...")
		tables := []string{
			"wallet_transactions", "payments", "contracts", "bids",
			"project_milestones", "projects", "reviews",
			"developer_portfolio", "developer_skills", "developer_profiles",
			"client_profiles", "notifications", "platform_wallets",
			"refresh_tokens", "oauth_accounts", "users",
		}
		for _, t := range tables {
			if _, err := pool.Exec(ctx, fmt.Sprintf("DELETE FROM %s WHERE id::text LIKE 'a0000000-%%' OR id::text LIKE 'a1000000-%%' OR id::text LIKE 'a2000000-%%' OR id::text LIKE 'a3000000-%%' OR id::text LIKE 'a4000000-%%' OR id::text LIKE 'a5000000-%%' OR id::text LIKE 'a6000000-%%' OR id::text LIKE 'a7000000-%%' OR id::text LIKE 'a8000000-%%' OR id::text LIKE 'b0000000-%%' OR id::text LIKE 'c0000000-%%' OR id::text LIKE 'd0000000-%%' OR id::text LIKE 'e0000000-%%' OR id::text LIKE 'f0000000-%%'", t)); err != nil {
				fmt.Fprintf(os.Stderr, "  warning: clean %s: %v\n", t, err)
			}
		}
		fmt.Println("Test data cleaned")
	}

	fmt.Println("Seeding test data...")
	sql, err := os.ReadFile("migrations/seed_test_data.sql")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to read seed file: %v\n", err)
		os.Exit(1)
	}

	if _, err := pool.Exec(ctx, string(sql)); err != nil {
		fmt.Fprintf(os.Stderr, "failed to execute seed: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Test data seeded successfully!")
	fmt.Println("")
	fmt.Println("Test accounts (password: Test123456):")
	fmt.Println("  Clients:    client1@test.com, client2@test.com, client3@test.com")
	fmt.Println("  Developers: dev1@test.com, dev2@test.com, dev3@test.com, dev4@test.com, dev5@test.com")
	fmt.Println("  Both:       both1@test.com, both2@test.com")
}
