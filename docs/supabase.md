```shell
go get github.com/jackc/pgx/v5

package main

import (
	"context"
	"log"
	"os"
	"github.com/jackc/pgx/v5"
)

func main() {
	conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}
	defer conn.Close(context.Background())

	// Example query to test connection
	var version string
	if err := conn.QueryRow(context.Background(), "SELECT version()").Scan(&version); err != nil {
		log.Fatalf("Query failed: %v", err)
	}

	log.Println("Connected to:", version)
}

user=postgres.cnlbftaurppzamivofhb 
password=[YOUR-PASSWORD] 
host=aws-1-ap-northeast-2.pooler.supabase.com
port=6543
dbname=postgres
```