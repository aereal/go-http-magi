package main

import (
	"io"
	"os/exec"
)

func runCheckHTTP(url string, outStream, errorStream io.Writer) error {
	cmd := exec.Command("check-http", "--url", url)
	cmd.Stdout = outStream
	cmd.Stderr = errorStream
	err := cmd.Run()
	return err
}

type URLCheckResult struct {
	outStream   io.Writer
	errorStream io.Writer
	err         error
}

func (r *URLCheckResult) OK() bool {
	return r.err == nil
}
