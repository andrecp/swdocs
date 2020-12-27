package swdocs

import (
	"html/template"
	"io/ioutil"
	"net/http"
)

type (
	HomePage struct {
		SwDocs []SwDoc
	}
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	message, err := ioutil.ReadFile("templates/home.gohtml")
	if err != nil {
		panic(err)
	}

	t, err := template.New("HomePage").Parse(string(message))
	if err != nil {
		panic(err)
	}

	h := HomePage{}
	err = t.Execute(w, h)
	if err != nil {
		panic(err)
	}
}
