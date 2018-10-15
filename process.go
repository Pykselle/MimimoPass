package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/nanobox-io/golang-scribble"
)

// App represents the stored application
type App struct {
	AppName         string
	UseSpecialChars bool
	Increment       int
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
	// Check if an app was selected
	app, isApp := r.Form["app"]
	passphrase, isPassphrase := r.Form["passphrase"]
	if isApp && isPassphrase {
		err = processSelectedApp(app[0], passphrase[0], &templateDatas, db)
		if err != nil {
			return nil, fmt.Errorf("Error while processing selected app: %v", err)
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
	return &templateDatas, nil
}

// processSelectedApp retrieves the given app in scribble, computes the password
// to be displayed and updates the template datas
func processSelectedApp(appName string, passphrase string, tmplVars *TmplVars, db *scribble.Driver) error {
	tmplVars.SelectedApp = appName
	app := App{}
	if err := db.Read("app", appName, &app); err != nil {
		return fmt.Errorf("Error while getting app : %v", err)
	}
	tmplVars.Pass = computePassword(passphrase, app)
	return nil
}

// processAppToCreate creates the given app
func processAppToCreate(appName string, useSpecialChars bool, db *scribble.Driver) error {
	err := storeApp(db, App{AppName: appName, UseSpecialChars: useSpecialChars})
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
