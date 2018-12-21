package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/nanobox-io/golang-scribble"
)

// App represents the stored application
type App struct {
	AppName         string
	Login           string
	UseSpecialChars bool
	Versions        []time.Time
}

func storeApp(db *scribble.Driver, a *App) error {
	if err := db.Write("app", a.AppName, a); err != nil {
		return err
	}
	return nil
}

func getApps(db *scribble.Driver) ([]App, error) {
	records, err := db.ReadAll("app")
	if err != nil {
		return nil, err
	}
	apps := []App{}
	for _, f := range records {
		app := App{}
		if err := json.Unmarshal([]byte(f), &app); err != nil {
			return nil, err
		}
		apps = append(apps, app)
	}
	return apps, nil
}

func getApp(db *scribble.Driver, appName string) (*App, error) {
	app := &App{}
	if err := db.Read("app", appName, &app); err != nil {
		return nil, fmt.Errorf("Error while getting app : %v", err)
	}
	return app, nil
}

func deleteApp(db *scribble.Driver, a App) error {
	if err := db.Delete("app", a.AppName); err != nil {
		return err
	}
	return nil
}
