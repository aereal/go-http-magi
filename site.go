package main

import (
	"os/exec"
	"syscall"
)

type Site struct {
	name string
	urls []string
}

type SiteCheckResult struct {
	ok         bool
	urlResults map[string]*URLCheckResult
}

func (r *SiteCheckResult) errors() []error {
	var errors []error
	for _, res := range r.urlResults {
		if !res.OK() {
			errors = append(errors, res.err)
		}
	}
	return errors
}

func (r *SiteCheckResult) status() int {
	totalStatus := 0
	for _, res := range r.urlResults {
		if exit, ok := res.err.(*exec.ExitError); ok {
			if status, ok := exit.Sys().(syscall.WaitStatus); ok && status.ExitStatus() > totalStatus {
				totalStatus = status.ExitStatus()
			}
		}
	}
	return totalStatus
}
