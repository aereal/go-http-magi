package main

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
