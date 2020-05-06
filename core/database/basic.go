package database

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
	"github.com/sirupsen/logrus"
)

// GetInitializedDatabase - Return initialized db object with created tables if they wasn't created before
func GetInitializedDatabase() *sql.DB {
	db, err := sql.Open("sqlite3", "./app.db")

	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(UsersTable)
	if nil != err {
		logrus.Fatal("Can't create users table")
	}

	_, err = db.Exec(GroupsTable)
	if nil != err {
		logrus.Fatal("Can't create groups table")
	}

	_, err = db.Exec(AuthorizedTable)
	if nil != err {
		logrus.Fatal("Can't create authorized table")
	}

	_, err = db.Exec(GenresTable)
	if nil != err {
		logrus.Fatal("Can't create genres table")
	}

	_, err = db.Exec(TranslatorsTable)
	if nil != err {
		logrus.Fatal("Can't create translators table")
	}

	_, err = db.Exec(AuthorsTable)
	if nil != err {
		logrus.Fatal("Can't create authors table")
	}

	return db
}
