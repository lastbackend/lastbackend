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

package reference

import (
	"errors"
	"fmt"
	"net/url"
	"path"
	"regexp"
	"strings"

	"github.com/opencontainers/go-digest"
)

var (
	ErrInvalid          = errors.New("invalid reference")
	ErrObjectRequired   = errors.New("object required")
	ErrHostnameRequired = errors.New("hostname required")
)

type Spec struct {
	Locator string
	Object  string
}

var splitRe = regexp.MustCompile(`[:@]`)

func Parse(s string) (Spec, error) {
	u, err := url.Parse("dummy://" + s)
	if err != nil {
		return Spec{}, err
	}

	if u.Scheme != "dummy" {
		return Spec{}, ErrInvalid
	}

	if u.Host == "" {
		return Spec{}, ErrHostnameRequired
	}

	var object string

	if idx := splitRe.FindStringIndex(u.Path); idx != nil {
		// This allows us to retain the @ to signify digests or shortened digests in
		// the object.
		object = u.Path[idx[0]:]
		if object[:1] == ":" {
			object = object[1:]
		}
		u.Path = u.Path[:idx[0]]
	}

	return Spec{
		Locator: path.Join(u.Host, u.Path),
		Object:  object,
	}, nil
}

func (r Spec) Hostname() string {
	i := strings.Index(r.Locator, "/")

	if i < 0 {
		return r.Locator
	}
	return r.Locator[:i]
}

func (r Spec) Digest() digest.Digest {
	_, dgst := SplitObject(r.Object)
	return dgst
}

func (r Spec) String() string {
	if r.Object == "" {
		return r.Locator
	}
	if r.Object[:1] == "@" {
		return fmt.Sprintf("%v%v", r.Locator, r.Object)
	}

	return fmt.Sprintf("%v:%v", r.Locator, r.Object)
}

func SplitObject(obj string) (tag string, dgst digest.Digest) {
	parts := strings.SplitAfterN(obj, "@", 2)
	if len(parts) < 2 {
		return parts[0], ""
	}
	return parts[0], digest.Digest(parts[1])
}
