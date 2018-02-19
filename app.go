package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"runtime"
	"sync"

	"github.com/mackerelio/checkers"
)

// App represents execution context of CLI application.
type App struct {
	site           *Site
	maxConcurrency int
}

func newApp(args []string, outStream, errorStream io.Writer) (*App, error) {
	var (
		configFilePath string
	)
	flgs := flag.NewFlagSet("magi", flag.ContinueOnError)
	flgs.StringVar(&configFilePath, "config", "", "config file path")
	flgs.SetOutput(errorStream)
	if err := flgs.Parse(args[1:]); err != nil {
		return nil, err
	}

	if configFilePath == "" {
		return nil, fmt.Errorf("config required")
	}
	site, err := parseConfigFile(configFilePath)
	if err != nil {
		return nil, err
	}
	if err := validateSiteConfig(site); err != nil {
		return nil, err
	}

	app := new(App)
	app.maxConcurrency = runtime.NumCPU()
	runtime.GOMAXPROCS(app.maxConcurrency)
	app.site = site
	return app, nil
}

func (a *App) checkURLs() *sync.Map {
	semaphore := make(chan int, a.maxConcurrency)
	var wg sync.WaitGroup
	checkResults := new(sync.Map)
	urls := []string{a.site.PrimaryURL, a.site.SecondaryURL}
	for _, url := range urls {
		wg.Add(1)
		go func(url string) {
			defer wg.Done()
			semaphore <- 1

			out := new(bytes.Buffer)
			errOut := new(bytes.Buffer)
			if status, err := runCheckHTTP(url, out, errOut); err != nil {
				res := &URLCheckResult{
					err:    err,
					status: status,
				}
				checkResults.Store(url, res)
			} else {
				res := &URLCheckResult{
					err:    nil,
					status: status,
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
		statusCode: int(checkers.OK),
	}
	var primaryURLResult *URLCheckResult
	var secondaryURLResult *URLCheckResult
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
		switch url {
		case a.site.PrimaryURL:
			primaryURLResult = urlResult
		case a.site.SecondaryURL:
			secondaryURLResult = urlResult
		}
		return true
	})
	if (primaryURLResult.status != checkers.OK) && (secondaryURLResult.status == checkers.OK) {
		result.statusCode = int(primaryURLResult.status)
	}
	return result
}

func (a *App) run() *SiteCheckResult {
	urlResults := a.checkURLs()
	siteResult := a.accumulateResults(urlResults)
	return siteResult
}

func parseConfigFile(configPath string) (*Site, error) {
	var err error
	configFile, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}
	jsonIn := json.NewDecoder(configFile)
	var site *Site
	err = jsonIn.Decode(&site)
	if err != nil {
		return nil, err
	}
	return site, nil
}

func validateSiteConfig(site *Site) error {
	if site.Name == "" {
		return fmt.Errorf("name required")
	}
	if site.PrimaryURL == "" {
		return fmt.Errorf("primary_url required")
	}
	if site.SecondaryURL == "" {
		return fmt.Errorf("secondary_url required")
	}
	return nil
}
