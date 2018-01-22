package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"strings"
	"sync"
)

type App struct {
	site *Site
}

type URLList []string

func (f *URLList) String() string {
	return strings.Join(*f, " ")
}

func (f *URLList) Set(value string) error {
	*f = append(*f, value)
	return nil
}

func newApp(args []string, outStream, errorStream io.Writer) (*App, error) {
	var (
		siteName string
		urls     URLList
	)
	flgs := flag.NewFlagSet("magi", flag.ContinueOnError)
	flgs.StringVar(&siteName, "name", "", "site name")
	flgs.Var(&urls, "url", "URLs")
	flgs.SetOutput(errorStream)
	if err := flgs.Parse(args[1:]); err != nil {
		return nil, err
	}

	if siteName == "" {
		return nil, fmt.Errorf("name required")
	}

	if len(urls) == 0 {
		return nil, fmt.Errorf("URLs required")
	}

	app := new(App)
	app.site = &Site{
		name: siteName,
		urls: urls,
	}
	return app, nil
}

func (a *App) checkURLs() *sync.Map {
	checkResults := new(sync.Map)
	for _, url := range a.site.urls {
		out := new(bytes.Buffer)
		errOut := new(bytes.Buffer)
		if err := runCheckHTTP(url, out, errOut); err != nil {
			res := &URLCheckResult{
				outStream:   out,
				errorStream: errOut,
				err:         err,
			}
			checkResults.Store(url, res)
		} else {
			res := &URLCheckResult{
				outStream:   out,
				errorStream: errOut,
				err:         nil,
			}
			checkResults.Store(url, res)
		}
	}
	return checkResults
}

func (a *App) accumulateResults(results *sync.Map) *SiteCheckResult {
	result := &SiteCheckResult{
		ok:         true,
		urlResults: make(map[string]*URLCheckResult),
	}
	results.Range(func(key interface{}, value interface{}) bool {
		var (
			ok        bool
			url       string
			urlResult *URLCheckResult
		)
		if url, ok = key.(string); !ok {
			return true
		}
		if urlResult, ok = value.(*URLCheckResult); !ok {
			return true
		}
		result.urlResults[url] = urlResult
		if !urlResult.OK() {
			result.ok = false
		}
		return true
	})
	return result
}

func (a *App) run() *SiteCheckResult {
	urlResults := a.checkURLs()
	siteResult := a.accumulateResults(urlResults)
	return siteResult
}
