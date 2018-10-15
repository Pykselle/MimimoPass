package main

import (
	"encoding/json"
	"fmt"

	"github.com/nanobox-io/golang-scribble"
)

func storeApp(db *scribble.Driver, a App) error {
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

func getApp(db *scribble.Driver, appName string) (App, error) {
	app := App{}
	if err := db.Read("app", appName, &app); err != nil {
		return app, fmt.Errorf("Error while getting app : %v", err)
	}
	return app, nil
}

func deleteApp(db *scribble.Driver, a App) error {
	if err := db.Delete("app", a.AppName); err != nil {
		return err
	}
	return nil
}
