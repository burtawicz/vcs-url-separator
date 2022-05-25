package vcs_url_separator

import (
	"errors"
	"regexp"
	"strings"
)

// UrlParts represents the components that make up a VCS url.
//
// Example
//
// 	separator.SeparateVcsUrl("https://github.com/burtawicz/vcs-url-separator/go.mod")
// should result in
// 	UrlParts{
// 		Provider:       "GitHub",
//		Organization:   "burtawicz",
//		Project:        "vcs-url-separator",
//		SubDirectories: []string{},
//		FilePath:       "go.mod",
// 	}
type UrlParts struct {
	Provider       string
	Organization   string
	Project        string
	SubDirectories []string
	FilePath       string
}

var illegalCharPattern = regexp.MustCompile(`[\[\]{}\\|"%~#<>]+`)
var topLevelDomainPattern = regexp.MustCompile(`([a-zA-Z0-9\-]{3,})(\.[a-zA-Z0-9\-]+)`)
var fileNamePattern = regexp.MustCompile(`[a-zA-Z0-9\-_]+\.[a-zA-Z0-9]+`)

// stripHttpPrefix removes the `http://` or `https://` prefix from the url.
// If the prefix is not found, the url is returned as is.
func stripHttpPrefix(url string) string {
	if strings.HasPrefix(url, "http://") {
		return strings.Split(url, "http://")[1]
	}

	if strings.HasPrefix(url, "https://") {
		return strings.Split(url, "https://")[1]
	}

	return url
}

// stripTopLevelDomain tries to match the providerUrl for a top level domain
// If found, the remainder of the provider's domain is returned.
// Else, the providerUrl is returned as is.
func stripTopLevelDomain(providerUrl string) string {
	if topLevelDomainPattern.MatchString(providerUrl) {
		return topLevelDomainPattern.FindAllStringSubmatch(providerUrl, -1)[0][1]
	} else {
		return providerUrl
	}
}

// matchProvider tries to match the providerName for a known match.
// If found, the stylized form of the providerName is returned.
// Else, the providerName is returned as is.
func matchProvider(providerName string) string {
	switch providerName {
	case "github":
		return "GitHub"
	case "bitbucket":
		return "BitBucket"
	case "gitlab":
		return "GitLab"
	default:
		return providerName
	}
}

// SeparateVcsUrl isolates individual components of a VCS url.
// TODO: add a better function comment.
func SeparateVcsUrl(url string) (UrlParts, error) {
	// verify url is not empty
	if len(strings.TrimSpace(url)) < 1 {
		// FIXME: replace with proper error
		return UrlParts{}, errors.New("invalid string length")
	}

	// verify url does not contain illegal characters
	if illegalCharPattern.MatchString(url) {
		return UrlParts{}, errors.New("illegal characters included")
	}

	withoutHttp := stripHttpPrefix(url)

	parts := strings.Split(withoutHttp, "/")
	// verify there are at least 3 substrings after the split [provider/owner/project[subdirs?,]/[filepath?]
	if len(parts) <= 2 {
		return UrlParts{}, errors.New("invalid url, does not contain enough information")
	}
	provider := matchProvider(stripTopLevelDomain(parts[0]))
	owner := parts[1]
	project := parts[2]
	subDirs := make([]string, 0)
	filePath := ""

	if len(parts) > 3 {
		// collect the possible subDirs from [4:-1]
		for i := 3; i < len(parts)-1; i++ {
			// if the name matches a file path, we've made a mistake
			if fileNamePattern.MatchString(parts[i]) {
				return UrlParts{}, errors.New("multiple file paths")
			}
			subDirs = append(subDirs, parts[i])
		}

		potentialFilePath := parts[len(parts)-1]
		// verify potentialFilePath is not the same as the project name
		if potentialFilePath != project {
			if fileNamePattern.MatchString(potentialFilePath) {
				filePath = potentialFilePath
			} else {
				subDirs = append(subDirs, potentialFilePath)
			}
		}
	}

	return UrlParts{
		Provider:       provider,
		Organization:   owner,
		Project:        project,
		SubDirectories: subDirs,
		FilePath:       filePath,
	}, nil
}
