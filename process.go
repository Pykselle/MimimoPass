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
	// Check if there's a password to show
	if app, isPassAsked := r.Form["password"]; isPassAsked {
		if err = processPasswordToShow(app[0], r, templateDatas, db); err != nil {
			return nil, fmt.Errorf("Error while processing password to show: %v", err)
		}
	}
	// Check if there is a new app to create
	if newApp, isNewApp := r.Form["newAppName"]; isNewApp {
		_, useSpecials := r.Form["useSpecials"]
		if err = processAppToCreate(newApp[0], useSpecials, db); err != nil {
			return nil, fmt.Errorf("Error while processing app to create: %v", err)
		}
	}
	// Check if there is an app to delete
	if appToDel, isAppToDel := r.Form["delete"]; isAppToDel {
		if err = processAppToDelete(appToDel[0], db); err != nil {
			return nil, fmt.Errorf("Error while processing app to delete: %v", err)
		}
	}
	// Check if there is an app history to display
	if historyApp, isHistoryToShow := r.Form["history"]; isHistoryToShow {
		if err = processHistoryToShow(historyApp[0], r, templateDatas, db); err != nil {
			return nil, fmt.Errorf("Error while processing history to display: %v", err)
		}
	}
	// Check if there is a new password to generate
	if app, isNewPassAsked := r.Form["new"]; isNewPassAsked {
		if err = processNewPassword(app[0], templateDatas, db); err != nil {
			return nil, fmt.Errorf("Error while processing new password: %v", err)
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

// processPasswordToShow TODO
func processPasswordToShow(appName string, r *http.Request, tmplVars *TmplVars, db *scribble.Driver) (err error) {
	tmplVars.SelectedApp = appName
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

// processHistoryToShow TODO
func processHistoryToShow(appName string, r *http.Request, tmplVars *TmplVars, db *scribble.Driver) (err error) {
	tmplVars.SelectedApp = appName
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

// processNewPassword TODO
func processNewPassword(appName string, tmplVars *TmplVars, db *scribble.Driver) (err error) {
	tmplVars.SelectedApp = appName
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
