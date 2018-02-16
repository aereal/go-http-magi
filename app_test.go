package main

import (
	"bytes"
	"strings"
	"testing"
)

type newAppFailureTestCase struct {
	args                 []string
	message              string
	expectedErrorMessage string
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
