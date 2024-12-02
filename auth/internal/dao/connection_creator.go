package dao

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"time"
)

type Storage struct {
	db *sql.DB
}

func New(storagePath string) (*Storage, error) {
	const fn = "storage.mysql.New"

	db, err := sql.Open("mysql", storagePath)
	if err != nil {
		log.Fatalf("cannot open mysql connection: %s", err)
	}

	if db.Ping() != nil {
		log.Fatalf("cannot ping mysql connection: %s", err)
	}

	db.SetConnMaxLifetime(1 * time.Minute)

	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS refreshsessions (
		id INT AUTO_INCREMENT PRIMARY KEY,
		userId VARCHAR(40) NOT NULL,
		refreshToken VARCHAR(10000) NOT NULL,
		discordRefreshToken VARCHAR(10000) NOT NULL,
		expiresIn INT NOT NULL,
		createdAt INT NOT NULL
	) ENGINE=InnoDB DEFAULT CHARSET=utf8;
`)

	if err != nil {
		return nil, fmt.Errorf("cannot create tables: %s", err)
	}

	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS  user (
		id INT AUTO_INCREMENT PRIMARY KEY,
		userId varchar(40) NOT NULL,
		name varchar(50) NOT NULL,
		avatarURL varchar(500) NOT NULL,
		subscription enum('none','sub') NOT NULL
	  ) ENGINE=InnoDB DEFAULT CHARSET=utf8;
	`)

	if err != nil {
		return nil, fmt.Errorf("cannot create tables: %s", err)
	}

	return &Storage{db: db}, nil

}
