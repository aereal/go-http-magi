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
	err := fmt.Errorf("%s: %s on %s", chkr.Status.String(), chkr.String(), url)
	return chkr.Status, err
}

// URLCheckResult represents a result of external monitoring.
type URLCheckResult struct {
	err    error
	status checkers.Status
}

func (r *URLCheckResult) ok() bool {
	return r.err == nil
}
