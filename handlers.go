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
	LastCreated *swDocsSlice
	LastUpdated *swDocsSlice
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func respondWithJSONError(w http.ResponseWriter, code int, message string) {
	log.WithFields(log.Fields{
		"code": code,
	}).Error(message)
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	log.WithFields(log.Fields{
		"code": code,
	}).Error(message)
	w.WriteHeader(code)
}

// Templated HTML pages //

func (a *App) homeHandler(w http.ResponseWriter, r *http.Request) {
	message, err := ioutil.ReadFile("../../templates/home.gohtml")
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	t, err := template.New("HomePage").Parse(string(message))
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	createDocs, err := getMostRecentCreatedSwDocs(a.DB)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	updatedDocs, err := getMostRecentUpdatedSwDocs(a.DB)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	c := swDocsSlice{&createDocs}
	u := swDocsSlice{&updatedDocs}

	h := createdAndUpdatedHomePage{
		LastCreated: &c,
		LastUpdated: &u,
	}
	err = t.Execute(w, h)
	if err != nil {
		log.Error(err.Error())
	}
}

func (a *App) swDocHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	swdocName := params["swDocName"]
	message, err := ioutil.ReadFile("../../templates/swdoc.gohtml")
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	t, err := template.New("SwDoc").Parse(string(message))
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	doc, err := getSwDocByName(a.DB, swdocName)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if doc.Name == "" {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, "SwDoc with this name does not exist")
		return
	}

	err = t.Execute(w, doc)
	if err != nil {
		log.Error(err.Error())
	}
}

func (a *App) searchHandler(w http.ResponseWriter, r *http.Request) {
	searchParams := r.URL.Query().Get("swdocsearch")
	message, err := ioutil.ReadFile("../../templates/search.gohtml")
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	t, err := template.New("Search").Parse(string(message))
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	docs, err := searchSwDocsByName(a.DB, searchParams)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	h := swDocsSlice{&docs}
	err = t.Execute(w, h)
	if err != nil {
		log.Error(err.Error())
	}
}

// REST API //

func (a *App) deleteSwDocHandler(w http.ResponseWriter, r *http.Request) {
	// Sqlite only allows one writer at a time, handlers that change the state must execute once at a time.
	a.Mutex.Lock()
	defer a.Mutex.Unlock()

	params := mux.Vars(r)
	swdocName := params["swDocName"]

	if err := deleteSwDoc(a.DB, swdocName); err != nil {
		respondWithJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, nil)
}

func (a *App) applySwDocHandler(w http.ResponseWriter, r *http.Request) {
	// Sqlite only allows one writer at a time, handlers that change the state must execute once at a time.
	a.Mutex.Lock()
	defer a.Mutex.Unlock()

	var s swDoc
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&s); err != nil {
		respondWithJSONError(w, http.StatusBadRequest, "Invalid request payload.\n"+err.Error())
		return
	}

	defer r.Body.Close()

	if err := createOrUpdateSwDoc(a.DB, &s); err != nil {
		respondWithJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, s)
}
