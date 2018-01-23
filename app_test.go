package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestNewApp(t *testing.T) {
	var err error

	siteName := "aereal.org"
	app, err := newApp(strings.Split("magi -name "+siteName+" -url https://aereal.org/", " "), new(bytes.Buffer), new(bytes.Buffer))
	if err != nil {
		t.Errorf("app should be successfully initiated but got error: %s", err)
	}
	if app.site.name != siteName {
		t.Errorf("app.site.name should be %s but got %s", siteName, app.site.name)
	}

	_, err = newApp(strings.Split("magi -unknown", " "), new(bytes.Buffer), new(bytes.Buffer))
	if err == nil {
		t.Errorf("app should reject unknown parameter but got %#v", err)
	}

	_, err = newApp(strings.Split("magi", " "), new(bytes.Buffer), new(bytes.Buffer))
	if err == nil || err.Error() != "name required" {
		t.Errorf("app should require name parameter; error: %#v", err)
	}

	_, err = newApp(strings.Split("magi -name aereal.org", " "), new(bytes.Buffer), new(bytes.Buffer))
	if err == nil || err.Error() != "URLs required" {
		t.Errorf("app should require url parameter but got %#v", err)
	}
}
