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

func (a *App) initializeRoutes() {
	a.Router.HandleFunc("/", a.homeHandler)
	a.Router.HandleFunc("/api/v1/swdocs", a.createSwDoc).Methods("POST")
}

func (a *App) createDbIfNotExists(dbpath string) (bool, error) {
	_, err := os.Stat(dbpath)
	if err == nil {
		log.Info("Database " + dbpath + " already exists")
	} else if os.IsNotExist(err) {
		log.Info("Creating database " + dbpath)
		file, err := os.Create(dbpath)
		if err != nil {
			return true, err
		}
		return false, file.Close()
	}
	return true, nil
}

func (a *App) populateDb() error {
	dbSchema := `
    CREATE TABLE IF NOT EXISTS swdocs (
		id INTEGER PRIMARY KEY,
		name TEXT UNIQUE,
		description TEXT)
	`
	statement, err := a.DB.Prepare(dbSchema)
	if err != nil {
		return err
	}
	_, err = statement.Exec()
	if err != nil {
		return err
	}
	return nil
}

func (a *App) Initialize(dbpath string) {
	var err error

	// Create DB file if not exists.
	exists, err := a.createDbIfNotExists(dbpath)
	if err != nil {
		log.Fatal(err.Error())
	}

	// Open up a DB connection.
	a.DB, err = sql.Open("sqlite3", dbpath)
	if err != nil {
		log.Fatal(err)
	}

	// Create tables in DB if it didn't exist before.
	if !exists {
		err = a.populateDb()
		if err != nil {
			log.Fatal(err)
		}
	}

	// Initialize the web app routes.
	a.Router = mux.NewRouter()
	a.initializeRoutes()
}

func (a *App) Run(addr string) {
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", addr), a.Router))
}
