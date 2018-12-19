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
	templateDatas, err := processFormDatas(r, db)
	if err != nil {
		return nil, fmt.Errorf("Error while processing Form datas: %v", err)
	}
	// Get all apps to be displayed
	templateDatas.Apps, err = getApps(db)
	if err != nil {
		return nil, fmt.Errorf("Error while getting apps : %v", err)
	}
	return templateDatas, nil
}

// processFormDatas processes all datas sent by forms
func processFormDatas(r *http.Request, db *scribble.Driver) (templateDatas *TmplVars, err error) {
	templateDatas = &TmplVars{}
	if err = r.ParseForm(); err != nil {
		return nil, fmt.Errorf("param parsing error: %v", err)
	}
	// Read action
	action, isAction := r.Form["action"]
	// Read App concerned
	app, isApp := r.Form["app"]
	if isApp {
		templateDatas.SelectedApp = app[0]
		if isAction {
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
				return nil, fmt.Errorf("Error with action %v for app %v: %v", action, app, err)
			}
		}
	}
	return templateDatas, nil
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
		return fmt.Errorf("Error while adding new app : %v", err)
	}
	return nil
}

// processAppToDelete deletes the given app
func processAppToDelete(appName string, db *scribble.Driver) error {
	err := deleteApp(db, App{AppName: appName})
	if err != nil {
		return fmt.Errorf("Error while deleting app : %v", err)
	}
	return nil
}

// processPasswordToShow displays the password for the given app
func processPasswordToShow(appName string, r *http.Request, tmplVars *TmplVars, db *scribble.Driver) (err error) {
	tmplVars.AppPasswordToShow, err = getApp(db, appName)
	if err != nil {
		return fmt.Errorf("Error while retrieving app: %v", err)
	}
	passphrase, isPassphrase := r.Form["passphrase"]
	if isPassphrase {
		inc := len(tmplVars.AppPasswordToShow.Versions) - 1
		tmplVars.Pass = computePassword(passphrase[0], tmplVars.AppPasswordToShow, inc)
	}
	return nil
}

// processHistoryToShow displays the history for the given app
func processHistoryToShow(appName string, r *http.Request, tmplVars *TmplVars, db *scribble.Driver) (err error) {
	tmplVars.AppHistoryToShow, err = getApp(db, appName)
	if err != nil {
		return fmt.Errorf("Error while retrieving app: %v", err)
	}
	passphrase, isPassphrase := r.Form["passphrase"]
	if isPassphrase {
		var pass, ts string
		for i := len(tmplVars.AppHistoryToShow.Versions) - 1; i >= 0; i-- {
			pass = computePassword(passphrase[0], tmplVars.AppHistoryToShow, i)
			ts = tmplVars.AppHistoryToShow.Versions[i].Format("02 Jan 2006 15:04")
			tmplVars.AppHistory = append(tmplVars.AppHistory, TmplPassHistory{TS: ts, Pass: pass})
		}
	}
	return nil
}

// processNewPassword adds a version for the given app
func processNewPassword(appName string, tmplVars *TmplVars, db *scribble.Driver) (err error) {
	tmplVars.AppHistoryToShow, err = getApp(db, appName)
	if err != nil {
		return fmt.Errorf("Error while retrieving app: %v", err)
	}
	tmplVars.AppHistoryToShow.Versions = append(tmplVars.AppHistoryToShow.Versions, time.Now())
	err = storeApp(db, tmplVars.AppHistoryToShow)
	if err != nil {
		return fmt.Errorf("Error while saving app: %v", err)
	}
	return nil
}
