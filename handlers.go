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

type (
	SwDocsSlice struct {
		SwDocs *[]SwDoc
	}
)

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	log.Error(message)
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

	docs, err := GetMostRecentSwDocs(a.DB)
	if err != nil {
		panic(err)
	}

	h := SwDocsSlice{&docs}
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

func (a *App) createSwDocHandler(w http.ResponseWriter, r *http.Request) {
	var s SwDoc
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&s); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload.\n"+err.Error())
		return
	}

	defer r.Body.Close()

	if err := CreateSwDoc(a.DB, &s); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, s)
}

func (a *App) deleteSwDocHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	swdocName := params["swDocName"]

	if err := DeleteSwDoc(a.DB, swdocName); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, nil)
}
