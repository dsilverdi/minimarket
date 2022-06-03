package database

import (
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	migrate "github.com/rubenv/sql-migrate"
)

type Config struct {
	Host        string
	Port        string
	User        string
	Pass        string
	Name        string
	SSLMode     string
	SSLCert     string
	SSLKey      string
	SSLRootCert string
}

func Connect(cfg Config) (*sqlx.DB, error) {
	var db *sqlx.DB
	var err error
	url := cfg.User + ":" + cfg.Pass + "@tcp(" + cfg.Host + ":" + cfg.Port + ")/" + cfg.Name + "?parseTime=true&clientFoundRows=true"

	for {
		db, err = sqlx.Connect("mysql", url)
		if err == nil {
			break
		}

		if !strings.Contains(err.Error(), "connect: connection refused") {
			return nil, err
		}

		const retryDuration = 5 * time.Second
		time.Sleep(retryDuration)
	}

	if err := migrateDB(db); err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(5)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(10 * time.Minute)

	return db, nil
}

func migrateDB(db *sqlx.DB) error {
	migrations := &migrate.MemoryMigrationSource{
		Migrations: []*migrate.Migration{
			{
				Id: "products_1",
				Up: []string{
					`CREATE TABLE IF NOT EXISTS PRODUCT (
						id       	INT	NOT NULL AUTO_INCREMENT,
						name    	VARCHAR(255) NOT NULL,
						description	VARCHAR(255),
						category    VARCHAR(255),
						picture		VARCHAR(255),
						created_at  DATETIME,
						updated_at	DATETIME,
						PRIMARY KEY (id)
					);`,

					`CREATE TABLE IF NOT EXISTS PRODUCT_COMMENT (
						id       	INT	NOT NULL AUTO_INCREMENT,
						product_id  INT,
						message		VARCHAR(255),
						reply_id	INT,
						owner		VARCHAR(255),
						created_at  DATETIME,
						updated_at	DATETIME,
						PRIMARY KEY (id),
						FOREIGN KEY (product_id) REFERENCES PRODUCT(id)
					);`,
				},
				Down: []string{
					"DROP TABLE users",
				},
			},
		},
	}

	_, err := migrate.Exec(db.DB, "mysql", migrations, migrate.Up)
	return err
}
