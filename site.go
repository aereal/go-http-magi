package main

import (
	"github.com/mackerelio/checkers"
)

// Site represents a set of URLs
type Site struct {
	Name         string `json:"name"`
	PrimaryURL   string `json:"primary_url"`
	SecondaryURL string `json:"secondary_url"`
}

// SiteCheckResult represents a check result of the site.
type SiteCheckResult struct {
	urlResults map[string]*URLCheckResult
	statusCode int
}

func (r *SiteCheckResult) errors() []error {
	var errors []error
	for _, res := range r.urlResults {
		if !res.ok() {
			errors = append(errors, res.err)
		}
	}
	return errors
}

func (r *SiteCheckResult) ok() bool {
	return r.statusCode == int(checkers.OK)
}
