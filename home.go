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

// TmplVars contains vars to be used in the web template
type TmplVars struct {
	Pass        string
	Apps        []App
	SelectedApp string
	AppHistory  []TmplPassHistory
	Action      string
}

//TmplPassHistory contains vars to be used for the password history
type TmplPassHistory struct {
	TS   string
	Pass string
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
		tmplPath+"/modalSeePass.html",
		tmplPath+"/modalHistory.html",
		tmplPath+"/modalConfDelete.html",
	)
	if err != nil {
		log.Print("template parsing error: ", err)
	}
	if err = t.Execute(w, templateDatas); err != nil {
		log.Print("template executing error: ", err)
	}
}
