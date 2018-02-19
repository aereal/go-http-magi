# go-http-magi

Run [check-http][] concurrently and accumulate results without false-positive.

# Motivation

- primary URL: The URL that we want to check
  - Serving **our** environment such as CMS
- secondary URL: The URL that shares HTTP serving environment such as proxy with primary URL
  - Serving **their** environment

Sometimes we may get check failure on primary URL but actually that is caused by their enironment issues, so we want to do triage such errors.

# Implementations

Set `*SiteCheckResult.status` to `checkers.WARNING` (or else) if only primary URL is failed.

primary: `http://example.com/subdir/` | secondary:`http://example.com/` | Total status
------------ | ------------- | -------------
OK | OK | OK
OK | NG | OK
NG | OK | NG
NG | NG | OK

[check-http]: https://github.com/mackerelio/go-check-plugins/tree/master/check-http
