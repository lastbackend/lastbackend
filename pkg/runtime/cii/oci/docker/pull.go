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
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
	"github.com/pkg/errors"
)

// Reader handles the read of content from a content store
type Reader interface {
	// ReadCloser is the interface that groups the basic Read and Close methods.
	io.ReadCloser
}

type puller struct {
	*base
}

func (p puller) pull(ctx context.Context, descriptor ocispec.Descriptor) (Reader, error) {
	p.log("fetch for digest %s", descriptor.Digest)

	ctx, err := contextWithScope(ctx, p.spec, false)
	if err != nil {
		return nil, err
	}

	return &reader{
		size: descriptor.Size,
		read: func(offset int64) (io.ReadCloser, error) {
			if len(descriptor.URLs) != 0 {
				rc, err := p.fetchFromUrls(ctx, descriptor, offset)
				if err != nil {
					p.log("fetch from specifies urls err: %v", err)
					return nil, err
				}
				return rc, nil
			}

			hosts := p.base.filterRegistryHosts(GrantPull)
			if len(hosts) == 0 {
				return nil, errors.Wrap(ErrNotFound, "no pull RegistryHost hosts")
			}

			// Download manifest
			switch descriptor.MediaType {
			case MediaTypeDockerSchema2Manifest, MediaTypeDockerSchema2ManifestList, MediaTypeDockerSchema1Manifest,
				ocispec.MediaTypeImageManifest, ocispec.MediaTypeImageIndex:
				rc, err := p.fetchManifestsForMediaType(ctx, descriptor, offset)
				if err != nil {
					p.log("fetch manifest for media type err: %v", err)
					return nil, err
				}
				return rc, nil
			}

			// Download blob
			rc, err := p.fetchBlob(ctx, descriptor, offset)
			if err != nil {
				p.log("fetch blob err: %v", err)
				return nil, err
			}
			return rc, nil
		},
	}, nil
}

func (p *puller) fetchFromUrls(ctx context.Context, descriptor ocispec.Descriptor, offset int64) (io.ReadCloser, error) {
	var saveFirstError error
	for _, rawUrl := range descriptor.URLs {
		u, err := url.Parse(rawUrl)
		if err != nil {
			return nil, err
		}

		registry := newRegistry(http.DefaultClient, *u, GrantPull)
		req := p.request(http.MethodGet, registry)
		req.path = u.Path

		if u.RawQuery != "" {
			req.path = fmt.Sprintf("%s?%s", req.path, u.RawQuery)
		}
		rc, err := p.fetch(ctx, req, descriptor.MediaType, offset)
		if err != nil {
			if saveFirstError == nil {
				saveFirstError = err
			}
			continue
		}
		return rc, nil
	}
	if IsNotFound(saveFirstError) {
		saveFirstError = errors.Wrapf(ErrNotFound, "could not fetch content descriptor %v (%v) from remote", descriptor.Digest, descriptor.MediaType)
	}
	return nil, saveFirstError
}

func (p *puller) fetchManifestsForMediaType(ctx context.Context, descriptor ocispec.Descriptor, offset int64) (io.ReadCloser, error) {
	var saveFirstError error
	for _, host := range p.hosts {
		req := p.request(http.MethodGet, host, "manifests", descriptor.Digest.String())
		rc, err := p.fetch(ctx, req, descriptor.MediaType, offset)
		if err != nil {
			if saveFirstError == nil {
				saveFirstError = err
			}
			continue
		}
		return rc, nil
	}

	if IsNotFound(saveFirstError) {
		saveFirstError = errors.Wrapf(ErrNotFound, "could not fetch content descriptor %v (%v) from remote", descriptor.Digest, descriptor.MediaType)
	}

	return nil, saveFirstError
}

func (p *puller) fetchBlob(ctx context.Context, desc ocispec.Descriptor, offset int64) (io.ReadCloser, error) {
	var saveFirstError error
	for _, host := range p.hosts {
		req := p.request(http.MethodGet, host, "blobs", desc.Digest.String())
		rc, err := p.fetch(ctx, req, desc.MediaType, offset)
		if err != nil {
			if saveFirstError == nil {
				saveFirstError = err
			}
			continue
		}
		return rc, nil
	}

	if IsNotFound(saveFirstError) {
		saveFirstError = errors.Wrapf(ErrNotFound, "could not fetch content descriptor %v (%v) from remote", desc.Digest, desc.MediaType)
	}

	return nil, saveFirstError
}

func (puller) fetch(ctx context.Context, req *request, mediaType string, offset int64) (io.ReadCloser, error) {

	// Add the data for the Accept header of the form
	req.header.Set("Accept", strings.Join([]string{mediaType, `*/*`}, ", "))

	if offset > 0 {
		// Add the data for the Range header of the form
		rangeHeader := fmt.Sprintf("bytes=%d-", offset)
		req.header.Set("Range", rangeHeader)
	}

	resp, err := req.doWithRetries(ctx, nil)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode > 299 {
		defer resp.Body.Close()
		if resp.StatusCode == http.StatusNotFound {
			return nil, errors.Wrapf(ErrNotFound, "content at %v not found", req.String())
		}
		return nil, handleErrorResponse(resp)
	}

	if offset > 0 {
		cr := resp.Header.Get("content-range")
		if cr != "" {
			rangeHeader := fmt.Sprintf("bytes %d-", offset)
			if !strings.HasPrefix(cr, rangeHeader) {
				return nil, errors.Errorf("unhandled content range in response: %v", cr)
			}
		} else {
			// Discard up to offset
			// Could use buffer pool here but this case should be rare
			n, err := io.Copy(ioutil.Discard, io.LimitReader(resp.Body, offset))
			if err != nil {
				return nil, errors.Wrap(err, "failed to discard to offset")
			}
			if n != offset {
				return nil, errors.Errorf("unable to discard to offset")
			}

		}
	}

	return resp.Body, nil
}

type reader struct {
	size   int64
	offset int64
	rc     io.ReadCloser
	read   func(offset int64) (io.ReadCloser, error)
	closed bool
}

func (r *reader) Read(p []byte) (n int, err error) {
	if r.closed {
		return 0, io.EOF
	}

	rd, err := r.reader()
	if err != nil {
		return 0, err
	}

	n, err = rd.Read(p)
	r.offset += int64(n)
	return
}

func (r *reader) Close() error {
	if r.closed {
		return nil
	}
	r.closed = true
	if r.rc != nil {
		return r.rc.Close()
	}

	return nil
}

func (r *reader) Seek(offset int64, whence int) (int64, error) {
	if r.closed {
		return 0, errors.Wrap(ErrUnavailable, "closed")
	}

	abs := r.offset
	switch whence {
	case io.SeekStart:
		abs = offset
	case io.SeekCurrent:
		abs += offset
	case io.SeekEnd:
		if r.size == -1 {
			return 0, errors.Wrap(ErrUnavailable, "unknown size, cannot seek from end")
		}
		abs = r.size + offset
	default:
		return 0, errors.Wrap(ErrInvalidArgument, "invalid whence")
	}

	if abs < 0 {
		return 0, errors.Wrapf(ErrInvalidArgument, "negative offset")
	}

	if abs != r.offset {
		if r.rc != nil {
			if err := r.rc.Close(); err != nil {
				fmt.Println("failed to close ReadCloser")
			}

			r.rc = nil
		}

		r.offset = abs
	}

	return r.offset, nil
}

func (r *reader) reader() (io.Reader, error) {
	if r.rc != nil {
		return r.rc, nil
	}

	if r.size == -1 || r.offset < r.size {
		// only try to reopen the body request if we are seeking to a value
		// less than the actual size.
		if r.read == nil {
			return nil, errors.Wrapf(ErrNotImplemented, "cannot fetch")
		}

		rc, err := r.read(r.offset)
		if err != nil {
			return nil, errors.Wrapf(err, "httpReaderSeeker: failed fetch")
		}

		if r.rc != nil {
			if err := r.rc.Close(); err != nil {
				fmt.Println("failed to close ReadCloser")
			}
		}
		r.rc = rc
	} else {
		// There is an edge case here where offset == size of the content. If
		// we seek, we will probably get an error for content that cannot be
		// sought (?). In that case, we should err on committing the content,
		// as the length is already satisfied but we just return the empty
		// reader instead.

		r.rc = ioutil.NopCloser(bytes.NewReader([]byte{}))
	}

	return r.rc, nil
}
