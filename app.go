package swdocs

import (
	"database/sql"
	"fmt"

	log "github.com/sirupsen/logrus"

	"net/http"
	"os"

	_ "github.com/mattn/go-sqlite3"

	"github.com/gorilla/mux"
)

type App struct {
	Router *mux.Router
	DB     *sql.DB
}

func createDbIfNotExists(dbname string) error {
	dbFile := fmt.Sprintf("%s.db", dbname)
	_, err := os.Stat(dbFile)
	if err == nil {
		log.Info("Database " + dbFile + " already exists")
	} else if os.IsNotExist(err) {
		log.Info("Creating database " + dbFile)
		file, err := os.Create(dbFile)
		if err != nil {
			return err
		}
		return file.Close()
	}
	return nil
}

func (a *App) Initialize(dbname string) {
	var err error
	if err = createDbIfNotExists(dbname); err != nil {
		log.Fatal(err.Error())
	}

	a.DB, err = sql.Open("sqlite3", dbname)
	if err != nil {
		log.Fatal(err)
	}

	a.Router = mux.NewRouter()
}

func (a *App) Run(addr string) {

	a.Router.HandleFunc("/", YourHandler)

	// Bind to a port and pass our router in
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", addr), a.Router))
}
