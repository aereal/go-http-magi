package main

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
)

type newAppSuccessTestCase struct {
	args               []string
	expectedSiteConfig *Site
}

type newAppFailureTestCase struct {
	args                 []string
	message              string
	expectedErrorMessage string
}

func TestNewApp_Valid(t *testing.T) {
	cases := []newAppSuccessTestCase{
		newAppSuccessTestCase{
			args: strings.Split("magi -config ./testdata/valid.json", " "),
			expectedSiteConfig: &Site{
				Name: "aereal.org",
				URLs: []string{
					"https://aereal.org/",
					"https://aereal.org/subdir/",
				},
			},
		},
	}

	for _, testCase := range cases {
		app, err := newApp(testCase.args, new(bytes.Buffer), new(bytes.Buffer))
		if err != nil {
			t.Errorf("newApp(%v) should succeed but failure with error: %s", testCase.args, err)
			continue
		}
		err = eqSite(testCase.expectedSiteConfig, app.site)
		if err != nil {
			t.Errorf("newApp(%v) should build *App instance with\n\texpected *Site: %#v\n\tactual *Site: %#v\n\terror: %s", testCase.args, testCase.expectedSiteConfig, app.site, err)
			continue
		}
	}
}

func TestNewApp_Invalid(t *testing.T) {
	cases := []newAppFailureTestCase{
		newAppFailureTestCase{
			args:                 strings.Split("magi -unknown", " "),
			message:              "app should reject unknown parameter but got %#v",
			expectedErrorMessage: "flag provided but not defined: -unknown",
		},
		newAppFailureTestCase{
			args:                 strings.Split("magi -config testdata/invalid.json", " "),
			message:              "app should require urls but got %#v",
			expectedErrorMessage: "urls required",
		},
	}

	for _, testCase := range cases {
		_, err := newApp(testCase.args, new(bytes.Buffer), new(bytes.Buffer))
		if err == nil {
			t.Errorf("newApp(%v) should return some error but nothing got", testCase.args)
			continue
		}
		if err.Error() != testCase.expectedErrorMessage {
			t.Errorf(testCase.message, err)
		}
	}
}

func eqSite(expected *Site, actual *Site) error {
	if expected.Name != actual.Name {
		return fmt.Errorf("expected name: %s; actual: %s", expected.Name, actual.Name)
	}
	if len(expected.URLs) != len(actual.URLs) {
		return fmt.Errorf("expected # of URLs: %d; actual: %d", len(expected.URLs), len(actual.URLs))
	}
	for i, expectedURL := range expected.URLs {
		actualURL := actual.URLs[i]
		if expectedURL != actualURL {
			return fmt.Errorf("expected urls[%d]: %s; actual: %s", i, expectedURL, actualURL)
		}
	}
	return nil
}
