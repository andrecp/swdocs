package swdocs

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

type createdAndUpdatedHomePage struct {
	LastCreated *SwDocsSlice
	LastUpdated *SwDocsSlice
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	log.WithFields(log.Fields{
		"code": code,
	}).Error(message)
	respondWithJSON(w, code, map[string]string{"error": message})
}

// Templated HTML pages //

func (a *App) homeHandler(w http.ResponseWriter, r *http.Request) {
	message, err := ioutil.ReadFile("../../templates/home.gohtml")
	if err != nil {
		panic(err)
	}

	t, err := template.New("HomePage").Parse(string(message))
	if err != nil {
		panic(err)
	}

	createDocs, err := GetMostRecentCreatedSwDocs(a.DB)
	if err != nil {
		panic(err)
	}
	updatedDocs, err := GetMostRecentUpdatedSwDocs(a.DB)
	if err != nil {
		panic(err)
	}

	c := SwDocsSlice{&createDocs}
	u := SwDocsSlice{&updatedDocs}

	h := createdAndUpdatedHomePage{
		LastCreated: &c,
		LastUpdated: &u,
	}
	err = t.Execute(w, h)
	if err != nil {
		panic(err)
	}
}

func (a *App) swDocHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	swdocName := params["swDocName"]
	message, err := ioutil.ReadFile("../../templates/swdoc.gohtml")
	if err != nil {
		panic(err)
	}

	t, err := template.New("SwDoc").Parse(string(message))
	if err != nil {
		panic(err)
	}

	doc, err := GetSwDocByName(a.DB, swdocName)
	if err != nil {
		panic(err)
	}

	if doc.Name == "" {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, "SwDoc with this name does not exist")
		return
	}

	err = t.Execute(w, doc)
	if err != nil {
		panic(err)
	}
}

func (a *App) searchHandler(w http.ResponseWriter, r *http.Request) {
	searchParams := r.URL.Query().Get("swdocsearch")
	message, err := ioutil.ReadFile("../../templates/search.gohtml")
	if err != nil {
		panic(err)
	}

	t, err := template.New("Search").Parse(string(message))
	if err != nil {
		panic(err)
	}

	docs, err := SearchSwDocsByName(a.DB, searchParams)
	if err != nil {
		panic(err)
	}

	h := SwDocsSlice{&docs}
	err = t.Execute(w, h)
	if err != nil {
		panic(err)
	}
}

// REST API //

func (a *App) deleteSwDocHandler(w http.ResponseWriter, r *http.Request) {
	// Sqlite only allows one writer at a time, handlers that change the state must execute once at a time.
	a.Mutex.Lock()
	defer a.Mutex.Unlock()

	params := mux.Vars(r)
	swdocName := params["swDocName"]

	if err := DeleteSwDoc(a.DB, swdocName); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, nil)
}

func (a *App) applySwDocHandler(w http.ResponseWriter, r *http.Request) {
	// Sqlite only allows one writer at a time, handlers that change the state must execute once at a time.
	a.Mutex.Lock()
	defer a.Mutex.Unlock()

	var s SwDoc
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&s); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload.\n"+err.Error())
		return
	}

	defer r.Body.Close()

	if err := CreateOrUpdateSwDoc(a.DB, &s); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, s)
}
