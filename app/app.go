package app

import (
	"http-response-hasher/hasher"
	p "http-response-hasher/processor"
	"net/http"
	"net/url"
	"strings"
)

func parseUrl(rawUrl string, defaultScheme string) (*url.URL, error) {
	if strings.Index(rawUrl, "//") == 0 {
		rawUrl = defaultScheme + ":" + rawUrl
	}
	if !strings.Contains(rawUrl, "://") {
		rawUrl = defaultScheme + "://" + rawUrl
	}

	return url.ParseRequestURI(rawUrl)
}

// requestUrlAndHashResponse performs a GET request to url and returns
// a processor.Result instance, which Input is url and Output is an MD5
// hash of the response body.
func requestUrlAndHashResponse(value p.Value) p.Result {
	rawUrl := string(value)
	url, err := parseUrl(rawUrl, "http")
	if err != nil {
		return p.Result{
			Input:  p.Value(rawUrl),
			Output: p.NoValue,
			Error:  err,
		}
	}

	rawUrl = url.String()
	res, err := http.Get(rawUrl)
	if err != nil {
		return p.Result{
			Input:  p.Value(rawUrl),
			Output: p.NoValue,
			Error:  err,
		}
	}

	resHash, err := hasher.HashHttpResponseToString(res)
	return p.Result{
		Input:  p.Value(rawUrl),
		Output: p.Value(resHash),
		Error:  err,
	}
}

// ProcessUrls performs GET requests to urls in parallel and returns
// an unordered generator of process.Result for corresponding responses.
func ProcessUrls(urls []string, nWorkers uint) (<-chan p.Result, error) {
	values := make(chan p.Value)
	go func() {
		for _, url := range urls {
			values <- p.Value(url)
		}
		close(values)
	}()

	return p.Process(values, requestUrlAndHashResponse, nWorkers, 1)
}
