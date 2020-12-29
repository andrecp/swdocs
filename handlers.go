package swdocs

import (
	"encoding/json"
	"html/template"
	"io/ioutil"
	"net/http"

	log "github.com/sirupsen/logrus"
)

type (
	HomePage struct {
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

func (a *App) homeHandler(w http.ResponseWriter, r *http.Request) {
	message, err := ioutil.ReadFile("templates/home.gohtml")
	if err != nil {
		panic(err)
	}

	t, err := template.New("HomePage").Parse(string(message))
	if err != nil {
		panic(err)
	}

	docs := GetMostRecentSwDocs()
	h := HomePage{docs}
	err = t.Execute(w, h)
	if err != nil {
		panic(err)
	}
}

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
