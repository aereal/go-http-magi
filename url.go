package main

import (
	"fmt"
	"io"

	"github.com/mackerelio/checkers"
	checkhttp "github.com/mackerelio/go-check-plugins/check-http/lib"
)

func runCheckHTTP(url string, outStream, errorStream io.Writer) (checkers.Status, error) {
	chkr := checkhttp.Run([]string{"-u", url})
	chkr.Name = "HTTP"
	if chkr.Status == checkers.OK {
		return chkr.Status, nil
	}
	return chkr.Status, fmt.Errorf(chkr.String())
}

// URLCheckResult represents a result of external monitoring.
type URLCheckResult struct {
	outStream   io.Writer
	errorStream io.Writer
	err         error
	status      checkers.Status
}

func (r *URLCheckResult) ok() bool {
	return r.err == nil
}
