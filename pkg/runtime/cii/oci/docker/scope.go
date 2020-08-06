//
// Last.Backend LLC CONFIDENTIAL
// __________________
//
// [2014] - [2020] Last.Backend LLC
// All Rights Reserved.
//
// NOTICE:  All information contained herein is, and remains
// the property of Last.Backend LLC and its suppliers,
// if any.  The intellectual and technical concepts contained
// herein are proprietary to Last.Backend LLC
// and its suppliers and may be covered by Russian Federation and Foreign Patents,
// patents in process, and are protected by trade secret or copyright law.
// Dissemination of this information or reproduction of this material
// is strictly forbidden unless prior written permission is obtained
// from Last.Backend LLC.
//

package docker

import (
	"context"
	"fmt"
	"net/url"
	"sort"
	"strings"
)

// repositoryScope returns a repository scope string such as "repository:foo/bar:pull"
// for "host/foo/bar:baz".
// When push is true, both pull and push are added to the scope.
func repositoryScope(refspec Spec, push bool) (string, error) {
	u, err := url.Parse("dummy://" + refspec.Locator)
	if err != nil {
		return "", err
	}
	s := "repository:" + strings.TrimPrefix(u.Path, "/") + ":pull"
	if push {
		s += ",push"
	}
	return s, nil
}

// tokenScopesKey is used for the key for context.WithValue().
// value: []string (e.g. {"registry:foo/bar:pull"})
type tokenScopesKey struct{}

// contextWithScope returns a context with tokenScopesKey{} and the repository scope value.
func contextWithScope(ctx context.Context, refspec Spec, push bool) (context.Context, error) {
	s, err := repositoryScope(refspec, push)
	if err != nil {
		return nil, err
	}
	return WithScope(ctx, s), nil
}

// WithScope appends a custom registry auth scope to the context.
func WithScope(ctx context.Context, scope string) context.Context {
	var scopes []string
	if v := ctx.Value(tokenScopesKey{}); v != nil {
		scopes = v.([]string)
		scopes = append(scopes, scope)
	} else {
		scopes = []string{scope}
	}
	return context.WithValue(ctx, tokenScopesKey{}, scopes)
}

// contextWithAppendPullRepositoryScope is used to append repository pull
// scope into existing scopes indexed by the tokenScopesKey{}.
func contextWithAppendPullRepositoryScope(ctx context.Context, repo string) context.Context {
	return WithScope(ctx, fmt.Sprintf("repository:%s:pull", repo))
}

// getTokenScopes returns deduplicated and sorted scopes from ctx.Value(tokenScopesKey{}) and common scopes.
func getTokenScopes(ctx context.Context, common []string) []string {
	var scopes []string
	if x := ctx.Value(tokenScopesKey{}); x != nil {
		scopes = append(scopes, x.([]string)...)
	}

	scopes = append(scopes, common...)
	sort.Strings(scopes)

	l := 0
	for idx := 1; idx < len(scopes); idx++ {
		// Note: this comparison is unaware of the scope grammar (https://docs.docker.com/registry/spec/auth/scope/)
		// So, "repository:foo/bar:pull,push" != "repository:foo/bar:push,pull", although semantically they are equal.
		if scopes[l] == scopes[idx] {
			continue
		}

		l++
		scopes[l] = scopes[idx]
	}
	return scopes[:l+1]
}
