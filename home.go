package main

import (
	"html/template"
	"log"
	"net/http"
)

const (
	tmplPath      = "./template"
	minPassLength = 10
	maxPassLength = 15
)

// TmplVars is the base of a template variable.
type TmplVars struct {
	Pass        string
	Apps        []App
	SelectedApp string
	UseSpecial  bool
	NewApp      string
	AppToDelete string
}

// Home listens and serves the homepage.
func home(w http.ResponseWriter, r *http.Request) {
	// Process the datas to be stored in templateDatas
	templateDatas, err := process(r)
	if err != nil {
		log.Print("processing error: ", err)
	}

	// Display the page.
	t, err := template.New("home.html").ParseFiles(
		tmplPath+"/home.html",
		tmplPath+"/modalNewApp.html",
	)
	if err != nil {
		log.Print("template parsing error: ", err)
	}

	if err = t.Execute(w, templateDatas); err != nil {
		log.Print("template executing error: ", err)
	}
}
