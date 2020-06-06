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
	"io"
	"time"

	"github.com/opencontainers/go-digest"
	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
)

type ImageMeta struct {
	Digest    digest.Digest
	Size      int64
	CreatedAt time.Time
	UpdatedAt time.Time
	Labels    map[string]string
}

type Status struct {
	Ref       string
	Offset    int64
	Total     int64
	Expected  digest.Digest
	StartedAt time.Time
	UpdatedAt time.Time
}

type Opt func(*ImageMeta) error

// Writer handles the write of content into a content store
type Writer interface {
	// WriteCloser is the interface that groups the basic Write and Close methods.
	io.WriteCloser

	// Digest may return empty image digest until committed.
	Digest() digest.Digest

	// Commit commits the blob (size and expected can be zero-value when unknown)
	Commit(ctx context.Context, size int64, expected digest.Digest, opts ...Opt) error

	// Status returns the current state of write
	Status() (Status, error)

	// Truncate updates the size of the target blob
	Truncate(size int64) error
}

type pusher struct {
	*base
	object  string
	tracker StatusTracker
}

func (p pusher) Push(ctx context.Context, descriptor ocispec.Descriptor) (Writer, error) {
	return nil, nil
}
