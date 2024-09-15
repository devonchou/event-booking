package test

import (
	"context"
	"database/sql"
	"log"
	"os"
	"path/filepath"

	"github.com/testcontainers/testcontainers-go/modules/mysql"
)

func setupMysqlContainer() (string, *sql.DB, func()) {
	ctx := context.Background()

	dirPath, _ := os.Getwd()

	mysqlContainer, err := mysql.Run(ctx,
		"mysql:8.0.36",
		mysql.WithDatabase("testdb"),
		mysql.WithUsername("testuser"),
		mysql.WithPassword("testpass"),
		mysql.WithScripts(filepath.Join(dirPath, "test.sql")),
	)
	if err != nil {
		log.Fatalf("failed to start mysql container: %s", err)
	}

	terminate := func() {
		if err := mysqlContainer.Terminate(ctx); err != nil {
			log.Fatalf("failed to terminate mysql container: %s", err)
		}
	}

	dsn, err := mysqlContainer.ConnectionString(ctx, "charset=utf8mb4", "parseTime=true")
	if err != nil {
		terminate()
		log.Fatalf("could not get mysql dsn: %s", err)
	}

	client, err := sql.Open("mysql", dsn)
	if err != nil {
		terminate()
		log.Fatalf("could not connect to mysql: %s", err)
	}

	return dsn, client, terminate
}
