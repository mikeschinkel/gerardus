package persister

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

func panicf(msg string, args ...any) {
	panic(fmt.Sprintf(msg, args...))
}
func debugBreakpointHere(...any) {
	// just a function for debugging
}

func ComposeCodebaseSourceURL(repoURL, versionTag string) (url string, err error) {
	url = fmt.Sprintf(`%s/tree/%s/src`, repoURL, versionTag)
	_, err = checkURL(url)
	if err != nil {
		goto end
	}
end:
	return url, err
}

func Close(c io.Closer, f func(err error)) {
	f(c.Close())
}

func WarnOnError(err error) {
	if err != nil {
		_, _ = fmt.Fprint(os.Stderr, err.Error())
	}
}

func checkURL(url string) (ok bool, err error) {
	var resp *http.Response
	var status int

	resp, err = http.Get(url)
	if err != nil {
		goto end
	}
	defer Close(resp.Body, WarnOnError) // Make sure to close the response body.

	status = resp.StatusCode
	if status < 200 || status > 299 {
		err = fmt.Errorf("received HTTP status code %d from %s", status, url)
		goto end
	}
	ok = true
end:
	return ok, err
}

func checkGitHubRepoURL(repoURL string) (parts []string, err error) {
	if repoURL == "" {
		err = ErrValueCannotBeEmpty.Args("which_value", RepoURLArg)
		goto end
	}
	if repoURL == "." {
		err = ErrInvalidGitHubRepoURL.Args(RepoURLArg, repoURL)
		goto end
	}
	parts = strings.Split(repoURL, "/")
	if len(parts) < 5 {
		err = ErrInvalidGitHubRepoURL.Args(RepoURLArg, repoURL)
		goto end
	}
end:
	return parts, err
}

func RequestGitHubRepoInfo(repoURL string) (info *RepoInfo, err error) {
	var body []byte
	var owner, repo, apiURL string
	var resp *http.Response

	info = &RepoInfo{}

	parts, err := checkGitHubRepoURL(repoURL)
	if err != nil {
		goto end
	}
	owner = parts[3]
	repo = parts[4]

	apiURL = fmt.Sprintf("https://api.github.com/repos/%s/%s", owner, repo)

	// Make the HTTP GET request
	resp, err = http.Get(apiURL)
	if err != nil {
		err = ErrHTTPRequestFailed.Err(err,
			"status_code", resp.StatusCode,
			"request_url", apiURL,
		)
		goto end
	}
	defer Close(resp.Body, WarnOnError)

	// Read the response body
	body, err = io.ReadAll(resp.Body)
	if err != nil {
		err = fmt.Errorf("failed to read the response body; %w", err)
		goto end
	}

	// Parse the JSON response
	err = json.Unmarshal(body, &info)
	if err != nil {
		err = fmt.Errorf("failed to unmarshal the JSON response; %w", err)
		goto end
	}
end:
	return info, err
}
