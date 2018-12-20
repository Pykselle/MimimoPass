package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/nanobox-io/golang-scribble"
)

// process parses the form datas sent in POST or GET and returns the datas to
// be used in HTML template
func process(r *http.Request) (*TmplVars, error) {
	// Open connexion to scribble
	curDir, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("Error while getting current directory : %v", err)
	}
	db, err := scribble.New(curDir, nil)
	if err != nil {
		return nil, fmt.Errorf("Error with scribble driver : %v", err)
	}
	// Process datas sent by forms
	templateDatas := processFormDatas(r, db)

	// Get all apps to be displayed
	templateDatas.Apps, err = getApps(db)
	if err != nil {
		return nil, fmt.Errorf("Error while getting apps : %v", err)
	}
	return templateDatas, nil
}

// processFormDatas processes all datas sent by forms
func processFormDatas(r *http.Request, db *scribble.Driver) (templateDatas *TmplVars) {
	templateDatas = &TmplVars{}
	var err error
	if err = r.ParseForm(); err != nil {
		templateDatas.Action = "error"
		templateDatas.Error = fmt.Sprintf("param parsing error: %v", err)
	}
	// Read action
	action, isAction := r.Form["action"]
	// Read App concerned
	app, isApp := r.Form["app"]
	if isApp {
		templateDatas.SelectedApp = app[0]
		if isAction {
			templateDatas.Action = action[0]
			switch action[0] {
			case "delete":
				err = processAppToDelete(app[0], db)
			case "newpass":
				err = processNewPassword(app[0], templateDatas, db)
			case "history":
				err = processHistoryToShow(app[0], r, templateDatas, db)
			case "showpass":
				err = processPasswordToShow(app[0], r, templateDatas, db)
			case "newapp":
				_, useSpecials := r.Form["useSpecials"]
				err = processAppToCreate(app[0], useSpecials, db)
			}
			if err != nil {
				templateDatas.Action = "error"
				templateDatas.Error = fmt.Sprintf("Error with action %v for app %v: %v", action, app, err)
			}
		}
	}
	return templateDatas
}

// processAppToCreate creates the given app
func processAppToCreate(appName string, useSpecialChars bool, db *scribble.Driver) error {
	err := storeApp(
		db,
		App{
			AppName:         appName,
			UseSpecialChars: useSpecialChars,
			Versions:        []time.Time{time.Now()},
		},
	)
	if err != nil {
		return err
	}
	return nil
}

// processAppToDelete deletes the given app
func processAppToDelete(appName string, db *scribble.Driver) error {
	err := deleteApp(db, App{AppName: appName})
	if err != nil {
		return err
	}
	return nil
}

// processPasswordToShow displays the password for the given app
func processPasswordToShow(appName string, r *http.Request, tmplVars *TmplVars, db *scribble.Driver) error {
	app, err := getApp(db, appName)
	if err != nil {
		return fmt.Errorf("Error while retrieving app: %v", err)
	}
	passphrase, isPassphrase := r.Form["passphrase"]
	if isPassphrase {
		inc := len(app.Versions) - 1
		tmplVars.Pass = computePassword(passphrase[0], app, inc)
	}
	return nil
}

// processHistoryToShow displays the history for the given app
func processHistoryToShow(appName string, r *http.Request, tmplVars *TmplVars, db *scribble.Driver) error {
	app, err := getApp(db, appName)
	if err != nil {
		return fmt.Errorf("Error while retrieving app: %v", err)
	}
	passphrase, isPassphrase := r.Form["passphrase"]
	if isPassphrase {
		var pass, ts string
		for i := len(app.Versions) - 1; i >= 0; i-- {
			pass = computePassword(passphrase[0], app, i)
			ts = app.Versions[i].Format("02 Jan 2006 15:04")
			tmplVars.AppHistory = append(tmplVars.AppHistory, TmplPassHistory{TS: ts, Pass: pass})
		}
	}
	return nil
}

// processNewPassword adds a version for the given app
func processNewPassword(appName string, tmplVars *TmplVars, db *scribble.Driver) error {
	app, err := getApp(db, appName)
	if err != nil {
		return fmt.Errorf("Error while retrieving app: %v", err)
	}
	app.Versions = append(app.Versions, time.Now())
	err = storeApp(db, app)
	if err != nil {
		return fmt.Errorf("Error while saving app: %v", err)
	}
	tmplVars.Action = "history"
	return nil
}
