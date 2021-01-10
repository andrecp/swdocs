package swdocs

import (
	"database/sql"
	"fmt"
	"sync"

	log "github.com/sirupsen/logrus"

	"net/http"
	"os"

	// We're using sqlite implementation of the sql interface.
	_ "github.com/mattn/go-sqlite3"

	"github.com/gorilla/mux"
)

// App is the struct representing our web application containing
// the database connection, the router, a mutex and its configuration.
type App struct {
	Router *mux.Router
	DB     *sql.DB
	Mutex  sync.Mutex
	Config AppConfig
}

// AppConfig holds the configuration used by the application.
type AppConfig struct {
	Port          string
	TemplatesPath string
	DbPath        string
}

func (a *App) initializeRoutes() {
	// Web Pages
	a.Router.HandleFunc("/", a.homeHandler).Methods("GET")
	a.Router.HandleFunc("/search", a.searchHandler).Methods("GET")
	a.Router.HandleFunc("/{swDocName}", a.swDocHandler).Methods("GET")
	// REST API
	a.Router.HandleFunc("/api/v1/swdocs/", a.getSwDocsHandler).Methods("GET")
	a.Router.HandleFunc("/api/v1/swdocs/{swDocName}", a.getSwDocHandler).Methods("GET")
	a.Router.HandleFunc("/api/v1/swdocs/{swDocName}", a.deleteSwDocHandler).Methods("DELETE")
	a.Router.HandleFunc("/api/v1/swdocs/apply", a.applySwDocHandler).Methods("POST")

}

func (a *App) createDbIfNotExists() (bool, error) {
	_, err := os.Stat(a.Config.DbPath)
	if err == nil {
		log.Info("Database " + a.Config.DbPath + " already exists")
	} else if os.IsNotExist(err) {
		log.Info("Creating database " + a.Config.DbPath)
		file, err := os.Create(a.Config.DbPath)
		if err != nil {
			return true, err
		}
		return false, file.Close()
	}
	return true, nil
}

func (a *App) populateDb() error {
	_, err := a.DB.Exec("PRAGMA encoding = \"UTF-8\";")
	if err != nil {
		return err
	}

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

// Initialize the web app database and routes.
func (a *App) Initialize() {
	var err error

	// Create DB file if not exists.
	exists, err := a.createDbIfNotExists()
	if err != nil {
		log.Fatal(err.Error())
	}

	// Open up a DB connection.
	a.DB, err = sql.Open("sqlite3", a.Config.DbPath)
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

// Run the web application.
func (a *App) Run() {
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", a.Config.Port), a.Router))
}
