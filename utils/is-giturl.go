// Based on https://github.com/go-git/go-git/blob/master/internal/url/url.go

package utils

import (
	"regexp"
)

var (
	isSchemeRegExp = regexp.MustCompile(`^[^:]+://`)

	// Ref: https://github.com/git/git/blob/master/Documentation/urls.txt#L37
	scpLikeUrlRegExp = regexp.MustCompile(`^(?:(?P<user>[^@]+)@)?(?P<host>[^:\s]+):(?:(?P<port>[0-9]{1,5}):)?(?P<path>[^\\].*)$`)
)

// matchesScheme returns true if the given string matches a URL-like
// format scheme.
func matchesScheme(url string) bool {
	return isSchemeRegExp.MatchString(url)
}

// matchesScpLike returns true if the given string matches an SCP-like
// format scheme.
func matchesScpLike(url string) bool {
	return scpLikeUrlRegExp.MatchString(url)
}

// isLocalEndpoint returns true if the given URL string specifies a
// local file endpoint.  For example, on a Linux machine,
// `/home/user/src/go-git` would match as a local endpoint, but
// `https://github.com/src-d/go-git` would not.
func isLocalEndpoint(url string) bool {
	return !matchesScheme(url) && !matchesScpLike(url)
}

// IsGitUrl returns true if the given URL string specifies a git repository.
func IsGitUrl(url string) bool {
	return matchesScheme(url) || matchesScpLike(url) || isLocalEndpoint(url)
}

// FindScpLikeComponents returns the user, host, port and path of the
// given SCP-like URL.
func FindScpLikeComponents(url string) (user, host, port, path string) {
	m := scpLikeUrlRegExp.FindStringSubmatch(url)
	return m[1], m[2], m[3], m[4]
}
