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
	"io"
	"io/ioutil"
	"net/http"
	"path"
	"strings"

	"github.com/lastbackend/lastbackend/version"
	"github.com/opencontainers/go-digest"
	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
	"github.com/pkg/errors"
)

const (
	// MaxManifestSize represents the largest size accepted from a registry during resolution.
	MaxManifestSize int64 = 4 * 1048 * 1048
)

type Resolver interface {
	Pull(ctx context.Context, ref string) (Reader, error)
	Push(ctx context.Context, ref string) (Writer, error)
}

type Config struct {
	Debug         bool
	Hosts         RegistryHosts
	Headers       http.Header
	ResolveHeader http.Header
	Tracker       StatusTracker
}

func New(cfg Config) (Resolver, error) {
	d := new(resolver)
	d.debug = cfg.Debug
	if cfg.Tracker == nil {
		cfg.Tracker = NewInMemoryTracker()
	}

	if cfg.Headers == nil {
		cfg.Headers = make(http.Header)
	}
	if _, ok := cfg.Headers["User-Agent"]; !ok {
		cfg.Headers.Set("User-Agent", "lastbackend/"+version.Version)
	}

	resolveHeader := http.Header{}
	if _, ok := cfg.Headers["Accept"]; !ok {
		// set headers for all the types we support for resolution.
		resolveHeader.Set("Accept", strings.Join([]string{
			MediaTypeDockerSchema2Manifest,
			MediaTypeDockerSchema2ManifestList,
			ocispec.MediaTypeImageManifest,
			ocispec.MediaTypeImageIndex, "*/*"}, ", "))
	} else {
		resolveHeader["Accept"] = cfg.Headers["Accept"]
		delete(cfg.Headers, "Accept")
	}

	return &resolver{
		hosts:         cfg.Hosts,
		header:        cfg.Headers,
		resolveHeader: resolveHeader,
		tracker:       cfg.Tracker,
	}, nil
}

type resolver struct {
	debug         bool
	refspec       Spec
	namespace     string
	hosts         RegistryHosts
	header        http.Header
	resolveHeader http.Header
	tracker       StatusTracker
}

func (r *resolver) Pull(ctx context.Context, ref string) (Reader, error) {
	name, desc, err := r.resolve(ctx, ref)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to resolve reference %q", ref)
	}
	spec, err := Parse(name)
	if err != nil {
		return nil, err
	}
	base, err := r.base(spec)
	if err != nil {
		return nil, err
	}
	p := puller{
		base: base,
	}
	return p.pull(ctx, desc)
}

func (r *resolver) Push(ctx context.Context, ref string) (Writer, error) {
	name, desc, err := r.resolve(ctx, ref)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to resolve reference %q", ref)
	}
	spec, err := Parse(name)
	if err != nil {
		return nil, err
	}
	base, err := r.base(spec)
	if err != nil {
		return nil, err
	}
	p := pusher{
		base:    base,
		object:  spec.Object,
		tracker: r.tracker,
	}
	return p.Push(ctx, desc)
}

func (r *resolver) base(spec Spec) (*base, error) {
	host := spec.Hostname()
	hosts, err := r.hosts(host)
	if err != nil {
		return nil, err
	}
	return &base{
		spec:      spec,
		namespace: strings.TrimPrefix(spec.Locator, host+"/"),
		hosts:     hosts,
		header:    r.header,
	}, nil
}

func (r *resolver) resolve(ctx context.Context, ref string) (string, ocispec.Descriptor, error) {
	spec, err := Parse(ref)
	if err != nil {
		return "", ocispec.Descriptor{}, err
	}

	if spec.Object == "" {
		return "", ocispec.Descriptor{}, ErrObjectRequired
	}

	base, err := r.base(spec)
	if err != nil {
		return "", ocispec.Descriptor{}, err
	}

	var (
		saveLastError error
		paths         [][]string
		dgst          = spec.Digest()
		grant         = GrantPull
	)

	if dgst != "" {
		if err := dgst.Validate(); err != nil {
			return "", ocispec.Descriptor{}, err
		}
		paths = append(paths, []string{"manifests", dgst.String()})
		paths = append(paths, []string{"blobs", dgst.String()})
	} else {
		paths = append(paths, []string{"manifests", spec.Object})
		grant |= GrantResolve
	}

	hosts := base.filterRegistryHosts(grant)
	if len(hosts) == 0 {
		return "", ocispec.Descriptor{}, errors.Wrap(ErrNotFound, "no resolve hosts")
	}

	ctx, err = contextWithScope(ctx, spec, false)
	if err != nil {
		return "", ocispec.Descriptor{}, err
	}

	for _, u := range paths {
		for _, host := range hosts {

			req := base.request(http.MethodHead, host, u...)
			for key, value := range r.resolveHeader {
				req.header[key] = append(req.header[key], value...)
			}

			resp, err := req.doWithRetries(ctx, nil)
			if err != nil {
				if errors.Is(err, ErrInvalidAuthorization) {
					err = errors.Wrapf(err, "pull access denied, repository does not exist or may require authorization")
				}
				if saveLastError == nil {
					saveLastError = err
				}
				continue
			}
			resp.Body.Close()

			if resp.StatusCode > 299 {
				if resp.StatusCode == http.StatusNotFound {
					continue
				}
				return "", ocispec.Descriptor{}, errors.Errorf("unexpected status code %v: %v", u, resp.Status)
			}
			size := resp.ContentLength
			contentType := getManifestMediaType(resp)

			if dgst == "" {
				digestHeader := digest.Digest(resp.Header.Get("Docker-Content-Digest"))

				if digestHeader != "" && size != -1 {
					if err := digestHeader.Validate(); err != nil {
						return "", ocispec.Descriptor{}, errors.Wrapf(err, "%q in header not a valid dgst", digestHeader)
					}
					dgst = digestHeader
				}
			}
			if dgst == "" || size == -1 {
				req = base.request(http.MethodGet, host, u...)
				for key, value := range r.resolveHeader {
					req.header[key] = append(req.header[key], value...)
				}

				resp, err := req.doWithRetries(ctx, nil)
				if err != nil {
					return "", ocispec.Descriptor{}, err
				}
				defer resp.Body.Close()

				bodyReader := countingReader{reader: resp.Body}

				contentType = getManifestMediaType(resp)
				if dgst == "" {
					if contentType == MediaTypeDockerSchema1Manifest {
						b, err := ReadStripSignature(&bodyReader)
						if err != nil {
							return "", ocispec.Descriptor{}, err
						}

						dgst = digest.FromBytes(b)
					} else {
						dgst, err = digest.FromReader(&bodyReader)
						if err != nil {
							return "", ocispec.Descriptor{}, err
						}
					}
				} else if _, err := io.Copy(ioutil.Discard, &bodyReader); err != nil {
					return "", ocispec.Descriptor{}, err
				}
				size = bodyReader.bytesRead
			}

			if size > MaxManifestSize {
				if saveLastError == nil {
					saveLastError = errors.Wrapf(ErrNotFound, "rejecting %d byte manifest for %s", size, ref)
				}
				continue
			}

			desc := ocispec.Descriptor{
				Digest:    dgst,
				MediaType: contentType,
				Size:      size,
			}

			return ref, desc, nil
		}
	}

	if saveLastError == nil {
		saveLastError = errors.Wrap(ErrNotFound, ref)
	}

	return "", ocispec.Descriptor{}, saveLastError
}

type base struct {
	debug     bool
	namespace string
	spec      Spec
	header    http.Header
	hosts     []RegistryHost
}

func (b *base) log(format string, a ...interface{}) {
	if b.debug {
		fmt.Println(fmt.Sprintf(format, a))
	}
}

func (b *base) filterRegistryHosts(grants RegistryGrant) []RegistryHost {
	var hosts = make([]RegistryHost, 0)
	for _, host := range b.hosts {
		if host.Grant.Has(grants) {
			hosts = append(hosts, host)
		}
	}
	return hosts
}

func (b *base) request(method string, host RegistryHost, ps ...string) *request {
	header := http.Header{}
	for key, value := range b.header {
		header[key] = append(header[key], value...)
	}
	parts := append([]string{"/", host.Path, b.namespace}, ps...)
	p := path.Join(parts...)
	// Join strips trailing slash, re-add ending "/" if included
	if len(parts) > 0 && strings.HasSuffix(parts[len(parts)-1], "/") {
		p = p + "/"
	}
	return &request{
		method: method,
		path:   p,
		header: header,
		host:   host,
	}
}

func getManifestMediaType(resp *http.Response) string {
	// Strip encoding data (manifests should always be ascii JSON)
	contentType := resp.Header.Get("Content-Type")
	if sp := strings.IndexByte(contentType, ';'); sp != -1 {
		contentType = contentType[0:sp]
	}
	// As of Apr 30 2019 the registry.access.redhat.com registry does not specify
	// the content type of any data but uses schema1 manifests.
	if contentType == "text/plain" {
		contentType = MediaTypeDockerSchema1Manifest
	}
	return contentType
}

type countingReader struct {
	reader    io.Reader
	bytesRead int64
}

func (r *countingReader) Read(p []byte) (int, error) {
	n, err := r.reader.Read(p)
	r.bytesRead += int64(n)
	return n, err
}
