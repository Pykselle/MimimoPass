package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/nanobox-io/golang-scribble"
)

// App represents the stored application
type App struct {
	AppName         string
	UseSpecialChars bool
	Versions        []time.Time
}

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
func processFormDatas(r *http.Request, db *scribble.Driver) (*TmplVars, error) {
	templateDatas := TmplVars{}
	err := r.ParseForm()
	if err != nil {
		return nil, fmt.Errorf("param parsing error: %v", err)
	}
	// Check if there's a password to show
	if app, isPassAsked := r.Form["password"]; isPassAsked {
		templateDatas.SelectedApp = app[0]
		templateDatas.AppPasswordToShow, err = getApp(db, app[0])
		if err != nil {
			return nil, fmt.Errorf("Error while retrieving app: %v", err)
		}
		passphrase, isPassphrase := r.Form["passphrase"]
		if isPassAsked && isPassphrase {
			inc := len(templateDatas.AppPasswordToShow.Versions) - 1
			templateDatas.Pass = computePassword(passphrase[0], templateDatas.AppPasswordToShow, inc)
		}
	}

	// Check if there is a new app to create
	if newApp, isNewApp := r.Form["newAppName"]; isNewApp {
		_, useSpecials := r.Form["useSpecials"]
		err = processAppToCreate(newApp[0], useSpecials, db)
		if err != nil {
			return nil, fmt.Errorf("Error while processing app to create: %v", err)
		}
	}

	// Check if there is an app to delete
	if appToDel, isAppToDel := r.Form["delete"]; isAppToDel {
		err = processAppToDelete(appToDel[0], db)
		if err != nil {
			return nil, fmt.Errorf("Error while processing app to delete: %v", err)
		}
	}

	// Check if there is an app history to display
	if historyApp, isHistoryToShow := r.Form["history"]; isHistoryToShow {
		templateDatas.SelectedApp = historyApp[0]
		templateDatas.AppHistoryToShow, err = getApp(db, historyApp[0])
		if err != nil {
			return nil, fmt.Errorf("Error while retrieving app: %v", err)
		}
		passphrase, isPassphrase := r.Form["passphrase"]
		if isPassphrase {
			var pass, ts string
			for i := len(templateDatas.AppHistoryToShow.Versions) - 1; i >= 0; i-- {
				pass = computePassword(passphrase[0], templateDatas.AppHistoryToShow, i)
				ts = templateDatas.AppHistoryToShow.Versions[i].Format("02 Jan 2006 15:04")
				templateDatas.AppHistory = append(templateDatas.AppHistory, TmplPassHistory{TS: ts, Pass: pass})
			}
		}
	}

	// Check if there is a new password to generate
	if app, isNewPassAsked := r.Form["new"]; isNewPassAsked {
		templateDatas.SelectedApp = app[0]
		templateDatas.AppHistoryToShow, err = getApp(db, app[0])
		if err != nil {
			return nil, fmt.Errorf("Error while retrieving app: %v", err)
		}
		templateDatas.AppHistoryToShow.Versions = append(templateDatas.AppHistoryToShow.Versions, time.Now())
		err = storeApp(db, templateDatas.AppHistoryToShow)
		if err != nil {
			return nil, fmt.Errorf("Error while saving app: %v", err)
		}
	}

	return &templateDatas, nil
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
