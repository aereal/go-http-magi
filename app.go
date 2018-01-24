package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"runtime"
	"strings"
	"sync"

	"github.com/mackerelio/checkers"
)

type AccumulationOp int

const (
	OP_OR AccumulationOp = iota
	OP_AND
)

func (f *AccumulationOp) String() string {
	return string(*f)
}

func (f *AccumulationOp) Set(value string) error {
	switch value {
	case "or":
		*f = OP_OR
		return nil
	case "and":
		*f = OP_AND
		return nil
	default:
		return fmt.Errorf("Unknown value: %s", value)
	}
}

// App represents execution context of CLI application.
type App struct {
	site           *Site
	maxConcurrency int
	accumulationOp AccumulationOp
}

type urlList []string

func (f *urlList) String() string {
	return strings.Join(*f, " ")
}

func (f *urlList) Set(value string) error {
	*f = append(*f, value)
	return nil
}

func newApp(args []string, outStream, errorStream io.Writer) (*App, error) {
	var (
		siteName string
		urls     urlList
		op       AccumulationOp
	)
	flgs := flag.NewFlagSet("magi", flag.ContinueOnError)
	flgs.StringVar(&siteName, "name", "", "site name")
	flgs.Var(&urls, "url", "URLs")
	flgs.Var(&op, "op", "operation")
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
	app.maxConcurrency = runtime.NumCPU()
	app.accumulationOp = op
	runtime.GOMAXPROCS(app.maxConcurrency)
	app.site = &Site{
		name: siteName,
		urls: urls,
	}
	return app, nil
}

func (a *App) checkURLs() *sync.Map {
	semaphore := make(chan int, a.maxConcurrency)
	var wg sync.WaitGroup
	checkResults := new(sync.Map)
	for _, url := range a.site.urls {
		wg.Add(1)
		go func(url string) {
			defer wg.Done()
			semaphore <- 1

			out := new(bytes.Buffer)
			errOut := new(bytes.Buffer)
			if status, err := runCheckHTTP(url, out, errOut); err != nil {
				res := &URLCheckResult{
					outStream:   out,
					errorStream: errOut,
					err:         err,
					status:      status,
				}
				checkResults.Store(url, res)
			} else {
				res := &URLCheckResult{
					outStream:   out,
					errorStream: errOut,
					err:         nil,
					status:      status,
				}
				checkResults.Store(url, res)
			}

			<-semaphore
		}(url)
	}
	wg.Wait()
	return checkResults
}

func (a *App) accumulateResults(results *sync.Map) *SiteCheckResult {
	result := &SiteCheckResult{
		urlResults: make(map[string]*URLCheckResult),
		statusCode: 0,
	}
	results.Range(func(key interface{}, value interface{}) bool {
		var (
			url       string
			urlResult *URLCheckResult
		)
		if k, casted := key.(string); casted {
			url = k
		} else {
			return true
		}
		if res, casted := value.(*URLCheckResult); casted {
			urlResult = res
		} else {
			return true
		}
		result.urlResults[url] = urlResult
		switch a.accumulationOp {
		case OP_AND:
			if result.statusCode < int(urlResult.status) {
				result.statusCode = int(urlResult.status)
			}
		case OP_OR:
			if urlResult.status == checkers.OK {
				result.statusCode = int(checkers.OK)
			}
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
