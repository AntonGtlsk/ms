package dao

import (
	"database/sql"
	"fmt"
	"logging"

	"github.com/sirupsen/logrus"

	_ "github.com/go-sql-driver/mysql"
	// "github.com/sirupsen/logrus"
)

type Storage struct {
	db *sql.DB
}

func MustConnect(storagePath string, logging *logging.Logger) *Storage {
	db, err := sql.Open("mysql", storagePath)

	if err != nil {
		fmt.Println(err)
		// logging.WithFields(logrus.Fields{
		// 	"error": err,
		// }).Fatalf("Cannot open mysql connection")
	}

	if db.Ping() != nil {
		fmt.Println(err)
		// logging.WithFields(logrus.Fields{
		// 	"error": err,
		// }).Fatalf("Cannot ping mysql connection")
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS message (
	  name varchar(50) NOT NULL,
	  avatarURL varchar(500) NOT NULL,
	  contractAddress varchar(42) NOT NULL,
	  guild varchar(150) NOT NULL,
	  body varchar(300) NOT NULL,
	  time datetime NOT NULL,
	  repliedMessageId int(11) NOT NULL,
	  id INT AUTO_INCREMENT PRIMARY KEY,
	  pinned tinyint(1) NOT NULL
	) ENGINE=InnoDB DEFAULT CHARSET=utf8;
		`)

	if err != nil {
		logging.WithFields(logrus.Fields{
			"error": err,
		}).Fatalf("Cannot create table message")
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS reaction (
	  userId varchar(25) NOT NULL,
	  name varchar(50) NOT NULL,
	  avatarURL varchar(500) NOT NULL,
	  reaction int NOT NULL,
	  guild varchar(100) NOT NULL,
	  contractAddress varchar(55) NOT NULL
	) ENGINE=InnoDB DEFAULT CHARSET=utf8;
		`)

	if err != nil {
		logging.WithFields(logrus.Fields{
			"error": err,
		}).Fatalf("Cannot create table reaction")

	}

	return &Storage{db: db}

}
