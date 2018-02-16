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
			args: strings.Split("magi -name aereal.org -url http://aereal.org/", " "),
			expectedSiteConfig: &Site{
				name: "aereal.org",
				urls: []string{
					"http://aereal.org/",
				},
			},
		},
	}

	for _, testCase := range cases {
		app, err := newApp(testCase.args, new(bytes.Buffer), new(bytes.Buffer))
		if err != nil {
			t.Errorf("newApp(%v) should succeed but failure with error %s", testCase.args, err)
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
			args:                 strings.Split("magi", " "),
			message:              "app should require name parameter; error: %#v",
			expectedErrorMessage: "name required",
		},
		newAppFailureTestCase{
			args:                 strings.Split("magi -name aereal.org", " "),
			message:              "app should require url parameter but got %#v",
			expectedErrorMessage: "URLs required",
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
	if expected.name != actual.name {
		return fmt.Errorf("expected name: %s; actual: %s", expected.name, actual.name)
	}
	if len(expected.urls) != len(actual.urls) {
		return fmt.Errorf("expected # of URLs: %d; actual: %d", len(expected.urls), len(actual.urls))
	}
	for i, expectedURL := range expected.urls {
		actualURL := actual.urls[i]
		if expectedURL != actualURL {
			return fmt.Errorf("expected urls[%d]: %s; actual: %s", i, expectedURL, actualURL)
		}
	}
	return nil
}
