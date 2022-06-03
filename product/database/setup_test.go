package database_test

import (
	"database/sql"
	"fmt"
	"log"
	"minimarket/product/database"
	"os"
	"testing"

	"github.com/jmoiron/sqlx"
	dockertest "github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
)

var db *sqlx.DB

func TestMain(m *testing.M) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	cfg := []string{
		"MYSQL_DATABASE=test",
		"MYSQL_ROOT_PASSWORD=test",
		"MYSQL_ROOT_HOST=%",
	}

	container, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "mysql",
		Tag:        "latest",
		Env:        cfg,
	}, func(config *docker.HostConfig) {
		// set AutoRemove to true so that stopped container goes away by itself
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{
			Name: "no",
		}
	})
	if err != nil {
		log.Fatalf("Could not start container: %s", err)
	}

	port := container.GetPort("3306/tcp")

	if err := pool.Retry(func() error {
		var err error
		db, err := sql.Open("mysql", fmt.Sprintf("root:test@(localhost:%s)/mysql", port))
		if err != nil {
			return err
		}
		return db.Ping()
	}); err != nil {
		log.Fatalf("Could not connect to database: %s", err)
	}

	dbConfig := database.Config{
		Host:        "localhost",
		Port:        port,
		User:        "root",
		Pass:        "test",
		Name:        "test",
		SSLMode:     "disable",
		SSLCert:     "",
		SSLKey:      "",
		SSLRootCert: "",
	}

	if db, err = database.Connect(dbConfig); err != nil {
		log.Fatalf("Could not setup test DB connection: %s", err)
	}

	code := m.Run()

	// Defers will not be run when using os.Exit
	db.Close()
	if err := pool.Purge(container); err != nil {
		log.Fatalf("Could not purge container: %s", err)
	}

	os.Exit(code)
}
