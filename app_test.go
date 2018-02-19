package main

import (
	"bytes"
	"fmt"
	"strings"
	"sync"
	"testing"

	"github.com/mackerelio/checkers"
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

type siteCheckResultTestCase struct {
	primaryCheckStatus   checkers.Status
	secondaryCheckStatus checkers.Status
	expectedStatusCode   int
}

func TestAccumulateResults(t *testing.T) {
	testCases := []*siteCheckResultTestCase{
		&siteCheckResultTestCase{
			primaryCheckStatus:   checkers.OK,
			secondaryCheckStatus: checkers.OK,
			expectedStatusCode:   int(checkers.OK),
		},
		&siteCheckResultTestCase{
			primaryCheckStatus:   checkers.CRITICAL,
			secondaryCheckStatus: checkers.OK,
			expectedStatusCode:   int(checkers.CRITICAL),
		},
		&siteCheckResultTestCase{
			primaryCheckStatus:   checkers.OK,
			secondaryCheckStatus: checkers.CRITICAL,
			expectedStatusCode:   int(checkers.CRITICAL),
		},
		&siteCheckResultTestCase{
			primaryCheckStatus:   checkers.CRITICAL,
			secondaryCheckStatus: checkers.CRITICAL,
			expectedStatusCode:   int(checkers.CRITICAL),
		},
	}
	for _, testCase := range testCases {
		out, errOut := new(bytes.Buffer), new(bytes.Buffer)
		app, _ := newApp(strings.Split("magi -config ./testdata/valid.json", " "), out, errOut)
		results := new(sync.Map)
		results.Store("https://aereal.org/subdir/", &URLCheckResult{
			status: testCase.primaryCheckStatus,
		})
		results.Store("https://aereal.org/", &URLCheckResult{
			status: testCase.secondaryCheckStatus,
		})
		siteCheckResult := app.accumulateResults(results)
		if siteCheckResult.statusCode != testCase.expectedStatusCode {
			t.Errorf(
				"site (primary:%v secondary:%v) check must be %#v but got %#v",
				testCase.primaryCheckStatus,
				testCase.secondaryCheckStatus,
				testCase.expectedStatusCode,
				siteCheckResult.statusCode,
			)
		}
	}
}

func TestNewApp_Valid(t *testing.T) {
	cases := []newAppSuccessTestCase{
		newAppSuccessTestCase{
			args: strings.Split("magi -config ./testdata/valid.json", " "),
			expectedSiteConfig: &Site{
				Name:         "aereal.org",
				PrimaryURL:   "https://aereal.org/subdir/",
				SecondaryURL: "https://aereal.org/",
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
			expectedErrorMessage: "primary_url required",
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
	if expected.PrimaryURL != actual.PrimaryURL {
		return fmt.Errorf("expected PrimaryURL: %s; actual: %s", expected.PrimaryURL, actual.PrimaryURL)
	}
	if expected.SecondaryURL != actual.SecondaryURL {
		return fmt.Errorf("expected SecondaryURL: %s; actual: %s", expected.SecondaryURL, actual.SecondaryURL)
	}
	return nil
}
